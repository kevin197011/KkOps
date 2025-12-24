package repository

import (
	"github.com/kkops/backend/internal/models"
	"gorm.io/gorm"
)

type SSHKeyRepository interface {
	Create(key *models.SSHKey) error
	GetByID(id uint64) (*models.SSHKey, error)
	GetByUserID(userID uint64) ([]models.SSHKey, error)
	List(userID uint64, offset, limit int) ([]models.SSHKey, int64, error)
	Update(key *models.SSHKey) error
	Delete(id uint64) error
	CheckKeyInUse(keyID uint64) (bool, error) // 检查密钥是否被主机使用
}

type sshKeyRepository struct {
	db *gorm.DB
}

func NewSSHKeyRepository(db *gorm.DB) SSHKeyRepository {
	return &sshKeyRepository{db: db}
}

func (r *sshKeyRepository) Create(key *models.SSHKey) error {
	return r.db.Create(key).Error
}

func (r *sshKeyRepository) GetByID(id uint64) (*models.SSHKey, error) {
	var key models.SSHKey
	err := r.db.First(&key, id).Error
	return &key, err
}

func (r *sshKeyRepository) GetByUserID(userID uint64) ([]models.SSHKey, error) {
	var keys []models.SSHKey
	err := r.db.Where("user_id = ?", userID).Find(&keys).Error
	return keys, err
}

func (r *sshKeyRepository) List(userID uint64, offset, limit int) ([]models.SSHKey, int64, error) {
	var keys []models.SSHKey
	var total int64

	query := r.db.Model(&models.SSHKey{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&keys).Error
	return keys, total, err
}

func (r *sshKeyRepository) Update(key *models.SSHKey) error {
	return r.db.Model(key).Updates(key).Error
}

func (r *sshKeyRepository) Delete(id uint64) error {
	return r.db.Delete(&models.SSHKey{}, id).Error
}

func (r *sshKeyRepository) CheckKeyInUse(keyID uint64) (bool, error) {
	var count int64
	err := r.db.Model(&models.Host{}).Where("ssh_key_id = ?", keyID).Count(&count).Error
	return count > 0, err
}

