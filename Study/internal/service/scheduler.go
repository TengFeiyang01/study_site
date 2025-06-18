package service

import (
	"context"
	"log"
	"regexp"
	"strings"
	"time"

	"Training/Study/internal/repository"
)

type Scheduler struct {
	crawler           *LeetCodeCrawler
	codingProblemRepo repository.CodingProblemRepository
}

func NewScheduler(crawler *LeetCodeCrawler, codingProblemRepo repository.CodingProblemRepository) *Scheduler {
	return &Scheduler{
		crawler:           crawler,
		codingProblemRepo: codingProblemRepo,
	}
}

// StartScheduler 启动定时任务 - 只处理每日一题更新
func (s *Scheduler) StartScheduler() {
	// 启动定时任务
	go s.scheduleTask()
}

// scheduleTask 定时任务主循环 - 每天0点更新每日一题
func (s *Scheduler) scheduleTask() {
	for {
		now := time.Now()
		next := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, now.Location())
		duration := next.Sub(now)

		log.Printf("⏰ 下次每日一题更新时间: %s (还有 %v)", next.Format("2006-01-02 15:04:05"), duration)

		// 等待到下一个0点
		time.Sleep(duration)

		// 获取每日一题
		log.Printf("🕛 到达0点，开始获取新的每日一题...")
		if err := s.fetchDailyProblem(); err != nil {
			log.Printf("❌ 获取每日一题失败: %v", err)
		} else {
			log.Printf("✅ 每日一题更新完成")
		}
	}
}

// fetchDailyProblem 获取每日一题
func (s *Scheduler) fetchDailyProblem() error {
	ctx := context.Background()
	return s.crawler.CrawlAndSaveDailyProblem(ctx)
}

// convertDifficultyToChinese 将英文难度转换为中文
func (s *Scheduler) convertDifficultyToChinese(difficulty string) string {
	switch strings.ToLower(difficulty) {
	case "easy":
		return "简单"
	case "medium":
		return "中等"
	case "hard":
		return "困难"
	default:
		return difficulty
	}
}

// cleanHTMLContent 清理HTML内容
func (s *Scheduler) cleanHTMLContent(content string) string {
	// 简单的HTML标签清理
	content = strings.ReplaceAll(content, "<p>", "\n")
	content = strings.ReplaceAll(content, "</p>", "\n")
	content = strings.ReplaceAll(content, "<br>", "\n")
	content = strings.ReplaceAll(content, "<br/>", "\n")
	content = strings.ReplaceAll(content, "<strong>", "**")
	content = strings.ReplaceAll(content, "</strong>", "**")
	content = strings.ReplaceAll(content, "<em>", "*")
	content = strings.ReplaceAll(content, "</em>", "*")

	// 移除其他HTML标签
	re := regexp.MustCompile(`<[^>]*>`)
	content = re.ReplaceAllString(content, "")

	// 清理多余的换行
	content = regexp.MustCompile(`\n\s*\n`).ReplaceAllString(content, "\n\n")
	content = strings.TrimSpace(content)

	return content
}
