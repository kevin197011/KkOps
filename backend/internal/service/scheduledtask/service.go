// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package scheduledtask

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/kkops/backend/internal/model"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

// Service 定时任务服务
type Service struct {
	db        *gorm.DB
	scheduler *Scheduler
}

// NewService 创建定时任务服务
func NewService(db *gorm.DB) *Service {
	return &Service{db: db}
}

// SetScheduler 设置调度器引用
func (s *Service) SetScheduler(scheduler *Scheduler) {
	s.scheduler = scheduler
}

// CreateScheduledTaskRequest 创建定时任务请求
type CreateScheduledTaskRequest struct {
	Name           string `json:"name" binding:"required"`
	Description    string `json:"description"`
	CronExpression string `json:"cron_expression" binding:"required"`
	TemplateID     *uint  `json:"template_id"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	AssetIDs       []uint `json:"asset_ids"`
	Timeout        int    `json:"timeout"`
	Enabled        bool   `json:"enabled"`
	UpdateAssets   bool   `json:"update_assets"` // 是否更新资产信息
}

// UpdateScheduledTaskRequest 更新定时任务请求
type UpdateScheduledTaskRequest struct {
	Name           string `json:"name"`
	Description    string `json:"description"`
	CronExpression string `json:"cron_expression"`
	TemplateID     *uint  `json:"template_id"`
	Content        string `json:"content"`
	Type           string `json:"type"`
	AssetIDs       []uint `json:"asset_ids"`
	Timeout        int    `json:"timeout"`
	Enabled        *bool  `json:"enabled"`
	UpdateAssets   *bool  `json:"update_assets"` // 是否更新资产信息
}

// ScheduledTaskResponse 定时任务响应
type ScheduledTaskResponse struct {
	ID             uint       `json:"id"`
	Name           string     `json:"name"`
	Description    string     `json:"description"`
	CronExpression string     `json:"cron_expression"`
	TemplateID     *uint      `json:"template_id,omitempty"`
	TemplateName   string     `json:"template_name,omitempty"`
	Content        string     `json:"content"`
	Type           string     `json:"type"`
	AssetIDs       []uint     `json:"asset_ids"`
	Timeout        int        `json:"timeout"`
	Enabled        bool       `json:"enabled"`
	UpdateAssets   bool       `json:"update_assets"` // 是否更新资产信息
	LastRunAt      *time.Time `json:"last_run_at,omitempty"`
	NextRunAt      *time.Time `json:"next_run_at,omitempty"`
	LastStatus     string     `json:"last_status,omitempty"`
	CreatedBy      uint       `json:"created_by"`
	CreatorName    string     `json:"creator_name,omitempty"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ListScheduledTasksResponse 定时任务列表响应
type ListScheduledTasksResponse struct {
	Data  []ScheduledTaskResponse `json:"data"`
	Total int64                   `json:"total"`
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
}

// ValidateCronExpression 验证 Cron 表达式（必须是 6 字段格式：秒 分 时 日 月 周）
func ValidateCronExpression(expr string) error {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(expr)
	if err != nil {
		return fmt.Errorf("无效的 Cron 表达式，请使用 6 字段格式（秒 分 时 日 月 周），例如: 0 * * * * * 表示每分钟执行")
	}
	return nil
}

// GetNextRunTime 获取下次执行时间（使用 6 字段格式）
func GetNextRunTime(cronExpr string) (*time.Time, error) {
	parser := cron.NewParser(cron.Second | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		return nil, err
	}
	nextTime := schedule.Next(time.Now())
	return &nextTime, nil
}

// CreateScheduledTask 创建定时任务
func (s *Service) CreateScheduledTask(userID uint, req CreateScheduledTaskRequest) (*ScheduledTaskResponse, error) {
	// 验证 Cron 表达式
	if err := ValidateCronExpression(req.CronExpression); err != nil {
		return nil, err
	}

	// 如果指定了模板，获取模板内容
	if req.TemplateID != nil && *req.TemplateID > 0 {
		var template model.TaskTemplate
		if err := s.db.First(&template, *req.TemplateID).Error; err != nil {
			return nil, fmt.Errorf("模板不存在")
		}
		if req.Content == "" {
			req.Content = template.Content
		}
		if req.Type == "" {
			req.Type = template.Type
		}
	}

	// 设置默认值
	if req.Type == "" {
		req.Type = "shell"
	}
	if req.Timeout == 0 {
		req.Timeout = 300
	}

	// 计算下次执行时间
	var nextRunAt *time.Time
	if req.Enabled {
		nextRunAt, _ = GetNextRunTime(req.CronExpression)
	}

	// 将 AssetIDs 转换为字符串
	assetIDStrs := make([]string, len(req.AssetIDs))
	for i, id := range req.AssetIDs {
		assetIDStrs[i] = strconv.FormatUint(uint64(id), 10)
	}

	task := &model.ScheduledTask{
		Name:           req.Name,
		Description:    req.Description,
		CronExpression: req.CronExpression,
		TemplateID:     req.TemplateID,
		Content:        req.Content,
		Type:           req.Type,
		AssetIDs:       strings.Join(assetIDStrs, ","),
		Timeout:        req.Timeout,
		Enabled:        req.Enabled,
		UpdateAssets:   req.UpdateAssets,
		NextRunAt:      nextRunAt,
		CreatedBy:      userID,
	}

	if err := s.db.Create(task).Error; err != nil {
		return nil, fmt.Errorf("创建定时任务失败: %w", err)
	}

	// 如果任务启用且调度器存在，添加到调度器
	if task.Enabled && s.scheduler != nil {
		if err := s.scheduler.AddTask(task); err != nil {
			// 记录错误但不影响创建
			fmt.Printf("添加任务到调度器失败: %v\n", err)
		}
	}

	return s.GetScheduledTask(task.ID)
}

// GetScheduledTask 获取定时任务
func (s *Service) GetScheduledTask(id uint) (*ScheduledTaskResponse, error) {
	var task model.ScheduledTask
	if err := s.db.Preload("Template").Preload("Creator").First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("定时任务不存在")
		}
		return nil, err
	}

	return s.taskToResponse(&task), nil
}

