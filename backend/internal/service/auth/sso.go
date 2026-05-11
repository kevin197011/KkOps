// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package auth

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/coreos/go-oidc/v3/oidc"
	"golang.org/x/oauth2"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

const (
	ssoStateValidity = 10 * time.Minute
	ssoProviderOIDC  = "oidc"
	// placeholder password hash for SSO-only users (valid bcrypt, no password will match)
	ssoNoPasswordBcrypt = "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZRGdjGj/n3.4RlKZP2tMp2qKqKqKqKq"
)

// SSOConfigResponse is returned to frontend for SSO login entry
type SSOConfigResponse struct {
	Enabled bool   `json:"enabled"`
	Label   string `json:"label,omitempty"` // e.g. "SSO 登录"
}

// GetSSOConfig returns whether SSO is enabled (no auth required)
func (s *Service) GetSSOConfig() SSOConfigResponse {
	return SSOConfigResponse{
		Enabled: s.config.SSO.Enabled && s.config.SSO.OIDC.IssuerURL != "" && s.config.SSO.OIDC.ClientID != "",
		Label:   "SSO 登录",
	}
}

// GetSSOLoginURL returns the IdP authorization URL with state (caller must store state for verification in callback)
func (s *Service) GetSSOLoginURL(state string) (string, error) {
	if !s.config.SSO.Enabled || s.config.SSO.OIDC.IssuerURL == "" || s.config.SSO.OIDC.ClientID == "" {
		return "", errors.New("SSO is not configured")
	}

	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, s.config.SSO.OIDC.IssuerURL)
	if err != nil {
		return "", fmt.Errorf("oidc provider: %w", err)
	}

	scopes := s.config.SSO.OIDC.Scopes
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	oauth2Config := oauth2.Config{
		ClientID:     s.config.SSO.OIDC.ClientID,
		ClientSecret: s.config.SSO.OIDC.ClientSecret,
		RedirectURL:  s.config.SSO.OIDC.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	return oauth2Config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// GenerateSSOState returns a signed state value and the raw state string for redirect
func (s *Service) GenerateSSOState() (state string, err error) {
	b := make([]byte, 32)
	if _, err = rand.Read(b); err != nil {
		return "", err
	}
	payload := base64.RawURLEncoding.EncodeToString(b) + "." + strconv.FormatInt(time.Now().Add(ssoStateValidity).Unix(), 10)
	mac := hmac.New(sha256.New, []byte(s.config.JWT.Secret))
	mac.Write([]byte(payload))
	sig := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
	return payload + "." + sig, nil
}

// VerifySSOState checks state signature and expiry
func (s *Service) VerifySSOState(state string) bool {
	parts := strings.SplitN(state, ".", 3)
	if len(parts) != 3 {
		return false
	}
	payload := parts[0] + "." + parts[1]
	sig, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return false
	}
	mac := hmac.New(sha256.New, []byte(s.config.JWT.Secret))
	mac.Write([]byte(payload))
	if !hmac.Equal(mac.Sum(nil), sig) {
		return false
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return false
	}
	return time.Now().Unix() < exp
}

