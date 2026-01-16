// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package environment

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles environment management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new environment service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateEnvironmentRequest represents a request to create an environment
type CreateEnvironmentRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateEnvironmentRequest represents a request to update an environment
type UpdateEnvironmentRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// EnvironmentResponse represents an environment response
type EnvironmentResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateEnvironment creates a new environment
func (s *Service) CreateEnvironment(req *CreateEnvironmentRequest) (*EnvironmentResponse, error) {
	// Check if name already exists
	var existingEnvironment model.Environment
	if err := s.db.Where("name = ?", req.Name).First(&existingEnvironment).Error; err == nil {
		return nil, errors.New("environment name already exists")
	}

	environment := model.Environment{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.Create(&environment).Error; err != nil {
		return nil, err
	}

	return &EnvironmentResponse{
		ID:          environment.ID,
		Name:        environment.Name,
		Description: environment.Description,
		CreatedAt:   environment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   environment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetEnvironment retrieves an environment by ID
func (s *Service) GetEnvironment(id uint) (*EnvironmentResponse, error) {
	var environment model.Environment
	if err := s.db.First(&environment, id).Error; err != nil {
		return nil, err
	}

	return &EnvironmentResponse{
		ID:          environment.ID,
		Name:        environment.Name,
		Description: environment.Description,
		CreatedAt:   environment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   environment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListEnvironments retrieves all environments
func (s *Service) ListEnvironments() ([]EnvironmentResponse, error) {
	var environments []model.Environment
	if err := s.db.Find(&environments).Error; err != nil {
		return nil, err
	}

	result := make([]EnvironmentResponse, len(environments))
	for i, environment := range environments {
		result[i] = EnvironmentResponse{
			ID:          environment.ID,
			Name:        environment.Name,
			Description: environment.Description,
			CreatedAt:   environment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   environment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, nil
}

// UpdateEnvironment updates an environment
func (s *Service) UpdateEnvironment(id uint, req *UpdateEnvironmentRequest) (*EnvironmentResponse, error) {
	var environment model.Environment
	if err := s.db.First(&environment, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		// Check if new name already exists (excluding current environment)
		var existingEnvironment model.Environment
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingEnvironment).Error; err == nil {
			return nil, errors.New("environment name already exists")
		}
		environment.Name = req.Name
	}

	if req.Description != "" {
		environment.Description = req.Description
	}

	if err := s.db.Save(&environment).Error; err != nil {
		return nil, err
	}

	return &EnvironmentResponse{
		ID:          environment.ID,
		Name:        environment.Name,
		Description: environment.Description,
		CreatedAt:   environment.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   environment.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteEnvironment deletes an environment
func (s *Service) DeleteEnvironment(id uint) error {
	return s.db.Delete(&model.Environment{}, id).Error
}
