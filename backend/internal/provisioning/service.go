// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provisioning

import (
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// APIService backs HTTP handlers for targets and run history.
type APIService struct {
	db    *gorm.DB
	reg   *Registry
	coord *Coordinator
}

// NewAPIService creates the provisioning API service.
func NewAPIService(db *gorm.DB, reg *Registry, coord *Coordinator) *APIService {
	return &APIService{db: db, reg: reg, coord: coord}
}

// CreateTargetRequest binds an integration to a provider kind.
type CreateTargetRequest struct {
	IntegrationID uint   `json:"integration_id" binding:"required"`
	ProviderKind  string `json:"provider_kind" binding:"required"`
	Enabled       *bool  `json:"enabled"`
}

// ListTargets returns all provisioning targets with integration preloaded.
func (s *APIService) ListTargets() ([]model.ProvisioningTarget, error) {
	var rows []model.ProvisioningTarget
	if err := s.db.Preload("Integration").Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// CreateTarget validates and inserts a target row.
func (s *APIService) CreateTarget(req *CreateTargetRequest) (*model.ProvisioningTarget, error) {
	if _, err := s.reg.Create(req.ProviderKind, []byte(`{"base_url":""}`)); err != nil {
		return nil, fmt.Errorf("invalid provider_kind: %w", err)
	}
	var integ model.Integration
	if err := s.db.First(&integ, req.IntegrationID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("integration not found")
		}
		return nil, err
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	t := model.ProvisioningTarget{
		IntegrationID: req.IntegrationID,
		ProviderKind:  req.ProviderKind,
		Status:        "idle",
		Enabled:       enabled,
	}
	if err := s.db.Create(&t).Error; err != nil {
		return nil, err
	}
	_ = s.db.Preload("Integration").First(&t, t.ID).Error
	return &t, nil
}

// ListRuns returns recent runs for a target (newest first).
func (s *APIService) ListRuns(targetID uint, limit int) ([]model.ProvisioningRun, error) {
	if limit <= 0 || limit > 500 {
		limit = 50
	}
	var runs []model.ProvisioningRun
	q := s.db.Where("target_id = ?", targetID).Order("id desc").Limit(limit)
	if err := q.Find(&runs).Error; err != nil {
		return nil, err
	}
	return runs, nil
}

// RequestManualSync enqueues a full target sync.
func (s *APIService) RequestManualSync(targetID uint) error {
	var t model.ProvisioningTarget
	if err := s.db.First(&t, targetID).Error; err != nil {
		return err
	}
	s.coord.EnqueueTargetSync(targetID)
	return nil
}
