// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package sshkey

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"golang.org/x/crypto/ssh"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// Service handles SSH key management operations
type Service struct {
	db     *gorm.DB
	config *config.Config
}

// NewService creates a new SSH key service
func NewService(db *gorm.DB, cfg *config.Config) *Service {
	return &Service{
		db:     db,
		config: cfg,
	}
}

// CreateSSHKeyRequest represents a request to create an SSH key
type CreateSSHKeyRequest struct {
	Name        string `json:"name" binding:"required"`
	Type        string `json:"type"` // rsa, ed25519, ecdsa, etc.
	PublicKey   string `json:"public_key"`  // Optional - auto-extracted from private key if empty
	PrivateKey  string `json:"private_key" binding:"required"`
	SSHUser     string `json:"ssh_user"`   // Default SSH username
	Passphrase  string `json:"passphrase"` // Passphrase for encrypted private key (if provided)
	Description string `json:"description"`
}

// UpdateSSHKeyRequest represents a request to update an SSH key
type UpdateSSHKeyRequest struct {
	Name        string `json:"name"`
	SSHUser     string `json:"ssh_user"`
	Description string `json:"description"`
}

// TestSSHKeyRequest represents a request to test an SSH key connection
type TestSSHKeyRequest struct {
	Host       string `json:"host" binding:"required"`
	Port       int    `json:"port"`
	Username   string `json:"username"`   // Optional, uses key's SSHUser if not provided
	Passphrase string `json:"passphrase"` // Optional, if key has password
}

// SSHKeyResponse represents an SSH key response (without private key)
type SSHKeyResponse struct {
	ID          uint       `json:"id"`
	UserID      uint       `json:"user_id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	PublicKey   string     `json:"public_key"`
	Fingerprint string     `json:"fingerprint"`
	SSHUser     string     `json:"ssh_user"`
	Description string     `json:"description"`
	LastUsedAt  *time.Time `json:"last_used_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// ListSSHKeys lists all SSH keys for a user
func (s *Service) ListSSHKeys(userID uint) ([]SSHKeyResponse, error) {
	var keys []model.SSHKey
	if err := s.db.Where("user_id = ?", userID).Find(&keys).Error; err != nil {
		return nil, err
	}

	responses := make([]SSHKeyResponse, len(keys))
	for i, key := range keys {
		responses[i] = SSHKeyResponse{
			ID:          key.ID,
			UserID:      key.UserID,
			Name:        key.Name,
			Type:        key.Type,
			PublicKey:   key.PublicKey,
			Fingerprint: key.Fingerprint,
			SSHUser:     key.SSHUser,
			Description: key.Description,
			LastUsedAt:  key.LastUsedAt,
			CreatedAt:   key.CreatedAt,
			UpdatedAt:   key.UpdatedAt,
		}
	}

	return responses, nil
}

// GetSSHKey gets an SSH key by ID (only if owned by user)
func (s *Service) GetSSHKey(userID, keyID uint) (*SSHKeyResponse, error) {
	var key model.SSHKey
	if err := s.db.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		return nil, err
	}

	return &SSHKeyResponse{
		ID:          key.ID,
		UserID:      key.UserID,
		Name:        key.Name,
		Type:        key.Type,
		PublicKey:   key.PublicKey,
		Fingerprint: key.Fingerprint,
		SSHUser:     key.SSHUser,
		Description: key.Description,
		LastUsedAt:  key.LastUsedAt,
		CreatedAt:   key.CreatedAt,
		UpdatedAt:   key.UpdatedAt,
	}, nil
}

