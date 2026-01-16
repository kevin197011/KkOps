// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package role

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles role and permission management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new role service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateRoleRequest represents a request to create a role
type CreateRoleRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

// UpdateRoleRequest represents a request to update a role
type UpdateRoleRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// RoleResponse represents a role response
type RoleResponse struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// CreateRole creates a new role
func (s *Service) CreateRole(req *CreateRoleRequest) (*RoleResponse, error) {
	// Check if name already exists
	var existingRole model.Role
	if err := s.db.Where("name = ?", req.Name).First(&existingRole).Error; err == nil {
		return nil, errors.New("role name already exists")
	}

	role := model.Role{
		Name:        req.Name,
		Description: req.Description,
	}

	if err := s.db.Create(&role).Error; err != nil {
		return nil, err
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetRole retrieves a role by ID
func (s *Service) GetRole(id uint) (*RoleResponse, error) {
	var role model.Role
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListRoles retrieves all roles
func (s *Service) ListRoles() ([]RoleResponse, error) {
	var roles []model.Role
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}

	result := make([]RoleResponse, len(roles))
	for i, role := range roles {
		result[i] = RoleResponse{
			ID:          role.ID,
			Name:        role.Name,
			Description: role.Description,
			CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, nil
}

// UpdateRole updates a role
func (s *Service) UpdateRole(id uint, req *UpdateRoleRequest) (*RoleResponse, error) {
	var role model.Role
	if err := s.db.First(&role, id).Error; err != nil {
		return nil, err
	}

	if req.Name != "" {
		// Check if new name already exists (excluding current role)
		var existingRole model.Role
		if err := s.db.Where("name = ? AND id != ?", req.Name, id).First(&existingRole).Error; err == nil {
			return nil, errors.New("role name already exists")
		}
		role.Name = req.Name
	}

	if req.Description != "" {
		role.Description = req.Description
	}

	if err := s.db.Save(&role).Error; err != nil {
		return nil, err
	}

	return &RoleResponse{
		ID:          role.ID,
		Name:        role.Name,
		Description: role.Description,
		CreatedAt:   role.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   role.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteRole deletes a role
func (s *Service) DeleteRole(id uint) error {
	return s.db.Delete(&model.Role{}, id).Error
}

// AssignRoleToUser assigns a role to a user
func (s *Service) AssignRoleToUser(userID, roleID uint) error {
	// Check if role exists
	var role model.Role
	if err := s.db.First(&role, roleID).Error; err != nil {
		return errors.New("role not found")
	}

	// Check if assignment already exists
	var userRole model.UserRole
	if err := s.db.Where("user_id = ? AND role_id = ?", userID, roleID).First(&userRole).Error; err == nil {
		return nil // Already assigned
	}

	userRole = model.UserRole{
		UserID: userID,
		RoleID: roleID,
	}

	return s.db.Create(&userRole).Error
}

// RemoveRoleFromUser removes a role from a user
func (s *Service) RemoveRoleFromUser(userID, roleID uint) error {
	return s.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&model.UserRole{}).Error
}

// AssignPermissionToRole assigns a permission to a role
func (s *Service) AssignPermissionToRole(roleID, permissionID uint) error {
	// Check if permission exists
	var permission model.Permission
	if err := s.db.First(&permission, permissionID).Error; err != nil {
		return errors.New("permission not found")
	}

	// Check if assignment already exists
	var rolePerm model.RolePermission
	if err := s.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).First(&rolePerm).Error; err == nil {
		return nil // Already assigned
	}

	rolePerm = model.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}

	return s.db.Create(&rolePerm).Error
}

// RemovePermissionFromRole removes a permission from a role
func (s *Service) RemovePermissionFromRole(roleID, permissionID uint) error {
	return s.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).Delete(&model.RolePermission{}).Error
}

// ListPermissions retrieves all permissions
func (s *Service) ListPermissions() ([]model.Permission, error) {
	var permissions []model.Permission
	if err := s.db.Find(&permissions).Error; err != nil {
		return nil, err
	}

	return permissions, nil
}

// GetRolePermissions retrieves permissions for a role
func (s *Service) GetRolePermissions(roleID uint) ([]model.Permission, error) {
	var role model.Role
	if err := s.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}

	return role.Permissions, nil
}
