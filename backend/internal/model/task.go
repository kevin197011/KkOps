// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// TaskTemplate represents a task template
type TaskTemplate struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Content     string         `gorm:"type:text" json:"content"` // Script or command content
	Type        string         `gorm:"size:50" json:"type"`      // shell, python, etc.
	CreatedBy   uint           `json:"created_by"`
	Creator     User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Tasks []Task `gorm:"foreignKey:TemplateID" json:"tasks,omitempty"`
}

// Task represents an execution task
type Task struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TemplateID  *uint          `json:"template_id"`
	Template    *TaskTemplate  `gorm:"foreignKey:TemplateID" json:"template,omitempty"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	Content     string         `gorm:"type:text" json:"content"`                    // Script or command content
	Type        string         `gorm:"size:50" json:"type"`                         // shell, python, etc.
	Timeout     int            `gorm:"default:600" json:"timeout"`                  // Execution timeout in seconds (default 10 minutes)
	Status      string         `gorm:"default:pending;size:20;index" json:"status"` // pending, running, success, failed, cancelled
	AssetIDs    string         `gorm:"type:text" json:"asset_ids_str"`              // Comma-separated asset IDs for execution
	CreatedBy   uint           `json:"created_by"`
	Creator     User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	StartedAt   *time.Time     `json:"started_at"`
	FinishedAt  *time.Time     `json:"finished_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Executions []TaskExecution `gorm:"foreignKey:TaskID" json:"executions,omitempty"`
}

// TaskExecution represents a task execution on a specific host
type TaskExecution struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	TaskID          *uint          `gorm:"index" json:"task_id,omitempty"`                      // 关联的执行任务 ID（可为空）
	Task            *Task          `gorm:"foreignKey:TaskID" json:"task,omitempty"`
	ScheduledTaskID *uint          `gorm:"index" json:"scheduled_task_id,omitempty"`            // 关联的定时任务 ID（可为空）
	AssetID         uint           `gorm:"not null;index" json:"asset_id"`
	Asset           Asset          `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
	TriggerType     string         `gorm:"size:20;default:manual" json:"trigger_type"`          // manual, scheduled
	Status          string         `gorm:"default:pending;size:20;index" json:"status"`         // pending, running, success, failed, cancelled
	ExitCode        *int           `json:"exit_code"`
	Output          string         `gorm:"type:text" json:"output"`                             // Command output
	Error           string         `gorm:"type:text" json:"error"`                              // Error message
	StartedAt       *time.Time     `json:"started_at"`
	FinishedAt      *time.Time     `json:"finished_at"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
