// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package cloudplatform

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles cloud platform management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new cloud platform service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateCloudPlatformRequest represents a request to create a cloud platform
type CreateCloudPlatformRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateCloudPlatformRequest represents a request to update a cloud platform
type UpdateCloudPlatformRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CloudPlatformResponse represents a cloud platform response
type CloudPlatformResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateCloudPlatform creates a new cloud platform
func (s *Service) CreateCloudPlatform(req *CreateCloudPlatformRequest) (*CloudPlatformResponse, error) {
	// Check if name already exists
	var existingCloudPlatform model.CloudPlatform
	if err := s.db.Where("name = ?", req.Name).First(&existingCloudPlatform).Error; err == nil {
		return nil, errors.New("cloud platform name already exists")
	}

	cloudPlatform := model.CloudPlatform{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.Create(&cloudPlatform).Error; err != nil {
		return nil, err
	}

	return &CloudPlatformResponse{
		ID:          cloudPlatform.ID,
		Name:        cloudPlatform.Name,
		Description: cloudPlatform.Description,
		CreatedAt:   cloudPlatform.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   cloudPlatform.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetCloudPlatform retrieves a cloud platform by ID
func (s *Service) GetCloudPlatform(id uint) (*CloudPlatformResponse, error) {
	var cloudPlatform model.CloudPlatform
	if err := s.db.First(&cloudPlatform, id).Error; err != nil {
		return nil, err
	}

	return &CloudPlatformResponse{
		ID:          cloudPlatform.ID,
		Name:        cloudPlatform.Name,
		Description: cloudPlatform.Description,
		CreatedAt:   cloudPlatform.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   cloudPlatform.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListCloudPlatforms retrieves all cloud platforms
func (s *Service) ListCloudPlatforms() ([]CloudPlatformResponse, error) {
	var cloudPlatforms []model.CloudPlatform
	if err := s.db.Find(&cloudPlatforms).Error; err != nil {
		return nil, err
	}

	result := make([]CloudPlatformResponse, len(cloudPlatforms))
	for i, cloudPlatform := range cloudPlatforms {
		result[i] = CloudPlatformResponse{
			ID:          cloudPlatform.ID,
			Name:        cloudPlatform.Name,
			Description: cloudPlatform.Description,
			CreatedAt:   cloudPlatform.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   cloudPlatform.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, nil
}

// UpdateCloudPlatform updates a cloud platform
func (s *Service) UpdateCloudPlatform(id uint, req *UpdateCloudPlatformRequest) (*CloudPlatformResponse, error) {
	var cloudPlatform model.CloudPlatform
	if err := s.db.First(&cloudPlatform, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		// Check if new name already exists (excluding current cloud platform)
		var existingCloudPlatform model.CloudPlatform
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingCloudPlatform).Error; err == nil {
			return nil, errors.New("cloud platform name already exists")
		}
		cloudPlatform.Name = req.Name
	}

	if req.Description != "" {
		cloudPlatform.Description = req.Description
	}

	if err := s.db.Save(&cloudPlatform).Error; err != nil {
		return nil, err
	}

	return &CloudPlatformResponse{
		ID:          cloudPlatform.ID,
		Name:        cloudPlatform.Name,
		Description: cloudPlatform.Description,
		CreatedAt:   cloudPlatform.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   cloudPlatform.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteCloudPlatform deletes a cloud platform
func (s *Service) DeleteCloudPlatform(id uint) error {
	return s.db.Delete(&model.CloudPlatform{}, id).Error
}
