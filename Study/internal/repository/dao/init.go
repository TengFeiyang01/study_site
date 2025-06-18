package dao

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfig 系统配置表
type SystemConfig struct {
	Id        int64     `gorm:"primaryKey,autoIncrement" json:"id"`
	Key       string    `gorm:"type:varchar(100);uniqueIndex;not null" json:"key"`
	Value     string    `gorm:"type:text" json:"value"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (SystemConfig) TableName() string {
	return "system_config"
}

func InitTables(db *gorm.DB) error {
	err := db.AutoMigrate(&Question{}, &CodingProblem{}, &DailyProblem{}, &SystemConfig{})
	if err != nil {
		return err
	}

	return nil
}
