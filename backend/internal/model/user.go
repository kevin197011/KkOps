// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// User represents a user account
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:100" json:"username"`
	PasswordHash string         `gorm:"not null;size:255" json:"-"`
	Email        string         `gorm:"uniqueIndex;not null;size:255" json:"email"`
	Phone        string         `gorm:"size:20" json:"phone"`
	RealName     string         `gorm:"size:100" json:"real_name"`
	DepartmentID *uint          `json:"department_id"`
	Department   *Department    `gorm:"foreignKey:DepartmentID" json:"department,omitempty"`
	Status       string         `gorm:"default:active;size:20" json:"status"` // active, disabled
	AvatarURL    string         `gorm:"size:500" json:"avatar_url"`
	LastLoginAt  *time.Time     `json:"last_login_at"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Roles     []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
	APITokens []APIToken `gorm:"foreignKey:UserID" json:"api_tokens,omitempty"`
}

// Department represents an organizational department
type Department struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	ParentID    *uint          `json:"parent_id"`
	Parent      *Department    `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Code        string         `gorm:"size:50" json:"code"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Children []Department `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Users    []User       `gorm:"foreignKey:DepartmentID" json:"users,omitempty"`
}

// APIToken represents an API token for programmatic access
type APIToken struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"`
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	TokenHash   string         `gorm:"uniqueIndex;not null;size:255" json:"-"`
	ExpiresAt   *time.Time     `json:"expires_at"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	Status      string         `gorm:"default:active;size:20" json:"status"` // active, disabled
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}
