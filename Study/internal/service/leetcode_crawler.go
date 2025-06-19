package service

import (
	"Training/Study/internal/domain"
	"Training/Study/internal/repository"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"
)

// LeetCodeCrawler LeetCode爬虫服务
type LeetCodeCrawler struct {
	client     *http.Client
	repository repository.CodingProblemRepository
	logger     *log.Logger
}

// LeetCodeAPI相关的数据结构
type LeetCodeDailyResponse struct {
	Data struct {
		TodayRecord []struct {
			Question struct {
				QuestionFrontendId string `json:"questionFrontendId"`
				QuestionTitleSlug  string `json:"questionTitleSlug"`
				Title              string `json:"title"`
				TranslatedTitle    string `json:"translatedTitle"`
				Difficulty         string `json:"difficulty"`
			} `json:"question"`
		} `json:"todayRecord"`
	} `json:"data"`
}

type LeetCodeProblemResponse struct {
	Data struct {
		Question struct {
			QuestionId         string `json:"questionId"`
			QuestionFrontendId string `json:"questionFrontendId"`
			Title              string `json:"title"`
			TitleSlug          string `json:"titleSlug"`
			Content            string `json:"content"`
			Difficulty         string `json:"difficulty"`
			TopicTags          []struct {
				Name string `json:"name"`
				Slug string `json:"slug"`
			} `json:"topicTags"`
		} `json:"question"`
	} `json:"data"`
}

// NewLeetCodeCrawler 创建LeetCode爬虫
func NewLeetCodeCrawler(repository repository.CodingProblemRepository) *LeetCodeCrawler {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	logger := log.New(log.Writer(), "[LeetCode Crawler] ", log.LstdFlags)

	return &LeetCodeCrawler{
		client:     client,
		repository: repository,
		logger:     logger,
	}
}

