// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// LaunchType: sso_link = same SSO (e.g. Keycloak), just open URL; jwt_token = KkOps issues JWT with user/perms
const (
	LaunchTypeSSOLink  = "sso_link"  // 同 SSO 应用：与 KkOps 共用同一 IdP，打开即已登录
	LaunchTypeJWTToken = "jwt_token" // Token 跳转：目标系统接收 KkOps 签发的 JWT
)

// ExternalSystem represents an external ops system for SSO portal (open after one login)
type ExternalSystem struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	Name        string `gorm:"not null;size:100" json:"name"`
	Description string `gorm:"type:text" json:"description"`
	// LaunchType: sso_link = open base_url only (same IdP); jwt_token = build JWT and redirect
	LaunchType string `gorm:"size:32;default:jwt_token" json:"launch_type"`
	// For both: base_url (sso_link = full URL to open; jwt_token = scheme+host)
	BaseURL string `gorm:"not null;size:512" json:"base_url"`
	// For jwt_token only: path that accepts token, e.g. /api/sso/consume
	LoginPath   string         `gorm:"size:255" json:"login_path"`
	Secret      string         `gorm:"size:255" json:"-"` // for jwt_token: shared secret
	RoleMapping string         `gorm:"type:text" json:"role_mapping"`
	Icon        string         `gorm:"size:64" json:"icon"`
	OrderIndex  int            `gorm:"default:0" json:"order_index"`
	Enabled     bool           `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies table name
func (ExternalSystem) TableName() string {
	return "external_systems"
}
