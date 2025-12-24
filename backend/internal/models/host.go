package models

import (
	"time"

	"gorm.io/gorm"
)

// Host 主机模型
type Host struct {
	ID           uint64         `gorm:"primaryKey" json:"id"`
	ProjectID    uint64         `gorm:"not null;index" json:"project_id"`
	Hostname     string         `gorm:"size:255;not null;index" json:"hostname"`
	IPAddress    string         `gorm:"size:45;not null;index" json:"ip_address"`
	SaltMinionID string         `gorm:"uniqueIndex;size:255" json:"salt_minion_id"`
	OSType       string         `gorm:"size:50" json:"os_type"`
	OSVersion    string         `gorm:"size:100" json:"os_version"`
	CPUCores     *int           `json:"cpu_cores"`
	MemoryGB     *float64       `gorm:"type:decimal(10,2)" json:"memory_gb"`
	DiskGB       *float64       `gorm:"type:decimal(10,2)" json:"disk_gb"`
	Status       string         `gorm:"size:20;not null;default:'unknown';index" json:"status"`
	Environment  string         `gorm:"size:20;index" json:"environment"`
	SSHPort      int            `gorm:"not null;default:22" json:"ssh_port"`
	SSHKeyID     *uint64        `gorm:"index" json:"ssh_key_id"` // 可选的默认SSH密钥ID
	LastSeenAt   *time.Time     `json:"last_seen_at"`
	SaltVersion  string         `gorm:"size:50" json:"salt_version"`
	Metadata     string         `gorm:"type:jsonb" json:"metadata"`
	Description  string         `json:"description"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Project *Project    `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	SSHKey  *SSHKey     `gorm:"foreignKey:SSHKeyID" json:"ssh_key,omitempty"`
	Groups  []HostGroup `gorm:"many2many:host_group_members;" json:"groups,omitempty"`
	Tags    []HostTag   `gorm:"many2many:host_tag_assignments;" json:"tags,omitempty"`
}

// TableName 指定表名
func (Host) TableName() string {
	return "hosts"
}

// HostGroup 主机组模型
type HostGroup struct {
	ID          uint64         `gorm:"primaryKey" json:"id"`
	ProjectID   uint64         `gorm:"not null;index" json:"project_id"`
	Name        string         `gorm:"size:100;not null;index" json:"name"`
	Description string         `json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// 关联
	Project *Project `gorm:"foreignKey:ProjectID" json:"project,omitempty"`
	Hosts   []Host   `gorm:"many2many:host_group_members;" json:"hosts,omitempty"`
}

// TableName 指定表名
func (HostGroup) TableName() string {
	return "host_groups"
}

// HostTag 主机标签模型
type HostTag struct {
	ID          uint64    `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"uniqueIndex;size:50;not null" json:"name"`
	Color       string    `gorm:"size:7" json:"color"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`

	// 关联
	Hosts []Host `gorm:"many2many:host_tag_assignments;" json:"hosts,omitempty"`
}

// TableName 指定表名
func (HostTag) TableName() string {
	return "host_tags"
}

// HostGroupMember 主机组成员关联模型
type HostGroupMember struct {
	ID        uint64    `gorm:"primaryKey"`
	HostID    uint64    `gorm:"not null;index"`
	GroupID   uint64    `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (HostGroupMember) TableName() string {
	return "host_group_members"
}

// HostTagAssignment 主机标签关联模型
type HostTagAssignment struct {
	ID        uint64    `gorm:"primaryKey"`
	HostID    uint64    `gorm:"not null;index"`
	TagID     uint64    `gorm:"not null;index"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

// TableName 指定表名
func (HostTagAssignment) TableName() string {
	return "host_tag_assignments"
}
