// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"github.com/kkops/backend/internal/model"
)

// Provider syncs users to an external system using decrypted integration JSON config.
type Provider interface {
	Kind() string
	Sync(ctx context.Context, user *model.User) error
	Delete(ctx context.Context, externalID string) error
	Verify(ctx context.Context) error
}

// Factory builds a Provider from decrypted integration config JSON.
type Factory func(cfgJSON []byte) (Provider, error)

// Built-in provider kind identifiers (must match registry keys).
const (
	KindSCIM        = "scim"
	KindGitLab      = "gitlab"
	KindJenkins     = "jenkins"
	KindGrafana     = "grafana"
	KindHarbor      = "harbor"
	KindArgoCD      = "argocd"
	KindJumpserver  = "jumpserver"
	KindNightingale = "nightingale"
)
