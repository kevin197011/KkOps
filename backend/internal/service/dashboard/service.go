// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package dashboard

import (
	"time"

	"github.com/kkops/backend/internal/model"
	"gorm.io/gorm"
)

// Service 仪表板服务
type Service struct {
	db *gorm.DB
}

// NewService 创建仪表板服务实例
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// StatsResponse 统计数据响应
type StatsResponse struct {
	// 基础统计
	TotalAssets   int64 `json:"total_assets"`
	TotalUsers    int64 `json:"total_users"`
	TotalTasks    int64 `json:"total_tasks"`
	TotalProjects int64 `json:"total_projects"`

	// 资产状态统计
	AssetsByStatus map[string]int64 `json:"assets_by_status"`

	// 任务执行统计
	TaskExecutionStats TaskExecutionStats `json:"task_execution_stats"`

	// 资产按环境分布
	AssetsByEnvironment []EnvironmentAssetCount `json:"assets_by_environment"`

	// 资产按项目分布
	AssetsByProject []ProjectAssetCount `json:"assets_by_project"`

	// 最近活动
	RecentActivities []RecentActivity `json:"recent_activities"`
}

// TaskExecutionStats 任务执行统计
type TaskExecutionStats struct {
	Total     int64 `json:"total"`
	Success   int64 `json:"success"`
	Failed    int64 `json:"failed"`
	Running   int64 `json:"running"`
	Pending   int64 `json:"pending"`
	Cancelled int64 `json:"cancelled"`
}

// EnvironmentAssetCount 环境资产计数
type EnvironmentAssetCount struct {
	EnvironmentID   uint   `json:"environment_id"`
	EnvironmentName string `json:"environment_name"`
	Count           int64  `json:"count"`
}

// ProjectAssetCount 项目资产计数
type ProjectAssetCount struct {
	ProjectID   uint   `json:"project_id"`
	ProjectName string `json:"project_name"`
	Count       int64  `json:"count"`
}

// RecentActivity 最近活动
type RecentActivity struct {
	ID          uint      `json:"id"`
	Type        string    `json:"type"` // "task_execution", "asset_created", "user_login"
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
}

// GetStats 获取仪表板统计数据
func (s *Service) GetStats() (*StatsResponse, error) {
	stats := &StatsResponse{
		AssetsByStatus: make(map[string]int64),
	}

	// 基础统计
	s.db.Model(&model.Asset{}).Count(&stats.TotalAssets)
	s.db.Model(&model.User{}).Count(&stats.TotalUsers)
	s.db.Model(&model.Task{}).Count(&stats.TotalTasks)
	s.db.Model(&model.Project{}).Count(&stats.TotalProjects)

	// 资产状态统计
	type StatusCount struct {
		Status string
		Count  int64
	}
	var statusCounts []StatusCount
	s.db.Model(&model.Asset{}).
		Select("status, count(*) as count").
		Group("status").
		Scan(&statusCounts)
	for _, sc := range statusCounts {
		if sc.Status == "" {
			stats.AssetsByStatus["unknown"] = sc.Count
		} else {
			stats.AssetsByStatus[sc.Status] = sc.Count
		}
	}

	// 任务执行统计
	s.db.Model(&model.TaskExecution{}).Count(&stats.TaskExecutionStats.Total)
	s.db.Model(&model.TaskExecution{}).Where("status = ?", "success").Count(&stats.TaskExecutionStats.Success)
	s.db.Model(&model.TaskExecution{}).Where("status = ?", "failed").Count(&stats.TaskExecutionStats.Failed)
	s.db.Model(&model.TaskExecution{}).Where("status = ?", "running").Count(&stats.TaskExecutionStats.Running)
	s.db.Model(&model.TaskExecution{}).Where("status = ?", "pending").Count(&stats.TaskExecutionStats.Pending)
	s.db.Model(&model.TaskExecution{}).Where("status = ?", "cancelled").Count(&stats.TaskExecutionStats.Cancelled)

	// 资产按环境分布
	type EnvCount struct {
		EnvironmentID uint
		Name          string
		Count         int64
	}
	var envCounts []EnvCount
	s.db.Model(&model.Asset{}).
		Select("assets.environment_id, environments.name, count(*) as count").
		Joins("LEFT JOIN environments ON environments.id = assets.environment_id").
		Where("assets.environment_id IS NOT NULL").
		Group("assets.environment_id, environments.name").
		Scan(&envCounts)
	for _, ec := range envCounts {
		stats.AssetsByEnvironment = append(stats.AssetsByEnvironment, EnvironmentAssetCount{
			EnvironmentID:   ec.EnvironmentID,
			EnvironmentName: ec.Name,
			Count:           ec.Count,
		})
	}

	// 资产按项目分布
	type ProjCount struct {
		ProjectID uint
		Name      string
		Count     int64
	}
	var projCounts []ProjCount
	s.db.Model(&model.Asset{}).
		Select("assets.project_id, projects.name, count(*) as count").
		Joins("LEFT JOIN projects ON projects.id = assets.project_id").
		Where("assets.project_id IS NOT NULL").
		Group("assets.project_id, projects.name").
		Scan(&projCounts)
	for _, pc := range projCounts {
		stats.AssetsByProject = append(stats.AssetsByProject, ProjectAssetCount{
			ProjectID:   pc.ProjectID,
			ProjectName: pc.Name,
			Count:       pc.Count,
		})
	}

	// 最近活动 (最近10条任务执行记录)
	var recentExecutions []model.TaskExecution
	s.db.Model(&model.TaskExecution{}).
		Preload("Task").
		Preload("Asset").
		Order("created_at DESC").
		Limit(10).
		Find(&recentExecutions)

	// 获取定时任务名称映射
	scheduledTaskNames := make(map[uint]string)
	var scheduledTasks []model.ScheduledTask
	s.db.Model(&model.ScheduledTask{}).Select("id, name").Find(&scheduledTasks)
	for _, st := range scheduledTasks {
		scheduledTaskNames[st.ID] = st.Name
	}

	for _, exec := range recentExecutions {
		taskName := "未知任务"
		// 优先检查普通任务
		if exec.Task != nil && exec.Task.ID != 0 {
			taskName = exec.Task.Name
		} else if exec.ScheduledTaskID != nil {
			// 检查定时任务
			if name, ok := scheduledTaskNames[*exec.ScheduledTaskID]; ok {
				taskName = name
			}
		}
		assetName := "未知主机"
		if exec.Asset.ID != 0 {
			assetName = exec.Asset.HostName
		}

		stats.RecentActivities = append(stats.RecentActivities, RecentActivity{
			ID:          exec.ID,
			Type:        "task_execution",
			Title:       taskName,
			Description: "执行于 " + assetName,
			Status:      exec.Status,
			CreatedAt:   exec.CreatedAt,
		})
	}

	return stats, nil
}
