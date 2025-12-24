package models

import (
	"time"

	"gorm.io/gorm"
)

// Project 项目模型
type Project struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"size:100;not null;index" json:"name"`
	Description string         `json:"description"`
	Status      string         `gorm:"size:20;not null;default:'active';index" json:"status"`
	CreatedBy   uint64         `gorm:"not null;index" json:"created_by"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Creator *User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName 指定表名
func (Project) TableName() string {
	return "projects"
}

// ProjectMember 项目成员关联模型
type ProjectMember struct {
	ID        uint64    `gorm:"primaryKey"`
	ProjectID uint64    `gorm:"not null;index"`
	UserID    uint64    `gorm:"not null;index"`
	Role      string    `gorm:"size:50;not null" json:"role"` // owner, admin, member, viewer
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (ProjectMember) TableName() string {
	return "project_members"
}
