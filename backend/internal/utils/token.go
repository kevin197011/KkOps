// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"crypto/rand"
	"encoding/base64"
	"time"
)

// GenerateAPIToken generates a secure random API token
func GenerateAPIToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashAPIToken hashes an API token using bcrypt
func HashAPIToken(token string) (string, error) {
	return HashPassword(token) // Reuse password hashing
}

// CheckAPITokenHash compares an API token with a hash
func CheckAPITokenHash(token, hash string) bool {
	return CheckPasswordHash(token, hash)
}

// GetTokenPrefix returns the first 8 characters of a token for display
func GetTokenPrefix(token string) string {
	if len(token) < 8 {
		return token
	}
	return token[:8] + "..."
}

// IsTokenExpired checks if a token is expired
func IsTokenExpired(expiresAt *time.Time) bool {
	if expiresAt == nil {
		return false // No expiration
	}
	return time.Now().After(*expiresAt)
}
