// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package aisvc orchestrates AI chat, RCA, and anomaly workflows.
package aisvc

import (
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	alertsvc "github.com/kkops/backend/internal/service/alert"
	incidentsvc "github.com/kkops/backend/internal/service/incident"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Service wires persistence and tool execution for AI features.
type Service struct {
	DB          *gorm.DB
	Log         *zap.Logger
	Tools       *ToolBridge
	Integration *integrationsvc.Service
	Incidents   *incidentsvc.Service
	Alerts      *alertsvc.Service
}

// NewService constructs the AI orchestration service.
func NewService(db *gorm.DB, log *zap.Logger, integ *integrationsvc.Service, alerts *alertsvc.Service, inc *incidentsvc.Service, bridge *ToolBridge) *Service {
	return &Service{
		DB:          db,
		Log:         log,
		Integration: integ,
		Alerts:      alerts,
		Incidents:   inc,
		Tools:       bridge,
	}
}

// --- Sessions ---

// ListSessions returns chat sessions for a user.
func (s *Service) ListSessions(userID uint, limit int) ([]model.AIChatSession, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	var rows []model.AIChatSession
	err := s.DB.Where("user_id = ?", userID).Order("updated_at DESC").Limit(limit).Find(&rows).Error
	return rows, err
}

// CreateSession creates an empty session titled by first message snippet.
func (s *Service) CreateSession(userID uint, title string) (*model.AIChatSession, error) {
	if title == "" {
		title = "Chat"
	}
	row := model.AIChatSession{UserID: userID, Title: title}
	if err := s.DB.Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// GetSession loads session if owned by user.
func (s *Service) GetSession(userID, sessionID uint) (*model.AIChatSession, []model.AIChatMessage, error) {
	var sess model.AIChatSession
	if err := s.DB.Where("id = ? AND user_id = ?", sessionID, userID).First(&sess).Error; err != nil {
		return nil, nil, err
	}
	var msgs []model.AIChatMessage
	if err := s.DB.Where("session_id = ?", sessionID).Order("id asc").Find(&msgs).Error; err != nil {
		return nil, nil, err
	}
	return &sess, msgs, nil
}

// DeleteSession removes session and messages.
func (s *Service) DeleteSession(userID, sessionID uint) error {
	res := s.DB.Where("id = ? AND user_id = ?", sessionID, userID).Delete(&model.AIChatSession{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("session not found")
	}
	return s.DB.Where("session_id = ?", sessionID).Delete(&model.AIChatMessage{}).Error
}

// AppendMessage stores one chat message.
func (s *Service) AppendMessage(sessionID uint, role, content string) (*model.AIChatMessage, error) {
	m := model.AIChatMessage{SessionID: sessionID, Role: role, Content: content}
	if err := s.DB.Create(&m).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

// TouchSessionTitle updates session title from first message snippet.
func (s *Service) TouchSessionTitle(sessionID uint, snippet string) {
	snippet = strings.TrimSpace(snippet)
	if idx := strings.IndexAny(snippet, "\r\n"); idx >= 0 {
		snippet = snippet[:idx]
	}
	snippet = strings.TrimSpace(snippet)
	if snippet == "" {
		return
	}
	if len([]rune(snippet)) > 80 {
		rs := []rune(snippet)
		snippet = string(rs[:80]) + "…"
	}
	_ = s.DB.Model(&model.AIChatSession{}).Where("id = ?", sessionID).Update("title", snippet).Error
}

// CitedCall records one tool invocation for RCA payloads.
type CitedCall struct {
	Tool string          `json:"tool"`
	Args json.RawMessage `json:"args,omitempty"`
}
