package service

import (
	"fmt"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
)

type CloudPlatformService interface {
	Create(cp *models.CloudPlatform) error
	GetByID(id uint64) (*models.CloudPlatform, error)
	GetByName(name string) (*models.CloudPlatform, error)
	List(page, pageSize int) ([]models.CloudPlatform, int64, error)
	Update(id uint64, cp *models.CloudPlatform) error
	Delete(id uint64) error
}

type cloudPlatformService struct {
	repo repository.CloudPlatformRepository
}

func NewCloudPlatformService(repo repository.CloudPlatformRepository) CloudPlatformService {
	return &cloudPlatformService{repo: repo}
}

func (s *cloudPlatformService) Create(cp *models.CloudPlatform) error {
	// 检查名称是否已存在
	existing, err := s.repo.GetByName(cp.Name)
	if err == nil && existing != nil {
		return fmt.Errorf("cloud platform with name '%s' already exists", cp.Name)
	}
	return s.repo.Create(cp)
}

func (s *cloudPlatformService) GetByID(id uint64) (*models.CloudPlatform, error) {
	return s.repo.GetByID(id)
}

func (s *cloudPlatformService) GetByName(name string) (*models.CloudPlatform, error) {
	return s.repo.GetByName(name)
}

func (s *cloudPlatformService) List(page, pageSize int) ([]models.CloudPlatform, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.List(offset, pageSize)
}

func (s *cloudPlatformService) Update(id uint64, cp *models.CloudPlatform) error {
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	// 检查名称是否与其他云平台冲突
	if cp.Name != "" && cp.Name != existing.Name {
		other, err := s.repo.GetByName(cp.Name)
		if err == nil && other != nil && other.ID != id {
			return fmt.Errorf("cloud platform with name '%s' already exists", cp.Name)
		}
		existing.Name = cp.Name
	}

	if cp.DisplayName != "" {
		existing.DisplayName = cp.DisplayName
	}
	if cp.Icon != "" {
		existing.Icon = cp.Icon
	}
	if cp.Color != "" {
		existing.Color = cp.Color
	}
	if cp.SortOrder != 0 {
		existing.SortOrder = cp.SortOrder
	}
	if cp.Description != "" {
		existing.Description = cp.Description
	}

	return s.repo.Update(existing)
}

func (s *cloudPlatformService) Delete(id uint64) error {
	return s.repo.Delete(id)
}
