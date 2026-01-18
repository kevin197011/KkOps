// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package task

import (
	"errors"
	"strconv"
	"strings"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/service/authorization"
)

// Service handles task and template management business logic
type Service struct {
	db       *gorm.DB
	authzSvc *authorization.Service
}

// NewService creates a new task service
func NewService(db *gorm.DB, authzSvc *authorization.Service) *Service {
	return &Service{
		db:       db,
		authzSvc: authzSvc,
	}
}

// CreateTemplateRequest represents a request to create a task template
type CreateTemplateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Type        string `json:"type"` // shell, python, etc.
}

// UpdateTemplateRequest represents a request to update a task template
type UpdateTemplateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Type        string `json:"type"`
}

// TemplateResponse represents a task template response
type TemplateResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	CreatedBy   uint   `json:"created_by"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	TemplateID  *uint  `json:"template_id"`
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Content     string `json:"content" binding:"required"`
	Type        string `json:"type"`      // shell, python, etc.
	Timeout     int    `json:"timeout"`   // Execution timeout in seconds (default 600)
	AssetIDs    []uint `json:"asset_ids"` // Target assets for execution (stored for later use)
}

// UpdateTaskRequest represents a request to update a task
type UpdateTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Content     string `json:"content"`
	Type        string `json:"type"`
	Timeout     int    `json:"timeout"` // Execution timeout in seconds
	Status      string `json:"status"`
}

// TaskResponse represents a task response
type TaskResponse struct {
	ID          uint    `json:"id"`
	TemplateID  *uint   `json:"template_id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Content     string  `json:"content"`
	Type        string  `json:"type"`
	Timeout     int     `json:"timeout"` // Execution timeout in seconds
	Status      string  `json:"status"`
	AssetIDs    []uint  `json:"asset_ids"`
	CreatedBy   uint    `json:"created_by"`
	StartedAt   *string `json:"started_at"`
	FinishedAt  *string `json:"finished_at"`
	CreatedAt   string  `json:"created_at"`
	UpdatedAt   string  `json:"updated_at"`
}

