// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// Service handles authentication business logic
type Service struct {
	db     *gorm.DB
	config *config.Config
}

// NewService creates a new authentication service
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{
		db:     db,
		config: cfg,
	}
}

// LoginRequest represents a login request
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// LoginResponse represents a login response
type LoginResponse struct {
	Token        string   `json:"token"`
	RefreshToken string   `json:"refresh_token"`
	ExpiresIn    int      `json:"expires_in"`
	User         UserInfo `json:"user"`
}

// UserInfo represents user information in login response
type UserInfo struct {
	ID       uint     `json:"id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	RealName string   `json:"real_name"`
	Roles    []string `json:"roles"`
}

// Login authenticates a user and returns a JWT token
func (s *Service) Login(req *LoginRequest) (*LoginResponse, error) {
	var user model.User
	if err := s.db.Where("username = ?", req.Username).Preload("Roles").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	// Check if user is active
	if user.Status != "active" {
		return nil, errors.New("user account is disabled")
	}

	// Verify password
	if !utils.CheckPasswordHash(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Get user roles (use role names instead of codes)
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	// Generate JWT token
	token, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		roleNames,
		s.config.JWT.Secret,
		s.config.JWT.ExpiresIn,
	)
	if err != nil {
		return nil, err
	}

	// Generate refresh token
	refreshToken, err := utils.GenerateJWT(
		user.ID,
		user.Username,
		roleNames,
		s.config.JWT.Secret,
		s.config.JWT.RefreshExpiresIn,
	)
	if err != nil {
		return nil, err
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	s.db.Save(&user)

	return &LoginResponse{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresIn:    s.config.JWT.ExpiresIn,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RealName: user.RealName,
			Roles:    roleNames,
		},
	}, nil
}

// GetCurrentUser returns the current user information
func (s *Service) GetCurrentUser(userID uint) (*UserInfo, error) {
	var user model.User
	if err := s.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}

	// Use role names instead of codes
	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	return &UserInfo{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		RealName: user.RealName,
		Roles:    roleNames,
	}, nil
}

// ChangePasswordRequest represents a change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

// ChangePassword changes the password for the current user
func (s *Service) ChangePassword(userID uint, req *ChangePasswordRequest) error {
	var user model.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return errors.New("user not found")
	}

	// Verify old password
	if !utils.CheckPasswordHash(req.OldPassword, user.PasswordHash) {
		return errors.New("原密码不正确")
	}

	// Hash new password
	newHash, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return errors.New("密码加密失败")
	}

	// Update password
	user.PasswordHash = newHash
	if err := s.db.Save(&user).Error; err != nil {
		return errors.New("密码更新失败")
	}

	return nil
}
