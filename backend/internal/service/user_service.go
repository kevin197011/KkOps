package service

import (
	"errors"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/utils"
)

type UserService interface {
	CreateUser(username, email, password, displayName string) (*models.User, error)
	GetUser(id uint64) (*models.User, error)
	ListUsers(page, pageSize int) ([]models.User, int64, error)
	UpdateUser(id uint64, displayName, email string) (*models.User, error)
	DeleteUser(id uint64) error
	ChangePassword(id uint64, oldPassword, newPassword string) error
	ResetPassword(id uint64, newPassword string) error
	AssignRole(userID, roleID uint64) error
	RemoveRole(userID, roleID uint64) error
	GetUserPermissions(userID uint64) ([]models.Permission, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (s *userService) CreateUser(username, email, password, displayName string) (*models.User, error) {
	// 检查用户名是否已存在
	_, err := s.userRepo.GetByUsername(username)
	if err == nil {
		return nil, errors.New("username already exists")
	}

	// 检查邮箱是否已存在
	_, err = s.userRepo.GetByEmail(email)
	if err == nil {
		return nil, errors.New("email already exists")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
		DisplayName:  displayName,
		Status:       "active",
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) GetUser(id uint64) (*models.User, error) {
	return s.userRepo.GetByID(id)
}

func (s *userService) ListUsers(page, pageSize int) ([]models.User, int64, error) {
	offset := (page - 1) * pageSize
	return s.userRepo.List(offset, pageSize)
}

func (s *userService) UpdateUser(id uint64, displayName, email string) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if email != "" && email != user.Email {
		// 检查邮箱是否已被其他用户使用
		existingUser, err := s.userRepo.GetByEmail(email)
		if err == nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = email
	}

	if displayName != "" {
		user.DisplayName = displayName
	}

	err = s.userRepo.Update(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint64) error {
	return s.userRepo.Delete(id)
}

func (s *userService) ChangePassword(id uint64, oldPassword, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if !utils.CheckPassword(oldPassword, user.PasswordHash) {
		return errors.New("invalid old password")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

// ResetPassword 管理员重置用户密码（无需验证旧密码）
func (s *userService) ResetPassword(id uint64, newPassword string) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.PasswordHash = hashedPassword
	return s.userRepo.Update(user)
}

func (s *userService) AssignRole(userID, roleID uint64) error {
	return s.userRepo.AssignRole(userID, roleID)
}

func (s *userService) RemoveRole(userID, roleID uint64) error {
	return s.userRepo.RemoveRole(userID, roleID)
}

func (s *userService) GetUserPermissions(userID uint64) ([]models.Permission, error) {
	return s.userRepo.GetUserPermissions(userID)
}

