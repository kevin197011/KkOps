// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// OAuth2Client represents an OIDC/OAuth2 client (application) that can use KkOps as IdP
type OAuth2Client struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	ClientID     string         `gorm:"uniqueIndex;not null;size:128" json:"client_id"`
	ClientSecret string         `gorm:"not null;size:255" json:"-"` // hashed
	Name         string         `gorm:"not null;size:100" json:"name"`
	Protocol     string         `gorm:"not null;size:16;default:oidc" json:"protocol"` // oidc | saml | ldap
	RedirectURIs string         `gorm:"type:text;not null" json:"-"`                   // JSON array of allowed redirect_uri
	Scopes       string         `gorm:"type:text" json:"-"`                            // space-separated, e.g. "openid profile email"
	Enabled      bool           `gorm:"default:true" json:"enabled"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName specifies table name
func (OAuth2Client) TableName() string {
	return "oauth2_clients"
}
