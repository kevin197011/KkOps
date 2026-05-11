// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// Alert source / status constants for normalized alerts.
const (
	AlertSourcePrometheusAM = "prometheus_alertmanager"
	AlertSourceNightingale  = "nightingale"
	AlertSourceWebhook      = "webhook"

	AlertStatusFiring       = "firing"
	AlertStatusResolved     = "resolved"
	AlertStatusAcknowledged = "acknowledged"
	AlertStatusDismissed    = "dismissed"
)

// AlertRecord stores a normalized alert row (Alertmanager, Nightingale, or webhook).
type AlertRecord struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	IntegrationID uint           `gorm:"index;default:0" json:"integration_id"` // 0 = webhook / unknown
	Source        string         `gorm:"size:64;index;not null" json:"source"`  // prometheus_alertmanager, nightingale, webhook
	Fingerprint   string         `gorm:"size:512;index;not null" json:"fingerprint"`
	Severity      string         `gorm:"size:32;not null" json:"severity"`
	Title         string         `gorm:"size:512;not null" json:"title"`
	Status        string         `gorm:"size:32;not null;index" json:"status"` // firing, resolved, acknowledged, dismissed
	LabelsJSON    string         `gorm:"type:text" json:"labels_json,omitempty"`
	StartsAt      *time.Time     `json:"starts_at,omitempty"`
	EndsAt        *time.Time     `json:"ends_at,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
	Integration   *Integration   `gorm:"foreignKey:IntegrationID" json:"-"`
}

// TableName specifies table name.
func (AlertRecord) TableName() string {
	return "alert_records"
}
