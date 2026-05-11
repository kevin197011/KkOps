// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package alert

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/kkops/backend/internal/integration/monitoring"
	"github.com/kkops/backend/internal/integration/provider"
	"github.com/kkops/backend/internal/model"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Service persists and aggregates alerts.
type Service struct {
	db      *gorm.DB
	integra *integrationsvc.Service
}

// NewService constructs the alert service.
func NewService(db *gorm.DB, integ *integrationsvc.Service) *Service {
	return &Service{db: db, integra: integ}
}

// AlertView is an API-safe row with integration display name.
type AlertView struct {
	model.AlertRecord
	IntegrationName string `json:"integration_name,omitempty"`
}

// ListParams filters list results.
type ListParams struct {
	Status        string
	IntegrationID uint
	Limit         int
	Offset        int
}

// List returns alerts ordered by starts_at desc, id desc.
func (s *Service) List(_ context.Context, p ListParams) ([]AlertView, int64, error) {
	if p.Limit <= 0 || p.Limit > 500 {
		p.Limit = 50
	}
	q := s.db.Model(&model.AlertRecord{})
	if p.Status != "" {
		q = q.Where("status = ?", p.Status)
	}
	if p.IntegrationID > 0 {
		q = q.Where("integration_id = ?", p.IntegrationID)
	}
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	var rows []model.AlertRecord
	if err := q.Order("starts_at DESC NULLS LAST, id DESC").Limit(p.Limit).Offset(p.Offset).Find(&rows).Error; err != nil {
		return nil, 0, err
	}
	out := make([]AlertView, 0, len(rows))
	for _, r := range rows {
		v := AlertView{AlertRecord: r}
		if r.IntegrationID > 0 {
			var in model.Integration
			if err := s.db.Select("name").First(&in, r.IntegrationID).Error; err == nil {
				v.IntegrationName = in.Name
			}
		}
		out = append(out, v)
	}
	return out, total, nil
}

// Sync pulls alerts from configured Prometheus (Alertmanager URL) and Nightingale (stub) integrations.
func (s *Service) Sync(ctx context.Context) (synced int, errs []string) {
	var integrations []model.Integration
	if err := s.db.Where("enabled = ?", true).Find(&integrations).Error; err != nil {
		return 0, []string{err.Error()}
	}
	for _, in := range integrations {
		k := provider.NormalizeKind(in.Kind)
		cfg, err := s.integra.DecryptConfigForWorker(in.ID)
		if err != nil {
			errs = append(errs, fmt.Sprintf("integration %d decrypt: %v", in.ID, err))
			continue
		}
		baseCfg, err := provider.ParseBase(cfg)
		if err != nil {
			errs = append(errs, fmt.Sprintf("integration %d config: %v", in.ID, err))
			continue
		}
		switch k {
		case provider.KindPrometheus:
			if baseCfg.AlertmanagerBaseURL == "" {
				continue
			}
			alerts, err := monitoring.FetchAlertmanagerAlerts(ctx, baseCfg.AlertmanagerBaseURL, baseCfg.Token)
			if err != nil {
				errs = append(errs, fmt.Sprintf("integration %d alertmanager: %v", in.ID, err))
				continue
			}
			for _, a := range alerts {
				if err := s.upsertFromAlertmanager(ctx, in.ID, a); err != nil {
					errs = append(errs, fmt.Sprintf("integration %d upsert: %v", in.ID, err))
					continue
				}
				synced++
			}
		case provider.KindNightingale:
			alerts, err := monitoring.FetchNightingaleAlerts(ctx, baseCfg.BaseURL, baseCfg.Token)
			if err != nil {
				errs = append(errs, fmt.Sprintf("integration %d nightingale: %v", in.ID, err))
				continue
			}
			for _, a := range alerts {
				if err := s.upsertFromAlertmanager(ctx, in.ID, a); err != nil {
					errs = append(errs, fmt.Sprintf("integration %d upsert: %v", in.ID, err))
					continue
				}
				synced++
			}
		default:
			continue
		}
	}
	return synced, errs
}

func (s *Service) upsertFromAlertmanager(_ context.Context, integrationID uint, a monitoring.AlertmanagerV2Alert) error {
	labels := a.Labels
	if labels == nil {
		labels = map[string]string{}
	}
	fp := strings.TrimSpace(a.Fingerprint)
	if fp == "" {
		fp = fingerprintFromLabels(labels)
	}
	sev := strings.ToLower(labels["severity"])
	if sev == "" {
		sev = "unknown"
	}
	title := pickAnnotation(a.Annotations, "summary", "description", "message")
	if title == "" {
		title = labels["alertname"]
	}
	if title == "" {
		title = "alert"
	}
	status := model.AlertStatusFiring
	switch strings.ToLower(a.State) {
	case "suppressed", "active":
		status = model.AlertStatusFiring
	default:
		if strings.EqualFold(a.State, "resolved") {
			status = model.AlertStatusResolved
		}
	}
	startsAt := parseRFC3339Ptr(firstNonEmpty(ptrStr(a.ActiveAt), ptrStr(a.StartsAt)))
	endsAt := parseRFC3339Ptr(ptrStr(a.EndsAt))
	return s.upsertRecord(integrationID, model.AlertSourcePrometheusAM, fp, sev, title, status, labels, startsAt, endsAt)
}

