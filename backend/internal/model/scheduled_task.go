// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// ScheduledTask 定时任务模型
type ScheduledTask struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:100;not null" json:"name"`
	Description    string         `gorm:"type:text" json:"description"`
	CronExpression string         `gorm:"size:100;not null" json:"cron_expression"`
	TemplateID     *uint          `gorm:"index" json:"template_id,omitempty"`
	Content        string         `gorm:"type:text" json:"content"`
	Type           string         `gorm:"size:50;default:shell" json:"type"`
	AssetIDs       string         `gorm:"type:text" json:"asset_ids"`
	Timeout        int            `gorm:"default:300" json:"timeout"`
	Enabled        bool           `gorm:"default:false" json:"enabled"`
	UpdateAssets   bool           `gorm:"default:false" json:"update_assets"` // 是否更新资产信息
	LastRunAt      *time.Time     `json:"last_run_at,omitempty"`
	NextRunAt      *time.Time     `json:"next_run_at,omitempty"`
	LastStatus     string         `gorm:"size:50" json:"last_status,omitempty"`
	CreatedBy      uint           `gorm:"not null" json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Template *TaskTemplate `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Creator  *User         `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// ScheduledTaskExecution 定时任务执行记录（复用 TaskExecution 但额外添加关联字段）
// TaskExecution 表示每次执行的记录
// 增加 ScheduledTaskID 字段关联到定时任务
