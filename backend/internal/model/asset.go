// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// Project represents a project
type Project struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100;uniqueIndex" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedBy   uint           `json:"created_by"`
	Creator     User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Assets []Asset `gorm:"foreignKey:ProjectID" json:"assets,omitempty"`
}

// AssetCategory represents an asset category
type AssetCategory struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	ParentID    *uint          `json:"parent_id"`
	Parent      *AssetCategory `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Children []AssetCategory `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

// Asset represents an IT asset
type Asset struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	HostName        string         `gorm:"not null;size:100;column:host_name;uniqueIndex" json:"hostName"`
	ProjectID       *uint          `json:"project_id"`
	Project         *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	CloudPlatformID *uint          `json:"cloud_platform_id"`
	CloudPlatform   *CloudPlatform `gorm:"foreignKey:CloudPlatformID" json:"cloud_platform,omitempty"`
	EnvironmentID   *uint          `json:"environment_id"`
	Environment     *Environment   `gorm:"foreignKey:EnvironmentID" json:"environment,omitempty"`
	IP              string         `gorm:"size:45;index" json:"ip"`
	SSHPort         int            `gorm:"default:22" json:"ssh_port"`
	SSHKeyID        *uint          `json:"ssh_key_id"`
	SSHKey          *SSHKey        `gorm:"foreignKey:SSHKeyID" json:"ssh_key,omitempty"`
	SSHUser         string         `gorm:"size:50" json:"ssh_user"`
	CPU             string         `gorm:"size:50" json:"cpu"`
	Memory          string         `gorm:"size:50" json:"memory"`
	Disk            string         `gorm:"size:50" json:"disk"`
	Status          string         `gorm:"default:active;size:20;index" json:"status"` // active, disabled
	Description     string         `gorm:"type:text" json:"description"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Tags []Tag `gorm:"many2many:asset_tags;" json:"tags,omitempty"`
}

// Tag represents a tag for categorizing assets
type Tag struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Color       string         `gorm:"size:20" json:"color"` // hex color code
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Assets []Asset `gorm:"many2many:asset_tags;" json:"assets,omitempty"`
}

// AssetTag represents the many-to-many relationship between assets and tags
type AssetTag struct {
	AssetID   uint      `gorm:"primaryKey" json:"asset_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
}

// SSHKey represents an SSH key for asset connections
type SSHKey struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	UserID      uint           `gorm:"not null;index" json:"user_id"` // Owner user ID
	User        User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Type        string         `gorm:"size:50" json:"type"` // rsa, ed25519, ecdsa, etc.
	PrivateKey  string         `gorm:"type:text" json:"-"`  // Encrypted private key
	PublicKey   string         `gorm:"type:text" json:"public_key"`
	Fingerprint string         `gorm:"size:100;index" json:"fingerprint"`
	SSHUser     string         `gorm:"size:50" json:"ssh_user"` // Default SSH user (username for SSH connection)
	Passphrase  string         `gorm:"type:text" json:"-"`      // Encrypted passphrase (if key has password)
	Description string         `gorm:"type:text" json:"description"`
	LastUsedAt  *time.Time     `json:"last_used_at"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Assets []Asset `gorm:"foreignKey:SSHKeyID" json:"assets,omitempty"`
}
