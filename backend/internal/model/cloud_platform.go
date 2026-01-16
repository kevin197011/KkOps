// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// CloudPlatform represents a cloud platform (e.g., Aliyun, Tencent Cloud, AWS, Azure)
type CloudPlatform struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Assets []Asset `gorm:"foreignKey:CloudPlatformID" json:"assets,omitempty"`
}
