package repository

import (
	"github.com/kronos/backend/internal/models"
	"gorm.io/gorm"
)

type SSHConnectionRepository interface {
	Create(conn *models.SSHConnection) error
	GetByID(id uint64) (*models.SSHConnection, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.SSHConnection, int64, error)
	Update(conn *models.SSHConnection) error
	Delete(id uint64) error
}

type SSHKeyRepository interface {
	Create(key *models.SSHKey) error
	GetByID(id uint64) (*models.SSHKey, error)
	GetByFingerprint(fingerprint string) (*models.SSHKey, error)
	List(offset, limit int) ([]models.SSHKey, int64, error)
	Update(key *models.SSHKey) error
	Delete(id uint64) error
}

type SSHSessionRepository interface {
	Create(session *models.SSHSession) error
	GetByID(id uint64) (*models.SSHSession, error)
	GetByToken(token string) (*models.SSHSession, error)
	List(offset, limit int, filters map[string]interface{}) ([]models.SSHSession, int64, error)
	Update(session *models.SSHSession) error
	CloseSession(id uint64) error
}

type sshConnectionRepository struct {
	db *gorm.DB
}

func NewSSHConnectionRepository(db *gorm.DB) SSHConnectionRepository {
	return &sshConnectionRepository{db: db}
}

func (r *sshConnectionRepository) Create(conn *models.SSHConnection) error {
	return r.db.Create(conn).Error
}

func (r *sshConnectionRepository) GetByID(id uint64) (*models.SSHConnection, error) {
	var conn models.SSHConnection
	err := r.db.Preload("Project").Preload("Host").Preload("Key").First(&conn, id).Error
	return &conn, err
}

func (r *sshConnectionRepository) List(offset, limit int, filters map[string]interface{}) ([]models.SSHConnection, int64, error) {
	var conns []models.SSHConnection
	var total int64

	query := r.db.Model(&models.SSHConnection{})
	
	if projectID, ok := filters["project_id"].(uint64); ok && projectID > 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}
	if hostID, ok := filters["host_id"].(uint64); ok && hostID > 0 {
		query = query.Where("host_id = ?", hostID)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Project").Preload("Host").Preload("Key").Offset(offset).Limit(limit).Find(&conns).Error
	return conns, total, err
}

func (r *sshConnectionRepository) Update(conn *models.SSHConnection) error {
	return r.db.Save(conn).Error
}

func (r *sshConnectionRepository) Delete(id uint64) error {
	return r.db.Delete(&models.SSHConnection{}, id).Error
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

func (r *sshKeyRepository) GetByFingerprint(fingerprint string) (*models.SSHKey, error) {
	var key models.SSHKey
	err := r.db.Where("fingerprint = ?", fingerprint).First(&key).Error
	return &key, err
}

func (r *sshKeyRepository) List(offset, limit int) ([]models.SSHKey, int64, error) {
	var keys []models.SSHKey
	var total int64

	err := r.db.Model(&models.SSHKey{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Offset(offset).Limit(limit).Find(&keys).Error
	return keys, total, err
}

func (r *sshKeyRepository) Update(key *models.SSHKey) error {
	return r.db.Save(key).Error
}

func (r *sshKeyRepository) Delete(id uint64) error {
	return r.db.Delete(&models.SSHKey{}, id).Error
}

type sshSessionRepository struct {
	db *gorm.DB
}

func NewSSHSessionRepository(db *gorm.DB) SSHSessionRepository {
	return &sshSessionRepository{db: db}
}

func (r *sshSessionRepository) Create(session *models.SSHSession) error {
	return r.db.Create(session).Error
}

func (r *sshSessionRepository) GetByID(id uint64) (*models.SSHSession, error) {
	var session models.SSHSession
	err := r.db.Preload("Connection").Preload("User").First(&session, id).Error
	return &session, err
}

func (r *sshSessionRepository) GetByToken(token string) (*models.SSHSession, error) {
	var session models.SSHSession
	err := r.db.Preload("Connection").Preload("User").
		Where("session_token = ?", token).First(&session).Error
	return &session, err
}

func (r *sshSessionRepository) List(offset, limit int, filters map[string]interface{}) ([]models.SSHSession, int64, error) {
	var sessions []models.SSHSession
	var total int64

	query := r.db.Model(&models.SSHSession{})
	
	if userID, ok := filters["user_id"].(uint64); ok && userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if connectionID, ok := filters["connection_id"].(uint64); ok && connectionID > 0 {
		query = query.Where("connection_id = ?", connectionID)
	}
	if status, ok := filters["status"].(string); ok && status != "" {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("Connection").Preload("User").
		Order("started_at DESC").Offset(offset).Limit(limit).Find(&sessions).Error
	return sessions, total, err
}

func (r *sshSessionRepository) Update(session *models.SSHSession) error {
	return r.db.Save(session).Error
}

func (r *sshSessionRepository) CloseSession(id uint64) error {
	return r.db.Model(&models.SSHSession{}).Where("id = ?", id).
		Updates(map[string]interface{}{
			"status": "closed",
		}).Error
}

