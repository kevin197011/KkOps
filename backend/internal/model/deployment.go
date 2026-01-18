// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package model

import (
	"time"

	"gorm.io/gorm"
)

// DeploymentModule represents a deployment module configuration
type DeploymentModule struct {
	ID               uint           `gorm:"primaryKey" json:"id"`
	ProjectID        uint           `gorm:"not null;index" json:"project_id"`
	Project          *Project       `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	EnvironmentID    *uint          `gorm:"index" json:"environment_id"`
	Environment      *Environment   `gorm:"foreignKey:EnvironmentID" json:"environment,omitempty"`
	TemplateID       *uint          `gorm:"index" json:"template_id"`                             // 关联执行模板
	Template         *TaskTemplate  `gorm:"foreignKey:TemplateID" json:"template,omitempty"`      // 执行模板
	Name             string         `gorm:"not null;size:100" json:"name"`
	Description      string         `gorm:"type:text" json:"description"`
	VersionSourceURL string         `gorm:"size:500;column:version_source_url" json:"version_source_url"`
	DeployScript     string         `gorm:"type:text" json:"deploy_script"`
	ScriptType       string         `gorm:"size:50;default:shell" json:"script_type"` // shell/python
	Timeout          int            `gorm:"default:600" json:"timeout"`               // Timeout in seconds
	AssetIDs         string         `gorm:"type:text" json:"asset_ids"`               // Comma-separated asset IDs
	CreatedBy        uint           `json:"created_by"`
	Creator          User           `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Deployments []Deployment `gorm:"foreignKey:ModuleID" json:"deployments,omitempty"`
}

// Deployment represents a deployment execution record
type Deployment struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	ModuleID   uint              `gorm:"not null;index" json:"module_id"`
	Module     *DeploymentModule `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
	Version    string            `gorm:"size:100" json:"version"`
	Status     string            `gorm:"default:pending;size:20;index" json:"status"` // pending/running/success/failed/cancelled
	AssetIDs   string            `gorm:"type:text" json:"asset_ids"`                  // Comma-separated asset IDs for this deployment
	Output     string            `gorm:"type:text" json:"output"`
	Error      string            `gorm:"type:text" json:"error"`
	CreatedBy  uint              `json:"created_by"`
	Creator    User              `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	StartedAt  *time.Time        `json:"started_at"`
	FinishedAt *time.Time        `json:"finished_at"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `gorm:"index" json:"-"`
}
