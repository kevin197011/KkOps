// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package idp

import (
	"crypto/rsa"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// idTokenClaims for OIDC id_token (signed with RS256)
type idTokenClaims struct {
	PreferredUsername string `json:"preferred_username"`
	Name              string `json:"name,omitempty"`
	Email             string `json:"email,omitempty"`
	jwt.RegisteredClaims
}

// accessTokenClaims for OAuth2 access_token (used to call userinfo)
type accessTokenClaims struct {
	ClientID string `json:"client_id,omitempty"`
	Scope    string `json:"scope,omitempty"`
	jwt.RegisteredClaims
}

func signIDToken(issuer, clientID, userID string, username, name, email string, key *rsa.PrivateKey, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := idTokenClaims{
		PreferredUsername: username,
		Name:              name,
		Email:             email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			Audience:  jwt.ClaimStrings{clientID},
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = getKeyID()
	return tok.SignedString(key)
}

func signAccessToken(issuer, clientID, userID, scope string, key *rsa.PrivateKey, expiry time.Duration) (string, error) {
	now := time.Now()
	claims := accessTokenClaims{
		ClientID: clientID,
		Scope:    scope,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    issuer,
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(now.Add(expiry)),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	tok.Header["kid"] = getKeyID()
	return tok.SignedString(key)
}

func verifyAccessToken(tokenStr string, pubKey *rsa.PublicKey) (*accessTokenClaims, error) {
	tok, err := jwt.ParseWithClaims(tokenStr, &accessTokenClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return pubKey, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := tok.Claims.(*accessTokenClaims)
	if !ok || !tok.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}

// userInfoResponse for GET /userinfo
func userInfoResponse(sub, username, name, email string) ([]byte, error) {
	m := map[string]interface{}{
		"sub":                sub,
		"preferred_username": username,
		"name":               name,
		"email":              email,
	}
	return json.Marshal(m)
}

func parseRedirectURIs(jsonStr string) ([]string, error) {
	if jsonStr == "" {
		return nil, nil
	}
	var uris []string
	if err := json.Unmarshal([]byte(jsonStr), &uris); err != nil {
		return nil, fmt.Errorf("redirect_uris: %w", err)
	}
	return uris, nil
}
