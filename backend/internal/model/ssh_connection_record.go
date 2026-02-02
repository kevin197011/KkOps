// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"
)

// SSHConnectionRecord WebSSH 连线记录（审计连线）
type SSHConnectionRecord struct {
	ID                  uint      `gorm:"primaryKey" json:"id"`
	UserID              uint      `gorm:"index" json:"user_id"`
	Username            string    `gorm:"size:100;index" json:"username"`
	AssetID             uint      `gorm:"index" json:"asset_id"`
	AssetHostname       string    `gorm:"size:100" json:"asset_hostname"`
	StartedAt           time.Time `gorm:"index" json:"started_at"`
	EndedAt             time.Time `gorm:"index" json:"ended_at"`
	DurationSeconds     int64     `json:"duration_seconds"`
	Transcript          string    `gorm:"type:text" json:"transcript,omitempty"` // JSON 数组 [{ "t": 0, "d": "base64..." }, ...]，PostgreSQL 用 text（无长度限制）
	TranscriptTruncated bool      `gorm:"default:false" json:"transcript_truncated"`
	CreatedAt           time.Time `gorm:"index" json:"created_at"`
}

// TableName 指定表名
func (SSHConnectionRecord) TableName() string {
	return "ssh_connection_records"
}
