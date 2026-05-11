// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package aisvc

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"

	"github.com/kkops/backend/internal/ai"
	"github.com/kkops/backend/internal/integration/monitoring"
	"github.com/kkops/backend/internal/integration/provider"
	"github.com/kkops/backend/internal/model"
	incidentsvc "github.com/kkops/backend/internal/service/incident"
)

// ScheduleAnomalyCron registers enabled anomaly rules on the given cron scheduler.
func (s *Service) ScheduleAnomalyCron(c *cron.Cron, reg *ai.Registry) error {
	var rules []model.AIAnomalyRule
	if err := s.DB.Where("enabled = ?", true).Find(&rules).Error; err != nil {
		return err
	}
	for _, r := range rules {
		rule := r
		spec := strings.TrimSpace(rule.ScheduleCron)
		if spec == "" {
			spec = "*/15 * * * *"
		}
		if _, err := c.AddFunc(spec, func() {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
			defer cancel()
			if err := s.EvaluateAnomalyRule(ctx, reg, rule.ID); err != nil && s.Log != nil {
				s.Log.Warn("anomaly rule run failed", zap.Uint("rule_id", rule.ID), zap.Error(err))
			}
		}); err != nil {
			return fmt.Errorf("cron rule %d: %w", rule.ID, err)
		}
	}
	return nil
}

// EvaluateAnomalyRule pulls metric samples and asks the LLM for severity + summary.
func (s *Service) EvaluateAnomalyRule(ctx context.Context, reg *ai.Registry, ruleID uint) error {
	var rule model.AIAnomalyRule
	if err := s.DB.First(&rule, ruleID).Error; err != nil {
		return err
	}
	if !rule.Enabled {
		return nil
	}
	pub, err := s.Integration.Get(rule.IntegrationID)
	if err != nil {
		return err
	}
	cfg, err := s.Integration.DecryptConfigForWorker(rule.IntegrationID)
	if err != nil {
		return err
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindPrometheus {
		return fmt.Errorf("anomaly rule integration must be prometheus")
	}
	cli, err := monitoring.NewPrometheusClientFromConfig(cfg)
	if err != nil {
		return err
	}
	end := time.Now()
	start := end.Add(-2 * time.Hour)
	res, err := cli.QueryRange(ctx, strings.TrimSpace(rule.Query), start, end, time.Minute)
	if err != nil {
		return err
	}
	seriesJSON, _ := json.MarshalIndent(res, "", "  ")
	template := strings.TrimSpace(rule.PromptTemplate)
	if template == "" {
		template = "You are an SRE assistant. Given Prometheus query_range JSON samples, reply ONLY with compact JSON: {\"severity\":\"info|warning|critical\",\"summary\":\"one sentence\"}. severity=critical only if there is clear outage-level anomaly."
	}
	msgs := []ai.Message{
		{Role: "user", Content: template + "\n\nDATA:\n" + string(seriesJSON)},
	}
	out, err := s.DrainChatNoTools(ctx, reg, nil, msgs)
	if err != nil {
		return err
	}
	sev, summ := parseSeveritySummary(out)
	raw := map[string]interface{}{
		"llm_raw": out,
		"series":  json.RawMessage(seriesJSON),
	}
	rawBytes, _ := json.Marshal(raw)

	find := model.AIAnomalyFinding{
		RuleID:   rule.ID,
		Ts:       time.Now().UTC(),
		Severity: sev,
		Summary:  summ,
		Raw:      string(rawBytes),
	}
	if err := s.DB.Create(&find).Error; err != nil {
		return err
	}

	if strings.EqualFold(sev, "critical") {
		uid, err := s.firstUserID()
		if err != nil {
			return err
		}
		title := fmt.Sprintf("[AI Anomaly] %s", rule.Name)
		if summ != "" {
			title = truncateRunes(title+": "+summ, 500)
		}
		_, err = s.Incidents.Create(ctx, incidentsvc.CreateInput{
			Title:     title,
			Severity:  "critical",
			CreatedBy: uid,
		})
		if err != nil && s.Log != nil {
			s.Log.Warn("failed to create incident from anomaly", zap.Error(err))
		}
	}
	return nil
}

func parseSeveritySummary(out string) (string, string) {
	type verdict struct {
		Severity string `json:"severity"`
		Summary  string `json:"summary"`
	}
	var v verdict
	if err := json.Unmarshal([]byte(strings.TrimSpace(out)), &v); err == nil && v.Severity != "" {
		return strings.ToLower(strings.TrimSpace(v.Severity)), strings.TrimSpace(v.Summary)
	}
	return "info", strings.TrimSpace(out)
}

func truncateRunes(s string, n int) string {
	rs := []rune(s)
	if len(rs) <= n {
		return s
	}
	return string(rs[:n]) + "…"
}

func (s *Service) firstUserID() (uint, error) {
	var u model.User
	if err := s.DB.Order("id asc").First(&u).Error; err != nil {
		return 0, err
	}
	return u.ID, nil
}
