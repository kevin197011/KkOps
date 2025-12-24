package models

import (
	"time"

	"gorm.io/gorm"
)

// DeploymentConfig 部署配置模型
type DeploymentConfig struct {
	ID             uint64         `gorm:"primaryKey" json:"id"`
	Name           string         `gorm:"size:100;not null;index" json:"name"`
	ApplicationName string        `gorm:"size:100;not null;index" json:"application_name"`
	Description    string         `json:"description"`
	SaltStateFiles string         `gorm:"type:text[]" json:"salt_state_files"` // PostgreSQL array
	TargetGroups   string         `gorm:"type:bigint[]" json:"target_groups"` // PostgreSQL array
	TargetHosts    string         `gorm:"type:bigint[]" json:"target_hosts"`   // PostgreSQL array
	Environment    string         `gorm:"size:50" json:"environment"`          // dev, test, prod
	ConfigData     string         `gorm:"type:jsonb" json:"config_data"`        // JSON格式
	CreatedBy      uint64         `gorm:"not null;index" json:"created_by"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Creator *User        `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
	Deployments []Deployment `gorm:"foreignKey:ConfigID" json:"deployments,omitempty"`
}

// TableName 指定表名
func (DeploymentConfig) TableName() string {
	return "deployment_configs"
}

// Deployment 部署执行记录模型
type Deployment struct {
	ID                      uint64         `gorm:"primaryKey" json:"id"`
	ConfigID                uint64         `gorm:"not null;index" json:"config_id"`
	Version                 string         `gorm:"size:50;not null" json:"version"`
	Status                  string         `gorm:"size:20;not null;default:'pending';index" json:"status"` // pending, running, completed, failed, cancelled
	StartedBy               uint64         `gorm:"not null;index" json:"started_by"`
	StartedAt               time.Time      `gorm:"not null;default:now()" json:"started_at"`
	CompletedAt             *time.Time     `json:"completed_at"`
	DurationSeconds         *int           `json:"duration_seconds"`
	TargetHosts             string         `gorm:"type:bigint[];not null" json:"target_hosts"` // PostgreSQL array
	SaltJobID               string         `gorm:"size:255;index" json:"salt_job_id"`
	Results                 string         `gorm:"type:jsonb" json:"results"` // JSON格式，按主机存储结果
	ErrorMessage            string         `gorm:"type:text" json:"error_message"`
	IsRollback              bool           `gorm:"not null;default:false" json:"is_rollback"`
	RollbackFromDeploymentID *uint64       `gorm:"index" json:"rollback_from_deployment_id"`
	CreatedAt               time.Time      `json:"created_at"`
	UpdatedAt               time.Time      `json:"updated_at"`

	// 关联
	Config                  *DeploymentConfig `gorm:"foreignKey:ConfigID" json:"config,omitempty"`
	Starter                 *User             `gorm:"foreignKey:StartedBy" json:"starter,omitempty"`
	RollbackFromDeployment  *Deployment       `gorm:"foreignKey:RollbackFromDeploymentID" json:"rollback_from_deployment,omitempty"`
}

// TableName 指定表名
func (Deployment) TableName() string {
	return "deployments"
}

// DeploymentVersion 部署版本模型
type DeploymentVersion struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	ApplicationName string   `gorm:"size:100;not null;index" json:"application_name"`
	Version        string    `gorm:"size:50;not null" json:"version"`
	ReleaseNotes   string    `json:"release_notes"`
	CreatedBy      uint64    `gorm:"not null;index" json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`

	// 关联
	Creator *User `gorm:"foreignKey:CreatedBy" json:"creator,omitempty"`
}

// TableName 指定表名
func (DeploymentVersion) TableName() string {
	return "deployment_versions"
}

