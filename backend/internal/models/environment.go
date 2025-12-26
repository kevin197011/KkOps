package models

import (
	"time"

	"gorm.io/gorm"
)

// Environment 环境模型
type Environment struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;uniqueIndex" json:"name"`
	DisplayName string         `gorm:"size:100;not null" json:"display_name"`
	Color       string         `gorm:"size:20" json:"color"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (Environment) TableName() string {
	return "environments"
}
