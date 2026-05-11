// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

type grafanaProvider struct{ cfgJSON []byte }

// NewGrafana stub — Admin API user invite shape (TODO).
func NewGrafana(cfgJSON []byte) (Provider, error) {
	return &grafanaProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *grafanaProvider) Kind() string { return KindGrafana }

func (p *grafanaProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/api/health", cfg.Token)
	if err != nil {
		LogSoftError(log, KindGrafana, "verify", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindGrafana, "verify", req)
	return nil
}

func (p *grafanaProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"name":  user.Username,
		"email": user.Email,
		"login": user.Username,
	}
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, "/api/admin/users", cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindGrafana, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindGrafana, "sync", req)
	return nil
}

func (p *grafanaProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	req, err := BuildDELETE(ctx, cfg.BaseURL, "/api/admin/users/"+externalID, cfg.Token)
	if err != nil {
		LogSoftError(log, KindGrafana, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindGrafana, "delete", req)
	return nil
}