func ptrStr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func firstNonEmpty(vals ...string) string {
	for _, v := range vals {
		if strings.TrimSpace(v) != "" {
			return v
		}
	}
	return ""
}

func pickAnnotation(a map[string]string, keys ...string) string {
	if a == nil {
		return ""
	}
	for _, k := range keys {
		if v := strings.TrimSpace(a[k]); v != "" {
			return v
		}
	}
	return ""
}

type webhookAlert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:"annotations"`
	StartsAt    string            `json:"startsAt"`
	EndsAt      string            `json:"endsAt"`
	Fingerprint string            `json:"fingerprint"`
}

// IngestWebhook parses an Alertmanager notification payload (or compatible single-alert JSON).
func (s *Service) IngestWebhook(body []byte) (int, error) {
	var envelope struct {
		Status string         `json:"status"`
		Alerts []webhookAlert `json:"alerts"`
	}
	alerts := []webhookAlert{}
	globalStatus := ""
	if err := json.Unmarshal(body, &envelope); err == nil && len(envelope.Alerts) > 0 {
		alerts = envelope.Alerts
		globalStatus = envelope.Status
	} else {
		var single webhookAlert
		if err := json.Unmarshal(body, &single); err != nil {
			return 0, fmt.Errorf("decode webhook: %w", err)
		}
		alerts = []webhookAlert{single}
	}
	n := 0
	for _, a := range alerts {
		labels := a.Labels
		if labels == nil {
			labels = map[string]string{}
		}
		fp := strings.TrimSpace(a.Fingerprint)
		if fp == "" {
			fp = fingerprintFromLabels(labels)
		}
		sev := strings.ToLower(labels["severity"])
		if sev == "" {
			sev = "unknown"
		}
		title := pickAnnotation(a.Annotations, "summary", "description", "message")
		if title == "" {
			title = labels["alertname"]
		}
		if title == "" {
			title = "alert"
		}
		status := model.AlertStatusFiring
		if strings.EqualFold(a.Status, "resolved") || strings.EqualFold(globalStatus, "resolved") {
			status = model.AlertStatusResolved
		}
		startsAt := parseRFC3339Ptr(a.StartsAt)
		endsAt := parseRFC3339Ptr(a.EndsAt)
		if err := s.upsertRecord(0, model.AlertSourceWebhook, fp, sev, title, status, labels, startsAt, endsAt); err != nil {
			return n, err
		}
		n++
	}
	return n, nil
}

func fingerprintFromLabels(labels map[string]string) string {
	keys := make([]string, 0, len(labels))
	for k := range labels {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(labels[k])
		b.WriteByte('|')
	}
	sum := sha256.Sum256([]byte(b.String()))
	return hex.EncodeToString(sum[:])
}

func parseRFC3339Ptr(s string) *time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	for _, layout := range []string{time.RFC3339, time.RFC3339Nano} {
		if t, err := time.Parse(layout, s); err == nil {
			return &t
		}
	}
	return nil
}

func (s *Service) upsertRecord(integrationID uint, source, fp, severity, title, status string, labels map[string]string, startsAt, endsAt *time.Time) error {
	labelsJSON, err := json.Marshal(labels)
	if err != nil {
		return err
	}
	var existing model.AlertRecord
	err = s.db.Where("integration_id = ? AND fingerprint = ?", integrationID, fp).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}
	if err == gorm.ErrRecordNotFound {
		rec := model.AlertRecord{
			IntegrationID: integrationID,
			Source:        source,
			Fingerprint:   fp,
			Severity:      severity,
			Title:         title,
			Status:        status,
			LabelsJSON:    string(labelsJSON),
			StartsAt:      startsAt,
			EndsAt:        endsAt,
		}
		return s.db.Create(&rec).Error
	}
	// Preserve user lifecycle states unless remote resolves.
	if existing.Status == model.AlertStatusAcknowledged || existing.Status == model.AlertStatusDismissed {
		if status == model.AlertStatusResolved {
			existing.Status = model.AlertStatusResolved
			existing.EndsAt = endsAt
			existing.UpdatedAt = time.Now().UTC()
		}
		return s.db.Save(&existing).Error
	}
	existing.Severity = severity
	existing.Title = title
	existing.Status = status
	existing.LabelsJSON = string(labelsJSON)
	if startsAt != nil {
		existing.StartsAt = startsAt
	}
	existing.EndsAt = endsAt
	existing.Source = source
	return s.db.Save(&existing).Error
}

// Acknowledge marks an alert as acknowledged.
func (s *Service) Acknowledge(id uint) error {
	var r model.AlertRecord
	if err := s.db.First(&r, id).Error; err != nil {
		return err
	}
	if r.Status == model.AlertStatusDismissed {
		return fmt.Errorf("cannot acknowledge dismissed alert")
	}
	r.Status = model.AlertStatusAcknowledged
	return s.db.Save(&r).Error
}

// Dismiss marks an alert as dismissed.
func (s *Service) Dismiss(id uint) error {
	var r model.AlertRecord
	if err := s.db.First(&r, id).Error; err != nil {
		return err
	}
	r.Status = model.AlertStatusDismissed
	return s.db.Save(&r).Error
}
