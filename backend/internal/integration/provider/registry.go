// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provider

import (
	"fmt"
)

// Registry maps integration kind strings to Provider factories.
type Registry struct {
	factories map[string]Factory
}

// NewRegistry builds an empty registry; call RegisterDefaults to populate built-in kinds.
func NewRegistry() *Registry {
	return &Registry{factories: make(map[string]Factory)}
}

// Register adds or replaces a factory for a canonical kind.
func (r *Registry) Register(kind string, fn Factory) {
	k := NormalizeKind(kind)
	r.factories[k] = fn
}

// New builds a Provider for the given kind and decrypted JSON config.
func (r *Registry) New(kind string, cfgJSON []byte) (Provider, error) {
	k := NormalizeKind(kind)
	fn, ok := r.factories[k]
	if !ok {
		return nil, fmt.Errorf("unknown integration kind: %s", kind)
	}
	return fn(cfgJSON)
}

// RegisterDefaults registers stub HTTP providers for all Phase 2 hub kinds.
func (r *Registry) RegisterDefaults() {
	r.Register(KindPrometheus, NewPrometheusStub)
	r.Register(KindNightingale, NewNightingaleStub)
	r.Register(KindLoki, NewLokiStub)
	r.Register(KindElasticsearch, NewElasticsearchStub)
	r.Register(KindGrafana, NewGrafanaStub)
	r.Register(KindJenkins, NewJenkinsStub)
	r.Register(KindGitLab, NewGitLabStub)
	r.Register(KindHarbor, NewHarborStub)
	r.Register(KindArgoCD, NewArgoCDStub)
	r.Register(KindKubernetes, NewKubernetesStub)
}
