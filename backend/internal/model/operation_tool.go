// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import "time"

// OperationTool 运维工具
type OperationTool struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"not null;size:255" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Category    string    `gorm:"size:100" json:"category"`
	Icon        string    `gorm:"size:255" json:"icon"`
	URL         string    `gorm:"not null;size:512" json:"url"`
	OrderIndex  int       `gorm:"default:0" json:"order"`
	Enabled     bool      `gorm:"default:true" json:"enabled"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 指定表名
func (OperationTool) TableName() string {
	return "operation_tools"
}