// GetDailyProblem 获取今日每日一题（使用新的 GraphQL 查询）
func (c *LeetCodeCrawler) GetDailyProblem(ctx context.Context) (*domain.DailyProblem, error) {
	c.logger.Println("开始通过 GraphQL 接口爬取 LeetCode 每日一题...")

	url := "https://leetcode.cn/graphql/"
	payload := `{"operationName":"questionOfToday","query":"query questionOfToday { todayRecord { date userStatus question { questionFrontendId questionTitleSlug title translatedTitle difficulty } } }","variables":{}}`

	req, err := http.NewRequestWithContext(ctx, "POST", url, strings.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("构造请求失败: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求 GraphQL 接口失败: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	var result LeetCodeDailyResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("解析 JSON 失败: %w", err)
	}

	if len(result.Data.TodayRecord) == 0 {
		return nil, errors.New("未获取到每日一题")
	}

	q := result.Data.TodayRecord[0].Question
	title := c.fallbackTitle(q.TranslatedTitle, q.Title)

	// 获取题目详细信息
	problem, err := c.CrawlProblemBySlug(ctx, q.QuestionTitleSlug)
	if err != nil {
		c.logger.Printf("获取题目详情失败，使用基本信息: %v", err)
		// 如果获取详情失败，使用基本信息
		dailyProblem := &domain.DailyProblem{
			Date:       time.Now(),
			Title:      title,
			Difficulty: q.Difficulty,
			Tags:       []string{"算法"},
			Source:     "leetcode",
			SourceId:   q.QuestionFrontendId,
			SourceUrl:  fmt.Sprintf("https://leetcode.cn/problems/%s/", q.QuestionTitleSlug),
			Ctime:      time.Now(),
			Utime:      time.Now(),
		}
		return dailyProblem, nil
	}

	// 转换为 DailyProblem
	dailyProblem := &domain.DailyProblem{
		Date:       time.Now(),
		Title:      problem.Title,
		Difficulty: problem.Difficulty,
		Tags:       problem.Tags,
		Source:     problem.Source,
		SourceId:   problem.SourceId,
		SourceUrl:  problem.SourceUrl,
		Ctime:      time.Now(),
		Utime:      time.Now(),
	}

	c.logger.Printf("成功获取每日一题: %s (%s)", dailyProblem.Title, dailyProblem.SourceUrl)
	return dailyProblem, nil
}

// fallbackTitle 选择标题（优先使用翻译后的标题）
func (c *LeetCodeCrawler) fallbackTitle(translatedTitle, originalTitle string) string {
	if translatedTitle != "" {
		return translatedTitle
	}
	return originalTitle
}

// CrawlProblemBySlug 根据题目slug爬取题目详情
func (c *LeetCodeCrawler) CrawlProblemBySlug(ctx context.Context, titleSlug string) (*domain.CodingProblem, error) {
	c.logger.Printf("开始爬取题目: %s", titleSlug)

	query := fmt.Sprintf(`
	{
		question(titleSlug: "%s") {
			questionId
			questionFrontendId
			title
			titleSlug
			content
			difficulty
			topicTags {
				name
				slug
			}
		}
	}`, titleSlug)

	response, err := c.makeGraphQLRequest(ctx, query)
	if err != nil {
		c.logger.Printf("爬取题目失败: %v", err)
		return nil, fmt.Errorf("爬取题目失败: %w", err)
	}

	var problemResponse LeetCodeProblemResponse
	if err := json.Unmarshal(response, &problemResponse); err != nil {
		c.logger.Printf("解析题目数据失败: %v", err)
		return nil, fmt.Errorf("解析题目数据失败: %w", err)
	}

	question := problemResponse.Data.Question
	if question.Title == "" {
		return nil, errors.New("获取的题目数据为空")
	}

	// 提取标签
	tags := make([]string, 0, len(question.TopicTags))
	for _, tag := range question.TopicTags {
		tags = append(tags, tag.Name)
	}

	// 清理HTML内容

	problem := &domain.CodingProblem{
		Title:       question.Title,
		Difficulty:  question.Difficulty,
		Tags:        tags,
		Source:      "leetcode",
		SourceId:    question.QuestionFrontendId,
		SourceUrl:   fmt.Sprintf("https://leetcode.cn/problems/%s/", question.TitleSlug),
		StudyStatus: "not_started",
		Ctime:       time.Now(),
		Utime:       time.Now(),
	}

	c.logger.Printf("成功爬取题目: %s [%s]", problem.Title, problem.Difficulty)
	return problem, nil
}

// CrawlAndSaveDailyProblem 爬取并保存每日一题
func (c *LeetCodeCrawler) CrawlAndSaveDailyProblem(ctx context.Context) error {
	c.logger.Println("开始爬取并保存每日一题...")

	dailyProblem, err := c.GetDailyProblem(ctx)
	if err != nil {
		return fmt.Errorf("获取每日一题失败: %w", err)
	}

	// 检查是否已经存在该题目
	existingProblems, err := c.repository.FindBySource(ctx, "leetcode")
	if err != nil {
		c.logger.Printf("查询已有题目失败: %v", err)
	}

	var existingProblem *domain.CodingProblem
	for _, problem := range existingProblems {
		if problem.SourceId == dailyProblem.SourceId {
			existingProblem = &problem
			break
		}
	}

	// 转换为CodingProblem并保存
	var codingProblem domain.CodingProblem
	var problemId int64

	if existingProblem != nil {
		// 更新已有题目
		codingProblem = *existingProblem

		codingProblem.Tags = dailyProblem.Tags
		codingProblem.IsDailyProblem = true
		codingProblem.DailyDate = &dailyProblem.Date
		codingProblem.Utime = time.Now()

		if err := c.repository.Update(ctx, codingProblem); err != nil {
			return fmt.Errorf("更新题目失败: %w", err)
		}
		problemId = codingProblem.Id
		c.logger.Printf("已更新题目: %s", codingProblem.Title)
	} else {
		// 创建新题目
		codingProblem = domain.CodingProblem{
			Title:          dailyProblem.Title,
			Difficulty:     dailyProblem.Difficulty,
			Tags:           dailyProblem.Tags,
			Source:         dailyProblem.Source,
			SourceId:       dailyProblem.SourceId,
			SourceUrl:      dailyProblem.SourceUrl,
			StudyStatus:    "not_started",
			IsDailyProblem: true,
			DailyDate:      &dailyProblem.Date,
			Ctime:          time.Now(),
			Utime:          time.Now(),
		}

		if err := c.repository.Create(ctx, codingProblem); err != nil {
			return fmt.Errorf("创建题目失败: %w", err)
		}

		// 获取新创建题目的ID - 需要重新查询以获取自动生成的ID
		createdProblems, err := c.repository.FindBySourceId(ctx, dailyProblem.SourceId)
		if err != nil || len(createdProblems) == 0 {
			c.logger.Printf("无法获取新创建题目的ID: %v", err)
			return fmt.Errorf("无法获取新创建题目的ID: %w", err)
		}
		problemId = createdProblems[0].Id
		c.logger.Printf("已创建题目: %s (ID: %d)", codingProblem.Title, problemId)
	}

	// 使用 MarkAsDailyProblem 方法来正确保存每日一题记录
	if err := c.repository.MarkAsDailyProblem(ctx, problemId, dailyProblem.Date); err != nil {
		c.logger.Printf("标记每日一题失败: %v", err)
		// 不中断流程，只记录错误
	} else {
		c.logger.Printf("已标记为每日一题: %s", codingProblem.Title)
	}

	c.logger.Println("每日一题保存完成")
	return nil
}

// makeGraphQLRequest 发送GraphQL请求
func (c *LeetCodeCrawler) makeGraphQLRequest(ctx context.Context, query string) ([]byte, error) {
	// LeetCode CN的GraphQL端点
	url := "https://leetcode.cn/graphql/"

	// 构建请求体
	requestBody := map[string]interface{}{
		"query":         query,
		"variables":     map[string]interface{}{},
		"operationName": nil,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("构建请求体失败: %w", err)
	}

	c.logger.Printf("发送请求到: %s", url)
	c.logger.Printf("请求体: %s", string(jsonData))

	// 创建请求
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置请求头，模拟浏览器请求
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9,en;q=0.8")
	req.Header.Set("Accept-Encoding", "gzip, deflate, br")
	req.Header.Set("Referer", "https://leetcode.cn/problemset/all/")
	req.Header.Set("Origin", "https://leetcode.cn")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")

	// 发送请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	c.logger.Printf("响应状态码: %d", resp.StatusCode)
	c.logger.Printf("响应内容: %s", string(body))

	// 检查响应状态
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	return body, nil
}

// cleanHTMLContent 清理HTML内容
func (c *LeetCodeCrawler) cleanHTMLContent(content string) string {
	// 移除HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	content = re.ReplaceAllString(content, "")

	// 移除多余的空白字符
	re = regexp.MustCompile(`\s+`)
	content = re.ReplaceAllString(content, " ")

	// 处理HTML实体
	content = strings.ReplaceAll(content, "&nbsp;", " ")
	content = strings.ReplaceAll(content, "&lt;", "<")
	content = strings.ReplaceAll(content, "&gt;", ">")
	content = strings.ReplaceAll(content, "&amp;", "&")
	content = strings.ReplaceAll(content, "&quot;", "\"")

	return strings.TrimSpace(content)
}

// StartDailyCrawler 启动每日一题定时爬取
func (c *LeetCodeCrawler) StartDailyCrawler(ctx context.Context) error {
	c.logger.Println("启动每日一题定时爬取器...")

	ticker := time.NewTicker(24 * time.Hour) // 每24小时执行一次
	defer ticker.Stop()

	// 立即执行一次
	if err := c.CrawlAndSaveDailyProblem(ctx); err != nil {
		c.logger.Printf("初始爬取每日一题失败: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			c.logger.Println("每日一题爬取器已停止")
			return ctx.Err()
		case <-ticker.C:
			c.logger.Println("开始定时爬取每日一题...")
			if err := c.CrawlAndSaveDailyProblem(ctx); err != nil {
				c.logger.Printf("定时爬取每日一题失败: %v", err)
			}
		}
	}
}
