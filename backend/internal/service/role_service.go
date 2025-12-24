package service

import (
	"errors"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
)

type RoleService interface {
	CreateRole(name, description string) (*models.Role, error)
	GetRole(id uint64) (*models.Role, error)
	ListRoles(page, pageSize int) ([]models.Role, int64, error)
	UpdateRole(id uint64, name, description string) (*models.Role, error)
	DeleteRole(id uint64) error
	AssignPermission(roleID, permissionID uint64) error
	RemovePermission(roleID, permissionID uint64) error
}

type roleService struct {
	roleRepo repository.RoleRepository
}

func NewRoleService(roleRepo repository.RoleRepository) RoleService {
	return &roleService{roleRepo: roleRepo}
}

func (s *roleService) CreateRole(name, description string) (*models.Role, error) {
	role := &models.Role{
		Name:        name,
		Description: description,
	}

	err := s.roleRepo.Create(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) GetRole(id uint64) (*models.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *roleService) ListRoles(page, pageSize int) ([]models.Role, int64, error) {
	offset := (page - 1) * pageSize
	return s.roleRepo.List(offset, pageSize)
}

func (s *roleService) UpdateRole(id uint64, name, description string) (*models.Role, error) {
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if name != "" {
		role.Name = name
	}
	if description != "" {
		role.Description = description
	}

	err = s.roleRepo.Update(role)
	if err != nil {
		return nil, err
	}

	return role, nil
}

func (s *roleService) DeleteRole(id uint64) error {
	// 检查是否有用户使用该角色
	role, err := s.roleRepo.GetByID(id)
	if err != nil {
		return err
	}

	if len(role.Users) > 0 {
		return errors.New("cannot delete role: role is assigned to users")
	}

	return s.roleRepo.Delete(id)
}

func (s *roleService) AssignPermission(roleID, permissionID uint64) error {
	return s.roleRepo.AssignPermission(roleID, permissionID)
}

func (s *roleService) RemovePermission(roleID, permissionID uint64) error {
	return s.roleRepo.RemovePermission(roleID, permissionID)
}

