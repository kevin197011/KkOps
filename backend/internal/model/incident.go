// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// IncidentStatus enumerates lifecycle states stored as string.
const (
	IncidentStatusOpen         = "open"
	IncidentStatusAcknowledged = "acknowledged"
	IncidentStatusResolved     = "resolved"
)

// Incident is a minimal operational incident ticket.
type Incident struct {
	ID                 uint           `gorm:"primaryKey" json:"id"`
	Title              string         `gorm:"size:512;not null" json:"title"`
	Severity           string         `gorm:"size:32;not null;index" json:"severity"`
	Status             string         `gorm:"size:32;not null;index" json:"status"`
	LinkedAlertIDsJSON string         `gorm:"type:text" json:"-"` // JSON array of uint64 alert IDs
	AssigneeUserID     *uint          `gorm:"index" json:"assignee_user_id,omitempty"`
	CreatedBy          uint           `gorm:"not null;index" json:"created_by"`
	CreatedAt          time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies table name.
func (Incident) TableName() string {
	return "incidents"
}
