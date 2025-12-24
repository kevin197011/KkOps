package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type HostTagRepository interface {
	Create(tag *models.HostTag) error
	GetByID(id uint64) (*models.HostTag, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.HostTag, int64, error)
	Update(tag *models.HostTag) error
	Delete(id uint64) error
}

type hostTagRepository struct {
	db *gorm.DB
}

func NewHostTagRepository(db *gorm.DB) HostTagRepository {
	return &hostTagRepository{db: db}
}

func (r *hostTagRepository) Create(tag *models.HostTag) error {
	return r.db.Create(tag).Error
}

func (r *hostTagRepository) GetByID(id uint64) (*models.HostTag, error) {
	var tag models.HostTag
	err := r.db.Preload("Hosts").First(&tag, id).Error
	return &tag, err
}

func (r *hostTagRepository) List(offset, limit int, filters map[string]interface{}) ([]models.HostTag, int64, error) {
	var tags []models.HostTag
	var total int64

	query := r.db.Model(&models.HostTag{})

	// 应用过滤器
	if name, ok := filters["name"].(string); ok && name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Hosts").Offset(offset).Limit(limit).Find(&tags).Error
	return tags, total, err
}

func (r *hostTagRepository) Update(tag *models.HostTag) error {
	return r.db.Save(tag).Error
}

func (r *hostTagRepository) Delete(id uint64) error {
	return r.db.Delete(&models.HostTag{}, id).Error
}

