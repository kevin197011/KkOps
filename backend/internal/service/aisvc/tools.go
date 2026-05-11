// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package aisvc

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	k8sreg "github.com/kkops/backend/internal/integration/k8s"
	"github.com/kkops/backend/internal/integration/logging"
	"github.com/kkops/backend/internal/integration/monitoring"
	"github.com/kkops/backend/internal/integration/provider"
	alertsvc "github.com/kkops/backend/internal/service/alert"
	gitopsviewsvc "github.com/kkops/backend/internal/service/gitopsview"
	incidentsvc "github.com/kkops/backend/internal/service/incident"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// ToolBridge executes read-only assistant tools against platform services.
type ToolBridge struct {
	Integration *integrationsvc.Service
	Alerts      *alertsvc.Service
	Incidents   *incidentsvc.Service
	GitOpsView  *gitopsviewsvc.Service
	K8sRegistry *k8sreg.ClusterRegistry
}

// Execute runs a tool by name with JSON args object string (may be empty {}).
func (b *ToolBridge) Execute(ctx context.Context, name string, argsJSON string) (string, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	var args map[string]interface{}
	if strings.TrimSpace(argsJSON) == "" {
		args = map[string]interface{}{}
	} else if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
		return "", fmt.Errorf("invalid tool args JSON: %w", err)
	}

	switch name {
	case "list_alerts":
		return b.toolListAlerts(ctx, args)
	case "get_incident":
		return b.toolGetIncident(ctx, args)
	case "query_metric":
		return b.toolQueryMetric(ctx, args)
	case "search_logs":
		return b.toolSearchLogs(ctx, args)
	case "list_pods":
		return b.toolListPods(ctx, args)
	case "pipeline_status":
		return b.toolPipelineStatus(ctx, args)
	case "list_integrations":
		return b.toolListIntegrations(ctx)
	default:
		return "", fmt.Errorf("unknown tool %q", name)
	}
}

func uintArg(args map[string]interface{}, key string) (uint, bool) {
	v, ok := args[key]
	if !ok || v == nil {
		return 0, false
	}
	switch t := v.(type) {
	case float64:
		return uint(t), true
	case int:
		return uint(t), true
	case json.Number:
		n, err := t.Int64()
		if err != nil {
			return 0, false
		}
		return uint(n), true
	case string:
		n, err := strconv.ParseUint(strings.TrimSpace(t), 10, 32)
		if err != nil {
			return 0, false
		}
		return uint(n), true
	default:
		return 0, false
	}
}

func strArg(args map[string]interface{}, key string) string {
	v, ok := args[key]
	if !ok || v == nil {
		return ""
	}
	switch t := v.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(t)
	}
}

func (b *ToolBridge) toolListAlerts(ctx context.Context, args map[string]interface{}) (string, error) {
	status := strArg(args, "status")
	rows, _, err := b.Alerts.List(ctx, alertsvc.ListParams{Status: status, Limit: 50, Offset: 0})
	if err != nil {
		return "", err
	}
	raw, _ := json.MarshalIndent(rows, "", "  ")
	return string(raw), nil
}

func (b *ToolBridge) toolGetIncident(ctx context.Context, args map[string]interface{}) (string, error) {
	id, ok := uintArg(args, "id")
	if !ok {
		return "", fmt.Errorf("id is required")
	}
	v, err := b.Incidents.Get(ctx, id)
	if err != nil {
		return "", err
	}
	raw, _ := json.MarshalIndent(v, "", "  ")
	return string(raw), nil
}

func (b *ToolBridge) toolQueryMetric(ctx context.Context, args map[string]interface{}) (string, error) {
	iid, ok := uintArg(args, "integration_id")
	if !ok {
		return "", fmt.Errorf("integration_id is required")
	}
	q := strArg(args, "query")
	if q == "" {
		return "", fmt.Errorf("query is required")
	}
	pub, err := b.Integration.Get(iid)
	if err != nil {
		return "", err
	}
	cfg, err := b.Integration.DecryptConfigForWorker(iid)
	if err != nil {
		return "", err
	}
	k := provider.NormalizeKind(pub.Kind)
	end := time.Now()
	start := end.Add(-1 * time.Hour)
	step := 30 * time.Second
	switch k {
	case provider.KindPrometheus:
		cli, err := monitoring.NewPrometheusClientFromConfig(cfg)
		if err != nil {
			return "", err
		}
		res, err := cli.QueryRange(ctx, q, start, end, step)
		if err != nil {
			return "", err
		}
		raw, _ := json.MarshalIndent(res, "", "  ")
		return string(raw), nil
	default:
		return "", fmt.Errorf("integration %d is not prometheus", iid)
	}
}

