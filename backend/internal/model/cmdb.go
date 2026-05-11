// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// CMDB asset kinds (logical configuration items, distinct from infrastructure hosts).
const (
	CMDBKindService  = "service"
	CMDBKindHost     = "host"
	CMDBKindDatabase = "database"
	CMDBKindCluster  = "cluster"
	CMDBKindOther    = "other"
)

// Relation types between CMDB assets.
const (
	RelationDependsOn = "depends_on"
	RelationRunsOn    = "runs_on"
	RelationCalls     = "calls"
)

// CMDBAsset is a logical ops configuration item (CMDB), not the SSH host Asset model.
type CMDBAsset struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	Name          string         `gorm:"not null;size:256;index" json:"name"`
	Kind          string         `gorm:"not null;size:64;index" json:"kind"`
	Env           string         `gorm:"size:64;index" json:"env"`
	OwnerUserID   *uint          `gorm:"index" json:"owner_user_id,omitempty"`
	LabelsJSON    string         `gorm:"type:jsonb;default:'{}'" json:"labels,omitempty"`
	IntegrationID *uint          `gorm:"index" json:"integration_id,omitempty"`
	ExternalRef   string         `gorm:"size:512" json:"external_ref,omitempty"`
	Notes         string         `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`

	Owner       *User        `gorm:"foreignKey:OwnerUserID" json:"owner,omitempty"`
	Integration *Integration `gorm:"foreignKey:IntegrationID" json:"integration,omitempty"`
}

// TableName maps to cmdb_assets.
func (CMDBAsset) TableName() string {
	return "cmdb_assets"
}

// AssetRelation links two CMDB assets with a typed dependency or placement edge.
type AssetRelation struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	FromAssetID  uint      `gorm:"not null;index:idx_asset_rel_from,priority:1" json:"from_asset_id"`
	ToAssetID    uint      `gorm:"not null;index:idx_asset_rel_to,priority:1" json:"to_asset_id"`
	RelationType string    `gorm:"not null;size:32;index" json:"relation_type"`
	MetaJSON     string    `gorm:"type:jsonb;default:'{}'" json:"meta,omitempty"`
	CreatedAt    time.Time `json:"created_at"`

	FromAsset *CMDBAsset `gorm:"foreignKey:FromAssetID" json:"from_asset,omitempty"`
	ToAsset   *CMDBAsset `gorm:"foreignKey:ToAssetID" json:"to_asset,omitempty"`
}

// TableName maps to asset_relations.
func (AssetRelation) TableName() string {
	return "asset_relations"
}
