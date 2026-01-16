// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package tag

import (
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles tag management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new tag service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateTagRequest represents a request to create a tag
type CreateTagRequest struct {
	Name        string `json:"name" binding:"required"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// UpdateTagRequest represents a request to update a tag
type UpdateTagRequest struct {
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
}

// TagResponse represents a tag response
type TagResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Color       string `json:"color"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTag creates a new tag
func (s *Service) CreateTag(req *CreateTagRequest) (*TagResponse, error) {
	tag := model.Tag{
		Name:        req.Name,
		Color:       req.Color,
		Description: req.Description,
	}

	if err := s.db.Create(&tag).Error; err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetTag retrieves a tag by ID
func (s *Service) GetTag(id uint) (*TagResponse, error) {
	var tag model.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListTags retrieves all tags
func (s *Service) ListTags() ([]TagResponse, error) {
	var tags []model.Tag
	if err := s.db.Find(&tags).Error; err != nil {
		return nil, err
	}

	result := make([]TagResponse, len(tags))
	for i, tag := range tags {
		result[i] = TagResponse{
			ID:          tag.ID,
			Name:        tag.Name,
			Color:       tag.Color,
			Description: tag.Description,
			CreatedAt:   tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, nil
}

// UpdateTag updates a tag
func (s *Service) UpdateTag(id uint, req *UpdateTagRequest) (*TagResponse, error) {
	var tag model.Tag
	if err := s.db.First(&tag, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		tag.Name = req.Name
	}

	if req.Color != "" {
		tag.Color = req.Color
	}

	if req.Description != "" {
		tag.Description = req.Description
	}

	if err := s.db.Save(&tag).Error; err != nil {
		return nil, err
	}

	return &TagResponse{
		ID:          tag.ID,
		Name:        tag.Name,
		Color:       tag.Color,
		Description: tag.Description,
		CreatedAt:   tag.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   tag.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteTag deletes a tag
func (s *Service) DeleteTag(id uint) error {
	return s.db.Delete(&model.Tag{}, id).Error
}
