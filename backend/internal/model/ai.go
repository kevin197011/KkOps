// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// AIChatSession stores a user's assistant conversation thread.
type AIChatSession struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	UserID    uint           `gorm:"not null;index" json:"user_id"`
	Title     string         `gorm:"size:256" json:"title"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies table name.
func (AIChatSession) TableName() string {
	return "ai_chat_sessions"
}

// AIChatMessage is one message in a chat session.
type AIChatMessage struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	SessionID uint      `gorm:"not null;index" json:"session_id"`
	Role      string    `gorm:"size:32;not null" json:"role"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies table name.
func (AIChatMessage) TableName() string {
	return "ai_chat_messages"
}

// AIAnomalyRule schedules metric anomaly checks via LLM.
type AIAnomalyRule struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:128;not null" json:"name"`
	IntegrationID  uint           `gorm:"not null;index" json:"integration_id"`
	Query          string         `gorm:"type:text;not null" json:"query"`
	ScheduleCron   string         `gorm:"size:64;not null" json:"schedule_cron"`
	Enabled        bool           `gorm:"default:true" json:"enabled"`
	PromptTemplate string         `gorm:"type:text" json:"prompt_template"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies table name.
func (AIAnomalyRule) TableName() string {
	return "ai_anomaly_rules"
}

// AIAnomalyFinding is one evaluated anomaly result.
type AIAnomalyFinding struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RuleID    uint      `gorm:"not null;index" json:"rule_id"`
	Ts        time.Time `gorm:"not null;index" json:"ts"`
	Severity  string    `gorm:"size:32;not null" json:"severity"`
	Summary   string    `gorm:"type:text" json:"summary"`
	Raw       string    `gorm:"type:text" json:"raw"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName specifies table name.
func (AIAnomalyFinding) TableName() string {
	return "ai_anomaly_findings"
}

// AIRcaReport stores generated root-cause analysis for an incident.
type AIRcaReport struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	IncidentID  uint      `gorm:"not null;index" json:"incident_id"`
	GeneratedAt time.Time `gorm:"not null" json:"generated_at"`
	ReportMD    string    `gorm:"type:text;not null" json:"report_md"`
	RawToolLog  string    `gorm:"type:text" json:"raw_tool_log"`
	CreatedAt   time.Time `json:"created_at"`
}

// TableName specifies table name.
func (AIRcaReport) TableName() string {
	return "ai_rca_reports"
}
