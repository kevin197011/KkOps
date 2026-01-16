// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package category

import (
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles asset category management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new category service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateCategoryRequest represents a request to create a category
type CreateCategoryRequest struct {
	Name        string `json:"name" binding:"required"`
	ParentID    *uint  `json:"parent_id"`
	Description string `json:"description"`
}

// UpdateCategoryRequest represents a request to update a category
type UpdateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CategoryResponse represents a category response
type CategoryResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	ParentID    *uint  `json:"parent_id"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateCategory creates a new category
func (s *Service) CreateCategory(req *CreateCategoryRequest) (*CategoryResponse, error) {
	category := model.AssetCategory{
		Name:        req.Name,
		ParentID:    req.ParentID,
		Description: req.Description,
	}

	if err := s.db.Create(&category).Error; err != nil {
		return nil, err
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetCategory retrieves a category by ID
func (s *Service) GetCategory(id uint) (*CategoryResponse, error) {
	var category model.AssetCategory
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListCategories retrieves all categories
func (s *Service) ListCategories() ([]CategoryResponse, error) {
	var categories []model.AssetCategory
	if err := s.db.Find(&categories).Error; err != nil {
		return nil, err
	}

	result := make([]CategoryResponse, len(categories))
	for i, category := range categories {
		result[i] = CategoryResponse{
			ID:          category.ID,
			Name:        category.Name,
			ParentID:    category.ParentID,
			Description: category.Description,
			CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, nil
}

// UpdateCategory updates a category
func (s *Service) UpdateCategory(id uint, req *UpdateCategoryRequest) (*CategoryResponse, error) {
	var category model.AssetCategory
	if err := s.db.First(&category, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		category.Name = req.Name
	}

	if req.Description != "" {
		category.Description = req.Description
	}

	if err := s.db.Save(&category).Error; err != nil {
		return nil, err
	}

	return &CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
		ParentID:    category.ParentID,
		Description: category.Description,
		CreatedAt:   category.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   category.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteCategory deletes a category
func (s *Service) DeleteCategory(id uint) error {
	return s.db.Delete(&model.AssetCategory{}, id).Error
}
