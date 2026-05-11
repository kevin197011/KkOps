// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

type harborProvider struct{ cfgJSON []byte }

// NewHarbor stub — Harbor users API (TODO: robot accounts vs humans).
func NewHarbor(cfgJSON []byte) (Provider, error) {
	return &harborProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *harborProvider) Kind() string { return KindHarbor }

func (p *harborProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/api/v2.0/health", cfg.Token)
	if err != nil {
		LogSoftError(log, KindHarbor, "verify", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindHarbor, "verify", req)
	return nil
}

func (p *harborProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"username": user.Username,
		"email":    user.Email,
		"realname": user.RealName,
	}
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, "/api/v2.0/users", cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindHarbor, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindHarbor, "sync", req)
	return nil
}

func (p *harborProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	req, err := BuildDELETE(ctx, cfg.BaseURL, "/api/v2.0/users/"+externalID, cfg.Token)
	if err != nil {
		LogSoftError(log, KindHarbor, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindHarbor, "delete", req)
	return nil
}
