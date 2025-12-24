package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type UserRepository interface {
	Create(user *models.User) error
	GetByID(id uint64) (*models.User, error)
	GetByUsername(username string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	List(offset, limit int) ([]models.User, int64, error)
	Update(user *models.User) error
	Delete(id uint64) error
	AssignRole(userID, roleID uint64) error
	RemoveRole(userID, roleID uint64) error
	GetUserRoles(userID uint64) ([]models.Role, error)
	GetUserPermissions(userID uint64) ([]models.Permission, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *userRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").First(&user, id).Error
	return &user, err
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").Where("username = ?", username).First(&user).Error
	return &user, err
}

func (r *userRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Preload("Roles.Permissions").Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *userRepository) List(offset, limit int) ([]models.User, int64, error) {
	var users []models.User
	var total int64

	err := r.db.Model(&models.User{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Roles").Offset(offset).Limit(limit).Find(&users).Error
	return users, total, err
}

func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

func (r *userRepository) Delete(id uint64) error {
	return r.db.Delete(&models.User{}, id).Error
}

func (r *userRepository) AssignRole(userID, roleID uint64) error {
	return r.db.Create(&models.UserRole{
		UserID: userID,
		RoleID: roleID,
	}).Error
}

func (r *userRepository) RemoveRole(userID, roleID uint64) error {
	return r.db.Where("user_id = ? AND role_id = ?", userID, roleID).Delete(&models.UserRole{}).Error
}

func (r *userRepository) GetUserRoles(userID uint64) ([]models.Role, error) {
	var user models.User
	err := r.db.Preload("Roles").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return user.Roles, nil
}

func (r *userRepository) GetUserPermissions(userID uint64) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Table("permissions").
		Joins("INNER JOIN role_permissions ON permissions.id = role_permissions.permission_id").
		Joins("INNER JOIN user_roles ON role_permissions.role_id = user_roles.role_id").
		Where("user_roles.user_id = ?", userID).
		Distinct().
		Find(&permissions).Error
	return permissions, err
}

