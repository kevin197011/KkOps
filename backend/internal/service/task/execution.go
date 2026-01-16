// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package task

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/config"
	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/service/sshkey"
	"github.com/kkops/backend/internal/utils"
)

// ExecutionService handles task execution
type ExecutionService struct {
	db        *gorm.DB
	config    *config.Config
	sshkeySvc *sshkey.Service
}

// NewExecutionService creates a new task execution service
func NewExecutionService(db *gorm.DB, cfg *config.Config, sshkeySvc *sshkey.Service) *ExecutionService {
	return &ExecutionService{
		db:        db,
		config:    cfg,
		sshkeySvc: sshkeySvc,
	}
}

// CreateTaskExecutions creates execution records for a task
func (s *ExecutionService) CreateTaskExecutions(taskID uint, assetIDs []uint) error {
	executions := make([]model.TaskExecution, len(assetIDs))
	for i, assetID := range assetIDs {
		executions[i] = model.TaskExecution{
			TaskID:  taskID,
			AssetID: assetID,
			Status:  "pending",
		}
	}

	return s.db.Create(&executions).Error
}

// parseTaskAssetIDs converts comma-separated string to uint slice
func parseTaskAssetIDs(assetIDsStr string) []uint {
	if assetIDsStr == "" {
		return []uint{}
	}
	parts := strings.Split(assetIDsStr, ",")
	ids := make([]uint, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		id, err := strconv.ParseUint(p, 10, 32)
		if err == nil {
			ids = append(ids, uint(id))
		}
	}
	return ids
}

// ExecuteTask executes a task on all target assets
func (s *ExecutionService) ExecuteTask(taskID uint, executionType string) error {
	// Get the task
	var task model.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return fmt.Errorf("task not found: %w", err)
	}

	// Parse asset IDs from the task
	assetIDs := parseTaskAssetIDs(task.AssetIDs)
	if len(assetIDs) == 0 {
		return fmt.Errorf("no assets configured for this task")
	}

	// Update task status
	task.Status = "running"
	now := time.Now()
	task.StartedAt = &now
	if err := s.db.Save(&task).Error; err != nil {
		return err
	}

	// Create new execution records for this run
	executions := make([]model.TaskExecution, len(assetIDs))
	for i, assetID := range assetIDs {
		executions[i] = model.TaskExecution{
			TaskID:  task.ID,
			AssetID: assetID,
			Status:  "pending",
		}
	}
	if err := s.db.Create(&executions).Error; err != nil {
		return fmt.Errorf("failed to create execution records: %w", err)
	}

	if executionType == "async" {
		// Execute asynchronously in goroutines
		for _, exec := range executions {
			go s.executeTaskOnAsset(context.Background(), task, exec.ID)
		}
		return nil
	}

	// Execute synchronously
	for _, exec := range executions {
		if err := s.executeTaskOnAsset(context.Background(), task, exec.ID); err != nil {
			// Continue with other executions even if one fails
			continue
		}
	}

	// Update task status based on execution results
	s.updateTaskStatus(taskID)
	return nil
}

// executeTaskOnAsset executes a task on a specific asset
func (s *ExecutionService) executeTaskOnAsset(ctx context.Context, task model.Task, executionID uint) error {
	var execution model.TaskExecution
	if err := s.db.Preload("Asset").Preload("Asset.SSHKey").First(&execution, executionID).Error; err != nil {
		return err
	}

	asset := execution.Asset
	if asset.Status != "active" {
		execution.Status = "failed"
		execution.Error = "asset is not active"
		now := time.Now()
		execution.FinishedAt = &now
		s.db.Save(&execution)
		return fmt.Errorf("asset %d is not active", asset.ID)
	}

	// Update execution status
	execution.Status = "running"
	now := time.Now()
	execution.StartedAt = &now
	if err := s.db.Save(&execution).Error; err != nil {
		return err
	}

	// Connect to asset via SSH
	sshClient, err := s.connectToAsset(asset)
	if err != nil {
		execution.Status = "failed"
		execution.Error = fmt.Sprintf("SSH connection failed: %v", err)
		now := time.Now()
		execution.FinishedAt = &now
		s.db.Save(&execution)
		// Update overall task status after execution fails (important for async)
		s.updateTaskStatus(task.ID)
		return err
	}
	defer sshClient.Close()

	// Execute the task content with timeout
	command := s.buildCommand(task.Content, task.Type)
	
	// Use task timeout or default to 600 seconds (10 minutes)
	timeout := time.Duration(task.Timeout) * time.Second
	if timeout <= 0 {
		timeout = 600 * time.Second
	}
	execCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	
	output, exitCode, err := sshClient.ExecuteCommandWithTimeout(execCtx, command)

	// Update execution result
	now = time.Now()
	execution.FinishedAt = &now
	execution.Output = output
	execution.ExitCode = &exitCode

	if err != nil {
		execution.Status = "failed"
		if err == context.DeadlineExceeded {
			execution.Error = fmt.Sprintf("execution timed out after %d seconds", task.Timeout)
		} else if err == context.Canceled {
			execution.Error = "execution was cancelled"
		} else {
			execution.Error = err.Error()
		}
	} else if exitCode == 0 {
		execution.Status = "success"
	} else {
		execution.Status = "failed"
		execution.Error = fmt.Sprintf("command exited with code %d", exitCode)
	}

	if err := s.db.Save(&execution).Error; err != nil {
		return err
	}

	// Update overall task status after each execution completes (important for async)
	s.updateTaskStatus(task.ID)

	return nil
}

