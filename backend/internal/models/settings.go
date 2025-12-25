package models

import (
	"time"
)

// SystemSettings 系统设置模型
type SystemSettings struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Key         string    `gorm:"uniqueIndex;size:100;not null" json:"key"`
	Value       string    `gorm:"type:text;not null" json:"value"`
	Category    string    `gorm:"size:50;index;not null" json:"category"` // e.g., "salt", "email", "notification"
	Description string    `gorm:"type:text" json:"description"`
	UpdatedAt   time.Time `json:"updated_at"`
	UpdatedBy   uint64    `gorm:"index" json:"updated_by"` // User ID who last updated

	// 关联
	Updater *User `gorm:"foreignKey:UpdatedBy" json:"updater,omitempty"`
}

// TableName 指定表名
func (SystemSettings) TableName() string {
	return "system_settings"
}

