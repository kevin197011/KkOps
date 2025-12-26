package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/kkops/backend/internal/models"
	"github.com/kkops/backend/internal/repository"
	"github.com/kkops/backend/internal/salt"
)

type BatchOperationService interface {
	CreateOperation(operation *models.BatchOperation) error
	ExecuteOperation(operationID uint) error
	GetOperation(id uint) (*models.BatchOperation, error)
	GetOperationStatus(id uint) (*models.BatchOperation, error)
	GetOperationResults(id uint) (*models.BatchOperation, error)
	CancelOperation(id uint) error
	ListOperations(page, pageSize int, filters map[string]interface{}) ([]models.BatchOperation, int64, error)
	CleanupOldOperations(beforeTime time.Time) (int64, error)
}

type batchOperationService struct {
	operationRepo *repository.BatchOperationRepository
	hostRepo      repository.HostRepository
	saltClient    *salt.Client
}

func NewBatchOperationService(
	operationRepo *repository.BatchOperationRepository,
	hostRepo repository.HostRepository,
	saltClient *salt.Client,
) BatchOperationService {
	return &batchOperationService{
		operationRepo: operationRepo,
		hostRepo:      hostRepo,
		saltClient:    saltClient,
	}
}

func (s *batchOperationService) CreateOperation(operation *models.BatchOperation) error {
	operation.Status = "pending"
	return s.operationRepo.Create(operation)
}

func (s *batchOperationService) ExecuteOperation(operationID uint) error {
	// 获取操作
	operation, err := s.operationRepo.Get(operationID)
	if err != nil {
		return err
	}

	// 检查状态
	if operation.Status != "pending" {
		return errors.New("operation is not in pending status")
	}

	// 异步执行
	go s.executeOperationAsync(operationID)

	return nil
}

