// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package rbac

import (
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
)

// Service handles RBAC business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new RBAC service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// GetUserPermissions returns all permissions for a user (from all roles)
func (s *Service) GetUserPermissions(userID uint) ([]model.Permission, error) {
	var user model.User
	if err := s.db.Preload("Roles.Permissions").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Collect unique permissions from all roles
	permissionMap := make(map[uint]model.Permission)
	for _, role := range user.Roles {
		for _, perm := range role.Permissions {
			permissionMap[perm.ID] = perm
		}
	}

	permissions := make([]model.Permission, 0, len(permissionMap))
	for _, perm := range permissionMap {
		permissions = append(permissions, perm)
	}

	return permissions, nil
}

// HasPermission checks if a user has a specific permission
func (s *Service) HasPermission(userID uint, resource, action string) (bool, error) {
	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		if perm.Resource == resource && perm.Action == action {
			return true, nil
		}
	}

	return false, nil
}

// HasPermissionByCode is deprecated. Use HasPermission(resource, action) instead.
// This method is kept for backward compatibility but will be removed in a future version.
func (s *Service) HasPermissionByCode(userID uint, permissionCode string) (bool, error) {
	// Permission code is no longer used. This method always returns false.
	// Use HasPermission(resource, action) instead.
	return false, nil
}

// GetUserRoles returns all roles for a user
func (s *Service) GetUserRoles(userID uint) ([]model.Role, error) {
	var user model.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	return user.Roles, nil
}