// ExchangeSSOCode exchanges authorization code for tokens, verifies id_token, and finds or creates user; returns JWT and user info
func (s *Service) ExchangeSSOCode(ctx context.Context, code, state string) (*LoginResponse, error) {
	if !s.VerifySSOState(state) {
		return nil, errors.New("invalid or expired state")
	}

	if !s.config.SSO.Enabled || s.config.SSO.OIDC.IssuerURL == "" || s.config.SSO.OIDC.ClientID == "" {
		return nil, errors.New("SSO is not configured")
	}

	provider, err := oidc.NewProvider(ctx, s.config.SSO.OIDC.IssuerURL)
	if err != nil {
		return nil, fmt.Errorf("oidc provider: %w", err)
	}

	scopes := s.config.SSO.OIDC.Scopes
	if len(scopes) == 0 {
		scopes = []string{oidc.ScopeOpenID, "profile", "email"}
	}

	oauth2Config := oauth2.Config{
		ClientID:     s.config.SSO.OIDC.ClientID,
		ClientSecret: s.config.SSO.OIDC.ClientSecret,
		RedirectURL:  s.config.SSO.OIDC.RedirectURL,
		Endpoint:     provider.Endpoint(),
		Scopes:       scopes,
	}

	token, err := oauth2Config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("token exchange: %w", err)
	}

	rawIDToken, ok := token.Extra("id_token").(string)
	if !ok {
		return nil, errors.New("no id_token in response")
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: s.config.SSO.OIDC.ClientID})
	idToken, err := verifier.Verify(ctx, rawIDToken)
	if err != nil {
		return nil, fmt.Errorf("id_token verify: %w", err)
	}

	var claims struct {
		Sub           string `json:"sub"`
		Email         string `json:"email"`
		EmailVerified bool   `json:"email_verified"`
		Name          string `json:"name"`
		PreferredName string `json:"preferred_username"`
	}
	if err := idToken.Claims(&claims); err != nil {
		return nil, fmt.Errorf("parse claims: %w", err)
	}

	username := claims.PreferredName
	if username == "" {
		username = claims.Sub
	}
	if username == "" {
		return nil, errors.New("oidc: no subject or preferred_username")
	}

	email := claims.Email
	if email == "" {
		email = username + "@sso.local"
	}

	user, err := s.findOrCreateSSOUser(claims.Sub, username, email, claims.Name)
	if err != nil {
		return nil, err
	}

	if user.Status != "active" {
		return nil, errors.New("user account is disabled")
	}

	roleNames := make([]string, len(user.Roles))
	for i, role := range user.Roles {
		roleNames[i] = role.Name
	}

	jwtToken, err := utils.GenerateJWT(user.ID, user.Username, roleNames, s.config.JWT.Secret, s.config.JWT.ExpiresIn, "")
	if err != nil {
		return nil, err
	}
	refreshToken, err := utils.GenerateJWT(user.ID, user.Username, roleNames, s.config.JWT.Secret, s.config.JWT.RefreshExpiresIn, "refresh")
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user.LastLoginAt = &now
	s.db.Model(&user).Update("last_login_at", now)

	return &LoginResponse{
		Token:        jwtToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.config.JWT.ExpiresIn,
		User: UserInfo{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			RealName: user.RealName,
			Roles:    roleNames,
		},
	}, nil
}

// findOrCreateSSOUser finds user by external_id (sub) or creates one; returns with Roles preloaded
func (s *Service) findOrCreateSSOUser(sub, username, email, realName string) (*model.User, error) {
	var user model.User
	err := s.db.Where("external_id = ? AND sso_provider = ?", sub, ssoProviderOIDC).Preload("Roles").First(&user).Error
	if err == nil {
		// Update profile from IdP
		updates := map[string]interface{}{"email": email, "real_name": realName, "updated_at": time.Now()}
		s.db.Model(&user).Updates(updates)
		return &user, nil
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Username might already exist as local user - then we must use a different internal username or fail
	var existing model.User
	if err := s.db.Where("username = ?", username).First(&existing).Error; err == nil {
		// Prefer linking by external_id; if this is same user (e.g. same email) we could update to SSO. For simplicity: require unique username.
		// Option: use sub as username suffix when conflict. Here we require admin to pre-create or use unique preferred_username.
		return nil, fmt.Errorf("username %s already exists; use a different IdP username or contact admin", username)
	}

	user = model.User{
		Username:     username,
		PasswordHash: ssoNoPasswordBcrypt,
		Email:        email,
		RealName:     realName,
		Status:       "active",
		Source:       "sso",
		ExternalID:   sub,
		SSOProvider:  ssoProviderOIDC,
	}
	if err := s.db.Create(&user).Error; err != nil {
		return nil, err
	}
	// Load roles (empty for new user)
	s.db.Preload("Roles").First(&user, user.ID)
	return &user, nil
}

// FrontendCallbackURL returns the URL to redirect the browser after SSO callback (with token in query)
func (s *Service) FrontendCallbackURL(token string, expiresIn int) string {
	path := "/auth/callback?token=" + url.QueryEscape(token) + "&expires_in=" + strconv.Itoa(expiresIn)
	if s.config.SSO.OIDC.FrontendBaseURL != "" {
		return strings.TrimSuffix(s.config.SSO.OIDC.FrontendBaseURL, "/") + path
	}
	return path
}
