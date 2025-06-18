package domain

import "time"

type Question struct {
	Id           int64     `json:"id"`
	Category     string    `json:"category"`
	Content      string    `json:"content"`
	Answer       string    `json:"answer"`
	MasteryLevel int       `json:"mastery_level"` // 0: 未学习, 1: 学习中, 2: 已掌握
	Ctime        time.Time `json:"ctime"`
	Utime        time.Time `json:"utime"`
}
