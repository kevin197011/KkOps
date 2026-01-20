// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// CreateAPITokenRequest represents a request to create an API token
type CreateAPITokenRequest struct {
	Name        string     `json:"name" binding:"required"`
	Description string     `json:"description"`
	ExpiresAt   *time.Time `json:"expires_at"`
}

// APITokenResponse represents an API token response
type APITokenResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	Token       string     `json:"token"` // Only shown once on creation
	Prefix      string     `json:"prefix"`
	ExpiresAt   *time.Time `json:"expires_at"`
	Status      string     `json:"status"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"created_at"`
}

// CreateAPIToken creates a new API token for a user
func (s *Service) CreateAPIToken(userID uint, req *CreateAPITokenRequest) (*APITokenResponse, error) {
	// Generate token
	token, err := utils.GenerateAPIToken()
	if err != nil {
		return nil, err
	}

	// Hash token for validation
	tokenHash, err := utils.HashAPIToken(token)
	if err != nil {
		return nil, err
	}

	// Encrypt token for later retrieval
	encryptedToken, err := utils.Encrypt([]byte(token), s.config.Encryption.Key)
	if err != nil {
		return nil, err
	}

	// Create API token record
	apiToken := model.APIToken{
		UserID:         userID,
		Name:           req.Name,
		TokenHash:      tokenHash,
		TokenEncrypted: encryptedToken,
		ExpiresAt:      req.ExpiresAt,
		Status:         "active",
		Description:    req.Description,
	}

	if err := s.db.Create(&apiToken).Error; err != nil {
		return nil, err
	}

	return &APITokenResponse{
		ID:          apiToken.ID,
		Name:        apiToken.Name,
		Token:       token, // Return full token
		Prefix:      utils.GetTokenPrefix(token),
		ExpiresAt:   apiToken.ExpiresAt,
		Status:      apiToken.Status,
		Description: apiToken.Description,
		CreatedAt:   apiToken.CreatedAt,
	}, nil
}

// ValidateAPIToken validates an API token and returns the user ID
func (s *Service) ValidateAPIToken(token string) (uint, error) {
	// Get all active tokens and check each one
	var apiTokens []model.APIToken
	if err := s.db.Where("status = ?", "active").Find(&apiTokens).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return 0, errors.New("invalid token")
		}
		return 0, err
	}

	// Check each token's hash
	for _, apiToken := range apiTokens {
		// Check expiration
		if utils.IsTokenExpired(apiToken.ExpiresAt) {
			continue
		}

		// Verify token hash
		if utils.CheckAPITokenHash(token, apiToken.TokenHash) {
			// Update last used time
			now := time.Now()
			apiToken.LastUsedAt = &now
			s.db.Save(&apiToken)
			return apiToken.UserID, nil
		}
	}

	return 0, errors.New("invalid token")
}

// ListAPITokens lists all API tokens for a user
func (s *Service) ListAPITokens(userID uint) ([]APITokenResponse, error) {
	var tokens []model.APIToken
	if err := s.db.Where("user_id = ?", userID).Find(&tokens).Error; err != nil {
		return nil, err
	}

	result := make([]APITokenResponse, len(tokens))
	for i, token := range tokens {
		// Decrypt token to get prefix for display
		prefix := "****"
		if token.TokenEncrypted != "" {
			decrypted, err := utils.Decrypt(token.TokenEncrypted, s.config.Encryption.Key)
			if err == nil {
				tokenStr := string(decrypted)
				prefix = utils.GetTokenPrefix(tokenStr)
			}
		}

		result[i] = APITokenResponse{
			ID:          token.ID,
			Name:        token.Name,
			Prefix:      prefix,
			ExpiresAt:   token.ExpiresAt,
			Status:      token.Status,
			Description: token.Description,
			CreatedAt:   token.CreatedAt,
		}
	}

	return result, nil
}

// GetAPIToken retrieves the full API token by ID (for the token owner)
func (s *Service) GetAPIToken(userID, tokenID uint) (*APITokenResponse, error) {
	var token model.APIToken
	if err := s.db.Where("id = ? AND user_id = ?", tokenID, userID).First(&token).Error; err != nil {
		return nil, errors.New("token not found")
	}

	// Decrypt token
	var fullToken string
	if token.TokenEncrypted != "" {
		decrypted, err := utils.Decrypt(token.TokenEncrypted, s.config.Encryption.Key)
		if err != nil {
			return nil, errors.New("failed to decrypt token")
		}
		fullToken = string(decrypted)
	} else {
		// For old tokens without encrypted version, return error
		return nil, errors.New("token cannot be retrieved (created before encryption was enabled)")
	}

	return &APITokenResponse{
		ID:          token.ID,
		Name:        token.Name,
		Token:       fullToken,
		Prefix:      utils.GetTokenPrefix(fullToken),
		ExpiresAt:   token.ExpiresAt,
		Status:      token.Status,
		Description: token.Description,
		CreatedAt:   token.CreatedAt,
	}, nil
}

// DeleteAPIToken deletes an API token
func (s *Service) DeleteAPIToken(userID, tokenID uint) error {
	var token model.APIToken
	if err := s.db.Where("id = ? AND user_id = ?", tokenID, userID).First(&token).Error; err != nil {
		return err
	}

	return s.db.Delete(&token).Error
}
