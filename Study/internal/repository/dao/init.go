package dao

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

func InitTables(db *gorm.DB) error {
	err := db.AutoMigrate(&Question{}, &CodingProblem{}, &DailyProblem{})
	if err != nil {
		return err
	}

	return nil
}

// StringSlice 自定义类型用于处理数组字段
type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
	if value == nil {
		*s = StringSlice{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, s)
}

// CodingProblem 刷题问题数据库模型
type CodingProblem struct {
	Id             int64       `gorm:"primaryKey,autoIncrement" json:"id"`
	Title          string      `gorm:"type:varchar(255);not null" json:"title"`
	Difficulty     string      `gorm:"type:varchar(50)" json:"difficulty"`
	Tags           StringSlice `gorm:"type:json" json:"tags"`
	Source         string      `gorm:"type:varchar(100)" json:"source"`
	SourceId       string      `gorm:"type:varchar(100)" json:"source_id"`
	SourceUrl      string      `gorm:"type:varchar(500)" json:"source_url"`
	StudyStatus    string      `gorm:"type:varchar(50);default:'not_started'" json:"study_status"` // 学习状态
	LastStudied    *time.Time  `gorm:"type:datetime(3)" json:"last_studied"`                       // 最后学习时间
	IsDailyProblem bool        `gorm:"type:boolean;default:false" json:"is_daily_problem"`         // 是否为每日一题
	DailyDate      *time.Time  `gorm:"type:datetime(3)" json:"daily_date"`                         // 每日一题日期
	IsHot100       bool        `gorm:"type:boolean;default:false" json:"is_hot100"`                // 是否为 Hot 100 题目
	Ctime          time.Time   `gorm:"type:datetime(3)" json:"ctime"`
	Utime          time.Time   `gorm:"type:datetime(3)" json:"utime"`
}

func (c CodingProblem) TableName() string {
	return "coding_problems"
}

// DailyProblem 每日一题模型
type DailyProblem struct {
	Id         int64         `gorm:"primarykey,autoIncrement"`
	Date       string        `gorm:"column:date;unique;not null"` // 格式: 2024-06-17
	Title      string        `gorm:"type:varchar(255);not null"`
	Difficulty string        `gorm:"type:varchar(50)"`
	Tags       StringSlice   `gorm:"type:json"`
	Source     string        `gorm:"type:varchar(100)"`
	SourceId   string        `gorm:"type:varchar(100)"`
	SourceUrl  string        `gorm:"type:varchar(500)"`
	ProblemId  int64         `gorm:"column:problem_id;not null"`
	Problem    CodingProblem `gorm:"foreignKey:ProblemId"`
	Ctime      time.Time     `gorm:"type:datetime(3)"`
	Utime      time.Time     `gorm:"type:datetime(3)"`
}

func (d DailyProblem) TableName() string {
	return "daily_problems"
}

type Question struct {
	Id           int64  `gorm:"primary_key,autoIncrement"`
	Category     string `gorm:"type:varchar(100);not null"`
	Content      string `gorm:"type:text;not null"`
	Answer       string `gorm:"type:text;not null"`
	MasteryLevel int    `gorm:"type:int;default:0"` // 0: 未学习, 1: 学习中, 2: 已掌握
	Count        int64  `gorm:"type:bigint;not null"`
	Ctime        int64  `gorm:"type:bigint;not null"`
	Utime        int64  `gorm:"type:bigint;not null"`
}

func (q Question) TableName() string {
	return "questions"
}
