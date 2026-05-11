// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// Integration is a third-party connector definition with encrypted JSON credentials.
type Integration struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	Name            string         `gorm:"not null;size:128" json:"name"`
	Kind            string         `gorm:"not null;size:64;index" json:"kind"`
	ConfigEncrypted string         `gorm:"type:text;not null" json:"-"`
	Enabled         bool           `gorm:"default:true" json:"enabled"`
	Description     string         `gorm:"size:512" json:"description,omitempty"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies table name.
func (Integration) TableName() string {
	return "integrations"
}

// IntegrationPublic is the API-safe representation (no secrets).
type IntegrationPublic struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Kind        string    `json:"kind"`
	Enabled     bool      `json:"enabled"`
	Description string    `json:"description,omitempty"`
	HasConfig   bool      `json:"has_config"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