// connectToAsset establishes an SSH connection to an asset
func (s *ExecutionService) connectToAsset(asset model.Asset) (*utils.SSHClient, error) {
	timeout := 30 * time.Second
	port := asset.SSHPort
	if port == 0 {
		port = 22
	}

	username := asset.SSHUser

	// Try SSH key authentication first
	if asset.SSHKeyID != nil {
		// Load SSH key to get owner user ID
		var sshKey model.SSHKey
		if err := s.db.First(&sshKey, *asset.SSHKeyID).Error; err != nil {
			return nil, fmt.Errorf("SSH key not found: %w", err)
		}

		// Get and decrypt private key
		privateKeyBytes, err := s.sshkeySvc.GetDecryptedPrivateKey(sshKey.UserID, *asset.SSHKeyID)
		if err != nil {
			return nil, fmt.Errorf("failed to get private key: %w", err)
		}

		// Get passphrase if exists
		passphraseBytes, err := s.sshkeySvc.GetDecryptedPassphrase(sshKey.UserID, *asset.SSHKeyID)
		if err != nil {
			return nil, fmt.Errorf("failed to get passphrase: %w", err)
		}

		// Use SSH key username if asset doesn't have one
		if username == "" {
			username = sshKey.SSHUser
		}
		if username == "" {
			return nil, fmt.Errorf("SSH username is required for asset %d", asset.ID)
		}

		if len(passphraseBytes) > 0 {
			return utils.NewSSHClientWithPassphrase(asset.IP, port, username, privateKeyBytes, passphraseBytes, timeout)
		}
		return utils.NewSSHClient(asset.IP, port, username, privateKeyBytes, timeout)
	}

	// Password authentication is not supported for task execution
	// (we don't store passwords in the asset model)
	return nil, fmt.Errorf("SSH key is required for asset %d", asset.ID)
}

// buildCommand builds the command to execute based on task type
func (s *ExecutionService) buildCommand(content, taskType string) string {
	switch taskType {
	case "python":
		return fmt.Sprintf("python3 -c %q", content)
	case "bash", "shell":
		return content
	default:
		return content
	}
}

// updateTaskStatus updates the task status based on execution results
func (s *ExecutionService) updateTaskStatus(taskID uint) {
	var task model.Task
	if err := s.db.Preload("Executions").First(&task, taskID).Error; err != nil {
		return
	}

	// Count execution statuses
	statusCounts := make(map[string]int)
	for _, exec := range task.Executions {
		statusCounts[exec.Status]++
	}

	// Determine overall task status
	if statusCounts["running"] > 0 {
		task.Status = "running"
	} else if statusCounts["failed"] > 0 {
		task.Status = "failed"
	} else if statusCounts["success"] == len(task.Executions) {
		task.Status = "success"
	} else {
		task.Status = "failed"
	}

	// Update finished time if all executions are done
	if statusCounts["running"] == 0 {
		now := time.Now()
		task.FinishedAt = &now
	}

	s.db.Save(&task)
}

// GetTaskExecutions retrieves all executions for a task
func (s *ExecutionService) GetTaskExecutions(taskID uint) ([]model.TaskExecution, error) {
	var executions []model.TaskExecution
	if err := s.db.Where("task_id = ?", taskID).Preload("Asset").Find(&executions).Error; err != nil {
		return nil, err
	}
	return executions, nil
}

// GetTaskExecution retrieves a single execution
func (s *ExecutionService) GetTaskExecution(executionID uint) (*model.TaskExecution, error) {
	var execution model.TaskExecution
	if err := s.db.Preload("Task").Preload("Asset").First(&execution, executionID).Error; err != nil {
		return nil, err
	}
	return &execution, nil
}

// CancelTaskExecution cancels a running task execution
func (s *ExecutionService) CancelTaskExecution(executionID uint) error {
	var execution model.TaskExecution
	if err := s.db.First(&execution, executionID).Error; err != nil {
		return err
	}

	if execution.Status != "running" {
		return fmt.Errorf("execution is not running")
	}

	execution.Status = "cancelled"
	now := time.Now()
	execution.FinishedAt = &now
	return s.db.Save(&execution).Error
}

// CancelTask cancels all running executions for a task
func (s *ExecutionService) CancelTask(taskID uint) error {
	var executions []model.TaskExecution
	if err := s.db.Where("task_id = ? AND status = ?", taskID, "running").Find(&executions).Error; err != nil {
		return err
	}

	now := time.Now()
	for _, exec := range executions {
		exec.Status = "cancelled"
		exec.FinishedAt = &now
		s.db.Save(&exec)
	}

	// Update task status
	var task model.Task
	if err := s.db.First(&task, taskID).Error; err != nil {
		return err
	}
	task.Status = "cancelled"
	task.FinishedAt = &now
	return s.db.Save(&task).Error
}