// CreateSSHKey creates a new SSH key
func (s *Service) CreateSSHKey(userID uint, req CreateSSHKeyRequest) (*SSHKeyResponse, error) {
	publicKey := req.PublicKey

	// If public key is not provided, extract it from private key
	if publicKey == "" {
		extracted, err := extractPublicKeyFromPrivateKey(req.PrivateKey, req.Passphrase)
		if err != nil {
			return nil, fmt.Errorf("failed to extract public key from private key: %w", err)
		}
		publicKey = extracted
	}

	// Validate that at least one key is provided
	if publicKey == "" {
		return nil, fmt.Errorf("public key is required (provide either public_key or private_key)")
	}

	// Calculate fingerprint from public key
	fingerprint, err := calculateFingerprint(publicKey)
	if err != nil {
		return nil, fmt.Errorf("invalid public key: %w", err)
	}

	// Detect key type if not provided
	keyType := req.Type
	if keyType == "" {
		keyType = detectKeyType(publicKey)
	}

	// Encrypt private key
	encryptedPrivateKey, err := utils.Encrypt([]byte(req.PrivateKey), s.config.Encryption.Key)
	if err != nil {
		return nil, fmt.Errorf("failed to encrypt private key: %w", err)
	}

	// Encrypt passphrase if provided
	var encryptedPassphrase string
	if req.Passphrase != "" {
		encryptedPassphrase, err = utils.Encrypt([]byte(req.Passphrase), s.config.Encryption.Key)
		if err != nil {
			return nil, fmt.Errorf("failed to encrypt passphrase: %w", err)
		}
	}

	// Check for duplicate fingerprint
	var existingKey model.SSHKey
	if err := s.db.Where("fingerprint = ?", fingerprint).First(&existingKey).Error; err == nil {
		return nil, fmt.Errorf("SSH key with this fingerprint already exists")
	}

	// Create SSH key
	key := model.SSHKey{
		UserID:      userID,
		Name:        req.Name,
		Type:        keyType,
		PrivateKey:  encryptedPrivateKey,
		PublicKey:   publicKey,
		Fingerprint: fingerprint,
		SSHUser:     req.SSHUser,
		Passphrase:  encryptedPassphrase,
		Description: req.Description,
	}

	if err := s.db.Create(&key).Error; err != nil {
		return nil, err
	}

	return &SSHKeyResponse{
		ID:          key.ID,
		UserID:      key.UserID,
		Name:        key.Name,
		Type:        key.Type,
		PublicKey:   key.PublicKey,
		Fingerprint: key.Fingerprint,
		SSHUser:     key.SSHUser,
		Description: key.Description,
		LastUsedAt:  key.LastUsedAt,
		CreatedAt:   key.CreatedAt,
		UpdatedAt:   key.UpdatedAt,
	}, nil
}

// UpdateSSHKey updates an SSH key (only name, SSHUser, and description)
func (s *Service) UpdateSSHKey(userID, keyID uint, req UpdateSSHKeyRequest) (*SSHKeyResponse, error) {
	var key model.SSHKey
	if err := s.db.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		return nil, err
	}

	// Update allowed fields
	if req.Name != "" {
		key.Name = req.Name
	}
	if req.SSHUser != "" {
		key.SSHUser = req.SSHUser
	}
	if req.Description != "" {
		key.Description = req.Description
	}

	if err := s.db.Save(&key).Error; err != nil {
		return nil, err
	}

	return &SSHKeyResponse{
		ID:          key.ID,
		UserID:      key.UserID,
		Name:        key.Name,
		Type:        key.Type,
		PublicKey:   key.PublicKey,
		Fingerprint: key.Fingerprint,
		SSHUser:     key.SSHUser,
		Description: key.Description,
		LastUsedAt:  key.LastUsedAt,
		CreatedAt:   key.CreatedAt,
		UpdatedAt:   key.UpdatedAt,
	}, nil
}

// DeleteSSHKey deletes an SSH key (only if not in use)
func (s *Service) DeleteSSHKey(userID, keyID uint) error {
	var key model.SSHKey
	if err := s.db.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		return err
	}

	// Check if key is in use
	var count int64
	if err := s.db.Model(&model.Asset{}).Where("ssh_key_id = ?", keyID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("SSH key is in use by %d asset(s)", count)
	}

	// Delete the key
	return s.db.Delete(&key).Error
}

// TestSSHKey tests an SSH key connection
func (s *Service) TestSSHKey(userID, keyID uint, req TestSSHKeyRequest) error {
	// Get the SSH key
	var key model.SSHKey
	if err := s.db.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		return err
	}

	// Decrypt private key
	privateKeyBytes, err := utils.Decrypt(key.PrivateKey, s.config.Encryption.Key)
	if err != nil {
		return fmt.Errorf("failed to decrypt private key: %w", err)
	}

	// Decrypt passphrase if exists
	var passphraseBytes []byte
	if key.Passphrase != "" {
		if req.Passphrase != "" {
			// Use provided passphrase
			passphraseBytes = []byte(req.Passphrase)
		} else {
			// Decrypt stored passphrase
			decrypted, err := utils.Decrypt(key.Passphrase, s.config.Encryption.Key)
			if err != nil {
				return fmt.Errorf("failed to decrypt passphrase: %w", err)
			}
			passphraseBytes = decrypted
		}
	}

	// Determine username
	username := req.Username
	if username == "" {
		username = key.SSHUser
	}
	if username == "" {
		return fmt.Errorf("username is required")
	}

	// Determine port
	port := req.Port
	if port == 0 {
		port = 22
	}

	// Test connection
	timeout := 10 * time.Second
	var client *utils.SSHClient
	if len(passphraseBytes) > 0 {
		client, err = utils.NewSSHClientWithPassphrase(req.Host, port, username, privateKeyBytes, passphraseBytes, timeout)
	} else {
		client, err = utils.NewSSHClient(req.Host, port, username, privateKeyBytes, timeout)
	}
	if err != nil {
		return fmt.Errorf("SSH connection failed: %w", err)
	}
	defer client.Close()

	// Update last used time
	now := time.Now()
	key.LastUsedAt = &now
	s.db.Save(&key)

	return nil
}

