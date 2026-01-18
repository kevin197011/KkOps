// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"
)

// RoleAsset represents the many-to-many relationship between roles and assets for authorization
// 角色资产授权关系表，允许为角色批量授权资产
type RoleAsset struct {
	RoleID    uint      `gorm:"primaryKey" json:"role_id"`
	AssetID   uint      `gorm:"primaryKey" json:"asset_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint      `json:"created_by"` // The user who granted this authorization
}

// TableName specifies the table name for GORM
func (RoleAsset) TableName() string {
	return "role_assets"
}
