package models

import (
	"time"

	"gorm.io/gorm"
)

// SSHKey SSH密钥模型
type SSHKey struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	UserID      uint64         `gorm:"not null;index" json:"user_id"`
	Name        string         `gorm:"size:100;not null" json:"name"`
	Username    string         `gorm:"size:100" json:"username"` // SSH用户名（可选）
	KeyType     string         `gorm:"size:20;not null" json:"key_type"` // rsa, ed25519, ecdsa, etc.
	PrivateKey  string         `gorm:"type:text;not null" json:"-"`     // 加密存储，不返回给前端
	PublicKey   string         `gorm:"type:text" json:"public_key"`
	Fingerprint string         `gorm:"size:100;index" json:"fingerprint"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	User *User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// TableName 指定表名
func (SSHKey) TableName() string {
	return "ssh_keys"
}