// calculateFingerprint calculates the MD5 fingerprint of an SSH public key
func calculateFingerprint(publicKey string) (string, error) {
	// Parse the public key
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return "", err
	}

	// Get the key bytes
	keyBytes := pubKey.Marshal()

	// Calculate MD5 hash
	hash := md5.Sum(keyBytes)
	fingerprint := hex.EncodeToString(hash[:])

	// Format as colon-separated hex (e.g., "a1:b2:c3:...")
	parts := make([]string, 0, len(fingerprint)/2)
	for i := 0; i < len(fingerprint); i += 2 {
		parts = append(parts, fingerprint[i:i+2])
	}
	return strings.Join(parts, ":"), nil
}

// detectKeyType detects the SSH key type from the public key
func detectKeyType(publicKey string) string {
	// Parse the public key to get its type
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey([]byte(publicKey))
	if err != nil {
		return "unknown"
	}

	keyType := pubKey.Type()
	switch keyType {
	case ssh.KeyAlgoRSA, ssh.KeyAlgoRSASHA256, ssh.KeyAlgoRSASHA512:
		return "rsa"
	case ssh.KeyAlgoECDSA256, ssh.KeyAlgoECDSA384, ssh.KeyAlgoECDSA521:
		return "ecdsa"
	case ssh.KeyAlgoED25519:
		return "ed25519"
	case ssh.KeyAlgoDSA:
		return "dsa"
	default:
		return strings.ToLower(keyType)
	}
}

// GetDecryptedPrivateKey gets and decrypts the private key (for internal use only)
func (s *Service) GetDecryptedPrivateKey(userID, keyID uint) ([]byte, error) {
	var key model.SSHKey
	if err := s.db.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		return nil, err
	}

	return utils.Decrypt(key.PrivateKey, s.config.Encryption.Key)
}

// GetDecryptedPassphrase gets and decrypts the passphrase (for internal use only)
func (s *Service) GetDecryptedPassphrase(userID, keyID uint) ([]byte, error) {
	var key model.SSHKey
	if err := s.db.Where("id = ? AND user_id = ?", keyID, userID).First(&key).Error; err != nil {
		return nil, err
	}

	if key.Passphrase == "" {
		return nil, nil
	}

	return utils.Decrypt(key.Passphrase, s.config.Encryption.Key)
}

// extractPublicKeyFromPrivateKey extracts the public key from a private key
func extractPublicKeyFromPrivateKey(privateKey, passphrase string) (string, error) {
	privateKeyBytes := []byte(privateKey)
	var signer ssh.Signer
	var err error

	// Parse private key (with or without passphrase)
	if passphrase != "" {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(privateKeyBytes, []byte(passphrase))
		if err != nil {
			// Try without passphrase in case the key is not encrypted
			signer, err = ssh.ParsePrivateKey(privateKeyBytes)
			if err != nil {
				return "", fmt.Errorf("failed to parse private key (with or without passphrase): %w", err)
			}
		}
	} else {
		signer, err = ssh.ParsePrivateKey(privateKeyBytes)
		if err != nil {
			// Check if error suggests the key is encrypted
			if strings.Contains(err.Error(), "password") || strings.Contains(err.Error(), "passphrase") {
				return "", fmt.Errorf("private key appears to be encrypted, passphrase required")
			}
			return "", fmt.Errorf("failed to parse private key: %w", err)
		}
	}

	// Get the public key from the signer
	publicKey := signer.PublicKey()

	// Format public key in OpenSSH format
	// Marshal the public key
	authorizedKey := ssh.MarshalAuthorizedKey(publicKey)
	// Remove trailing newline and convert to string
	publicKeyStr := strings.TrimSpace(string(authorizedKey))

	return publicKeyStr, nil
}
