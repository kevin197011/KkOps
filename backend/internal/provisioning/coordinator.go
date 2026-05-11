// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provisioning

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"

	"github.com/kkops/backend/internal/model"
	"github.com/kkops/backend/internal/provisioning/providers"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// JobKind identifies async provisioning work.
type JobKind string

const (
	JobUserSync   JobKind = "user_sync"
	JobUserDelete JobKind = "user_delete"
	JobTargetSync JobKind = "target_sync"
)

// Job is queued for the in-process worker pool (TODO: Redis-backed queue).
type Job struct {
	Kind     JobKind
	UserID   uint
	TargetID uint // only for JobTargetSync
}

// Coordinator owns the worker pool and executes provisioning against integrations.
type Coordinator struct {
	db           *gorm.DB
	integrations *integrationsvc.Service
	reg          *Registry
	log          *zap.Logger
	jobs         chan Job
	workerCount  int
}

// NewCoordinator starts background workers. Queue capacity is fixed for Phase 1.
func NewCoordinator(db *gorm.DB, integ *integrationsvc.Service, reg *Registry, log *zap.Logger, workers int, queueSize int) *Coordinator {
	if workers < 1 {
		workers = 4
	}
	if queueSize < 1 {
		queueSize = 256
	}
	c := &Coordinator{
		db:           db,
		integrations: integ,
		reg:          reg,
		log:          log,
		jobs:         make(chan Job, queueSize),
		workerCount:  workers,
	}
	for i := 0; i < workers; i++ {
		go c.workerLoop()
	}
	return c
}

func (c *Coordinator) workerLoop() {
	for j := range c.jobs {
		ctx := providers.WithLogger(context.Background(), c.log)
		switch j.Kind {
		case JobUserSync:
			c.runUserSync(ctx, j.UserID)
		case JobUserDelete:
			c.runUserDelete(ctx, j.UserID)
		case JobTargetSync:
			c.runTargetFullSync(ctx, j.TargetID)
		}
	}
}

// EnqueueUserSync schedules sync for one user across all enabled targets.
func (c *Coordinator) EnqueueUserSync(userID uint) {
	select {
	case c.jobs <- Job{Kind: JobUserSync, UserID: userID}:
	default:
		c.log.Warn("provisioning queue full, dropping user_sync", zap.Uint("user_id", userID))
	}
}

// EnqueueUserDelete schedules delete hooks for one user.
func (c *Coordinator) EnqueueUserDelete(userID uint) {
	select {
	case c.jobs <- Job{Kind: JobUserDelete, UserID: userID}:
	default:
		c.log.Warn("provisioning queue full, dropping user_delete", zap.Uint("user_id", userID))
	}
}

// EnqueueTargetSync schedules a manual full sync for one target.
func (c *Coordinator) EnqueueTargetSync(targetID uint) {
	select {
	case c.jobs <- Job{Kind: JobTargetSync, TargetID: targetID}:
	default:
		c.log.Warn("provisioning queue full, dropping target_sync", zap.Uint("target_id", targetID))
	}
}

func (c *Coordinator) runUserSync(ctx context.Context, userID uint) {
	var user model.User
	if err := c.db.First(&user, userID).Error; err != nil {
		c.log.Debug("provisioning skip: user not found", zap.Uint("user_id", userID))
		return
	}
	var targets []model.ProvisioningTarget
	if err := c.db.Where("enabled = ?", true).Find(&targets).Error; err != nil {
		c.log.Error("list provisioning targets", zap.Error(err))
		return
	}
	for i := range targets {
		c.syncOneUserTarget(ctx, &user, &targets[i])
	}
}