// ListScheduledTasks 获取定时任务列表
func (s *Service) ListScheduledTasks(page, pageSize int, enabled *bool) (*ListScheduledTasksResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var tasks []model.ScheduledTask
	var total int64

	query := s.db.Model(&model.ScheduledTask{})
	if enabled != nil {
		query = query.Where("enabled = ?", *enabled)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Template").Preload("Creator").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&tasks).Error; err != nil {
		return nil, err
	}

	data := make([]ScheduledTaskResponse, len(tasks))
	for i, task := range tasks {
		data[i] = *s.taskToResponse(&task)
	}

	return &ListScheduledTasksResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Size:  pageSize,
	}, nil
}

// UpdateScheduledTask 更新定时任务
func (s *Service) UpdateScheduledTask(id uint, req UpdateScheduledTaskRequest) (*ScheduledTaskResponse, error) {
	var task model.ScheduledTask
	if err := s.db.First(&task, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("定时任务不存在")
		}
		return nil, err
	}

	// 验证 Cron 表达式
	if req.CronExpression != "" {
		if err := ValidateCronExpression(req.CronExpression); err != nil {
			return nil, err
		}
		task.CronExpression = req.CronExpression
	}

	// 更新字段
	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.TemplateID != nil {
		task.TemplateID = req.TemplateID
	}
	if req.Content != "" {
		task.Content = req.Content
	}
	if req.Type != "" {
		task.Type = req.Type
	}
	if len(req.AssetIDs) > 0 {
		assetIDStrs := make([]string, len(req.AssetIDs))
		for i, id := range req.AssetIDs {
			assetIDStrs[i] = strconv.FormatUint(uint64(id), 10)
		}
		task.AssetIDs = strings.Join(assetIDStrs, ",")
	}
	if req.Timeout > 0 {
		task.Timeout = req.Timeout
	}
	if req.Enabled != nil {
		task.Enabled = *req.Enabled
		// 更新下次执行时间
		if *req.Enabled {
			task.NextRunAt, _ = GetNextRunTime(task.CronExpression)
		} else {
			task.NextRunAt = nil
		}
	}
	if req.UpdateAssets != nil {
		task.UpdateAssets = *req.UpdateAssets
	}

	if err := s.db.Save(&task).Error; err != nil {
		return nil, fmt.Errorf("更新定时任务失败: %w", err)
	}

	// 更新调度器中的任务
	if s.scheduler != nil {
		if err := s.scheduler.UpdateTask(&task); err != nil {
			fmt.Printf("更新调度器中的任务失败: %v\n", err)
		}
	}

	return s.GetScheduledTask(task.ID)
}

