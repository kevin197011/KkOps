// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package externalsystem

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/utils"
)

// Service handles external ops system CRUD and SSO launch
type Service struct {
	db *gorm.DB
}

// NewService creates a new service
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// List returns enabled external systems for display (no secret)
func (s *Service) List(enabledOnly bool) ([]model.ExternalSystem, error) {
	var list []model.ExternalSystem
	q := s.db.Model(&model.ExternalSystem{})
	if enabledOnly {
		q = q.Where("enabled = ?", true)
	}
	if err := q.Order("order_index ASC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// GetByID returns one system by id (no secret in response is handled by handler)
func (s *Service) GetByID(id uint) (*model.ExternalSystem, error) {
	var sys model.ExternalSystem
	if err := s.db.First(&sys, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("external system not found")
		}
		return nil, err
	}
	return &sys, nil
}

// CreateRequest for creating external system
type CreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	LaunchType  string `json:"launch_type"` // sso_link | jwt_token
	BaseURL     string `json:"base_url" binding:"required"`
	LoginPath   string `json:"login_path"`
	Secret      string `json:"secret"`
	RoleMapping string `json:"role_mapping"`
	Icon        string `json:"icon"`
	OrderIndex  int    `json:"order_index"`
	Enabled     bool   `json:"enabled"`
}

// UpdateRequest for updating (all optional)
type UpdateRequest struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	LaunchType  *string `json:"launch_type"`
	BaseURL     *string `json:"base_url"`
	LoginPath   *string `json:"login_path"`
	Secret      *string `json:"secret"`
	RoleMapping *string `json:"role_mapping"`
	Icon        *string `json:"icon"`
	OrderIndex  *int    `json:"order_index"`
	Enabled     *bool   `json:"enabled"`
}

func validateBaseURL(u string) error {
	u = strings.TrimSpace(u)
	if u == "" {
		return errors.New("base_url is required")
	}
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return errors.New("base_url must start with http:// or https://")
	}
	_, err := url.Parse(u)
	return err
}

// Create creates an external system
func (s *Service) Create(req *CreateRequest) (*model.ExternalSystem, error) {
	if err := validateBaseURL(req.BaseURL); err != nil {
		return nil, err
	}
	launchType := strings.TrimSpace(req.LaunchType)
	if launchType == "" {
		launchType = model.LaunchTypeJWTToken
	}
	if launchType != model.LaunchTypeSSOLink && launchType != model.LaunchTypeJWTToken {
		launchType = model.LaunchTypeJWTToken
	}

	sys := &model.ExternalSystem{
		Name:        req.Name,
		Description: req.Description,
		LaunchType:  launchType,
		BaseURL:     strings.TrimSpace(req.BaseURL),
		RoleMapping: req.RoleMapping,
		Icon:        req.Icon,
		OrderIndex:  req.OrderIndex,
		Enabled:     req.Enabled,
	}
	if launchType == model.LaunchTypeJWTToken {
		loginPath := strings.TrimSpace(req.LoginPath)
		if loginPath == "" {
			loginPath = "/sso/consume"
		}
		if !strings.HasPrefix(loginPath, "/") {
			loginPath = "/" + loginPath
		}
		sys.LoginPath = loginPath
		sys.Secret = req.Secret
		sys.BaseURL = strings.TrimSuffix(sys.BaseURL, "/")
	}

	if err := s.db.Create(sys).Error; err != nil {
		return nil, err
	}
	return sys, nil
}

// Update updates an external system
func (s *Service) Update(id uint, req *UpdateRequest) (*model.ExternalSystem, error) {
	sys, err := s.GetByID(id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		sys.Name = *req.Name
	}
	if req.Description != nil {
		sys.Description = *req.Description
	}
	if req.LaunchType != nil {
		t := strings.TrimSpace(*req.LaunchType)
		if t == model.LaunchTypeSSOLink || t == model.LaunchTypeJWTToken {
			sys.LaunchType = t
		}
	}
	if req.BaseURL != nil {
		if err := validateBaseURL(*req.BaseURL); err != nil {
			return nil, err
		}
		sys.BaseURL = strings.TrimSpace(*req.BaseURL)
		if sys.LaunchType == model.LaunchTypeJWTToken {
			sys.BaseURL = strings.TrimSuffix(sys.BaseURL, "/")
		}
	}
	if req.LoginPath != nil && sys.LaunchType == model.LaunchTypeJWTToken {
		p := strings.TrimSpace(*req.LoginPath)
		if p != "" && !strings.HasPrefix(p, "/") {
			p = "/" + p
		}
		if p != "" {
			sys.LoginPath = p
		}
	}
	if req.Secret != nil && *req.Secret != "" && sys.LaunchType == model.LaunchTypeJWTToken {
		sys.Secret = *req.Secret
	}
	if req.RoleMapping != nil {
		sys.RoleMapping = *req.RoleMapping
	}
	if req.Icon != nil {
		sys.Icon = *req.Icon
	}
	if req.OrderIndex != nil {
		sys.OrderIndex = *req.OrderIndex
	}
	if req.Enabled != nil {
		sys.Enabled = *req.Enabled
	}
	if err := s.db.Save(sys).Error; err != nil {
		return nil, err
	}
	return sys, nil
}

// Delete deletes an external system
func (s *Service) Delete(id uint) error {
	res := s.db.Delete(&model.ExternalSystem{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("external system not found")
	}
	return nil
}

// Launch builds redirect URL: for sso_link just return base_url; for jwt_token build JWT and append
func (s *Service) Launch(systemID uint, userID uint, username, email, realName string, roles, permissions []string) (redirectURL string, err error) {
	sys, err := s.GetByID(systemID)
	if err != nil {
		return "", err
	}
	if !sys.Enabled {
		return "", errors.New("external system is disabled")
	}

	// Same SSO: user already logged in to IdP (e.g. Keycloak), just open URL; target will use same IdP
	if sys.LaunchType == model.LaunchTypeSSOLink {
		return sys.BaseURL, nil
	}
	// jwt_token or legacy (empty): build JWT and redirect

	// jwt_token: build JWT and redirect to login_path
	mappedRoles := mapRoles(sys.RoleMapping, roles)
	token, err := utils.GenerateOutboundSSOToken(userID, username, email, realName, roles, permissions, mappedRoles, sys.Secret, 300)
	if err != nil {
		return "", fmt.Errorf("generate token: %w", err)
	}

	base := strings.TrimSuffix(sys.BaseURL, "/")
	path := sys.LoginPath
	if path == "" {
		path = "/sso/consume"
	}
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := base + path
	if strings.Contains(u, "?") {
		u += "&token=" + url.QueryEscape(token)
	} else {
		u += "?token=" + url.QueryEscape(token)
	}
	return u, nil
}

// mapRoles applies RoleMapping JSON to user roles and returns target system role names (unique)
func mapRoles(roleMappingJSON string, userRoles []string) []string {
	if roleMappingJSON == "" {
		return userRoles
	}
	var m map[string]string
	if err := json.Unmarshal([]byte(roleMappingJSON), &m); err != nil {
		return userRoles
	}
	seen := make(map[string]bool)
	var out []string
	for _, r := range userRoles {
		if mapped, ok := m[r]; ok && mapped != "" && !seen[mapped] {
			seen[mapped] = true
			out = append(out, mapped)
		} else if !seen[r] {
			seen[r] = true
			out = append(out, r)
		}
	}
	return out
}
