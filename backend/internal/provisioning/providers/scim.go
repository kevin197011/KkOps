// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"go.uber.org/zap"

	"github.com/kkops/backend/internal/model"
)

type scimProvider struct {
	cfgJSON []byte
}

// NewSCIM builds a generic SCIM 2.0 HTTP client stub (TODO: full Users endpoint CRUD).
func NewSCIM(cfgJSON []byte) (Provider, error) {
	return &scimProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *scimProvider) Kind() string { return KindSCIM }

func (p *scimProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		log.Warn("scim verify skipped: invalid config")
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/ServiceProviderConfig", cfg.Token)
	if err != nil {
		log.Warn("scim verify: build request", zap.Error(err))
		return nil
	}
	st, err := DoJSON(ctx, log, KindSCIM, "verify", req)
	if err != nil || st >= 400 {
		return nil
	}
	return nil
}

func (p *scimProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		log.Warn("scim sync skipped: invalid config")
		return nil
	}
	payload := map[string]interface{}{
		"schemas":     []string{"urn:ietf:params:scim:schemas:core:2.0:User"},
		"userName":    user.Username,
		"displayName": user.RealName,
		"active":      user.Status == "active",
		"emails": []map[string]interface{}{
			{"value": user.Email, "primary": true},
		},
	}
	path := "/Users"
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, path, cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindSCIM, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindSCIM, "sync", req)
	// TODO: upsert via PUT /Users/{id}; capture external id from Location header
	return nil
}

func (p *scimProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	path := "/Users/" + externalID
	req, err := BuildDELETE(ctx, cfg.BaseURL, path, cfg.Token)
	if err != nil {
		LogSoftError(log, KindSCIM, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindSCIM, "delete", req)
	return nil
}
