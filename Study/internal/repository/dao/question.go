package dao

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type QuestDao interface {
	Insert(ctx context.Context, quest Question) error
	FindByCategory(ctx context.Context, category string) ([]Question, error)
	FindAll(ctx context.Context) ([]Question, error)
	UpdateById(ctx context.Context, quest Question) error
	UpdateMasteryLevel(ctx context.Context, id int64, masteryLevel int) error
	DeleteById(ctx context.Context, id int64) error
	DeleteByCategory(ctx context.Context, category string) error
	FindAllCategories(ctx context.Context) ([]string, error)
	GetMasteryStats(ctx context.Context) (map[string]int, error)
}

type questionDao struct {
	db *gorm.DB
}

func NewQuestionDao(db *gorm.DB) QuestDao {
	return &questionDao{db: db}
}

func (dao *questionDao) Insert(ctx context.Context, quest Question) error {
	now := time.Now()
	quest.Ctime = now.Unix()
	quest.Utime = now.Unix()
	err := dao.db.WithContext(ctx).Create(&quest).Error
	return err
}

func (dao *questionDao) FindByCategory(ctx context.Context, category string) ([]Question, error) {
	var questions []Question
	err := dao.db.WithContext(ctx).Where("category = ?", category).Find(&questions).Error
	return questions, err
}

func (dao *questionDao) FindAll(ctx context.Context) ([]Question, error) {
	var questions []Question
	err := dao.db.WithContext(ctx).Find(&questions).Error
	return questions, err
}

func (dao *questionDao) UpdateById(ctx context.Context, quest Question) error {
	now := time.Now()
	quest.Utime = now.Unix()
	err := dao.db.WithContext(ctx).Where("id = ?", quest.Id).Updates(&quest).Error
	return err
}

func (dao *questionDao) UpdateMasteryLevel(ctx context.Context, id int64, masteryLevel int) error {
	now := time.Now()
	err := dao.db.WithContext(ctx).Model(&Question{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"mastery_level": masteryLevel,
			"utime":         now.Unix(),
		}).Error
	return err
}

func (dao *questionDao) DeleteById(ctx context.Context, id int64) error {
	err := dao.db.WithContext(ctx).Where("id = ?", id).Delete(&Question{}).Error
	return err
}

func (dao *questionDao) DeleteByCategory(ctx context.Context, category string) error {
	err := dao.db.WithContext(ctx).Where("category = ?", category).Delete(&Question{}).Error
	return err
}

func (dao *questionDao) FindAllCategories(ctx context.Context) ([]string, error) {
	var categories []string
	err := dao.db.WithContext(ctx).Model(&Question{}).Distinct("category").Pluck("category", &categories).Error
	return categories, err
}

func (dao *questionDao) GetMasteryStats(ctx context.Context) (map[string]int, error) {
	stats := make(map[string]int)

	// 统计总数
	var total int64
	if err := dao.db.WithContext(ctx).Model(&Question{}).Count(&total).Error; err != nil {
		return nil, err
	}
	stats["total"] = int(total)

	// 统计各掌握度级别的数量
	var unlearned, learning, mastered int64

	dao.db.WithContext(ctx).Model(&Question{}).Where("mastery_level = ?", 0).Count(&unlearned)
	dao.db.WithContext(ctx).Model(&Question{}).Where("mastery_level = ?", 1).Count(&learning)
	dao.db.WithContext(ctx).Model(&Question{}).Where("mastery_level = ?", 2).Count(&mastered)

	stats["unlearned"] = int(unlearned)
	stats["learning"] = int(learning)
	stats["mastered"] = int(mastered)

	return stats, nil
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
