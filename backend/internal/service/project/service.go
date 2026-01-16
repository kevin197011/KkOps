// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package project

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles project management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new project service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateProjectRequest represents a request to create a project
type CreateProjectRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateProjectRequest represents a request to update a project
type UpdateProjectRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// ProjectResponse represents a project response
type ProjectResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedBy   uint   `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateProject creates a new project
func (s *Service) CreateProject(userID uint, req *CreateProjectRequest) (*ProjectResponse, error) {
	// Check if name already exists
	var existingProject model.Project
	if err := s.db.Where("name = ?", req.Name).First(&existingProject).Error; err == nil {
		return nil, errors.New("project name already exists")
	}

	project := model.Project{
		Name:        req.Name,
		Description: req.Description,
		CreatedBy:   userID,
	}

	if err := s.db.Create(&project).Error; err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   project.CreatedBy,
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetProject retrieves a project by ID
func (s *Service) GetProject(id uint) (*ProjectResponse, error) {
	var project model.Project
	if err := s.db.First(&project, id).Error; err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   project.CreatedBy,
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListProjects retrieves all projects
func (s *Service) ListProjects() ([]ProjectResponse, error) {
	var projects []model.Project
	if err := s.db.Find(&projects).Error; err != nil {
		return nil, err
	}

	result := make([]ProjectResponse, len(projects))
	for i, project := range projects {
		result[i] = ProjectResponse{
			ID:          project.ID,
			Name:        project.Name,
			Description: project.Description,
			CreatedBy:   project.CreatedBy,
			CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, nil
}

// UpdateProject updates a project
func (s *Service) UpdateProject(id uint, req *UpdateProjectRequest) (*ProjectResponse, error) {
	var project model.Project
	if err := s.db.First(&project, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		// Check if new name already exists (excluding current project)
		var existingProject model.Project
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingProject).Error; err == nil {
			return nil, errors.New("project name already exists")
		}
		project.Name = req.Name
	}

	if req.Description != "" {
		project.Description = req.Description
	}

	if err := s.db.Save(&project).Error; err != nil {
		return nil, err
	}

	return &ProjectResponse{
		ID:          project.ID,
		Name:        project.Name,
		Description: project.Description,
		CreatedBy:   project.CreatedBy,
		CreatedAt:   project.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   project.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteProject deletes a project
func (s *Service) DeleteProject(id uint) error {
	return s.db.Delete(&model.Project{}, id).Error
}
