package service

import (
	"log"
	"time"
)

type SchedulerService interface {
	Start()
	Stop()
}

type schedulerService struct {
	settingsService       SettingsService
	hostService           HostService
	batchOperationService BatchOperationService
	stopChan              chan struct{}
}

func NewSchedulerService(settingsService SettingsService, hostService HostService, batchOperationService BatchOperationService) SchedulerService {
	return &schedulerService{
		settingsService:       settingsService,
		hostService:           hostService,
		batchOperationService: batchOperationService,
		stopChan:              make(chan struct{}),
	}
}

func (s *schedulerService) Start() {
	log.Println("启动定时任务调度器")

	// 启动审计日志清理任务（凌晨3点）
	go s.startAuditLogCleanupScheduler()

	// 启动主机状态同步任务（凌晨3点）
	go s.startHostStatusSyncScheduler()

	// 启动批量操作清理任务（每天凌晨4点）
	go s.startBatchOperationCleanupScheduler()
}

func (s *schedulerService) Stop() {
	log.Println("停止定时任务调度器")
	close(s.stopChan)
}

// startAuditLogCleanupScheduler 启动审计日志清理调度器
func (s *schedulerService) startAuditLogCleanupScheduler() {
	// 计算下一次凌晨3点的时间
	nextRun := s.getNext3AM()
	log.Printf("审计日志清理任务已调度，下次运行时间: %v", nextRun)

	for {
		// 计算到下次运行的等待时间
		waitDuration := time.Until(nextRun)
		
		select {
		case <-time.After(waitDuration):
			// 执行清理任务
			s.runAuditLogCleanup()
			// 计算下一次运行时间（24小时后的凌晨3点）
			nextRun = nextRun.Add(24 * time.Hour)
			log.Printf("下次审计日志清理时间: %v", nextRun)
			
		case <-s.stopChan:
			log.Println("审计日志清理调度器已停止")
			return
		}
	}
}

// getNext3AM 获取下一次凌晨3点的时间
func (s *schedulerService) getNext3AM() time.Time {
	now := time.Now()
	// 今天凌晨3点
	today3AM := time.Date(now.Year(), now.Month(), now.Day(), 3, 0, 0, 0, now.Location())
	
	// 如果现在已经过了今天凌晨3点，则返回明天凌晨3点
	if now.After(today3AM) {
		return today3AM.Add(24 * time.Hour)
	}
	
	// 否则返回今天凌晨3点
	return today3AM
}

// runAuditLogCleanup 执行审计日志清理
func (s *schedulerService) runAuditLogCleanup() {
	log.Println("开始执行审计日志清理任务")
	
	deletedCount, err := s.settingsService.CleanupAuditLogs()
	if err != nil {
		log.Printf("审计日志清理失败: %v", err)
		return
	}
	
	log.Printf("审计日志清理完成，删除了 %d 条记录", deletedCount)
}

// startHostStatusSyncScheduler 启动主机状态同步调度器（每天凌晨3点执行）
func (s *schedulerService) startHostStatusSyncScheduler() {
	// 计算下一次凌晨3点的时间（比审计日志清理晚5分钟，避免同时执行）
	nextRun := s.getNext3AM().Add(5 * time.Minute)
	log.Printf("主机状态同步任务已调度，下次运行时间: %v", nextRun)

	for {
		// 计算到下次运行的等待时间
		waitDuration := time.Until(nextRun)
		
		select {
		case <-time.After(waitDuration):
			// 执行同步任务
			s.runHostStatusSync()
			// 计算下一次运行时间（24小时后）
			nextRun = nextRun.Add(24 * time.Hour)
			log.Printf("下次主机状态同步时间: %v", nextRun)
			
		case <-s.stopChan:
			log.Println("主机状态同步调度器已停止")
			return
		}
	}
}

// runHostStatusSync 执行主机状态同步
func (s *schedulerService) runHostStatusSync() {
	log.Println("开始执行主机状态同步任务")
	
	if s.hostService == nil {
		log.Println("主机服务未初始化，跳过状态同步")
		return
	}
	
	err := s.hostService.SyncAllHostsStatus()
	if err != nil {
		log.Printf("主机状态同步失败: %v", err)
		return
	}
	
	log.Println("主机状态同步完成")
}

// startBatchOperationCleanupScheduler 启动批量操作清理调度器
func (s *schedulerService) startBatchOperationCleanupScheduler() {
	// 计算下一次凌晨4点的时间
	nextRun := s.getNext4AM()
	log.Printf("批量操作清理任务已调度，下次运行时间: %v", nextRun)

	for {
		// 计算到下次运行的等待时间
		waitDuration := time.Until(nextRun)

		select {
		case <-time.After(waitDuration):
			s.runBatchOperationCleanup()
			// 计算下一次运行时间
			nextRun = s.getNext4AM()
			log.Printf("批量操作清理任务已重新调度，下次运行时间: %v", nextRun)
		case <-s.stopChan:
			log.Println("批量操作清理调度器已停止")
			return
		}
	}
}

// getNext4AM 计算下一次凌晨4点的时间
func (s *schedulerService) getNext4AM() time.Time {
	now := time.Now()
	next := time.Date(now.Year(), now.Month(), now.Day(), 4, 0, 0, 0, now.Location())

	// 如果当前时间已经过了凌晨4点，则计算明天的凌晨4点
	if now.After(next) {
		next = next.AddDate(0, 0, 1)
	}

	return next
}

// runBatchOperationCleanup 执行批量操作清理
func (s *schedulerService) runBatchOperationCleanup() {
	log.Println("开始执行批量操作清理任务")

	if s.batchOperationService == nil {
		log.Println("批量操作服务未初始化，跳过清理任务")
		return
	}

	// 计算1个月前的时间
	oneMonthAgo := time.Now().AddDate(0, -1, 0)

	count, err := s.batchOperationService.CleanupOldOperations(oneMonthAgo)
	if err != nil {
		log.Printf("批量操作清理失败: %v", err)
		return
	}

	log.Printf("批量操作清理完成，已清理 %d 条1个月前的记录", count)
}