// DeleteScheduledTask 删除定时任务
func (s *Service) DeleteScheduledTask(id uint) error {
	// 先从调度器移除
	if s.scheduler != nil {
		s.scheduler.RemoveTask(id)
	}

	result := s.db.Delete(&model.ScheduledTask{}, id)
	if result.Error != nil {
		return fmt.Errorf("删除定时任务失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("定时任务不存在")
	}
	return nil
}

// EnableScheduledTask 启用定时任务
func (s *Service) EnableScheduledTask(id uint) error {
	enabled := true
	_, err := s.UpdateScheduledTask(id, UpdateScheduledTaskRequest{Enabled: &enabled})
	return err
}

// DisableScheduledTask 禁用定时任务
func (s *Service) DisableScheduledTask(id uint) error {
	enabled := false
	_, err := s.UpdateScheduledTask(id, UpdateScheduledTaskRequest{Enabled: &enabled})
	return err
}

// GetEnabledScheduledTasks 获取所有启用的定时任务
func (s *Service) GetEnabledScheduledTasks() ([]model.ScheduledTask, error) {
	var tasks []model.ScheduledTask
	if err := s.db.Where("enabled = ?", true).Find(&tasks).Error; err != nil {
		return nil, err
	}
	return tasks, nil
}

// UpdateTaskLastRun 更新任务最后执行时间和状态
func (s *Service) UpdateTaskLastRun(id uint, status string) error {
	now := time.Now()
	updates := map[string]interface{}{
		"last_run_at": now,
		"last_status": status,
	}

	// 计算下次执行时间
	var task model.ScheduledTask
	if err := s.db.First(&task, id).Error; err == nil && task.Enabled {
		if nextRunAt, err := GetNextRunTime(task.CronExpression); err == nil {
			updates["next_run_at"] = nextRunAt
		}
	}

	return s.db.Model(&model.ScheduledTask{}).Where("id = ?", id).Updates(updates).Error
}

// GetScheduledTaskExecutions 获取定时任务的执行历史
func (s *Service) GetScheduledTaskExecutions(taskID uint, page, pageSize int) ([]model.TaskExecution, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	var executions []model.TaskExecution
	var total int64

	query := s.db.Model(&model.TaskExecution{}).Where("scheduled_task_id = ?", taskID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Preload("Asset").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&executions).Error; err != nil {
		return nil, 0, err
	}

	return executions, total, nil
}

// taskToResponse 将模型转换为响应
func (s *Service) taskToResponse(task *model.ScheduledTask) *ScheduledTaskResponse {
	resp := &ScheduledTaskResponse{
		ID:             task.ID,
		Name:           task.Name,
		Description:    task.Description,
		CronExpression: task.CronExpression,
		TemplateID:     task.TemplateID,
		Content:        task.Content,
		Type:           task.Type,
		Timeout:        task.Timeout,
		Enabled:        task.Enabled,
		UpdateAssets:   task.UpdateAssets,
		LastRunAt:      task.LastRunAt,
		NextRunAt:      task.NextRunAt,
		LastStatus:     task.LastStatus,
		CreatedBy:      task.CreatedBy,
		CreatedAt:      task.CreatedAt,
		UpdatedAt:      task.UpdatedAt,
	}

	// 解析 AssetIDs
	if task.AssetIDs != "" {
		ids := strings.Split(task.AssetIDs, ",")
		resp.AssetIDs = make([]uint, 0, len(ids))
		for _, idStr := range ids {
			if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
				resp.AssetIDs = append(resp.AssetIDs, uint(id))
			}
		}
	}

	// 获取模板名称
	if task.Template != nil {
		resp.TemplateName = task.Template.Name
	}

	// 获取创建者名称
	if task.Creator != nil {
		resp.CreatorName = task.Creator.Username
	}

	return resp
}

// ExportScheduledTaskConfig 导出定时任务配置结构
type ExportScheduledTaskConfig struct {
	Name           string   `json:"name"`
	Description    string   `json:"description"`
	CronExpression string   `json:"cron_expression"`
	TemplateName   string   `json:"template_name,omitempty"` // 模板名称（如果有）
	Content        string   `json:"content"`                  // 脚本内容（如果没有模板或自定义）
	Type           string   `json:"type"`
	Timeout        int      `json:"timeout"`
	Enabled        bool     `json:"enabled"`
	UpdateAssets   bool     `json:"update_assets"`
	TargetHosts    []string `json:"target_hosts"` // 主机名或 IP 列表
}

// ExportScheduledTasksConfig 导出定时任务配置根结构
type ExportScheduledTasksConfig struct {
	Version  string                    `json:"version"`
	ExportAt string                    `json:"export_at"`
	Tasks    []ExportScheduledTaskConfig `json:"tasks"`
}

