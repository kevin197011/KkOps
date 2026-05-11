// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// ProvisioningTarget binds an Integration record to a concrete provisioning provider kind.
type ProvisioningTarget struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	IntegrationID uint           `gorm:"not null;index" json:"integration_id"`
	ProviderKind  string         `gorm:"not null;size:64;index" json:"provider_kind"` // scim, gitlab, ...
	Status        string         `gorm:"not null;size:32;default:idle" json:"status"` // idle, syncing, ok, error
	LastError     string         `gorm:"type:text" json:"last_error,omitempty"`
	LastSyncAt    *time.Time     `json:"last_sync_at,omitempty"`
	Enabled       bool           `gorm:"default:true" json:"enabled"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Integration *Integration `gorm:"foreignKey:IntegrationID" json:"integration,omitempty"`
}

func (ProvisioningTarget) TableName() string {
	return "provisioning_targets"
}

// ProvisioningRun stores history for a provisioning attempt on a target.
type ProvisioningRun struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	TargetID  uint      `gorm:"not null;index" json:"target_id"`
	UserID    *uint     `gorm:"index" json:"user_id,omitempty"` // null for full-target sync
	Action    string    `gorm:"not null;size:32" json:"action"` // sync, delete, verify
	Status    string    `gorm:"not null;size:32" json:"status"` // success, partial, failed
	Message   string    `gorm:"type:text" json:"message,omitempty"`
	StartedAt time.Time `json:"started_at"`
	EndedAt   time.Time `json:"ended_at"`
}

func (ProvisioningRun) TableName() string {
	return "provisioning_runs"
}

// ProvisioningUserLink maps a local user to an external principal identifier per target.
type ProvisioningUserLink struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       uint      `gorm:"not null;uniqueIndex:idx_user_target;index" json:"user_id"`
	TargetID     uint      `gorm:"not null;uniqueIndex:idx_user_target;index" json:"target_id"`
	ExternalID   string    `gorm:"not null;size:512" json:"external_id"`
	LastSyncedAt time.Time `json:"last_synced_at"`
}

func (ProvisioningUserLink) TableName() string {
	return "provisioning_user_links"
}
