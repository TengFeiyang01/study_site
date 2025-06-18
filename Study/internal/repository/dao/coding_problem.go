package dao

import (
	"Training/Study/internal/domain"
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

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

type CodingProblemDAO interface {
	Insert(ctx context.Context, problem CodingProblem) error
	FindAll(ctx context.Context) ([]CodingProblem, error)
	FindById(ctx context.Context, id int64) (CodingProblem, error)
	FindBySource(ctx context.Context, source string) ([]CodingProblem, error)
	FindBySourceId(ctx context.Context, sourceId string) ([]CodingProblem, error)
	FindByDifficulty(ctx context.Context, difficulty string) ([]CodingProblem, error)
	UpdateById(ctx context.Context, problem CodingProblem) error
	DeleteById(ctx context.Context, id int64) error
	GetDailyProblem(ctx context.Context) (*CodingProblem, error)
	SetDailyProblem(ctx context.Context, problemId int64) error
	GetDailyProblemHistory(ctx context.Context) ([]CodingProblem, error)
	MarkAsDailyProblem(ctx context.Context, problemId int64, date time.Time) error
	UpdateStudyStatus(ctx context.Context, problemId int64, status string, lastStudied *time.Time) error
	SaveDailyProblem(ctx context.Context, dailyProblem *domain.DailyProblem) error
}

type GormCodingProblemDAO struct {
	db *gorm.DB
}

func NewGormCodingProblemDAO(db *gorm.DB) CodingProblemDAO {
	return &GormCodingProblemDAO{db: db}
}

// 转换函数
func (p *CodingProblem) toDomain() domain.CodingProblem {
	return domain.CodingProblem{
		Id:             p.Id,
		Title:          p.Title,
		Difficulty:     p.Difficulty,
		Tags:           []string(p.Tags),
		Source:         p.Source,
		SourceId:       p.SourceId,
		SourceUrl:      p.SourceUrl,
		StudyStatus:    p.StudyStatus,
		LastStudied:    p.LastStudied,
		IsDailyProblem: p.IsDailyProblem,
		DailyDate:      p.DailyDate,
		IsHot100:       p.IsHot100,
		Ctime:          p.Ctime,
		Utime:          p.Utime,
	}
}

func toDAOCodingProblem(p domain.CodingProblem) CodingProblem {
	return CodingProblem{
		Id:             p.Id,
		Title:          p.Title,
		Difficulty:     p.Difficulty,
		Tags:           StringSlice(p.Tags),
		Source:         p.Source,
		SourceId:       p.SourceId,
		SourceUrl:      p.SourceUrl,
		StudyStatus:    p.StudyStatus,
		LastStudied:    p.LastStudied,
		IsDailyProblem: p.IsDailyProblem,
		DailyDate:      p.DailyDate,
		IsHot100:       p.IsHot100,
		Ctime:          p.Ctime,
		Utime:          p.Utime,
	}
}

// CodingProblem CRUD 实现
func (g *GormCodingProblemDAO) Insert(ctx context.Context, problem CodingProblem) error {
	now := time.Now()
	if problem.Ctime.IsZero() {
		problem.Ctime = now
	}
	if problem.Utime.IsZero() {
		problem.Utime = now
	}
	return g.db.WithContext(ctx).Create(&problem).Error
}

func (g *GormCodingProblemDAO) FindAll(ctx context.Context) ([]CodingProblem, error) {
	var problems []CodingProblem
	err := g.db.WithContext(ctx).Find(&problems).Error
	return problems, err
}

func (g *GormCodingProblemDAO) FindById(ctx context.Context, id int64) (CodingProblem, error) {
	var problem CodingProblem
	err := g.db.WithContext(ctx).Where("id = ?", id).First(&problem).Error
	return problem, err
}

func (g *GormCodingProblemDAO) FindBySource(ctx context.Context, source string) ([]CodingProblem, error) {
	var problems []CodingProblem
	err := g.db.WithContext(ctx).Where("source = ?", source).Find(&problems).Error
	return problems, err
}

func (g *GormCodingProblemDAO) FindBySourceId(ctx context.Context, sourceId string) ([]CodingProblem, error) {
	var problems []CodingProblem
	err := g.db.WithContext(ctx).Where("source_id = ?", sourceId).Find(&problems).Error
	return problems, err
}

func (g *GormCodingProblemDAO) FindByDifficulty(ctx context.Context, difficulty string) ([]CodingProblem, error) {
	var problems []CodingProblem
	err := g.db.WithContext(ctx).Where("difficulty = ?", difficulty).Find(&problems).Error
	return problems, err
}

func (g *GormCodingProblemDAO) UpdateById(ctx context.Context, problem CodingProblem) error {
	problem.Utime = time.Now()
	return g.db.WithContext(ctx).Where("id = ?", problem.Id).Updates(&problem).Error
}

func (g *GormCodingProblemDAO) DeleteById(ctx context.Context, id int64) error {
	return g.db.WithContext(ctx).Where("id = ?", id).Delete(&CodingProblem{}).Error
}

// 每日一题相关方法
func (g *GormCodingProblemDAO) GetDailyProblem(ctx context.Context) (*CodingProblem, error) {
	today := time.Now()
	todayStr := today.Format("2006-01-02")

	// 查找今天标记为每日一题的题目
	var problem CodingProblem
	err := g.db.WithContext(ctx).Where("is_daily_problem = ? AND DATE(daily_date) = ?", true, todayStr).First(&problem).Error
	if err == nil {
		return &problem, nil
	}

	if err != gorm.ErrRecordNotFound {
		// 如果是其他错误，记录日志
		log.Printf("Error finding daily problem: %v", err)
	}

	// 如果没有找到今天的每日一题，尝试从 daily_problems 表获取
	var dailyProblem DailyProblem
	err = g.db.WithContext(ctx).Where("date = ?", todayStr).First(&dailyProblem).Error
	if err == nil {
		// 找到了每日一题记录，获取对应的题目
		err = g.db.WithContext(ctx).Where("id = ?", dailyProblem.ProblemId).First(&problem).Error
		if err == nil {
			return &problem, nil
		}
	}

	// 如果还是没有找到，随机选择一题作为每日一题
	err = g.db.WithContext(ctx).Order("RAND()").First(&problem).Error
	if err != nil {
		return nil, err
	}

	// 将选中的题目标记为每日一题
	err = g.MarkAsDailyProblem(ctx, problem.Id, today)
	if err != nil {
		log.Printf("Error marking problem as daily: %v", err)
		return nil, err
	}

	return &problem, nil
}

func (g *GormCodingProblemDAO) SetDailyProblem(ctx context.Context, problemId int64) error {
	return g.MarkAsDailyProblem(ctx, problemId, time.Now())
}

func (g *GormCodingProblemDAO) GetDailyProblemHistory(ctx context.Context) ([]CodingProblem, error) {
	// 从 daily_problems 表获取历史记录，按日期倒序排列
	var dailyProblems []DailyProblem
	err := g.db.WithContext(ctx).Order("date desc").Find(&dailyProblems).Error
	if err != nil {
		return nil, err
	}

	// 转换为 CodingProblem 列表
	var problems []CodingProblem
	for _, dp := range dailyProblems {
		problem := CodingProblem{
			Id:             dp.ProblemId, // 使用原始题目ID
			Title:          dp.Title,
			Difficulty:     dp.Difficulty,
			Tags:           dp.Tags,
			Source:         dp.Source,
			SourceId:       dp.SourceId,
			SourceUrl:      dp.SourceUrl,
			IsDailyProblem: true,               // 历史记录中的都是每日一题
			DailyDate:      parseDate(dp.Date), // 转换日期字符串为时间
			Ctime:          dp.Ctime,
			Utime:          dp.Utime,
		}
		problems = append(problems, problem)
	}

	return problems, nil
}

// parseDate 辅助函数，将日期字符串转换为时间指针
func parseDate(dateStr string) *time.Time {
	if dateStr == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return nil
	}

	return &t
}