func (c *Coordinator) syncOneUserTarget(ctx context.Context, user *model.User, target *model.ProvisioningTarget) {
	start := time.Now()
	cfg, err := c.integrations.DecryptConfigForWorker(target.IntegrationID)
	if err != nil {
		c.finishTargetRun(target.ID, &user.ID, "sync", "failed", err.Error(), start)
		return
	}
	p, err := c.reg.Create(target.ProviderKind, cfg)
	if err != nil {
		c.finishTargetRun(target.ID, &user.ID, "sync", "failed", err.Error(), start)
		return
	}
	ctx = providers.WithLogger(ctx, c.log)
	err = p.Sync(ctx, user)
	msg := "ok"
	st := "success"
	if err != nil {
		msg = err.Error()
		st = "partial"
	}
	extID := fmt.Sprintf("%s-%d", target.ProviderKind, user.ID)
	var row model.ProvisioningUserLink
	q := c.db.Where("user_id = ? AND target_id = ?", user.ID, target.ID).First(&row)
	if q.Error != nil && errors.Is(q.Error, gorm.ErrRecordNotFound) {
		row = model.ProvisioningUserLink{
			UserID:       user.ID,
			TargetID:     target.ID,
			ExternalID:   extID,
			LastSyncedAt: time.Now(),
		}
		_ = c.db.Create(&row).Error
	} else if q.Error == nil {
		row.ExternalID = extID
		row.LastSyncedAt = time.Now()
		_ = c.db.Save(&row).Error
	}
	c.finishTargetRun(target.ID, &user.ID, "sync", st, msg, start)
}

func (c *Coordinator) finishTargetRun(targetID uint, userID *uint, action, status, msg string, started time.Time) {
	run := model.ProvisioningRun{
		TargetID:  targetID,
		UserID:    userID,
		Action:    action,
		Status:    status,
		Message:   msg,
		StartedAt: started,
		EndedAt:   time.Now(),
	}
	_ = c.db.Create(&run).Error
	var t model.ProvisioningTarget
	if err := c.db.First(&t, targetID).Error; err == nil {
		t.Status = status
		if status == "success" {
			t.LastError = ""
			n := time.Now()
			t.LastSyncAt = &n
		} else {
			t.LastError = msg
		}
		_ = c.db.Save(&t).Error
	}
}

func (c *Coordinator) runUserDelete(ctx context.Context, userID uint) {
	var targets []model.ProvisioningTarget
	if err := c.db.Where("enabled = ?", true).Find(&targets).Error; err != nil {
		return
	}
	for i := range targets {
		start := time.Now()
		var link model.ProvisioningUserLink
		_ = c.db.Where("user_id = ? AND target_id = ?", userID, targets[i].ID).First(&link).Error
		cfg, err := c.integrations.DecryptConfigForWorker(targets[i].IntegrationID)
		if err != nil {
			c.finishTargetRun(targets[i].ID, &userID, "delete", "failed", err.Error(), start)
			continue
		}
		p, err := c.reg.Create(targets[i].ProviderKind, cfg)
		if err != nil {
			c.finishTargetRun(targets[i].ID, &userID, "delete", "failed", err.Error(), start)
			continue
		}
		ext := link.ExternalID
		if ext == "" {
			ext = fmt.Sprintf("%s-%d", targets[i].ProviderKind, userID)
		}
		ctx = providers.WithLogger(ctx, c.log)
		_ = p.Delete(ctx, ext)
		c.db.Where("user_id = ? AND target_id = ?", userID, targets[i].ID).Delete(&model.ProvisioningUserLink{})
		c.finishTargetRun(targets[i].ID, &userID, "delete", "success", "", start)
	}
}

func (c *Coordinator) runTargetFullSync(ctx context.Context, targetID uint) {
	var target model.ProvisioningTarget
	if err := c.db.First(&target, targetID).Error; err != nil {
		return
	}
	startAll := time.Now()
	cfg, err := c.integrations.DecryptConfigForWorker(target.IntegrationID)
	if err != nil {
		c.finishTargetRun(target.ID, nil, "verify", "failed", err.Error(), startAll)
		return
	}
	p, err := c.reg.Create(target.ProviderKind, cfg)
	if err != nil {
		c.finishTargetRun(target.ID, nil, "verify", "failed", err.Error(), startAll)
		return
	}
	ctx = providers.WithLogger(ctx, c.log)
	_ = p.Verify(ctx)

	var users []model.User
	if err := c.db.Find(&users).Error; err != nil {
		c.finishTargetRun(target.ID, nil, "sync", "failed", err.Error(), startAll)
		return
	}
	for i := range users {
		u := users[i]
		c.syncOneUserTarget(ctx, &u, &target)
	}
}
