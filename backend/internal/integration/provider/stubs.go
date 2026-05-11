// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

type simpleStub struct {
	kind       string
	cfg        BaseConfig
	verifyPath string
	meta       map[string]string
}

func (s *simpleStub) Kind() string { return s.kind }

func (s *simpleStub) Verify(ctx context.Context) error {
	return VerifyHTTPGET(ctx, s.cfg.BaseURL, s.cfg.Token, s.verifyPath)
}

func (s *simpleStub) Metadata() map[string]string {
	out := map[string]string{"kind": s.kind, "base_url": s.cfg.BaseURL}
	for k, v := range s.meta {
		out[k] = v
	}
	return out
}

// NewPrometheusStub verifies Prometheus-style metrics endpoint readiness.
func NewPrometheusStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindPrometheus, cfg: c, verifyPath: "/-/ready", meta: map[string]string{"probe": "prometheus_ready"}}, nil
}

// NewNightingaleStub probes Nightingale / FlashDuty API root (n9e commonly exposes version under /api/n9e).
func NewNightingaleStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindNightingale, cfg: c, verifyPath: "/ping", meta: map[string]string{"probe": "nightingale_ping"}}, nil
}

// NewLokiStub hits Loki readiness path.
func NewLokiStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindLoki, cfg: c, verifyPath: "/ready", meta: map[string]string{"probe": "loki_ready"}}, nil
}

// NewElasticsearchStub uses cluster health (may require roles on the token).
func NewElasticsearchStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindElasticsearch, cfg: c, verifyPath: "/_cluster/health", meta: map[string]string{"probe": "es_cluster_health"}}, nil
}

// NewGrafanaStub checks Grafana health endpoint.
func NewGrafanaStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindGrafana, cfg: c, verifyPath: "/api/health", meta: map[string]string{"probe": "grafana_health"}}, nil
}

// NewJenkinsStub requests Jenkins JSON API with optional API token as Bearer (many setups use basic auth in practice).
func NewJenkinsStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindJenkins, cfg: c, verifyPath: "/api/json", meta: map[string]string{"probe": "jenkins_api_json"}}, nil
}

// NewGitLabStub calls GitLab version API.
func NewGitLabStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindGitLab, cfg: c, verifyPath: "/api/v4/version", meta: map[string]string{"probe": "gitlab_version"}}, nil
}

// NewHarborStub checks Harbor health API v2.
func NewHarborStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindHarbor, cfg: c, verifyPath: "/api/v2.0/health", meta: map[string]string{"probe": "harbor_health"}}, nil
}

// NewArgoCDStub probes Argo CD health endpoint.
func NewArgoCDStub(cfgJSON []byte) (Provider, error) {
	c, err := ParseBase(cfgJSON)
	if err != nil {
		return nil, err
	}
	return &simpleStub{kind: KindArgoCD, cfg: c, verifyPath: "/healthz", meta: map[string]string{"probe": "argocd_healthz"}}, nil
}

type kubeStub struct {
	kind string
}

func (k *kubeStub) Kind() string { return k.kind }

func (k *kubeStub) Verify(ctx context.Context) error {
	_ = ctx
	return nil
}

func (k *kubeStub) Metadata() map[string]string {
	return map[string]string{"kind": k.kind, "credential": "kubeconfig"}
}

// NewKubernetesStub validates decrypted JSON contains non-empty kubeconfig (connectivity is checked on first API call).
func NewKubernetesStub(cfgJSON []byte) (Provider, error) {
	var cred struct {
		Kubeconfig string `json:"kubeconfig"`
	}
	if len(cfgJSON) == 0 {
		return nil, fmt.Errorf("empty integration config")
	}
	if err := json.Unmarshal(cfgJSON, &cred); err != nil {
		return nil, err
	}
	if strings.TrimSpace(cred.Kubeconfig) == "" {
		return nil, fmt.Errorf("kubeconfig is required")
	}
	return &kubeStub{kind: KindKubernetes}, nil
}
