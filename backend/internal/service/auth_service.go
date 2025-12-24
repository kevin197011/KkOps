package service

import (
	"errors"
	"time"

	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
	"github.com/kronos/backend/internal/utils"
)

type AuthService interface {
	Login(username, password, ip string) (string, *models.User, error)
	Register(username, email, password, displayName string) (*models.User, error)
}

type authService struct {
	userRepo repository.UserRepository
}

func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}

func (s *authService) Login(username, password, ip string) (string, *models.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	if user.Status != "active" {
		return "", nil, errors.New("user account is disabled")
	}

	if !utils.CheckPassword(password, user.PasswordHash) {
		return "", nil, errors.New("invalid credentials")
	}

	// 更新最后登录信息
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = ip
	s.userRepo.Update(user)

	// 生成JWT token
	roles := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roles[i] = role.Name
	}

	token, err := utils.GenerateToken(user.ID, user.Username, roles)
	if err != nil {
		return "", nil, err
	}

	return token, user, nil
}

func (s *authService) Register(username, email, password, displayName string) (*models.User, error) {
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

	// 加密密码
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

