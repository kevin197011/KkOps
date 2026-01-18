// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"
)

// AuditLog 审计日志
type AuditLog struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`                    // 操作用户 ID
	Username     string    `gorm:"size:100;index" json:"username"`          // 操作用户名（冗余存储）
	Action       string    `gorm:"size:50;index" json:"action"`             // 操作类型：create, update, delete, execute, login, logout
	Module       string    `gorm:"size:50;index" json:"module"`             // 模块：user, role, asset, task, deployment, auth
	ResourceID   *uint     `gorm:"index" json:"resource_id,omitempty"`      // 资源 ID（可选）
	ResourceName string    `gorm:"size:200" json:"resource_name,omitempty"` // 资源名称（冗余存储）
	Detail       string    `gorm:"type:text" json:"detail,omitempty"`       // 操作详情（JSON 格式）
	IPAddress    string    `gorm:"size:50" json:"ip_address"`               // 客户端 IP
	UserAgent    string    `gorm:"size:500" json:"user_agent,omitempty"`    // 用户代理
	Status       string    `gorm:"size:20;index" json:"status"`             // 操作结果：success, failed
	ErrorMsg     string    `gorm:"type:text" json:"error_msg,omitempty"`    // 错误信息（失败时）
	CreatedAt    time.Time `gorm:"index" json:"created_at"`                 // 操作时间
}

// TableName 指定表名
func (AuditLog) TableName() string {
	return "audit_logs"
}

// AuditAction 审计操作类型
type AuditAction string

const (
	AuditActionCreate  AuditAction = "create"
	AuditActionUpdate  AuditAction = "update"
	AuditActionDelete  AuditAction = "delete"
	AuditActionExecute AuditAction = "execute"
	AuditActionLogin   AuditAction = "login"
	AuditActionLogout  AuditAction = "logout"
	AuditActionEnable  AuditAction = "enable"
	AuditActionDisable AuditAction = "disable"
	AuditActionExport  AuditAction = "export"
	AuditActionConnect AuditAction = "connect"
)

// AuditModule 审计模块
type AuditModule string

const (
	AuditModuleAuth       AuditModule = "auth"
	AuditModuleUser       AuditModule = "user"
	AuditModuleRole       AuditModule = "role"
	AuditModuleAsset      AuditModule = "asset"
	AuditModuleTask       AuditModule = "task"
	AuditModuleTemplate   AuditModule = "template"
	AuditModuleScheduled  AuditModule = "scheduled_task"
	AuditModuleDeployment AuditModule = "deployment"
	AuditModuleSSH        AuditModule = "ssh"
	AuditModuleSSHKey     AuditModule = "ssh_key"
	AuditModuleProject    AuditModule = "project"
	AuditModuleEnv        AuditModule = "environment"
	AuditModuleCloud      AuditModule = "cloud_platform"
	AuditModuleTag        AuditModule = "tag"
)

// AuditStatus 审计状态
type AuditStatus string

const (
	AuditStatusSuccess AuditStatus = "success"
	AuditStatusFailed  AuditStatus = "failed"
)
