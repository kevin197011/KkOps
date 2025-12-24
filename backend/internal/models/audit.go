package models

import (
	"time"

	"gorm.io/gorm"
)

// AuditLog 审计日志模型
type AuditLog struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	UserID       *uint64        `gorm:"index" json:"user_id"` // 可为空，系统操作
	Username     string         `gorm:"size:50;index" json:"username"` // 冗余字段，便于查询
	Action       string         `gorm:"size:50;not null;index" json:"action"` // login, create, update, delete, execute
	ResourceType string         `gorm:"size:50;index" json:"resource_type"` // host, deployment, task, user, etc.
	ResourceID   *uint64        `gorm:"index" json:"resource_id"`
	ResourceName string         `gorm:"size:255" json:"resource_name"` // 冗余字段
	IPAddress    string         `gorm:"size:45" json:"ip_address"`
	UserAgent    string         `gorm:"type:text" json:"user_agent"`
	RequestData  string         `gorm:"type:jsonb" json:"request_data"` // JSON格式
	ResponseData string         `gorm:"type:jsonb" json:"response_data"` // JSON格式
	BeforeData   string         `gorm:"type:jsonb" json:"before_data"` // 变更前数据
	AfterData    string         `gorm:"type:jsonb" json:"after_data"` // 变更后数据
	Status       string         `gorm:"size:20;not null;default:'success';index" json:"status"` // success, failed
	ErrorMessage string         `gorm:"type:text" json:"error_message"`
	DurationMs   int            `json:"duration_ms"` // 操作耗时(毫秒)
	CreatedAt    time.Time      `gorm:"index" json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

