package models

import (
	"time"

	"gorm.io/datatypes"
)

// BatchOperation 批量操作模型
type BatchOperation struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"size:255;not null" json:"name"`
	Description     string         `gorm:"type:text" json:"description"`
	CommandType     string         `gorm:"size:50;not null" json:"command_type"` // custom, template, builtin
	CommandFunction string         `gorm:"size:255;not null" json:"command_function"`
	CommandArgs     datatypes.JSON `gorm:"type:jsonb" json:"command_args"`
	TargetHosts     datatypes.JSON `gorm:"type:jsonb;not null" json:"target_hosts"`
	TargetCount     int            `gorm:"not null" json:"target_count"`
	Status          string         `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, running, completed, failed, cancelled
	SaltJobID       string         `gorm:"size:255" json:"salt_job_id"`
	StartedBy       uint64         `gorm:"not null;index" json:"started_by"`
	StartedAt       time.Time      `gorm:"not null;default:now();index" json:"started_at"`
	CompletedAt     *time.Time     `json:"completed_at"`
	DurationSeconds *int           `json:"duration_seconds"`
	Results         datatypes.JSON `gorm:"type:jsonb" json:"results"`
	SuccessCount    int            `gorm:"default:0" json:"success_count"`
	FailedCount     int            `gorm:"default:0" json:"failed_count"`
	ErrorMessage    string         `gorm:"type:text" json:"error_message"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`

	// 关联
	StartedByUser User `gorm:"foreignKey:StartedBy" json:"started_by_user,omitempty"`
}

// TableName 指定表名
func (BatchOperation) TableName() string {
	return "batch_operations"
}

