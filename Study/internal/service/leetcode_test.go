package service

import (
	"context"
	"fmt"
	"testing"
)

// 测试每日一题爬取功能
func TestLeetCodeCrawler_GetDailyProblem(t *testing.T) {
	// 这里需要一个真实的 repository 实现，暂时跳过
	t.Skip("需要数据库连接")

	crawler := NewLeetCodeCrawler(nil)

	dailyProblem, err := crawler.GetDailyProblem(context.Background())
	if err != nil {
		t.Fatalf("获取每日一题失败: %v", err)
	}
	t.Logf("今日题目: %s\n", dailyProblem.Title)
	t.Logf("难度: %s\n", dailyProblem.Difficulty)
	t.Logf("链接: %s\n", dailyProblem.SourceUrl)
}

// 测试 Hot 100 题目列表
func TestLeetCodeCrawler_GetHot100Problems(t *testing.T) {
	crawler := NewLeetCodeCrawler(nil)

	hot100 := crawler.GetHot100Problems()

	if len(hot100) == 0 {
		t.Fatal("Hot 100 题目列表为空")
	}

	fmt.Printf("Hot 100 题目总数: %d\n", len(hot100))

	// 打印前5个题目作为示例
	for i, problem := range hot100[:5] {
		fmt.Printf("%d. %s (%s) - %s\n", i+1, problem.Title, problem.Difficulty, problem.TitleSlug)
	}
}

// 测试 fallbackTitle 函数
func TestLeetCodeCrawler_FallbackTitle(t *testing.T) {
	crawler := NewLeetCodeCrawler(nil)

	// 测试优先使用翻译标题
	title1 := crawler.fallbackTitle("两数之和", "Two Sum")
	if title1 != "两数之和" {
		t.Errorf("期望 '两数之和'，实际得到 '%s'", title1)
	}

	// 测试翻译标题为空时使用原标题
	title2 := crawler.fallbackTitle("", "Two Sum")
	if title2 != "Two Sum" {
		t.Errorf("期望 'Two Sum'，实际得到 '%s'", title2)
	}
}
