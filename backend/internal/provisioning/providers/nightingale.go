// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

type nightingaleProvider struct{ cfgJSON []byte }

// NewNightingale stub — n9e / FlashDuty style user bootstrap (TODO).
func NewNightingale(cfgJSON []byte) (Provider, error) {
	return &nightingaleProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *nightingaleProvider) Kind() string { return KindNightingale }

func (p *nightingaleProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/ping", cfg.Token)
	if err != nil {
		LogSoftError(log, KindNightingale, "verify", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindNightingale, "verify", req)
	return nil
}

func (p *nightingaleProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"nickname": user.RealName,
	}
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, "/api/n9e/users", cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindNightingale, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindNightingale, "sync", req)
	return nil
}

func (p *nightingaleProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	req, err := BuildDELETE(ctx, cfg.BaseURL, "/api/n9e/users/"+externalID, cfg.Token)
	if err != nil {
		LogSoftError(log, KindNightingale, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindNightingale, "delete", req)
	return nil
}
