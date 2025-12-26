package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type EnvironmentRepository interface {
	Create(env *models.Environment) error
	GetByID(id uint64) (*models.Environment, error)
	GetByName(name string) (*models.Environment, error)
	List(offset, limit int) ([]models.Environment, int64, error)
	Update(env *models.Environment) error
	Delete(id uint64) error
}

type environmentRepository struct {
	db *gorm.DB
}

func NewEnvironmentRepository(db *gorm.DB) EnvironmentRepository {
	return &environmentRepository{db: db}
}

func (r *environmentRepository) Create(env *models.Environment) error {
	return r.db.Create(env).Error
}

func (r *environmentRepository) GetByID(id uint64) (*models.Environment, error) {
	var env models.Environment
	err := r.db.First(&env, id).Error
	return &env, err
}

func (r *environmentRepository) GetByName(name string) (*models.Environment, error) {
	var env models.Environment
	err := r.db.Where("name = ?", name).First(&env).Error
	return &env, err
}

func (r *environmentRepository) List(offset, limit int) ([]models.Environment, int64, error) {
	var envs []models.Environment
	var total int64

	r.db.Model(&models.Environment{}).Count(&total)
	err := r.db.Order("sort_order ASC, id ASC").Offset(offset).Limit(limit).Find(&envs).Error
	return envs, total, err
}

func (r *environmentRepository) Update(env *models.Environment) error {
	return r.db.Save(env).Error
}

func (r *environmentRepository) Delete(id uint64) error {
	return r.db.Delete(&models.Environment{}, id).Error
}
