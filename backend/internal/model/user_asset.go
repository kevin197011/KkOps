// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"
)

// UserAsset 用户资产授权关联表
// 存储用户与资产的授权关系，用于细粒度的访问控制
type UserAsset struct {
	UserID    uint      `gorm:"primaryKey" json:"user_id"`
	AssetID   uint      `gorm:"primaryKey" json:"asset_id"`
	CreatedAt time.Time `json:"created_at"`
	CreatedBy uint      `json:"created_by"` // 授权操作人

	// 关联
	User  User  `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Asset Asset `gorm:"foreignKey:AssetID" json:"asset,omitempty"`
}

// TableName 指定表名
func (UserAsset) TableName() string {
	return "user_assets"
}
