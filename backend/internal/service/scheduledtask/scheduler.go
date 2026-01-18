// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package scheduledtask

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/model"
	sshUtils "github.com/kkops/backend/internal/utils"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// AssetInfo 资产信息结构（用于解析任务输出）
type AssetInfo struct {
	Hostname string `json:"hostname"`
	CPU      string `json:"cpu"`
	Memory   string `json:"memory"`
	Disk     string `json:"disk"`
}

// Scheduler Cron 调度器
type Scheduler struct {
	cron       *cron.Cron
	db         *gorm.DB
	cfg        *config.Config
	logger     *zap.Logger
	entries    map[uint]cron.EntryID // taskID -> cronEntryID
	entriesMux sync.RWMutex
	service    *Service
}

// NewScheduler 创建调度器
func NewScheduler(db *gorm.DB, cfg *config.Config, logger *zap.Logger) *Scheduler {
	return &Scheduler{
		cron:    cron.New(cron.WithSeconds(), cron.WithChain(cron.Recover(cron.DefaultLogger))),
		db:      db,
		cfg:     cfg,
		logger:  logger,
		entries: make(map[uint]cron.EntryID),
		service: NewService(db),
	}
}

// Start 启动调度器
func (s *Scheduler) Start() error {
	// 加载所有启用的定时任务
	tasks, err := s.service.GetEnabledScheduledTasks()
	if err != nil {
		return fmt.Errorf("加载定时任务失败: %w", err)
	}

	s.logger.Info("正在加载定时任务", zap.Int("count", len(tasks)))

	for _, task := range tasks {
		if err := s.AddTask(&task); err != nil {
			s.logger.Error("添加定时任务失败",
				zap.Uint("task_id", task.ID),
				zap.String("name", task.Name),
				zap.Error(err))
		} else {
			s.logger.Info("已添加定时任务",
				zap.Uint("task_id", task.ID),
				zap.String("name", task.Name),
				zap.String("cron", task.CronExpression))
		}
	}

	// 启动 Cron 调度器
	s.cron.Start()
	s.logger.Info("Cron 调度器已启动")

	return nil
}

// Stop 停止调度器
func (s *Scheduler) Stop() {
	ctx := s.cron.Stop()
	<-ctx.Done()
	s.logger.Info("Cron 调度器已停止")
}

// AddTask 添加任务到调度器
func (s *Scheduler) AddTask(task *model.ScheduledTask) error {
	s.entriesMux.Lock()
	defer s.entriesMux.Unlock()

	// 如果任务已存在，先移除
	if entryID, exists := s.entries[task.ID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, task.ID)
	}

	// 添加新任务
	taskID := task.ID // 捕获 task ID 以避免闭包问题
	entryID, err := s.cron.AddFunc(task.CronExpression, func() {
		s.executeTask(taskID)
	})
	if err != nil {
		return fmt.Errorf("添加 Cron 任务失败: %w", err)
	}

	s.entries[task.ID] = entryID
	return nil
}

// RemoveTask 从调度器移除任务
func (s *Scheduler) RemoveTask(taskID uint) {
	s.entriesMux.Lock()
	defer s.entriesMux.Unlock()

	if entryID, exists := s.entries[taskID]; exists {
		s.cron.Remove(entryID)
		delete(s.entries, taskID)
		s.logger.Info("已移除定时任务", zap.Uint("task_id", taskID))
	}
}

// UpdateTask 更新任务调度
func (s *Scheduler) UpdateTask(task *model.ScheduledTask) error {
	if task.Enabled {
		return s.AddTask(task)
	} else {
		s.RemoveTask(task.ID)
		return nil
	}
}

