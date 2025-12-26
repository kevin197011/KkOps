package service

import (
	"fmt"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
)

type EnvironmentService interface {
	Create(env *models.Environment) error
	GetByID(id uint64) (*models.Environment, error)
	GetByName(name string) (*models.Environment, error)
	List(page, pageSize int) ([]models.Environment, int64, error)
	Update(id uint64, env *models.Environment) error
	Delete(id uint64) error
}

type environmentService struct {
	repo repository.EnvironmentRepository
}

func NewEnvironmentService(repo repository.EnvironmentRepository) EnvironmentService {
	return &environmentService{repo: repo}
}

func (s *environmentService) Create(env *models.Environment) error {
	// 检查名称是否已存在
	existing, err := s.repo.GetByName(env.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("environment with name '%s' already exists", env.Name)
	}
	return s.repo.Create(env)
}

func (s *environmentService) GetByID(id uint64) (*models.Environment, error) {
	return s.repo.GetByID(id)
}

func (s *environmentService) GetByName(name string) (*models.Environment, error) {
	return s.repo.GetByName(name)
}

func (s *environmentService) List(page, pageSize int) ([]models.Environment, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

func (s *environmentService) Update(id uint64, env *models.Environment) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查名称是否与其他环境冲突
	if env.Name != "" && env.Name != existing.Name {
		other, err := s.repo.GetByName(env.Name)
		if err == nil && other != nil && other.ID != id {
			return fmt.Errorf("environment with name '%s' already exists", env.Name)
		}
		existing.Name = env.Name
	}

	if env.DisplayName != "" {
		existing.DisplayName = env.DisplayName
	}
	if env.Color != "" {
		existing.Color = env.Color
	}
	if env.SortOrder != 0 {
		existing.SortOrder = env.SortOrder
	}
	if env.Description != "" {
		existing.Description = env.Description
	}

	return s.repo.Update(existing)
}

func (s *environmentService) Delete(id uint64) error {
	return s.repo.Delete(id)
}
