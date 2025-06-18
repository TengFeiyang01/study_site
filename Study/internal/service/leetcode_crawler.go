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
	"strconv"
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

// Hot100题目列表（写死）
var hot100Problems = []struct {
	ID         string
	Title      string
	TitleSlug  string
	Difficulty string
}{
	{"1", "两数之和", "two-sum", "Easy"},
	{"2", "两数相加", "add-two-numbers", "Medium"},
	{"3", "无重复字符的最长子串", "longest-substring-without-repeating-characters", "Medium"},
	{"4", "寻找两个正序数组的中位数", "median-of-two-sorted-arrays", "Hard"},
	{"5", "最长回文子串", "longest-palindromic-substring", "Medium"},
	{"10", "正则表达式匹配", "regular-expression-matching", "Hard"},
	{"11", "盛最多水的容器", "container-with-most-water", "Medium"},
	{"15", "三数之和", "3sum", "Medium"},
	{"17", "电话号码的字母组合", "letter-combinations-of-a-phone-number", "Medium"},
	{"19", "删除链表的倒数第 N 个结点", "remove-nth-node-from-end-of-list", "Medium"},
	{"20", "有效的括号", "valid-parentheses", "Easy"},
	{"21", "合并两个有序链表", "merge-two-sorted-lists", "Easy"},
	{"22", "括号生成", "generate-parentheses", "Medium"},
	{"23", "合并K个升序链表", "merge-k-sorted-lists", "Hard"},
	{"25", "K 个一组翻转链表", "reverse-nodes-in-k-group", "Hard"},
	{"31", "下一个排列", "next-permutation", "Medium"},
	{"32", "最长有效括号", "longest-valid-parentheses", "Hard"},
	{"33", "搜索旋转排序数组", "search-in-rotated-sorted-array", "Medium"},
	{"34", "在排序数组中查找元素的第一个和最后一个位置", "find-first-and-last-position-of-element-in-sorted-array", "Medium"},
	{"35", "搜索插入位置", "search-insert-position", "Easy"},
	{"39", "组合总和", "combination-sum", "Medium"},
	{"42", "接雨水", "trapping-rain-water", "Hard"},
	{"46", "全排列", "permutations", "Medium"},
	{"48", "旋转图像", "rotate-image", "Medium"},
	{"49", "字母异位词分组", "group-anagrams", "Medium"},
	{"53", "最大子数组和", "maximum-subarray", "Medium"},
	{"55", "跳跃游戏", "jump-game", "Medium"},
	{"56", "合并区间", "merge-intervals", "Medium"},
	{"62", "不同路径", "unique-paths", "Medium"},
	{"64", "最小路径和", "minimum-path-sum", "Medium"},
	{"70", "爬楼梯", "climbing-stairs", "Easy"},
	{"72", "编辑距离", "edit-distance", "Medium"},
	{"75", "颜色分类", "sort-colors", "Medium"},
	{"76", "最小覆盖子串", "minimum-window-substring", "Hard"},
	{"78", "子集", "subsets", "Medium"},
	{"79", "单词搜索", "word-search", "Medium"},
	{"84", "柱状图中最大的矩形", "largest-rectangle-in-histogram", "Hard"},
	{"85", "最大矩形", "maximal-rectangle", "Hard"},
	{"94", "二叉树的中序遍历", "binary-tree-inorder-traversal", "Easy"},
	{"96", "不同的二叉搜索树", "unique-binary-search-trees", "Medium"},
	{"98", "验证二叉搜索树", "validate-binary-search-tree", "Medium"},
	{"101", "对称二叉树", "symmetric-tree", "Easy"},
	{"102", "二叉树的层序遍历", "binary-tree-level-order-traversal", "Medium"},
	{"104", "二叉树的最大深度", "maximum-depth-of-binary-tree", "Easy"},
	{"105", "从前序与中序遍历序列构造二叉树", "construct-binary-tree-from-preorder-and-inorder-traversal", "Medium"},
	{"114", "二叉树展开为链表", "flatten-binary-tree-to-linked-list", "Medium"},
	{"121", "买卖股票的最佳时机", "best-time-to-buy-and-sell-stock", "Easy"},
	{"124", "二叉树中的最大路径和", "binary-tree-maximum-path-sum", "Hard"},
	{"128", "最长连续序列", "longest-consecutive-sequence", "Medium"},
	{"136", "只出现一次的数字", "single-number", "Easy"},
	{"139", "单词拆分", "word-break", "Medium"},
	{"141", "环形链表", "linked-list-cycle", "Easy"},
	{"142", "环形链表 II", "linked-list-cycle-ii", "Medium"},
	{"146", "LRU 缓存", "lru-cache", "Medium"},
	{"148", "排序链表", "sort-list", "Medium"},
	{"152", "乘积最大子数组", "maximum-product-subarray", "Medium"},
	{"155", "最小栈", "min-stack", "Medium"},
	{"160", "相交链表", "intersection-of-two-linked-lists", "Easy"},
	{"169", "多数元素", "majority-element", "Easy"},
	{"198", "打家劫舍", "house-robber", "Medium"},
	{"200", "岛屿数量", "number-of-islands", "Medium"},
	{"206", "反转链表", "reverse-linked-list", "Easy"},
	{"207", "课程表", "course-schedule", "Medium"},
	{"208", "实现 Trie (前缀树)", "implement-trie-prefix-tree", "Medium"},
	{"215", "数组中的第K个最大元素", "kth-largest-element-in-an-array", "Medium"},
	{"221", "最大正方形", "maximal-square", "Medium"},
	{"226", "翻转二叉树", "invert-binary-tree", "Easy"},
	{"234", "回文链表", "palindrome-linked-list", "Easy"},
	{"236", "二叉树的最近公共祖先", "lowest-common-ancestor-of-a-binary-tree", "Medium"},
	{"238", "除自身以外数组的乘积", "product-of-array-except-self", "Medium"},
	{"239", "滑动窗口最大值", "sliding-window-maximum", "Hard"},
	{"240", "搜索二维矩阵 II", "search-a-2d-matrix-ii", "Medium"},
	{"253", "会议室 II", "meeting-rooms-ii", "Medium"},
	{"279", "完全平方数", "perfect-squares", "Medium"},
	{"283", "移动零", "move-zeroes", "Easy"},
	{"287", "寻找重复数", "find-the-duplicate-number", "Medium"},
	{"297", "二叉树的序列化与反序列化", "serialize-and-deserialize-binary-tree", "Hard"},
	{"300", "最长递增子序列", "longest-increasing-subsequence", "Medium"},
	{"301", "删除无效的括号", "remove-invalid-parentheses", "Hard"},
	{"309", "买卖股票的最佳时机含冷冻期", "best-time-to-buy-and-sell-stock-with-cooldown", "Medium"},
	{"312", "戳气球", "burst-balloons", "Hard"},
	{"322", "零钱兑换", "coin-change", "Medium"},
	{"337", "打家劫舍 III", "house-robber-iii", "Medium"},
	{"338", "比特位计数", "counting-bits", "Easy"},
	{"347", "前 K 个高频元素", "top-k-frequent-elements", "Medium"},
	{"394", "字符串解码", "decode-string", "Medium"},
	{"399", "除法求值", "evaluate-division", "Medium"},
	{"406", "根据身高重建队列", "queue-reconstruction-by-height", "Medium"},
	{"416", "分割等和子集", "partition-equal-subset-sum", "Medium"},
	{"437", "路径总和 III", "path-sum-iii", "Medium"},
	{"438", "找到字符串中所有字母异位词", "find-all-anagrams-in-a-string", "Medium"},
	{"448", "找到所有数组中消失的数字", "find-all-numbers-disappeared-in-an-array", "Easy"},
	{"461", "汉明距离", "hamming-distance", "Easy"},
	{"494", "目标和", "target-sum", "Medium"},
	{"538", "把二叉搜索树转换为累加树", "convert-bst-to-greater-tree", "Medium"},
	{"543", "二叉树的直径", "diameter-of-binary-tree", "Easy"},
	{"560", "和为 K 的子数组", "subarray-sum-equals-k", "Medium"},
	{"581", "最短无序连续子数组", "shortest-unsorted-continuous-subarray", "Medium"},
	{"617", "合并二叉树", "merge-two-binary-trees", "Easy"},
	{"621", "任务调度器", "task-scheduler", "Medium"},
	{"647", "回文子串", "palindromic-substrings", "Medium"},
	{"739", "每日温度", "daily-temperatures", "Medium"},
	{"763", "划分字母区间", "partition-labels", "Medium"},
	{"1143", "最长公共子序列", "longest-common-subsequence", "Medium"},
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

// GetHot100Problems 获取 Hot 100 题目列表
func (c *LeetCodeCrawler) GetHot100Problems() []struct {
	ID         string
	Title      string
	TitleSlug  string
	Difficulty string
} {
	return hot100Problems
}

// CrawlAndSaveHot100Problems 爬取并保存 Hot 100 题目
func (c *LeetCodeCrawler) CrawlAndSaveHot100Problems(ctx context.Context) error {
	c.logger.Printf("开始爬取并保存 Hot 100 题目，共 %d 个题目...", len(hot100Problems))

	for i, hotProblem := range hot100Problems {
		c.logger.Printf("正在处理第 %d/%d 个题目: %s", i+1, len(hot100Problems), hotProblem.Title)

		// 检查是否已经存在该题目
		existingProblems, err := c.repository.FindBySourceId(ctx, hotProblem.ID)
		if err != nil {
			c.logger.Printf("查询已有题目失败: %v", err)
		}

		if len(existingProblems) > 0 {
			c.logger.Printf("题目 %s 已存在，跳过", hotProblem.Title)
			continue
		}

		// 爬取题目详情
		problem, err := c.CrawlProblemBySlug(ctx, hotProblem.TitleSlug)
		if err != nil {
			c.logger.Printf("爬取题目 %s 失败: %v", hotProblem.Title, err)
			// 如果爬取失败，使用基本信息创建题目
			problem = &domain.CodingProblem{
				Title:       hotProblem.Title,
				Difficulty:  hotProblem.Difficulty,
				Tags:        []string{"Hot 100", "算法"},
				Source:      "leetcode",
				SourceId:    hotProblem.ID,
				SourceUrl:   fmt.Sprintf("https://leetcode.cn/problems/%s/", hotProblem.TitleSlug),
				StudyStatus: "not_started",
				IsHot100:    true,
				Ctime:       time.Now(),
				Utime:       time.Now(),
			}
		} else {
			// 标记为 Hot 100 题目
			problem.IsHot100 = true
			problem.Tags = append(problem.Tags, "Hot 100")
		}

		// 保存题目
		if err := c.repository.Create(ctx, *problem); err != nil {
			c.logger.Printf("保存题目 %s 失败: %v", problem.Title, err)
			continue
		}

		c.logger.Printf("已保存题目: %s", problem.Title)

		// 添加延迟避免请求过快
		time.Sleep(1 * time.Second)
	}

	c.logger.Println("Hot 100 题目爬取完成")
	return nil
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

// BatchCrawlProblems 批量爬取题目
func (c *LeetCodeCrawler) BatchCrawlProblems(ctx context.Context, titleSlugs []string) error {
	c.logger.Printf("开始批量爬取 %d 个题目...", len(titleSlugs))

	for i, slug := range titleSlugs {
		c.logger.Printf("正在爬取第 %d/%d 个题目: %s", i+1, len(titleSlugs), slug)

		problem, err := c.CrawlProblemBySlug(ctx, slug)
		if err != nil {
			c.logger.Printf("爬取题目 %s 失败: %v", slug, err)
			continue
		}

		// 检查是否已存在
		existingProblems, err := c.repository.FindBySource(ctx, "leetcode")
		if err != nil {
			c.logger.Printf("查询已有题目失败: %v", err)
		}

		exists := false
		for _, existing := range existingProblems {
			if existing.SourceId == problem.SourceId {
				exists = true
				break
			}
		}

		if !exists {
			if err := c.repository.Create(ctx, *problem); err != nil {
				c.logger.Printf("保存题目 %s 失败: %v", problem.Title, err)
				continue
			}
			c.logger.Printf("已保存题目: %s", problem.Title)
		} else {
			c.logger.Printf("题目 %s 已存在，跳过", problem.Title)
		}

		// 添加延迟避免请求过快
		time.Sleep(2 * time.Second)
	}

	c.logger.Println("批量爬取完成")
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

// GetProblemsInRange 获取指定范围的题目列表（用于批量爬取）
func (c *LeetCodeCrawler) GetProblemsInRange(ctx context.Context, start, end int) ([]string, error) {
	c.logger.Printf("获取题目范围: %d - %d", start, end)

	var titleSlugs []string

	// 这里可以通过LeetCode的API获取题目列表
	// 由于LeetCode的API限制，这里提供一个示例实现
	query := `
	{
		problemsetQuestionList: problemsetQuestionList(
			categorySlug: ""
			limit: 50
			skip: 0
			filters: {}
		) {
			total: totalNum
			questions: data {
				acRate
				difficulty
				freqBar
				frontendQuestionId: questionFrontendId
				isFavor
				paidOnly: isPaidOnly
				status
				title
				titleSlug
				topicTags {
					name
					id
					slug
				}
				hasSolution
				hasVideoSolution
			}
		}
	}`

	response, err := c.makeGraphQLRequest(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("获取题目列表失败: %w", err)
	}

	// 解析响应获取titleSlug列表
	var result map[string]interface{}
	if err := json.Unmarshal(response, &result); err != nil {
		return nil, fmt.Errorf("解析题目列表失败: %w", err)
	}

	// 提取题目slug
	if data, ok := result["data"].(map[string]interface{}); ok {
		if problemsetQuestionList, ok := data["problemsetQuestionList"].(map[string]interface{}); ok {
			if questions, ok := problemsetQuestionList["questions"].([]interface{}); ok {
				for _, q := range questions {
					if question, ok := q.(map[string]interface{}); ok {
						if titleSlug, ok := question["titleSlug"].(string); ok {
							if frontendId, ok := question["frontendQuestionId"].(string); ok {
								if id, err := strconv.Atoi(frontendId); err == nil {
									if id >= start && id <= end {
										titleSlugs = append(titleSlugs, titleSlug)
									}
								}
							}
						}
					}
				}
			}
		}
	}

	c.logger.Printf("找到 %d 个题目", len(titleSlugs))
	return titleSlugs, nil
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
