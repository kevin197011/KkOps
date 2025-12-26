package models

import (
	"time"

	"gorm.io/datatypes"
)

// CommandTemplate 命令模板模型
type CommandTemplate struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	Category        string         `gorm:"size:50;index" json:"category"` // system, network, disk, process, custom
	CommandFunction string         `gorm:"size:255;not null" json:"command_function"`
	CommandArgs     datatypes.JSON `gorm:"type:jsonb" json:"command_args"`
	Icon            string         `gorm:"size:50" json:"icon"`
	CreatedBy       uint64         `gorm:"not null;index" json:"created_by"`
	IsPublic        bool           `gorm:"default:false" json:"is_public"`
	UsageCount      int            `gorm:"default:0" json:"usage_count"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// 关联
	CreatedByUser User `gorm:"foreignKey:CreatedBy" json:"created_by_user,omitempty"`
}

// TableName 指定表名
func (CommandTemplate) TableName() string {
	return "command_templates"
}

