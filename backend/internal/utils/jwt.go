// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID   uint     `json:"user_id"`
	Username string   `json:"username"`
	Roles    []string `json:"roles"`
	// TokenUse is "refresh" for refresh tokens; empty or "access" for API access tokens.
	TokenUse string `json:"token_use,omitempty"`
	jwt.RegisteredClaims
}

// GenerateJWT generates a JWT token for a user.
// tokenUse should be "" or "access" for access tokens, "refresh" for refresh tokens.
func GenerateJWT(userID uint, username string, roles []string, secret string, expiresIn int, tokenUse string) (string, error) {
	claims := JWTClaims{
		UserID:   userID,
		Username: username,
		Roles:    roles,
		TokenUse: tokenUse,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresIn) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString, secret string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// OutboundSSOClaims is the JWT payload sent to external ops systems when launching SSO
type OutboundSSOClaims struct {
	UserID      uint     `json:"sub"` // subject = user id
	Username    string   `json:"username"`
	Email       string   `json:"email"`
	RealName    string   `json:"real_name"`
	Roles       []string `json:"roles"`
	Permissions []string `json:"permissions"`  // resource:action pairs
	MappedRoles []string `json:"mapped_roles"` // KkOps roles mapped to target system role names (by RoleMapping)
	jwt.RegisteredClaims
}

// GenerateOutboundSSOToken creates a short-lived JWT for redirecting to an external system (signed with shared secret)
func GenerateOutboundSSOToken(userID uint, username, email, realName string, roles, permissions, mappedRoles []string, secret string, expiresInSeconds int) (string, error) {
	if expiresInSeconds <= 0 {
		expiresInSeconds = 300
	}
	claims := OutboundSSOClaims{
		UserID:      userID,
		Username:    username,
		Email:       email,
		RealName:    realName,
		Roles:       roles,
		Permissions: permissions,
		MappedRoles: mappedRoles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expiresInSeconds) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
