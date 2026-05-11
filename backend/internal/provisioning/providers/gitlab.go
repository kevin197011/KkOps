// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

type gitlabProvider struct {
	cfgJSON []byte
}

// NewGitLab stub — sync uses GitLab Users API shape (TODO: PAT scopes, conflict handling).
func NewGitLab(cfgJSON []byte) (Provider, error) {
	return &gitlabProvider{cfgJSON: append([]byte(nil), cfgJSON...)}, nil
}

func (p *gitlabProvider) Kind() string { return KindGitLab }

func (p *gitlabProvider) Verify(ctx context.Context) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	req, err := BuildGET(ctx, cfg.BaseURL, "/api/v4/user", cfg.Token)
	if err != nil {
		LogSoftError(log, KindGitLab, "verify", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindGitLab, "verify", req)
	return nil
}

func (p *gitlabProvider) Sync(ctx context.Context, user *model.User) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" {
		return nil
	}
	payload := map[string]interface{}{
		"email":             user.Email,
		"username":          user.Username,
		"name":              user.RealName,
		"skip_confirmation": true,
	}
	req, err := BuildPOSTJSON(ctx, cfg.BaseURL, "/api/v4/users", cfg.Token, payload)
	if err != nil {
		LogSoftError(log, KindGitLab, "sync", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindGitLab, "sync", req)
	return nil
}

func (p *gitlabProvider) Delete(ctx context.Context, externalID string) error {
	log := Logger(ctx)
	cfg, err := parseBase(p.cfgJSON)
	if err != nil || cfg.BaseURL == "" || externalID == "" {
		return nil
	}
	req, err := BuildDELETE(ctx, cfg.BaseURL, "/api/v4/users/"+externalID, cfg.Token)
	if err != nil {
		LogSoftError(log, KindGitLab, "delete", err, 0)
		return nil
	}
	_, _ = DoJSON(ctx, log, KindGitLab, "delete", req)
	return nil
}
