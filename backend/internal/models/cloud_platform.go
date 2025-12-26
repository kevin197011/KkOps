package models

import (
	"time"

	"gorm.io/gorm"
)

// CloudPlatform 云平台模型
type CloudPlatform struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:50;not null;uniqueIndex" json:"name"`
	DisplayName string         `gorm:"size:100;not null" json:"display_name"`
	Icon        string         `gorm:"size:50" json:"icon"`
	Color       string         `gorm:"size:20" json:"color"`
	SortOrder   int            `gorm:"default:0" json:"sort_order"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (CloudPlatform) TableName() string {
	return "cloud_platforms"
}