func (g *GormCodingProblemDAO) MarkAsDailyProblem(ctx context.Context, problemId int64, date time.Time) error {
	// 先重置所有题目的每日一题标记
	err := g.db.WithContext(ctx).Model(&CodingProblem{}).Where("is_daily_problem = ?", true).Updates(map[string]interface{}{
		"is_daily_problem": false,
		"daily_date":       nil,
	}).Error
	if err != nil {
		return err
	}

	// 更新选中的题目为每日一题
	err = g.db.WithContext(ctx).Model(&CodingProblem{}).Where("id = ?", problemId).Updates(map[string]interface{}{
		"is_daily_problem": true,
		"daily_date":       date,
	}).Error
	if err != nil {
		return err
	}

	// 获取题目信息
	var problem CodingProblem
	err = g.db.WithContext(ctx).Where("id = ?", problemId).First(&problem).Error
	if err != nil {
		return err
	}

	// 同时在 daily_problems 表中记录
	dailyProblem := DailyProblem{
		Date:       date.Format("2006-01-02"),
		Title:      problem.Title,
		Difficulty: problem.Difficulty,
		Tags:       problem.Tags,
		Source:     problem.Source,
		SourceId:   problem.SourceId,
		SourceUrl:  problem.SourceUrl,
		ProblemId:  problem.Id,
		Ctime:      date,
		Utime:      date,
	}

	// 先删除可能存在的当天记录
	err = g.db.WithContext(ctx).Where("date = ?", dailyProblem.Date).Delete(&DailyProblem{}).Error
	if err != nil {
		return err
	}

	// 创建新记录
	return g.db.WithContext(ctx).Create(&dailyProblem).Error
}

