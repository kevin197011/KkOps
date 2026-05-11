// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

type argocdProvider struct{ cfgJSON []byte }

// NewArgoCD stub — account API (TODO: proper session token vs bearer).
func NewArgoCD(cfgJSON []byte) (Provider, error) {
	return &argocdProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *argocdProvider) Kind() string { return KindArgoCD }

func (p *argocdProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/healthz", cfg.Token)
	if err != nil {
		LogSoftError(log, KindArgoCD, "verify", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindArgoCD, "verify", req)
	return nil
}

func (p *argocdProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"username": user.Username,
		"password": "PLACEHOLDER_PROVISIONING_TODO",
	}
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, "/api/v1/account", cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindArgoCD, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindArgoCD, "sync", req)
	return nil
}

func (p *argocdProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	req, err := BuildDELETE(ctx, cfg.BaseURL, "/api/v1/account/"+externalID, cfg.Token)
	if err != nil {
		LogSoftError(log, KindArgoCD, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindArgoCD, "delete", req)
	return nil
}
