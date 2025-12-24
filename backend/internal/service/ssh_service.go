package service

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/kronos/backend/internal/models"
	"github.com/kronos/backend/internal/repository"
)

type SSHConnectionService interface {
	CreateConnection(conn *models.SSHConnection) error
	GetConnection(id uint64) (*models.SSHConnection, error)
	ListConnections(page, pageSize int, filters map[string]interface{}) ([]models.SSHConnection, int64, error)
	UpdateConnection(id uint64, conn *models.SSHConnection) error
	DeleteConnection(id uint64) error
}

type SSHKeyService interface {
	CreateKey(key *models.SSHKey) error
	GetKey(id uint64) (*models.SSHKey, error)
	ListKeys(page, pageSize int) ([]models.SSHKey, int64, error)
	UpdateKey(id uint64, key *models.SSHKey) error
	DeleteKey(id uint64) error
}

type SSHSessionService interface {
	CreateSession(session *models.SSHSession) error
	GetSession(id uint64) (*models.SSHSession, error)
	GetSessionByToken(token string) (*models.SSHSession, error)
	ListSessions(page, pageSize int, filters map[string]interface{}) ([]models.SSHSession, int64, error)
	CloseSession(id uint64) error
}

type sshConnectionService struct {
	connRepo repository.SSHConnectionRepository
}

func NewSSHConnectionService(connRepo repository.SSHConnectionRepository) SSHConnectionService {
	return &sshConnectionService{connRepo: connRepo}
}

func (s *sshConnectionService) CreateConnection(conn *models.SSHConnection) error {
	return s.connRepo.Create(conn)
}

func (s *sshConnectionService) GetConnection(id uint64) (*models.SSHConnection, error) {
	return s.connRepo.GetByID(id)
}

func (s *sshConnectionService) ListConnections(page, pageSize int, filters map[string]interface{}) ([]models.SSHConnection, int64, error) {
	offset := (page - 1) * pageSize
	return s.connRepo.List(offset, pageSize, filters)
}

func (s *sshConnectionService) UpdateConnection(id uint64, conn *models.SSHConnection) error {
	existingConn, err := s.connRepo.GetByID(id)
	if err != nil {
		return err
	}

	if conn.Name != "" {
		existingConn.Name = conn.Name
	}
	if conn.Hostname != "" {
		existingConn.Hostname = conn.Hostname
	}
	if conn.Port > 0 {
		existingConn.Port = conn.Port
	}
	if conn.Username != "" {
		existingConn.Username = conn.Username
	}
	if conn.AuthType != "" {
		existingConn.AuthType = conn.AuthType
	}
	if conn.Status != "" {
		existingConn.Status = conn.Status
	}

	return s.connRepo.Update(existingConn)
}

func (s *sshConnectionService) DeleteConnection(id uint64) error {
	return s.connRepo.Delete(id)
}

type sshKeyService struct {
	keyRepo repository.SSHKeyRepository
}

func NewSSHKeyService(keyRepo repository.SSHKeyRepository) SSHKeyService {
	return &sshKeyService{keyRepo: keyRepo}
}

func (s *sshKeyService) CreateKey(key *models.SSHKey) error {
	return s.keyRepo.Create(key)
}

func (s *sshKeyService) GetKey(id uint64) (*models.SSHKey, error) {
	return s.keyRepo.GetByID(id)
}

func (s *sshKeyService) ListKeys(page, pageSize int) ([]models.SSHKey, int64, error) {
	offset := (page - 1) * pageSize
	return s.keyRepo.List(offset, pageSize)
}

func (s *sshKeyService) UpdateKey(id uint64, key *models.SSHKey) error {
	existingKey, err := s.keyRepo.GetByID(id)
	if err != nil {
		return err
	}

	if key.Name != "" {
		existingKey.Name = key.Name
	}

	return s.keyRepo.Update(existingKey)
}

func (s *sshKeyService) DeleteKey(id uint64) error {
	return s.keyRepo.Delete(id)
}

type sshSessionService struct {
	sessionRepo repository.SSHSessionRepository
}

func NewSSHSessionService(sessionRepo repository.SSHSessionRepository) SSHSessionService {
	return &sshSessionService{sessionRepo: sessionRepo}
}

func (s *sshSessionService) CreateSession(session *models.SSHSession) error {
	// 生成会话令牌
	token, err := generateSessionToken()
	if err != nil {
		return err
	}
	session.SessionToken = token
	session.Status = "active"
	session.StartedAt = time.Now()

	return s.sessionRepo.Create(session)
}

func (s *sshSessionService) GetSession(id uint64) (*models.SSHSession, error) {
	return s.sessionRepo.GetByID(id)
}

func (s *sshSessionService) GetSessionByToken(token string) (*models.SSHSession, error) {
	return s.sessionRepo.GetByToken(token)
}

func (s *sshSessionService) ListSessions(page, pageSize int, filters map[string]interface{}) ([]models.SSHSession, int64, error) {
	offset := (page - 1) * pageSize
	return s.sessionRepo.List(offset, pageSize, filters)
}

func (s *sshSessionService) CloseSession(id uint64) error {
	session, err := s.sessionRepo.GetByID(id)
	if err != nil {
		return err
	}

	now := time.Now()
	session.EndedAt = &now
	if session.StartedAt.Unix() > 0 {
		duration := int(now.Sub(session.StartedAt).Seconds())
		session.DurationSeconds = &duration
	}
	session.Status = "closed"

	return s.sessionRepo.Update(session)
}

func generateSessionToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