func (s *batchOperationService) executeOperationAsync(operationID uint) {
	// 更新状态为 running
	_ = s.operationRepo.UpdateStatus(operationID, "running")

	// 获取操作
	operation, err := s.operationRepo.Get(operationID)
	if err != nil {
		_ = s.operationRepo.UpdateResults(operationID, map[string]interface{}{
			"status":       "failed",
			"error_message": fmt.Sprintf("Failed to get operation: %v", err),
		})
		return
	}

	if s.saltClient == nil {
		_ = s.operationRepo.UpdateResults(operationID, map[string]interface{}{
			"status":       "failed",
			"error_message": "Salt client not configured",
		})
		return
	}

	// 解析目标主机
	var targetHosts []map[string]interface{}
	if err := json.Unmarshal(operation.TargetHosts, &targetHosts); err != nil {
		_ = s.operationRepo.UpdateResults(operationID, map[string]interface{}{
			"status":       "failed",
			"error_message": fmt.Sprintf("Failed to parse target hosts: %v", err),
		})
		return
	}

	// 构建 Salt target（主机ID列表）
	var targetList []string
	hostMap := make(map[string]map[string]interface{})
	for _, hostData := range targetHosts {
		hostID := uint64(hostData["id"].(float64))
		host, err := s.hostRepo.GetByID(hostID)
		if err != nil {
			log.Printf("Failed to get host %d: %v", hostID, err)
			continue
		}
		if host.SaltMinionID == nil || *host.SaltMinionID == "" {
			log.Printf("Host %d does not have Salt Minion ID", hostID)
			continue
		}
		targetList = append(targetList, *host.SaltMinionID)
		hostMap[*host.SaltMinionID] = hostData
	}

	if len(targetList) == 0 {
		_ = s.operationRepo.UpdateResults(operationID, map[string]interface{}{
			"status":       "failed",
			"error_message": "No valid target hosts found",
		})
		return
	}

	// 构建 Salt target
	var target string
	if len(targetList) == 1 {
		target = targetList[0]
	} else {
		// 当有多个目标时，使用通配符以避免 Salt Master 不响应自己的问题
		target = "*"
		log.Printf("Using wildcard target '*' for %d minions to avoid Salt Master self-command issue", len(targetList))
	}

	// 解析命令参数
	var args []interface{}
	if len(operation.CommandArgs) > 0 {
		if err := json.Unmarshal(operation.CommandArgs, &args); err != nil {
			log.Printf("Warning: Failed to parse command args: %v", err)
		}
	}

	// 对于没有参数的命令，不要传递空数组
	var finalArgs []interface{}
	if len(args) > 0 {
		finalArgs = args
	}

	startTime := time.Now()

	// 创建带超时的上下文（默认30分钟超时）
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	// 执行 Salt 命令（这里需要检查 saltClient 是否支持 context，如果不支持，需要在 goroutine 中处理超时）
	// 由于 SaltClient.ExecuteCommand 可能不支持 context，我们使用 goroutine + select 来实现超时
	done := make(chan struct{})
	var result map[string]interface{}
	var cmdErr error

	go func() {
		result, cmdErr = s.saltClient.ExecuteCommand(target, operation.CommandFunction, finalArgs)
		close(done)
	}()

	select {
	case <-ctx.Done():
		// 超时
		_ = s.operationRepo.UpdateResults(operationID, map[string]interface{}{
			"status":       "failed",
			"error_message": "Operation timeout after 30 minutes",
			"completed_at": time.Now(),
			"duration_seconds": int(time.Since(startTime).Seconds()),
		})
		return
	case <-done:
		// 命令执行完成，继续处理结果
	}
	if cmdErr != nil {
		_ = s.operationRepo.UpdateResults(operationID, map[string]interface{}{
			"status":       "failed",
			"error_message": fmt.Sprintf("Failed to execute command: %v", cmdErr),
			"completed_at": time.Now(),
			"duration_seconds": int(time.Since(startTime).Seconds()),
		})
		return
	}

	// 处理结果
	results := make(map[string]interface{})
	successCount := 0
	failedCount := 0

	// 将结果按主机组织
	for minionID, output := range result {
		hostData, exists := hostMap[minionID]
		if !exists {
			// 如果使用通配符，跳过不在我们目标列表中的 minion
			log.Printf("Skipping minion %s as it's not in our target list", minionID)
			continue
		}
		hostID := uint64(hostData["id"].(float64))

		var isSuccess bool
		var outputStr string
		var errMsg string

		// 检查输出类型并判断成功/失败
		if output == true || output == "True" {
			isSuccess = true
			outputStr = "Success"
		} else if output == false || output == "False" {
			isSuccess = false
			errMsg = "Command returned False"
		} else {
			// 尝试序列化输出
			outputBytes, err := json.Marshal(output)
			if err != nil {
				outputStr = fmt.Sprintf("%v", output)
			} else {
				outputStr = string(outputBytes)
			}
			isSuccess = true // 默认认为成功（除非是明确的错误）
		}

		results[fmt.Sprintf("%d", hostID)] = map[string]interface{}{
			"success": isSuccess,
			"output":  outputStr,
			"error":   errMsg,
		}

		if isSuccess {
			successCount++
		} else {
			failedCount++
		}
	}

	// 更新操作结果
	resultsJSON, _ := json.Marshal(results)
	duration := int(time.Since(startTime).Seconds())

	_ = s.operationRepo.UpdateResults(operationID, map[string]interface{}{
		"status":          "completed",
		"results":         string(resultsJSON),
		"completed_at":    time.Now(),
		"duration_seconds": duration,
		"success_count":   successCount,
		"failed_count":    failedCount,
	})
}

func (s *batchOperationService) GetOperation(id uint) (*models.BatchOperation, error) {
	return s.operationRepo.Get(id)
}

func (s *batchOperationService) GetOperationStatus(id uint) (*models.BatchOperation, error) {
	return s.operationRepo.Get(id)
}

func (s *batchOperationService) GetOperationResults(id uint) (*models.BatchOperation, error) {
	return s.operationRepo.Get(id)
}

func (s *batchOperationService) CancelOperation(id uint) error {
	operation, err := s.operationRepo.Get(id)
	if err != nil {
		return err
	}

	if operation.Status != "running" {
		return errors.New("operation is not running, cannot cancel")
	}

	// 更新状态为取消
	return s.operationRepo.UpdateResults(id, map[string]interface{}{
		"status": "cancelled",
	})
}

func (s *batchOperationService) ListOperations(page, pageSize int, filters map[string]interface{}) ([]models.BatchOperation, int64, error) {
	return s.operationRepo.List(page, pageSize, filters)
}

// CleanupOldOperations 清理超过指定时间的老操作记录
func (s *batchOperationService) CleanupOldOperations(beforeTime time.Time) (int64, error) {
	// 使用原生SQL进行批量删除
	result := s.operationRepo.GetDB().Unscoped().Where("started_at < ?", beforeTime).Delete(&models.BatchOperation{})
	if result.Error != nil {
		return 0, result.Error
	}
	return result.RowsAffected, nil
}

