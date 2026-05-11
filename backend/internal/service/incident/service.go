// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package incident

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles incident CRUD.
type Service struct {
	db *gorm.DB
}

// NewService constructs the incident service.
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateInput holds fields for new incidents.
type CreateInput struct {
	Title          string `json:"title"`
	Severity       string `json:"severity"`
	LinkedAlertIDs []uint `json:"linked_alert_ids,omitempty"`
	AssigneeUserID *uint  `json:"assignee_user_id,omitempty"`
	CreatedBy      uint   `json:"-"`
}

// Create persists a new incident.
func (s *Service) Create(ctx context.Context, in CreateInput) (*IncidentView, error) {
	if in.Title == "" || in.Severity == "" {
		return nil, fmt.Errorf("title and severity are required")
	}
	raw, err := json.Marshal(in.LinkedAlertIDs)
	if err != nil {
		return nil, err
	}
	inc := model.Incident{
		Title:              in.Title,
		Severity:           in.Severity,
		Status:             model.IncidentStatusOpen,
		LinkedAlertIDsJSON: string(raw),
		AssigneeUserID:     in.AssigneeUserID,
		CreatedBy:          in.CreatedBy,
	}
	if err := s.db.Create(&inc).Error; err != nil {
		return nil, err
	}
	return viewFromModel(&inc)
}

// IncidentView includes decoded linked alert IDs for APIs.
type IncidentView struct {
	ID             uint      `json:"id"`
	Title          string    `json:"title"`
	Severity       string    `json:"severity"`
	Status         string    `json:"status"`
	LinkedAlertIDs []uint    `json:"linked_alert_ids,omitempty"`
	AssigneeUserID *uint     `json:"assignee_user_id,omitempty"`
	CreatedBy      uint      `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func viewFromModel(m *model.Incident) (*IncidentView, error) {
	v := &IncidentView{
		ID:             m.ID,
		Title:          m.Title,
		Severity:       m.Severity,
		Status:         m.Status,
		AssigneeUserID: m.AssigneeUserID,
		CreatedBy:      m.CreatedBy,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
	if strings.TrimSpace(m.LinkedAlertIDsJSON) != "" {
		var ids []uint
		if err := json.Unmarshal([]byte(m.LinkedAlertIDsJSON), &ids); err != nil {
			return nil, err
		}
		v.LinkedAlertIDs = ids
	}
	return v, nil
}

// List returns incidents with pagination.
func (s *Service) List(_ context.Context, status string, limit, offset int) ([]IncidentView, int64, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	q := s.db.Model(&model.Incident{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.Incident
	if err := q.Order("updated_at DESC").Limit(limit).Offset(offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]IncidentView, 0, len(rows))
	for i := range rows {
		v, err := viewFromModel(&rows[i])
		if err != nil {
			return nil, 0, err
		}
		out = append(out, *v)
	}
	return out, total, nil
}

// Get returns one incident by ID.
func (s *Service) Get(_ context.Context, id uint) (*IncidentView, error) {
	var m model.Incident
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return viewFromModel(&m)
}

// PatchInput contains optional updates.
type PatchInput struct {
	Title          *string `json:"title,omitempty"`
	Severity       *string `json:"severity,omitempty"`
	Status         *string `json:"status,omitempty"`
	LinkedAlertIDs *[]uint `json:"linked_alert_ids,omitempty"`
	AssigneeUserID *uint   `json:"assignee_user_id,omitempty"`
}

// Patch updates allowed fields.
func (s *Service) Patch(_ context.Context, id uint, p PatchInput) (*IncidentView, error) {
	var m model.Incident
	if err := s.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	if p.Title != nil {
		m.Title = *p.Title
	}
	if p.Severity != nil {
		m.Severity = *p.Severity
	}
	if p.Status != nil {
		if !validIncidentStatus(*p.Status) {
			return nil, fmt.Errorf("invalid status")
		}
		m.Status = *p.Status
	}
	if p.LinkedAlertIDs != nil {
		raw, err := json.Marshal(*p.LinkedAlertIDs)
		if err != nil {
			return nil, err
		}
		m.LinkedAlertIDsJSON = string(raw)
	}
	if p.AssigneeUserID != nil {
		m.AssigneeUserID = p.AssigneeUserID
	}
	if err := s.db.Save(&m).Error; err != nil {
		return nil, err
	}
	return viewFromModel(&m)
}

func validIncidentStatus(s string) bool {
	switch s {
	case model.IncidentStatusOpen, model.IncidentStatusAcknowledged, model.IncidentStatusResolved:
		return true
	default:
		return false
	}
}

// ErrNotFound wraps gorm not found for handlers.
func ErrNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
