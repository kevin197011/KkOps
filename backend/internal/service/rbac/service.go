// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package rbac

import (
	"errors"

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
// Admin users automatically have all permissions
func (s *Service) HasPermission(userID uint, resource, action string) (bool, error) {
	// First check if user is admin (admin has all permissions)
	isAdmin, err := s.IsAdmin(userID)
	if err != nil {
		return false, err
	}
	if isAdmin {
		return true, nil
	}

	// Check specific permission
	permissions, err := s.GetUserPermissions(userID)
	if err != nil {
		return false, err
	}

	for _, perm := range permissions {
		// Match exact resource and action, or wildcard action
		if perm.Resource == resource && (perm.Action == action || perm.Action == "*") {
			return true, nil
		}
	}

	return false, nil
}

// IsAdmin checks if a user is an admin (has a role with is_admin=true)
func (s *Service) IsAdmin(userID uint) (bool, error) {
	var user model.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return false, err
	}

	for _, role := range user.Roles {
		if role.IsAdmin {
			return true, nil
		}
	}

	return false, nil
}

// GetUserPermissionList returns a list of permissions for the user (for frontend)
// Returns resource:action pairs as strings
func (s *Service) GetUserPermissionList(userID uint) ([]string, error) {
	// Check if admin - admins have all permissions
	isAdmin, err := s.IsAdmin(userID)
	if err != nil {
		return nil, err
	}
	if isAdmin {
		// Return all menu permissions
		permissions := make([]string, len(model.AllMenuPermissions))
		for i, perm := range model.AllMenuPermissions {
			permissions[i] = perm.Resource + ":" + perm.Action
		}
		return permissions, nil
	}

	// Get user's actual permissions
	userPerms, err := s.GetUserPermissions(userID)
	if err != nil {
		return nil, err
	}

	permissions := make([]string, len(userPerms))
	for i, perm := range userPerms {
		permissions[i] = perm.Resource + ":" + perm.Action
	}

	return permissions, nil
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
