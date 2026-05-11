// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package integration

import (
	"encoding/json"
	"errors"
	"fmt"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// Service manages integration records.
type Service struct {
	db     *gorm.DB
	config *config.Config
}

// NewService creates an integration service.
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{db: db, config: cfg}
}

// CreateRequest holds plaintext config for create/update (never stored verbatim).
type CreateRequest struct {
	Name        string                 `json:"name"`
	Kind        string                 `json:"kind"`
	Enabled     *bool                  `json:"enabled"`
	Description *string                `json:"description"`
	Config      map[string]interface{} `json:"config"`
}

func (s *Service) toPublic(in *model.Integration) model.IntegrationPublic {
	hasConfig := in.ConfigEncrypted != ""
	return model.IntegrationPublic{
		ID:          in.ID,
		Name:        in.Name,
		Kind:        in.Kind,
		Enabled:     in.Enabled,
		Description: in.Description,
		HasConfig:   hasConfig,
		CreatedAt:   in.CreatedAt,
		UpdatedAt:   in.UpdatedAt,
	}
}

// List returns all integrations (public projection).
func (s *Service) List() ([]model.IntegrationPublic, error) {
	var rows []model.Integration
	if err := s.db.Order("id asc").Find(&rows).Error; err != nil {
		return nil, err
	}
	out := make([]model.IntegrationPublic, 0, len(rows))
	for i := range rows {
		out = append(out, s.toPublic(&rows[i]))
	}
	return out, nil
}

// Get returns one integration (public projection).
func (s *Service) Get(id uint) (*model.IntegrationPublic, error) {
	var row model.Integration
	if err := s.db.First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("integration not found")
		}
		return nil, err
	}
	p := s.toPublic(&row)
	return &p, nil
}

// Create stores integration with encrypted config JSON.
func (s *Service) Create(req *CreateRequest) (*model.IntegrationPublic, error) {
	if req.Name == "" || req.Kind == "" {
		return nil, fmt.Errorf("name and kind are required")
	}
	raw, err := json.Marshal(req.Config)
	if err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	enc, err := utils.Encrypt(raw, s.config.Encryption.Key)
	if err != nil {
		return nil, fmt.Errorf("encrypt config: %w", err)
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	desc := ""
	if req.Description != nil {
		desc = *req.Description
	}
	row := model.Integration{
		Name:            req.Name,
		Kind:            req.Kind,
		ConfigEncrypted: enc,
		Enabled:         enabled,
		Description:     desc,
	}
	if err := s.db.Create(&row).Error; err != nil {
		return nil, err
	}
	p := s.toPublic(&row)
	return &p, nil
}

// Update replaces mutable fields and optionally config.
func (s *Service) Update(id uint, req *CreateRequest) (*model.IntegrationPublic, error) {
	var row model.Integration
	if err := s.db.First(&row, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("integration not found")
		}
		return nil, err
	}
	if req.Name != "" {
		row.Name = req.Name
	}
	if req.Kind != "" {
		row.Kind = req.Kind
	}
	if req.Description != nil {
		row.Description = *req.Description
	}
	if req.Enabled != nil {
		row.Enabled = *req.Enabled
	}
	if req.Config != nil {
		raw, err := json.Marshal(req.Config)
		if err != nil {
			return nil, fmt.Errorf("invalid config: %w", err)
		}
		enc, err := utils.Encrypt(raw, s.config.Encryption.Key)
		if err != nil {
			return nil, fmt.Errorf("encrypt config: %w", err)
		}
		row.ConfigEncrypted = enc
	}
	if err := s.db.Save(&row).Error; err != nil {
		return nil, err
	}
	p := s.toPublic(&row)
	return &p, nil
}

// Delete removes an integration (soft delete).
func (s *Service) Delete(id uint) error {
	res := s.db.Delete(&model.Integration{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return fmt.Errorf("integration not found")
	}
	return nil
}

// DecryptConfigForWorker returns decrypted JSON for internal connector use only (not HTTP API).
func (s *Service) DecryptConfigForWorker(id uint) ([]byte, error) {
	var row model.Integration
	if err := s.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return utils.Decrypt(row.ConfigEncrypted, s.config.Encryption.Key)
}
