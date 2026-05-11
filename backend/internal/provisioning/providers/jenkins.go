// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

type jenkinsProvider struct{ cfgJSON []byte }

// NewJenkins stub — uses Jenkins REST with Basic from token or dedicated fields (TODO).
func NewJenkins(cfgJSON []byte) (Provider, error) {
	return &jenkinsProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *jenkinsProvider) Kind() string { return KindJenkins }

func (p *jenkinsProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/api/json", cfg.Token)
	if err != nil {
		LogSoftError(log, KindJenkins, "verify", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindJenkins, "verify", req)
	return nil
}

func (p *jenkinsProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"fullName": user.RealName,
		"username": user.Username,
	}
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, "/createUser/api/json", cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindJenkins, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindJenkins, "sync", req)
	return nil
}

func (p *jenkinsProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	req, err := BuildDELETE(ctx, cfg.BaseURL, "/user/"+externalID+"/delete", cfg.Token)
	if err != nil {
		LogSoftError(log, KindJenkins, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindJenkins, "delete", req)
	return nil
}
