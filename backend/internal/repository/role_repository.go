package repository

import (
	"github.com/kronos/backend/internal/models"
	"gorm.io/gorm"
)

type RoleRepository interface {
	Create(role *models.Role) error
	GetByID(id uint64) (*models.Role, error)
	List(offset, limit int) ([]models.Role, int64, error)
	Update(role *models.Role) error
	Delete(id uint64) error
	AssignPermission(roleID, permissionID uint64) error
	RemovePermission(roleID, permissionID uint64) error
}

type roleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{db: db}
}

func (r *roleRepository) Create(role *models.Role) error {
	return r.db.Create(role).Error
}

func (r *roleRepository) GetByID(id uint64) (*models.Role, error) {
	var role models.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (r *roleRepository) List(offset, limit int) ([]models.Role, int64, error) {
	var roles []models.Role
	var total int64

	err := r.db.Model(&models.Role{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Preload("Permissions").Offset(offset).Limit(limit).Find(&roles).Error
	return roles, total, err
}

func (r *roleRepository) Update(role *models.Role) error {
	return r.db.Save(role).Error
}

func (r *roleRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Role{}, id).Error
}

func (r *roleRepository) AssignPermission(roleID, permissionID uint64) error {
	return r.db.Create(&models.RolePermission{
		RoleID:       roleID,
		PermissionID: permissionID,
	}).Error
}

func (r *roleRepository) RemovePermission(roleID, permissionID uint64) error {
	return r.db.Where("role_id = ? AND permission_id = ?", roleID, permissionID).
		Delete(&models.RolePermission{}).Error
}