// ImportScheduledTaskConfig 导入定时任务配置结构
type ImportScheduledTaskConfig struct {
	Name           string   `json:"name" binding:"required"`
	Description    string   `json:"description"`
	CronExpression string   `json:"cron_expression" binding:"required"`
	TemplateName   string   `json:"template_name"`   // 模板名称（可选）
	Content        string   `json:"content"`         // 脚本内容（如果没有模板）
	Type           string   `json:"type"`
	Timeout        int      `json:"timeout"`
	Enabled        bool     `json:"enabled"`
	UpdateAssets   bool     `json:"update_assets"`
	TargetHosts    []string `json:"target_hosts"` // 主机名或 IP 列表
}

// ImportScheduledTasksConfig 导入定时任务配置根结构
type ImportScheduledTasksConfig struct {
	Version string                    `json:"version"`
	Tasks   []ImportScheduledTaskConfig `json:"tasks" binding:"required"`
}

// ImportScheduledTasksResult 导入定时任务结果
type ImportScheduledTasksResult struct {
	Total     int      `json:"total"`
	Success   int      `json:"success"`
	Failed    int      `json:"failed"`
	Errors    []string `json:"errors,omitempty"`
	Skipped   []string `json:"skipped,omitempty"`
}

// ExportScheduledTasks 导出所有定时任务
func (s *Service) ExportScheduledTasks() (*ExportScheduledTasksConfig, error) {
	var tasks []model.ScheduledTask
	if err := s.db.Preload("Template").Find(&tasks).Error; err != nil {
		return nil, err
	}

	exportTasks := make([]ExportScheduledTaskConfig, len(tasks))
	for i, t := range tasks {
		// 解析 AssetIDs 并转换为主机名/IP 列表
		assetIDs := []uint{}
		if t.AssetIDs != "" {
			ids := strings.Split(t.AssetIDs, ",")
			for _, idStr := range ids {
				if id, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 32); err == nil {
					assetIDs = append(assetIDs, uint(id))
				}
			}
		}
		targetHosts := s.getAssetHostnames(assetIDs)

		// 获取模板名称
		templateName := ""
		if t.Template != nil {
			templateName = t.Template.Name
		}

		exportTasks[i] = ExportScheduledTaskConfig{
			Name:           t.Name,
			Description:    t.Description,
			CronExpression: t.CronExpression,
			TemplateName:   templateName,
			Content:        t.Content,
			Type:           t.Type,
			Timeout:        t.Timeout,
			Enabled:        t.Enabled,
			UpdateAssets:   t.UpdateAssets,
			TargetHosts:    targetHosts,
		}
	}

	return &ExportScheduledTasksConfig{
		Version:  "1.0",
		ExportAt: time.Now().Format(time.RFC3339),
		Tasks:    exportTasks,
	}, nil
}

// getAssetHostnames 根据资产ID列表获取主机名或IP列表
func (s *Service) getAssetHostnames(assetIDs []uint) []string {
	if len(assetIDs) == 0 {
		return []string{}
	}

	var assets []model.Asset
	s.db.Where("id IN ?", assetIDs).Find(&assets)

	hostnames := make([]string, 0, len(assets))
	for _, asset := range assets {
		// 优先使用主机名，没有则使用 IP
		if asset.HostName != "" {
			hostnames = append(hostnames, asset.HostName)
		} else if asset.IP != "" {
			hostnames = append(hostnames, asset.IP)
		}
	}
	return hostnames
}

// findAssetsByHostnames 根据主机名或IP列表查找资产ID
func (s *Service) findAssetsByHostnames(hostnames []string) ([]uint, []string) {
	if len(hostnames) == 0 {
		return []uint{}, []string{}
	}

	var assets []model.Asset
	// 同时匹配主机名和 IP
	s.db.Where("host_name IN ? OR ip IN ?", hostnames, hostnames).Find(&assets)

	// 创建已找到的主机名/IP 映射
	foundSet := make(map[string]bool)
	assetIDs := make([]uint, 0, len(assets))
	for _, asset := range assets {
		assetIDs = append(assetIDs, asset.ID)
		if asset.HostName != "" {
			foundSet[asset.HostName] = true
		}
		if asset.IP != "" {
			foundSet[asset.IP] = true
		}
	}

	// 找出未找到的主机名
	missingHosts := make([]string, 0)
	for _, h := range hostnames {
		if !foundSet[h] {
			missingHosts = append(missingHosts, h)
		}
	}

	return assetIDs, missingHosts
}

// findTemplateByName 根据模板名称查找模板ID
func (s *Service) findTemplateByName(name string) (uint, error) {
	var template model.TaskTemplate
	if err := s.db.Where("name = ?", name).First(&template).Error; err != nil {
		return 0, fmt.Errorf("模板不存在: %s", name)
	}
	return template.ID, nil
}

