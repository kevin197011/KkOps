package models

import (
	"time"

	"gorm.io/datatypes"
)

// Formula Salt Formula 模型
type Formula struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	Name        string         `json:"name" gorm:"uniqueIndex;not null"`        // Formula 名称（如 nginx, redis）
	Description string         `json:"description"`                             // Formula 描述
	Category    string         `json:"category" gorm:"index"`                  // 分类（base, middleware, runtime, app）
	Version     string         `json:"version"`                                 // 版本号
	Path        string         `json:"path" gorm:"not null"`                   // 在仓库中的路径
	Repository  string         `json:"repository" gorm:"not null"`             // 仓库URL或本地路径
	Icon        string         `json:"icon"`                                   // 图标
	Tags        datatypes.JSON  `json:"tags" gorm:"type:jsonb"`                // 标签数组
	Metadata    datatypes.JSON  `json:"metadata" gorm:"type:jsonb"`            // 元数据（pillar schema, dependencies等）
	IsActive    bool           `json:"is_active" gorm:"default:true"`          // 是否启用
	CreatedBy   uint64         `json:"created_by" gorm:"not null;default:1"`   // 创建者ID，默认值为1
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

// FormulaParameter Formula 参数定义
type FormulaParameter struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	FormulaID   uint           `json:"formula_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`                              // 参数名称
	Type        string         `json:"type" gorm:"not null"`                              // 参数类型 (string, number, boolean, array, object)
	Default     datatypes.JSON `json:"default" gorm:"column:default_value;type:jsonb"`    // 默认值
	Required    bool           `json:"required" gorm:"default:false"`                     // 是否必填
	Label       string         `json:"label"`                                             // 显示标签
	Description string         `json:"description"`                                       // 参数描述
	Validation  datatypes.JSON `json:"validation" gorm:"type:jsonb"`                      // 验证规则
	Order       int            `json:"order" gorm:"column:order_index;default:0"`         // 显示顺序
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// FormulaDeployment Formula 部署记录
type FormulaDeployment struct {
	ID             uint           `json:"id" gorm:"primarykey"`
	FormulaID      uint           `json:"formula_id" gorm:"not null;index"`
	Name           string         `json:"name" gorm:"not null"`                // 部署名称
	Description    string         `json:"description"`                         // 部署描述
	TargetHosts    datatypes.JSON `json:"target_hosts" gorm:"type:jsonb"`      // 目标主机列表
	PillarData     datatypes.JSON `json:"pillar_data" gorm:"type:jsonb"`       // Pillar 参数数据
	Status         string         `json:"status" gorm:"default:pending"`        // 状态 (pending, running, completed, failed, cancelled)
	SaltJobID      string         `json:"salt_job_id"`                         // Salt Job ID
	Results        datatypes.JSON `json:"results" gorm:"type:jsonb"`           // 执行结果
	SuccessCount   int            `json:"success_count" gorm:"default:0"`      // 成功主机数
	FailedCount    int            `json:"failed_count" gorm:"default:0"`       // 失败主机数
	ErrorMessage   string         `json:"error_message"`                       // 错误信息
	StartedBy      uint64         `json:"started_by" gorm:"not null;default:1"`  // 执行者ID
	StartedAt      *time.Time     `json:"started_at"`                          // 开始时间
	CompletedAt    *time.Time     `json:"completed_at"`                       // 完成时间
	Duration       int            `json:"duration_seconds"`                    // 执行时长(秒)
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

// FormulaTemplate Formula 模板（预设配置）
type FormulaTemplate struct {
	ID          uint           `json:"id" gorm:"primarykey"`
	FormulaID   uint           `json:"formula_id" gorm:"not null;index"`
	Name        string         `json:"name" gorm:"not null"`                   // 模板名称
	Description string         `json:"description"`                            // 模板描述
	PillarData  datatypes.JSON `json:"pillar_data" gorm:"type:jsonb"`         // 预设的Pillar数据
	IsPublic    bool           `json:"is_public" gorm:"default:false"`         // 是否公开
	CreatedBy   uint64         `json:"created_by" gorm:"not null;default:1"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

// FormulaRepository Formula 仓库配置
type FormulaRepository struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"uniqueIndex;not null"`           // 仓库名称
	URL         string    `json:"url" gorm:"not null"`                        // Git仓库URL
	Branch      string    `json:"branch" gorm:"default:master"`               // 分支
	LocalPath   string    `json:"local_path"`                                 // 本地路径（如果不是Git仓库）
	IsActive    bool      `json:"is_active" gorm:"default:true"`              // 是否启用
	LastSyncAt  *time.Time `json:"last_sync_at"`                              // 最后同步时间
	CreatedBy   uint64    `json:"created_by" gorm:"not null;default:1"`       // 创建者ID，默认值为1
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// TableName 指定表名
func (Formula) TableName() string {
	return "formulas"
}

func (FormulaParameter) TableName() string {
	return "formula_parameters"
}

func (FormulaDeployment) TableName() string {
	return "formula_deployments"
}

func (FormulaTemplate) TableName() string {
	return "formula_templates"
}

func (FormulaRepository) TableName() string {
	return "formula_repositories"
}
