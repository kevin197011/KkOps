package models

import (
	"time"

	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;size:50;not null" json:"username"`
	Email        string         `gorm:"uniqueIndex;size:100;not null" json:"email"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	DisplayName  string         `gorm:"size:100" json:"display_name"`
	Status       string         `gorm:"size:20;not null;default:'active'" json:"status"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	LastLoginIP  string         `gorm:"size:45" json:"last_login_ip"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Roles []Role `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}

// Role 角色模型
type Role struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Users       []User        `gorm:"many2many:user_roles;" json:"users,omitempty"`
}

// TableName 指定表名
func (Role) TableName() string {
	return "roles"
}

// Permission 权限模型
type Permission struct {
	ID           uint64    `gorm:"primaryKey" json:"id"`
	Code         string    `gorm:"uniqueIndex;size:100;not null" json:"code"`
	Name         string    `gorm:"size:100;not null" json:"name"`
	ResourceType string    `gorm:"size:50;not null" json:"resource_type"`
	Action       string    `gorm:"size:50;not null" json:"action"`
	Description  string    `json:"description"`
	CreatedAt    time.Time `json:"created_at"`

	// 关联
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// TableName 指定表名
func (Permission) TableName() string {
	return "permissions"
}

// UserRole 用户角色关联模型
type UserRole struct {
	ID        uint64    `gorm:"primaryKey"`
	UserID    uint64    `gorm:"not null;index"`
	RoleID    uint64    `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (UserRole) TableName() string {
	return "user_roles"
}

// RolePermission 角色权限关联模型
type RolePermission struct {
	ID           uint64    `gorm:"primaryKey"`
	RoleID       uint64    `gorm:"not null;index"`
	PermissionID uint64    `gorm:"not null;index"`
	CreatedAt    time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (RolePermission) TableName() string {
	return "role_permissions"
}

