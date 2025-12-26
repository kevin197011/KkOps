package service

import (
	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
)

type CommandTemplateService interface {
	CreateTemplate(template *models.CommandTemplate) error
	UpdateTemplate(id uint, template *models.CommandTemplate) error
	DeleteTemplate(id uint) error
	ListTemplates(filters map[string]interface{}) ([]models.CommandTemplate, error)
	GetTemplate(id uint) (*models.CommandTemplate, error)
	IncrementUsageCount(id uint) error
}

type commandTemplateService struct {
	templateRepo *repository.CommandTemplateRepository
}

func NewCommandTemplateService(templateRepo *repository.CommandTemplateRepository) CommandTemplateService {
	return &commandTemplateService{
		templateRepo: templateRepo,
	}
}

func (s *commandTemplateService) CreateTemplate(template *models.CommandTemplate) error {
	return s.templateRepo.Create(template)
}

func (s *commandTemplateService) UpdateTemplate(id uint, template *models.CommandTemplate) error {
	existing, err := s.templateRepo.Get(id)
	if err != nil {
		return err
	}

	// 更新字段
	if template.Name != "" {
		existing.Name = template.Name
	}
	if template.Description != "" {
		existing.Description = template.Description
	}
	if template.Category != "" {
		existing.Category = template.Category
	}
	if template.CommandFunction != "" {
		existing.CommandFunction = template.CommandFunction
	}
	if len(template.CommandArgs) > 0 {
		existing.CommandArgs = template.CommandArgs
	}
	if template.Icon != "" {
		existing.Icon = template.Icon
	}
	existing.IsPublic = template.IsPublic

	return s.templateRepo.Update(existing)
}

func (s *commandTemplateService) DeleteTemplate(id uint) error {
	return s.templateRepo.Delete(id)
}

func (s *commandTemplateService) ListTemplates(filters map[string]interface{}) ([]models.CommandTemplate, error) {
	return s.templateRepo.List(filters)
}

func (s *commandTemplateService) GetTemplate(id uint) (*models.CommandTemplate, error) {
	return s.templateRepo.Get(id)
}

func (s *commandTemplateService) IncrementUsageCount(id uint) error {
	return s.templateRepo.IncrementUsageCount(id)
}

