// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provisioning

import (
	"github.com/kkops/backend/internal/provisioning/providers"
)

// Registry maps provider kind strings to factories.
type Registry struct {
	factories map[string]providers.Factory
}

// NewRegistry builds a registry with all built-in stubs registered.
func NewRegistry() *Registry {
	r := &Registry{factories: make(map[string]providers.Factory)}
	r.Register(providers.KindSCIM, providers.NewSCIM)
	r.Register(providers.KindGitLab, providers.NewGitLab)
	r.Register(providers.KindJenkins, providers.NewJenkins)
	r.Register(providers.KindGrafana, providers.NewGrafana)
	r.Register(providers.KindHarbor, providers.NewHarbor)
	r.Register(providers.KindArgoCD, providers.NewArgoCD)
	r.Register(providers.KindJumpserver, providers.NewJumpserver)
	r.Register(providers.KindNightingale, providers.NewNightingale)
	return r
}

// Register replaces or adds a factory for kind.
func (r *Registry) Register(kind string, f providers.Factory) {
	r.factories[kind] = f
}

// Create builds a Provider instance.
func (r *Registry) Create(kind string, cfgJSON []byte) (providers.Provider, error) {
	f, ok := r.factories[kind]
	if !ok || f == nil {
		return nil, ErrUnknownKind
	}
	return f(cfgJSON)
}