// taskExistsByName 检查同名任务是否存在
func (s *Service) taskExistsByName(name string) bool {
	var task model.ScheduledTask
	if err := s.db.Where("name = ?", name).First(&task).Error; err != nil {
		return false
	}
	return true
}

// ImportScheduledTasks 导入定时任务
func (s *Service) ImportScheduledTasks(config *ImportScheduledTasksConfig, userID uint) (*ImportScheduledTasksResult, error) {
	result := &ImportScheduledTasksResult{
		Total:   len(config.Tasks),
		Errors:  []string{},
		Skipped: []string{},
	}

	for _, t := range config.Tasks {
		// 验证必填字段
		if t.Name == "" {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("任务缺少必填字段 'name'"))
			continue
		}
		if t.CronExpression == "" {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("任务 '%s': 缺少必填字段 'cron_expression'", t.Name))
			continue
		}

		// 验证 Cron 表达式
		if err := ValidateCronExpression(t.CronExpression); err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("任务 '%s': %v", t.Name, err))
			continue
		}

		// 检查是否已存在同名任务
		if s.taskExistsByName(t.Name) {
			result.Skipped = append(result.Skipped, fmt.Sprintf("任务 '%s': 已存在，已跳过", t.Name))
			continue
		}

		// 查找模板（如果提供了模板名称）
		var templateID *uint
		if t.TemplateName != "" {
			tplID, err := s.findTemplateByName(t.TemplateName)
			if err != nil {
				// 模板不存在，检查是否有自定义内容
				if t.Content == "" {
					result.Failed++
					result.Errors = append(result.Errors, fmt.Sprintf("任务 '%s': 模板 '%s' 不存在且未提供自定义内容", t.Name, t.TemplateName))
					continue
				}
				// 有自定义内容，继续但给出警告
				result.Skipped = append(result.Skipped, fmt.Sprintf("任务 '%s': 模板 '%s' 不存在，已跳过模板关联", t.Name, t.TemplateName))
			} else {
				templateID = &tplID
			}
		}

		// 匹配主机
		var assetIDs []uint
		var missingHosts []string
		if len(t.TargetHosts) > 0 {
			assetIDs, missingHosts = s.findAssetsByHostnames(t.TargetHosts)
			if len(missingHosts) > 0 {
				result.Skipped = append(result.Skipped, fmt.Sprintf("任务 '%s': 主机 %v 不存在，已跳过", t.Name, missingHosts))
			}
		}

		// 如果没有模板且没有自定义内容，失败
		if templateID == nil && t.Content == "" {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("任务 '%s': 缺少脚本内容", t.Name))
			continue
		}

		// 设置默认值
		taskType := t.Type
		if taskType == "" {
			taskType = "shell"
		}
		timeout := t.Timeout
		if timeout == 0 {
			timeout = 300
		}

		// 如果有关联模板，尝试从模板获取内容
		content := t.Content
		taskTypeFromTemplate := taskType
		if templateID != nil {
			var template model.TaskTemplate
			if err := s.db.First(&template, *templateID).Error; err == nil {
				if content == "" {
					content = template.Content
				}
				if taskType == "shell" {
					taskTypeFromTemplate = template.Type
				}
			}
		}

		// 计算下次执行时间
		var nextRunAt *time.Time
		if t.Enabled {
			nextRunAt, _ = GetNextRunTime(t.CronExpression)
		}

		// 将 AssetIDs 转换为字符串
		assetIDStrs := make([]string, len(assetIDs))
		for i, id := range assetIDs {
			assetIDStrs[i] = strconv.FormatUint(uint64(id), 10)
		}

		// 创建定时任务
		task := &model.ScheduledTask{
			Name:           t.Name,
			Description:    t.Description,
			CronExpression: t.CronExpression,
			TemplateID:     templateID,
			Content:        content,
			Type:           taskTypeFromTemplate,
			AssetIDs:       strings.Join(assetIDStrs, ","),
			Timeout:        timeout,
			Enabled:        t.Enabled,
			UpdateAssets:   t.UpdateAssets,
			NextRunAt:      nextRunAt,
			CreatedBy:      userID,
		}

		if err := s.db.Create(task).Error; err != nil {
			result.Failed++
			result.Errors = append(result.Errors, fmt.Sprintf("任务 '%s': 创建失败: %v", t.Name, err))
			continue
		}

		// 如果任务启用且调度器存在，添加到调度器
		if task.Enabled && s.scheduler != nil {
			if err := s.scheduler.AddTask(task); err != nil {
				// 记录错误但不影响导入
				fmt.Printf("添加任务到调度器失败: %v\n", err)
			}
		}

		result.Success++
	}

	return result, nil
}