// CreateTemplate creates a new task template
func (s *Service) CreateTemplate(userID uint, req *CreateTemplateRequest) (*TemplateResponse, error) {
	template := model.TaskTemplate{
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
		Type:        req.Type,
		CreatedBy:   userID,
	}

	if template.Type == "" {
		template.Type = "shell"
	}

	if err := s.db.Create(&template).Error; err != nil {
		return nil, err
	}

	return &TemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Content:     template.Content,
		Type:        template.Type,
		CreatedBy:   template.CreatedBy,
		CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetTemplate retrieves a task template by ID
func (s *Service) GetTemplate(id uint) (*TemplateResponse, error) {
	var template model.TaskTemplate
	if err := s.db.First(&template, id).Error; err != nil {
		return nil, err
	}

	return &TemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Content:     template.Content,
		Type:        template.Type,
		CreatedBy:   template.CreatedBy,
		CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListTemplates retrieves all task templates
func (s *Service) ListTemplates() ([]TemplateResponse, error) {
	var templates []model.TaskTemplate
	if err := s.db.Find(&templates).Error; err != nil {
		return nil, err
	}

	result := make([]TemplateResponse, len(templates))
	for i, template := range templates {
		result[i] = TemplateResponse{
			ID:          template.ID,
			Name:        template.Name,
			Description: template.Description,
			Content:     template.Content,
			Type:        template.Type,
			CreatedBy:   template.CreatedBy,
			CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, nil
}

// UpdateTemplate updates a task template
func (s *Service) UpdateTemplate(id uint, req *UpdateTemplateRequest) (*TemplateResponse, error) {
	var template model.TaskTemplate
	if err := s.db.First(&template, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		template.Name = req.Name
	}
	if req.Description != "" {
		template.Description = req.Description
	}
	if req.Content != "" {
		template.Content = req.Content
	}
	if req.Type != "" {
		template.Type = req.Type
	}

	if err := s.db.Save(&template).Error; err != nil {
		return nil, err
	}

	return &TemplateResponse{
		ID:          template.ID,
		Name:        template.Name,
		Description: template.Description,
		Content:     template.Content,
		Type:        template.Type,
		CreatedBy:   template.CreatedBy,
		CreatedAt:   template.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   template.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteTemplate deletes a task template
func (s *Service) DeleteTemplate(id uint) error {
	return s.db.Delete(&model.TaskTemplate{}, id).Error
}

// CreateTask creates a new task (without execution records - those are created on Execute)
// 增加资产权限检查：用户只能在已授权的资产上创建任务
func (s *Service) CreateTask(userID uint, req *CreateTaskRequest) (*TaskResponse, error) {
	// If template ID is provided, load template content
	if req.TemplateID != nil {
		var template model.TaskTemplate
		if err := s.db.First(&template, *req.TemplateID).Error; err != nil {
			return nil, errors.New("template not found")
		}
		if req.Content == "" {
			req.Content = template.Content
		}
		if req.Type == "" {
			req.Type = template.Type
		}
	}

	// 检查用户对目标资产的访问权限
	if len(req.AssetIDs) > 0 && s.authzSvc != nil {
		authorizedAssets, err := s.authzSvc.HasMultipleAssetAccess(userID, req.AssetIDs)
		if err != nil {
			return nil, errors.New("failed to check asset permissions")
		}
		if len(authorizedAssets) != len(req.AssetIDs) {
			return nil, errors.New("no permission to execute on selected assets")
		}
	}

	// Convert asset IDs to comma-separated string
	assetIDsStr := ""
	if len(req.AssetIDs) > 0 {
		ids := make([]string, len(req.AssetIDs))
		for i, id := range req.AssetIDs {
			ids[i] = strconv.FormatUint(uint64(id), 10)
		}
		assetIDsStr = strings.Join(ids, ",")
	}

	// Set default timeout if not provided (600 seconds = 10 minutes)
	timeout := req.Timeout
	if timeout <= 0 {
		timeout = 600
	}

	task := model.Task{
		TemplateID:  req.TemplateID,
		Name:        req.Name,
		Description: req.Description,
		Content:     req.Content,
		Type:        req.Type,
		Timeout:     timeout,
		Status:      "pending",
		AssetIDs:    assetIDsStr,
		CreatedBy:   userID,
	}

	if task.Type == "" {
		task.Type = "shell"
	}

	if err := s.db.Create(&task).Error; err != nil {
		return nil, err
	}

	// Execution records will be created when task is executed
	return s.taskToResponse(task), nil
}

// GetTask retrieves a task by ID
func (s *Service) GetTask(id uint) (*TaskResponse, error) {
	var task model.Task
	if err := s.db.First(&task, id).Error; err != nil {
		return nil, err
	}

	return s.taskToResponse(task), nil
}

// ListTasks retrieves all tasks
func (s *Service) ListTasks(page, pageSize int) ([]TaskResponse, int64, error) {
	var tasks []model.Task
	var total int64

	offset := (page - 1) * pageSize

	if err := s.db.Model(&model.Task{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&tasks).Error; err != nil {
		return nil, 0, err
	}

	result := make([]TaskResponse, len(tasks))
	for i, task := range tasks {
		result[i] = *s.taskToResponse(task)
	}

	return result, total, nil
}

// UpdateTask updates a task
func (s *Service) UpdateTask(id uint, req *UpdateTaskRequest) (*TaskResponse, error) {
	var task model.Task
	if err := s.db.First(&task, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Content != "" {
		task.Content = req.Content
	}
	if req.Type != "" {
		task.Type = req.Type
	}
	if req.Timeout > 0 {
		task.Timeout = req.Timeout
	}
	if req.Status != "" {
		task.Status = req.Status
	}

	if err := s.db.Save(&task).Error; err != nil {
		return nil, err
	}

	return s.taskToResponse(task), nil
}

// DeleteTask deletes a task
func (s *Service) DeleteTask(id uint) error {
	return s.db.Delete(&model.Task{}, id).Error
}

// parseAssetIDs converts comma-separated string to uint slice
func parseAssetIDs(assetIDsStr string) []uint {
	if assetIDsStr == "" {
		return []uint{}
	}
	parts := strings.Split(assetIDsStr, ",")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseUint(p, 10, 32)
		if err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids
}

// taskToResponse converts a task model to response
func (s *Service) taskToResponse(task model.Task) *TaskResponse {
	resp := &TaskResponse{
		ID:          task.ID,
		TemplateID:  task.TemplateID,
		Name:        task.Name,
		Description: task.Description,
		Content:     task.Content,
		Type:        task.Type,
		Timeout:     task.Timeout,
		Status:      task.Status,
		AssetIDs:    parseAssetIDs(task.AssetIDs),
		CreatedBy:   task.CreatedBy,
		CreatedAt:   task.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   task.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if task.StartedAt != nil {
		startedAt := task.StartedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.StartedAt = &startedAt
	}
	if task.FinishedAt != nil {
		finishedAt := task.FinishedAt.Format("2006-01-02T15:04:05Z07:00")
		resp.FinishedAt = &finishedAt
	}

	return resp
}
