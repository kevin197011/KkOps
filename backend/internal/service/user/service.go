// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package user

import (
	"errors"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// Service handles user management business logic
type Service struct {
	db *gorm.DB
}

// NewService creates a new user service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// CreateUserRequest represents a request to create a user
type CreateUserRequest struct {
	Username     string `json:"username" binding:"required"`
	Email        string `json:"email" binding:"required,email"`
	Password     string `json:"password" binding:"required,min=6"`
	Phone        string `json:"phone"`
	RealName     string `json:"real_name"`
	DepartmentID *uint  `json:"department_id"`
}

// UpdateUserRequest represents a request to update a user
type UpdateUserRequest struct {
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	RealName     string `json:"real_name"`
	DepartmentID *uint  `json:"department_id"`
	Status       string `json:"status"`
}

// UserResponse represents a user response
type UserResponse struct {
	ID           uint   `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Phone        string `json:"phone"`
	RealName     string `json:"real_name"`
	DepartmentID *uint  `json:"department_id"`
	Status       string `json:"status"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

// CreateUser creates a new user
func (s *Service) CreateUser(req *CreateUserRequest) (*UserResponse, error) {
	// Check if username already exists
	var existingUser model.User
	if err := s.db.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		return nil, errors.New("username already exists")
	}

	// Check if email already exists
	if err := s.db.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		return nil, errors.New("email already exists")
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := model.User{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: passwordHash,
		Phone:        req.Phone,
		RealName:     req.RealName,
		DepartmentID: req.DepartmentID,
		Status:       "active",
	}

	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		RealName:     user.RealName,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// GetUser retrieves a user by ID
func (s *Service) GetUser(id uint) (*UserResponse, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		RealName:     user.RealName,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// ListUsers retrieves a paginated list of users
func (s *Service) ListUsers(page, pageSize int) ([]UserResponse, int64, error) {
	var users []model.User
	var total int64

	offset := (page - 1) * pageSize

	if err := s.db.Model(&model.User{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := s.db.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	result := make([]UserResponse, len(users))
	for i, user := range users {
		result[i] = UserResponse{
			ID:           user.ID,
			Username:     user.Username,
			Email:        user.Email,
			Phone:        user.Phone,
			RealName:     user.RealName,
			DepartmentID: user.DepartmentID,
			Status:       user.Status,
			CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	return result, total, nil
}

// UpdateUser updates a user
func (s *Service) UpdateUser(id uint, req *UpdateUserRequest) (*UserResponse, error) {
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return nil, err
	}

	if req.Email != "" {
		// Check if email already exists for another user
		var existingUser model.User
		if err := s.db.Where("email = ? AND id != ?", req.Email, id).First(&existingUser).Error; err == nil {
			return nil, errors.New("email already exists")
		}
		user.Email = req.Email
	}

	if req.Phone != "" {
		user.Phone = req.Phone
	}

	if req.RealName != "" {
		user.RealName = req.RealName
	}

	if req.DepartmentID != nil {
		user.DepartmentID = req.DepartmentID
	}

	if req.Status != "" {
		user.Status = req.Status
	}

	if err := s.db.Save(&user).Error; err != nil {
		return nil, err
	}

	return &UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Phone:        user.Phone,
		RealName:     user.RealName,
		DepartmentID: user.DepartmentID,
		Status:       user.Status,
		CreatedAt:    user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:    user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}, nil
}

// DeleteUser deletes a user
func (s *Service) DeleteUser(id uint) error {
	return s.db.Delete(&model.User{}, id).Error
}

// ResetPassword resets a user's password
func (s *Service) ResetPassword(id uint, newPassword string) error {
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	return s.db.Model(&model.User{}).Where("id = ?", id).Update("password_hash", passwordHash).Error
}
