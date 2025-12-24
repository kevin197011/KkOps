package models

import (
	"time"

	"gorm.io/gorm"
)

// SSHConnection SSH连接配置模型
type SSHConnection struct {
	ID                uint64         `gorm:"primaryKey" json:"id"`
	ProjectID         uint64         `gorm:"not null;index" json:"project_id"`
	Name              string         `gorm:"size:100;not null" json:"name"`
	HostID            *uint64         `json:"host_id"`
	Hostname          string         `gorm:"size:255;not null" json:"hostname"`
	Port              int            `gorm:"not null;default:22" json:"port"`
	Username          string         `gorm:"size:100;not null" json:"username"`
	AuthType          string         `gorm:"size:20;not null" json:"auth_type"`
	PasswordEncrypted string         `gorm:"type:text" json:"-"`
	KeyID             *uint64        `json:"key_id"`
	Status            string         `gorm:"size:20;not null;default:'active';index" json:"status"`
	LastConnectedAt   *time.Time     `json:"last_connected_at"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Host    *Host    `gorm:"foreignKey:HostID" json:"host,omitempty"`
	Key     *SSHKey  `gorm:"foreignKey:KeyID" json:"key,omitempty"`
}

// TableName 指定表名
func (SSHConnection) TableName() string {
	return "ssh_connections"
}

// SSHKey SSH密钥模型
type SSHKey struct {
	ID                  uint64         `gorm:"primaryKey" json:"id"`
	Name                string         `gorm:"size:100;not null" json:"name"`
	KeyType             string         `gorm:"size:20;not null" json:"key_type"`
	PrivateKeyEncrypted string         `gorm:"type:text;not null" json:"-"`
	PublicKey           string         `gorm:"type:text;not null" json:"public_key"`
	Fingerprint         string         `gorm:"uniqueIndex;size:100" json:"fingerprint"`
	PassphraseEncrypted string         `gorm:"type:text" json:"-"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt          time.Time      `json:"updated_at"`
	DeletedAt          gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 指定表名
func (SSHKey) TableName() string {
	return "ssh_keys"
}

// SSHSession SSH会话模型
type SSHSession struct {
	ID              uint64     `gorm:"primaryKey" json:"id"`
	ConnectionID    uint64     `gorm:"not null;index" json:"connection_id"`
	UserID          uint64     `gorm:"not null;index" json:"user_id"`
	SessionToken    string     `gorm:"uniqueIndex;size:255;not null" json:"session_token"`
	ClientIP        string     `gorm:"size:45" json:"client_ip"`
	StartedAt       time.Time  `gorm:"not null;default:now();index" json:"started_at"`
	EndedAt         *time.Time `json:"ended_at"`
	DurationSeconds *int       `json:"duration_seconds"`
	Status          string     `gorm:"size:20;not null;default:'active';index" json:"status"`
	CreatedAt       time.Time  `json:"created_at"`

	// 关联
	Connection *SSHConnection `gorm:"foreignKey:ConnectionID" json:"connection,omitempty"`
	User       *User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (SSHSession) TableName() string {
	return "ssh_sessions"
}

