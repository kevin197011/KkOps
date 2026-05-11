// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

type jumpserverProvider struct{ cfgJSON []byte }

// NewJumpserver stub — Jumpserver core API users resource (TODO).
func NewJumpserver(cfgJSON []byte) (Provider, error) {
	return &jumpserverProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *jumpserverProvider) Kind() string { return KindJumpserver }

func (p *jumpserverProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/api/v1/ping", cfg.Token)
	if err != nil {
		LogSoftError(log, KindJumpserver, "verify", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindJumpserver, "verify", req)
	return nil
}

func (p *jumpserverProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"name":     user.Username,
		"username": user.Username,
		"email":    user.Email,
	}
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, "/api/v1/users/", cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindJumpserver, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindJumpserver, "sync", req)
	return nil
}

func (p *jumpserverProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	req, err := BuildDELETE(ctx, cfg.BaseURL, "/api/v1/users/"+externalID+"/", cfg.Token)
	if err != nil {
		LogSoftError(log, KindJumpserver, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindJumpserver, "delete", req)
	return nil
}
