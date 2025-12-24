package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type PermissionRepository interface {
	Create(permission *models.Permission) error
	GetByID(id uint64) (*models.Permission, error)
	GetByCode(code string) (*models.Permission, error)
	List(offset, limit int) ([]models.Permission, int64, error)
	ListByResourceType(resourceType string) ([]models.Permission, error)
	Update(permission *models.Permission) error
	Delete(id uint64) error
}

type permissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository(db *gorm.DB) PermissionRepository {
	return &permissionRepository{db: db}
}

func (r *permissionRepository) Create(permission *models.Permission) error {
	return r.db.Create(permission).Error
}

func (r *permissionRepository) GetByID(id uint64) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.First(&permission, id).Error
	return &permission, err
}

func (r *permissionRepository) GetByCode(code string) (*models.Permission, error) {
	var permission models.Permission
	err := r.db.Where("code = ?", code).First(&permission).Error
	return &permission, err
}

func (r *permissionRepository) List(offset, limit int) ([]models.Permission, int64, error) {
	var permissions []models.Permission
	var total int64

	err := r.db.Model(&models.Permission{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&permissions).Error
	return permissions, total, err
}

func (r *permissionRepository) ListByResourceType(resourceType string) ([]models.Permission, error) {
	var permissions []models.Permission
	err := r.db.Where("resource_type = ?", resourceType).Find(&permissions).Error
	return permissions, err
}

func (r *permissionRepository) Update(permission *models.Permission) error {
	return r.db.Save(permission).Error
}

func (r *permissionRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Permission{}, id).Error
}