// executeTask 执行定时任务
func (s *Scheduler) executeTask(taskID uint) {
	s.logger.Info("开始执行定时任务", zap.Uint("task_id", taskID))

	// 获取任务详情
	var task model.ScheduledTask
	if err := s.db.First(&task, taskID).Error; err != nil {
		s.logger.Error("获取定时任务失败", zap.Uint("task_id", taskID), zap.Error(err))
		return
	}

	// 检查任务是否仍然启用
	if !task.Enabled {
		s.logger.Warn("定时任务已禁用，跳过执行", zap.Uint("task_id", taskID))
		return
	}

	// 解析目标主机
	assetIDs := s.parseAssetIDs(task.AssetIDs)
	if len(assetIDs) == 0 {
		s.logger.Warn("定时任务没有目标主机，跳过执行", zap.Uint("task_id", taskID))
		s.service.UpdateTaskLastRun(taskID, "skipped")
		return
	}

	// 获取主机信息
	var assets []model.Asset
	if err := s.db.Preload("SSHKey").Where("id IN ?", assetIDs).Find(&assets).Error; err != nil {
		s.logger.Error("获取主机信息失败", zap.Uint("task_id", taskID), zap.Error(err))
		s.service.UpdateTaskLastRun(taskID, "failed")
		return
	}

	// 在所有主机上执行任务
	var wg sync.WaitGroup
	results := make(chan bool, len(assets))

	for _, asset := range assets {
		wg.Add(1)
		go func(a model.Asset) {
			defer wg.Done()
			success := s.executeOnAsset(&task, &a)
			results <- success
		}(asset)
	}

	// 等待所有执行完成
	go func() {
		wg.Wait()
		close(results)
	}()

	// 统计执行结果
	var successCount, failCount int
	for success := range results {
		if success {
			successCount++
		} else {
			failCount++
		}
	}

	// 更新任务状态
	status := "success"
	if failCount > 0 {
		if successCount > 0 {
			status = "partial"
		} else {
			status = "failed"
		}
	}

	s.service.UpdateTaskLastRun(taskID, status)
	s.logger.Info("定时任务执行完成",
		zap.Uint("task_id", taskID),
		zap.Int("success", successCount),
		zap.Int("failed", failCount),
		zap.String("status", status))
}

// executeOnAsset 在单个主机上执行任务
func (s *Scheduler) executeOnAsset(task *model.ScheduledTask, asset *model.Asset) bool {
	// 创建执行记录
	now := time.Now()
	execution := &model.TaskExecution{
		ScheduledTaskID: &task.ID,
		AssetID:         asset.ID,
		TriggerType:     "scheduled",
		Status:          "running",
		StartedAt:       &now,
	}

	if err := s.db.Create(execution).Error; err != nil {
		s.logger.Error("创建执行记录失败",
			zap.Uint("task_id", task.ID),
			zap.Uint("asset_id", asset.ID),
			zap.Error(err))
		return false
	}

	// 执行命令
	output, exitCode, err := s.executeCommand(task, asset)

	// 更新执行记录
	finishedAt := time.Now()
	execution.FinishedAt = &finishedAt
	execution.Output = output
	execution.ExitCode = &exitCode

	if err != nil {
		execution.Status = "failed"
		execution.Error = err.Error()
		s.logger.Error("执行命令失败",
			zap.Uint("task_id", task.ID),
			zap.Uint("asset_id", asset.ID),
			zap.Error(err))
	} else if exitCode != 0 {
		execution.Status = "failed"
		execution.Error = fmt.Sprintf("退出码: %d", exitCode)
	} else {
		execution.Status = "success"

		// 如果启用了资产更新，解析输出并更新资产信息
		if task.UpdateAssets {
			s.updateAssetFromOutput(asset.ID, output)
		}
	}

	s.db.Save(execution)
	return execution.Status == "success"
}