func (g *GormCodingProblemDAO) UpdateStudyStatus(ctx context.Context, problemId int64, status string, lastStudied *time.Time) error {
	updates := map[string]interface{}{
		"study_status": status,
		"utime":        time.Now(),
	}
	if lastStudied != nil {
		updates["last_studied"] = *lastStudied
	}
	return g.db.WithContext(ctx).Model(&CodingProblem{}).Where("id = ?", problemId).Updates(updates).Error
}

func (g *GormCodingProblemDAO) SaveDailyProblem(ctx context.Context, dailyProblem *domain.DailyProblem) error {
	// 将domain.DailyProblem转换为dao.DailyProblem
	daoDailyProblem := DailyProblem{
		Date:        dailyProblem.Date.Format("2006-01-02"),
		Title:       dailyProblem.Title,
		Difficulty:  dailyProblem.Difficulty,
		Tags:        StringSlice(dailyProblem.Tags),
		Source:      dailyProblem.Source,
		SourceId:    dailyProblem.SourceId,
		SourceUrl:   dailyProblem.SourceUrl,
		Ctime:       dailyProblem.Ctime,
		Utime:       dailyProblem.Utime,
	}

	// 首先检查是否已存在相同日期的记录
	var existing DailyProblem
	err := g.db.WithContext(ctx).Where("date = ?", daoDailyProblem.Date).First(&existing).Error

	if err == nil {
		// 如果存在，更新记录
		daoDailyProblem.Id = existing.Id
		return g.db.WithContext(ctx).Save(&daoDailyProblem).Error
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		// 如果不存在，创建新记录
		return g.db.WithContext(ctx).Create(&daoDailyProblem).Error
	}

	return err
}

// 自定义JSON序列化
func (p *CodingProblem) BeforeCreate(tx *gorm.DB) error {
	return nil
}

func (p *CodingProblem) AfterFind(tx *gorm.DB) error {
	return nil
}

// 处理Tags字段的JSON序列化
func (p CodingProblem) Value() (interface{}, error) {
	return json.Marshal(p.Tags)
}

func (p *CodingProblem) Scan(value interface{}) error {
	if value == nil {
		p.Tags = []string{}
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &p.Tags)
	case string:
		return json.Unmarshal([]byte(v), &p.Tags)
	default:
		return fmt.Errorf("cannot scan %T into Tags", value)
	}
}
