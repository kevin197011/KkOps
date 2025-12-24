package service

import (
	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
)

type PermissionService interface {
	CreatePermission(code, name, resourceType, action, description string) (*models.Permission, error)
	GetPermission(id uint64) (*models.Permission, error)
	GetPermissionByCode(code string) (*models.Permission, error)
	ListPermissions(page, pageSize int) ([]models.Permission, int64, error)
	ListPermissionsByResourceType(resourceType string) ([]models.Permission, error)
	UpdatePermission(id uint64, name, description string) (*models.Permission, error)
	DeletePermission(id uint64) error
}

type permissionService struct {
	permissionRepo repository.PermissionRepository
}

func NewPermissionService(permissionRepo repository.PermissionRepository) PermissionService {
	return &permissionService{permissionRepo: permissionRepo}
}

func (s *permissionService) CreatePermission(code, name, resourceType, action, description string) (*models.Permission, error) {
	permission := &models.Permission{
		Code:         code,
		Name:         name,
		ResourceType: resourceType,
		Action:       action,
		Description:  description,
	}

	err := s.permissionRepo.Create(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) GetPermission(id uint64) (*models.Permission, error) {
	return s.permissionRepo.GetByID(id)
}

func (s *permissionService) GetPermissionByCode(code string) (*models.Permission, error) {
	return s.permissionRepo.GetByCode(code)
}

func (s *permissionService) ListPermissions(page, pageSize int) ([]models.Permission, int64, error) {
	offset := (page - 1) * pageSize
	return s.permissionRepo.List(offset, pageSize)
}

func (s *permissionService) ListPermissionsByResourceType(resourceType string) ([]models.Permission, error) {
	return s.permissionRepo.ListByResourceType(resourceType)
}

func (s *permissionService) UpdatePermission(id uint64, name, description string) (*models.Permission, error) {
	permission, err := s.permissionRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		permission.Name = name
	}
	if description != "" {
		permission.Description = description
	}

	err = s.permissionRepo.Update(permission)
	if err != nil {
		return nil, err
	}

	return permission, nil
}

func (s *permissionService) DeletePermission(id uint64) error {
	return s.permissionRepo.Delete(id)
}