// updateAssetFromOutput 从输出解析并更新资产信息
func (s *Scheduler) updateAssetFromOutput(assetID uint, output string) {
	// 尝试从输出中提取 JSON 部分
	jsonStart := strings.Index(output, "{")
	jsonEnd := strings.LastIndex(output, "}")
	if jsonStart == -1 || jsonEnd == -1 || jsonEnd < jsonStart {
		s.logger.Warn("无法从输出中提取 JSON", zap.Uint("asset_id", assetID))
		return
	}

	jsonStr := output[jsonStart : jsonEnd+1]

	var info AssetInfo
	if err := json.Unmarshal([]byte(jsonStr), &info); err != nil {
		s.logger.Error("解析资产信息 JSON 失败",
			zap.Uint("asset_id", assetID),
			zap.Error(err),
			zap.String("json", jsonStr))
		return
	}

	// 更新资产信息
	updates := map[string]interface{}{}
	if info.CPU != "" {
		updates["cpu"] = info.CPU
	}
	if info.Memory != "" {
		updates["memory"] = info.Memory
	}
	if info.Disk != "" {
		updates["disk"] = info.Disk
	}

	if len(updates) > 0 {
		if err := s.db.Model(&model.Asset{}).Where("id = ?", assetID).Updates(updates).Error; err != nil {
			s.logger.Error("更新资产信息失败",
				zap.Uint("asset_id", assetID),
				zap.Error(err))
		} else {
			s.logger.Info("已更新资产信息",
				zap.Uint("asset_id", assetID),
				zap.Any("updates", updates))
		}
	}
}

// executeCommand 执行命令
func (s *Scheduler) executeCommand(task *model.ScheduledTask, asset *model.Asset) (string, int, error) {
	if asset.SSHKey == nil {
		return "", -1, fmt.Errorf("主机 %s 没有配置 SSH 密钥", asset.HostName)
	}

	// 获取 SSH 用户
	sshUser := asset.SSHUser
	if sshUser == "" {
		sshUser = asset.SSHKey.SSHUser
	}
	if sshUser == "" {
		sshUser = "root"
	}

	// 获取 SSH 端口
	sshPort := asset.SSHPort
	if sshPort == 0 {
		sshPort = 22
	}

	// 解密 SSH 密钥
	privateKeyBytes, err := sshUtils.Decrypt(asset.SSHKey.PrivateKey, s.cfg.Encryption.Key)
	if err != nil {
		return "", -1, fmt.Errorf("解密 SSH 密钥失败: %w", err)
	}

	// 解密密钥密码
	var passphraseBytes []byte
	if asset.SSHKey.Passphrase != "" {
		passphraseBytes, err = sshUtils.Decrypt(asset.SSHKey.Passphrase, s.cfg.Encryption.Key)
		if err != nil {
			return "", -1, fmt.Errorf("解密密钥密码失败: %w", err)
		}
	}

	// 创建 SSH 客户端
	var client *sshUtils.SSHClient
	if len(passphraseBytes) > 0 {
		client, err = sshUtils.NewSSHClientWithPassphrase(
			asset.IP,
			sshPort,
			sshUser,
			privateKeyBytes,
			passphraseBytes,
			30*time.Second,
		)
	} else {
		client, err = sshUtils.NewSSHClient(
			asset.IP,
			sshPort,
			sshUser,
			privateKeyBytes,
			30*time.Second,
		)
	}
	if err != nil {
		return "", -1, fmt.Errorf("连接 SSH 失败: %w", err)
	}
	defer client.Close()

	// 创建超时上下文
	timeout := time.Duration(task.Timeout) * time.Second
	if timeout == 0 {
		timeout = 5 * time.Minute
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 执行命令
	return client.ExecuteCommandWithTimeout(ctx, task.Content)
}

// parseAssetIDs 解析主机 ID 字符串
func (s *Scheduler) parseAssetIDs(assetIDsStr string) []uint {
	if assetIDsStr == "" {
		return nil
	}

	parts := strings.Split(assetIDsStr, ",")
	ids := make([]uint, 0, len(parts))
	for _, part := range parts {
		if id, err := strconv.ParseUint(strings.TrimSpace(part), 10, 32); err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids
}