func (b *ToolBridge) toolSearchLogs(ctx context.Context, args map[string]interface{}) (string, error) {
	iid, ok := uintArg(args, "integration_id")
	if !ok {
		return "", fmt.Errorf("integration_id is required")
	}
	query := strArg(args, "query")
	if query == "" {
		return "", fmt.Errorf("query is required")
	}
	rng := strArg(args, "range")
	pub, err := b.Integration.Get(iid)
	if err != nil {
		return "", err
	}
	cfg, err := b.Integration.DecryptConfigForWorker(iid)
	if err != nil {
		return "", err
	}
	end := time.Now()
	start := end.Add(-1 * time.Hour)
	if rng != "" {
		if d, err := time.ParseDuration(rng); err == nil {
			start = end.Add(-d)
		}
	}
	k := provider.NormalizeKind(pub.Kind)
	switch k {
	case provider.KindLoki:
		cli, err := logging.NewLokiClientFromConfig(cfg)
		if err != nil {
			return "", err
		}
		lines, err := cli.Search(ctx, query, start, end, 80)
		if err != nil {
			return "", err
		}
		raw, _ := json.MarshalIndent(lines, "", "  ")
		return string(raw), nil
	case provider.KindElasticsearch:
		cli, err := logging.NewElasticsearchClientFromConfig(cfg)
		if err != nil {
			return "", err
		}
		lines, err := cli.Search(ctx, query, start, end, 80)
		if err != nil {
			return "", err
		}
		raw, _ := json.MarshalIndent(lines, "", "  ")
		return string(raw), nil
	default:
		return "", fmt.Errorf("integration is not loki/elasticsearch")
	}
}

func (b *ToolBridge) toolListPods(ctx context.Context, args map[string]interface{}) (string, error) {
	cid, ok := uintArg(args, "cluster_id")
	if !ok {
		return "", fmt.Errorf("cluster_id is required")
	}
	ns := strArg(args, "namespace")
	if ns == "" {
		ns = metav1.NamespaceAll
	}
	pub, err := b.Integration.Get(cid)
	if err != nil {
		return "", err
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindKubernetes {
		return "", fmt.Errorf("cluster_id must be a kubernetes integration")
	}
	raw, err := b.Integration.DecryptConfigForWorker(cid)
	if err != nil {
		return "", err
	}
	yamlBytes, err := k8sreg.ParseCredentials(raw)
	if err != nil {
		return "", err
	}
	cli, err := b.K8sRegistry.Clientset(cid, yamlBytes)
	if err != nil {
		return "", err
	}
	list, err := cli.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		return "", err
	}
	type slim struct {
		Name      string `json:"name"`
		Namespace string `json:"namespace"`
		Phase     string `json:"phase"`
	}
	out := make([]slim, 0, len(list.Items))
	for _, p := range list.Items {
		out = append(out, slim{Name: p.Name, Namespace: p.Namespace, Phase: string(p.Status.Phase)})
	}
	rawJ, _ := json.MarshalIndent(out, "", "  ")
	return string(rawJ), nil
}

func (b *ToolBridge) toolPipelineStatus(ctx context.Context, args map[string]interface{}) (string, error) {
	app := strArg(args, "app")
	if app == "" {
		return "", fmt.Errorf("app is required")
	}
	var argoID *uint
	if id, ok := uintArg(args, "argocd_integration_id"); ok {
		argoID = &id
	}
	ev, err := b.GitOpsView.PipelineView(ctx, app, argoID)
	if err != nil {
		return "", err
	}
	raw, _ := json.MarshalIndent(ev, "", "  ")
	return string(raw), nil
}

func (b *ToolBridge) toolListIntegrations(ctx context.Context) (string, error) {
	list, err := b.Integration.List()
	if err != nil {
		return "", err
	}
	raw, _ := json.MarshalIndent(list, "", "  ")
	return string(raw), nil
}
