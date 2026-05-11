// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package aisvc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/kkops/backend/internal/ai"
	"github.com/kkops/backend/internal/model"
)

// RcaResult is returned from GenerateRCA.
type RcaResult struct {
	ReportMD       string
	CitedToolCalls []CitedCall
	RawToolLog     string
}

// GenerateRCA builds a structured Markdown RCA for an incident.
func (s *Service) GenerateRCA(ctx context.Context, reg *ai.Registry, incidentID uint, integrationOverride *uint) (*RcaResult, error) {
	v, err := s.Incidents.Get(ctx, incidentID)
	if err != nil {
		return nil, err
	}
	cited := make([]CitedCall, 0, 8)
	log := make([]map[string]string, 0, 8)

	// Pre-fetch a few tool results for timeline context
	addTool := func(name string, args string, out string) {
		cited = append(cited, CitedCall{Tool: name, Args: json.RawMessage(args)})
		log = append(log, map[string]string{"tool": name, "args": args, "output_preview": truncateRunes(out, 1200)})
	}

	if out, err := s.Tools.Execute(ctx, "list_alerts", `{}`); err == nil {
		addTool("list_alerts", `{}`, out)
	}

	incJSON := fmt.Sprintf(`{"id":%d}`, incidentID)
	if out, err := s.Tools.Execute(ctx, "get_incident", incJSON); err == nil {
		addTool("get_incident", incJSON, out)
	}

	_ = v.LinkedAlertIDs

	logJSON, _ := json.MarshalIndent(log, "", "  ")
	prompt := fmt.Sprintf(`Produce a root cause analysis in Markdown with these sections:
## Symptoms
## Timeline
## Likely Cause
## Evidence
## Recommended Actions

Incident title: %s
Severity: %s
Status: %s

Correlated alert IDs (may be empty): %v

TOOL CONTEXT PREVIEW (JSON, truncated evidence pipeline):
%s

Use ONLY facts supported by tool outputs when citing specifics; otherwise state uncertainty.`,
		v.Title, v.Severity, v.Status, v.LinkedAlertIDs, string(logJSON))

	msgs := []ai.Message{
		{Role: "user", Content: prompt},
	}

	report, err := s.DrainChatNonStream(ctx, reg, integrationOverride, msgs)
	if err != nil {
		return nil, err
	}
	rawToolLog, _ := json.MarshalIndent(log, "", "  ")
	return &RcaResult{
		ReportMD:       report,
		CitedToolCalls: cited,
		RawToolLog:     string(rawToolLog),
	}, nil
}

// PersistRCA stores RCA report row.
func (s *Service) PersistRCA(incidentID uint, report string, rawLog string) (*model.AIRcaReport, error) {
	row := model.AIRcaReport{
		IncidentID:  incidentID,
		GeneratedAt: time.Now().UTC(),
		ReportMD:    report,
		RawToolLog:  rawLog,
	}
	if err := s.DB.Create(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

// ListRcaReports lists RCA rows optionally filtered by incident.
func (s *Service) ListRcaReports(incidentID uint, limit int) ([]model.AIRcaReport, error) {
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	q := s.DB.Order("generated_at DESC").Limit(limit)
	if incidentID > 0 {
		q = q.Where("incident_id = ?", incidentID)
	}
	var rows []model.AIRcaReport
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

// --- anomaly CRUD ---

// SaveRule creates or updates anomaly rule.
func (s *Service) SaveRule(in model.AIAnomalyRule) (*model.AIAnomalyRule, error) {
	if in.ID == 0 {
		if err := s.DB.Create(&in).Error; err != nil {
			return nil, err
		}
		return &in, nil
	}
	var existing model.AIAnomalyRule
	if err := s.DB.First(&existing, in.ID).Error; err != nil {
		return nil, err
	}
	existing.Name = in.Name
	existing.IntegrationID = in.IntegrationID
	existing.Query = in.Query
	existing.ScheduleCron = in.ScheduleCron
	existing.Enabled = in.Enabled
	existing.PromptTemplate = in.PromptTemplate
	if err := s.DB.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

// ListRules returns all anomaly rules.
func (s *Service) ListRules() ([]model.AIAnomalyRule, error) {
	var rows []model.AIAnomalyRule
	err := s.DB.Order("id asc").Find(&rows).Error
	return rows, err
}

// DeleteRule removes a rule.
func (s *Service) DeleteRule(id uint) error {
	return s.DB.Delete(&model.AIAnomalyRule{}, id).Error
}

// ListFindings returns findings with optional filters.
func (s *Service) ListFindings(ruleID uint, since *time.Time, limit int) ([]model.AIAnomalyFinding, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	q := s.DB.Order("ts DESC").Limit(limit)
	if ruleID > 0 {
		q = q.Where("rule_id = ?", ruleID)
	}
	if since != nil {
		q = q.Where("ts >= ?", *since)
	}
	var rows []model.AIAnomalyFinding
	if err := q.Find(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
