package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type CloudPlatformRepository interface {
	Create(cp *models.CloudPlatform) error
	GetByID(id uint64) (*models.CloudPlatform, error)
	GetByName(name string) (*models.CloudPlatform, error)
	List(offset, limit int) ([]models.CloudPlatform, int64, error)
	Update(cp *models.CloudPlatform) error
	Delete(id uint64) error
}

type cloudPlatformRepository struct {
	db *gorm.DB
}

func NewCloudPlatformRepository(db *gorm.DB) CloudPlatformRepository {
	return &cloudPlatformRepository{db: db}
}

func (r *cloudPlatformRepository) Create(cp *models.CloudPlatform) error {
	return r.db.Create(cp).Error
}

func (r *cloudPlatformRepository) GetByID(id uint64) (*models.CloudPlatform, error) {
	var cp models.CloudPlatform
	err := r.db.First(&cp, id).Error
	return &cp, err
}

func (r *cloudPlatformRepository) GetByName(name string) (*models.CloudPlatform, error) {
	var cp models.CloudPlatform
	err := r.db.Where("name = ?", name).First(&cp).Error
	return &cp, err
}

func (r *cloudPlatformRepository) List(offset, limit int) ([]models.CloudPlatform, int64, error) {
	var cps []models.CloudPlatform
	var total int64

	r.db.Model(&models.CloudPlatform{}).Count(&total)
	err := r.db.Order("sort_order ASC, id ASC").Offset(offset).Limit(limit).Find(&cps).Error
	return cps, total, err
}

func (r *cloudPlatformRepository) Update(cp *models.CloudPlatform) error {
	return r.db.Save(cp).Error
}

func (r *cloudPlatformRepository) Delete(id uint64) error {
	return r.db.Delete(&models.CloudPlatform{}, id).Error
}
