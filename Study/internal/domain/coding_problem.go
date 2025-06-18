package domain

import "time"

// CodingProblem 刷题问题
type CodingProblem struct {
	Id             int64      `json:"id"`
	Title          string     `json:"title"`
	Difficulty     string     `json:"difficulty"` // Easy, Medium, Hard
	Tags           []string   `json:"tags"`
	Source         string     `json:"source"`                 // leetcode, nowcoder
	SourceId       string     `json:"source_id"`              // 原网站的问题ID
	SourceUrl      string     `json:"source_url"`             // 原网站的链接
	StudyStatus    string     `json:"study_status"`           // 学习状态: not_started, in_progress, completed
	LastStudied    *time.Time `json:"last_studied,omitempty"` // 最后学习时间
	IsDailyProblem bool       `json:"is_daily_problem"`       // 是否为每日一题
	DailyDate      *time.Time `json:"daily_date,omitempty"`   // 每日一题的日期
	IsHot100       bool       `json:"is_hot100"`              // 是否为 Hot 100 题目
	Ctime          time.Time  `json:"ctime"`
	Utime          time.Time  `json:"utime"`
}

// DailyProblem 每日一题记录
type DailyProblem struct {
	Id         int64     `json:"id"`
	Date       time.Time `json:"date"`       // 每日一题的日期
	Title      string    `json:"title"`      // 题目标题
	Difficulty string    `json:"difficulty"` // 难度
	Tags       []string  `json:"tags"`       // 标签
	Source     string    `json:"source"`     // 来源
	SourceId   string    `json:"source_id"`  // 原网站的问题ID
	SourceUrl  string    `json:"source_url"` // 原网站的链接
	Ctime      time.Time `json:"ctime"`      // 创建时间
	Utime      time.Time `json:"utime"`      // 更新时间
}

// CodingProblemRequest 创建/更新刷题问题的请求
type CodingProblemRequest struct {
	Title       string   `json:"title" binding:"required"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty" binding:"required"`
	Tags        []string `json:"tags"`
	Source      string   `json:"source" binding:"required"`
	SourceId    string   `json:"source_id"`
	SourceUrl   string   `json:"source_url"`
}
