// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// Role represents a role in the RBAC system
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Users       []User       `gorm:"many2many:user_roles;" json:"users,omitempty"`
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}

// Permission represents a permission in the RBAC system
type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Resource    string         `gorm:"not null;size:100" json:"resource"` // e.g., "user", "asset", "task"
	Action      string         `gorm:"not null;size:50" json:"action"`    // e.g., "create", "read", "update", "delete"
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Roles []Role `gorm:"many2many:role_permissions;" json:"roles,omitempty"`
}

// Unique constraint on (resource, action) is handled by database migration

// UserRole represents the many-to-many relationship between users and roles
type UserRole struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
}

// RolePermission represents the many-to-many relationship between roles and permissions
type RolePermission struct {
	RoleID       uint      `gorm:"primaryKey" json:"role_id"`
	PermissionID uint      `gorm:"primaryKey" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
}
