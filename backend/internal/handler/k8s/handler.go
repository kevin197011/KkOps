// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package k8s

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	k8sreg "github.com/kkops/backend/internal/integration/k8s"
	"github.com/kkops/backend/internal/integration/provider"
	integrationsvc "github.com/kkops/backend/internal/service/integration"
)

// Handler exposes Kubernetes proxy APIs for integrations of kind kubernetes.
type Handler struct {
	svc      *integrationsvc.Service
	registry *k8sreg.ClusterRegistry
}

// NewHandler constructs the handler.
func NewHandler(svc *integrationsvc.Service, reg *k8sreg.ClusterRegistry) *Handler {
	return &Handler{svc: svc, registry: reg}
}

func (h *Handler) bindClientset(c *gin.Context) (*kubernetes.Clientset, bool) {
	id64, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil || id64 == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
		return nil, false
	}
	id := uint(id64)
	pub, err := h.svc.Get(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return nil, false
	}
	if provider.NormalizeKind(pub.Kind) != provider.KindKubernetes {
		c.JSON(http.StatusBadRequest, gin.H{"error": "integration kind must be kubernetes"})
		return nil, false
	}
	raw, err := h.svc.DecryptConfigForWorker(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "decrypt config failed"})
		return nil, false
	}
	yamlBytes, err := k8sreg.ParseCredentials(raw)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return nil, false
	}
	cli, err := h.registry.Clientset(id, yamlBytes)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return nil, false
	}
	return cli, true
}

func nsOrDefault(q string) string {
	if strings.TrimSpace(q) == "" {
		return metav1.NamespaceDefault
	}
	return q
}

// ListNamespaces GET /k8s/clusters/:id/namespaces
func (h *Handler) ListNamespaces(c *gin.Context) {
	cli, ok := h.bindClientset(c)
	if !ok {
		return
	}
	ctx := c.Request.Context()
	list, err := cli.CoreV1().Namespaces().List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	out := make([]gin.H, 0, len(list.Items))
	for _, n := range list.Items {
		out = append(out, gin.H{"name": n.Name, "status": string(n.Status.Phase)})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// WorkloadRow is a normalized workload for the UI.
type WorkloadRow struct {
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Replicas  int32  `json:"replicas"`
	Available int32  `json:"available"`
	Image     string `json:"image"`
}

func firstImage(template *corev1.PodTemplateSpec) string {
	if template == nil {
		return ""
	}
	for _, c := range template.Spec.Containers {
		if c.Image != "" {
			return c.Image
		}
	}
	return ""
}

// ListWorkloads GET /k8s/clusters/:id/workloads?namespace=
func (h *Handler) ListWorkloads(c *gin.Context) {
	cli, ok := h.bindClientset(c)
	if !ok {
		return
	}
	ns := nsOrDefault(c.Query("namespace"))
	ctx := c.Request.Context()
	out := make([]WorkloadRow, 0)

	dl, err := cli.AppsV1().Deployments(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	for _, d := range dl.Items {
		avail := d.Status.AvailableReplicas
		out = append(out, WorkloadRow{
			Kind:      "Deployment",
			Name:      d.Name,
			Replicas:  ptrReplicaCount(d.Spec.Replicas),
			Available: avail,
			Image:     firstImage(&d.Spec.Template),
		})
	}

	sl, err := cli.AppsV1().StatefulSets(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	for _, s := range sl.Items {
		out = append(out, WorkloadRow{
			Kind:      "StatefulSet",
			Name:      s.Name,
			Replicas:  ptrReplicaCount(s.Spec.Replicas),
			Available: s.Status.ReadyReplicas,
			Image:     firstImage(&s.Spec.Template),
		})
	}

	dsl, err := cli.AppsV1().DaemonSets(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	for _, ds := range dsl.Items {
		des := ds.Status.DesiredNumberScheduled
		ready := ds.Status.NumberReady
		out = append(out, WorkloadRow{
			Kind:      "DaemonSet",
			Name:      ds.Name,
			Replicas:  des,
			Available: ready,
			Image:     firstImage(&ds.Spec.Template),
		})
	}

	c.JSON(http.StatusOK, gin.H{"data": out})
}

func ptrReplicaCount(r *int32) int32 {
	if r == nil {
		return 0
	}
	return *r
}

// PodRow is a normalized pod row.
type PodRow struct {
	Name     string `json:"name"`
	Status   string `json:"status"`
	Node     string `json:"node"`
	Restarts int32  `json:"restarts"`
	Age      string `json:"age"`
}

func podRestartSum(pod *corev1.Pod) int32 {
	var n int32
	for _, cs := range pod.Status.ContainerStatuses {
		n += cs.RestartCount
	}
	for _, cs := range pod.Status.InitContainerStatuses {
		n += cs.RestartCount
	}
	return n
}

func podPhaseMessage(pod *corev1.Pod) string {
	if pod.Status.Phase != "" {
		return string(pod.Status.Phase)
	}
	for _, c := range pod.Status.Conditions {
		if c.Type == corev1.PodReady && c.Status == corev1.ConditionFalse && c.Message != "" {
			return c.Message
		}
	}
	return "Unknown"
}

func humanAge(t metav1.Time) string {
	if t.IsZero() {
		return ""
	}
	d := time.Since(t.Time).Truncate(time.Second)
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d.Minutes()))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d.Hours()))
	}
	return fmt.Sprintf("%dd", int(d.Hours()/24))
}

// ListPods GET /k8s/clusters/:id/pods?namespace=
func (h *Handler) ListPods(c *gin.Context) {
	cli, ok := h.bindClientset(c)
	if !ok {
		return
	}
	ns := nsOrDefault(c.Query("namespace"))
	ctx := c.Request.Context()
	list, err := cli.CoreV1().Pods(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	out := make([]PodRow, 0, len(list.Items))
	for _, p := range list.Items {
		node := p.Spec.NodeName
		out = append(out, PodRow{
			Name:     p.Name,
			Status:   podPhaseMessage(&p),
			Node:     node,
			Restarts: podRestartSum(&p),
			Age:      humanAge(p.CreationTimestamp),
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// PodLogs GET /k8s/clusters/:id/pods/:name/logs?namespace=&container=&tail=
func (h *Handler) PodLogs(c *gin.Context) {
	cli, ok := h.bindClientset(c)
	if !ok {
		return
	}
	ns := nsOrDefault(c.Query("namespace"))
	name := c.Param("name")
	if name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "pod name is required"})
		return
	}
	tail := int64(200)
	if t := strings.TrimSpace(c.Query("tail")); t != "" {
		n, err := strconv.ParseInt(t, 10, 64)
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid tail"})
			return
		}
		if n > 10000 {
			n = 10000
		}
		tail = n
	}
	opts := &corev1.PodLogOptions{TailLines: &tail}
	if ctr := strings.TrimSpace(c.Query("container")); ctr != "" {
		opts.Container = ctr
	}
	ctx := c.Request.Context()
	req := cli.CoreV1().Pods(ns).GetLogs(name, opts)
	stream, err := req.Stream(ctx)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	defer stream.Close()
	body, err := io.ReadAll(io.LimitReader(stream, 4<<20))
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.Data(http.StatusOK, "text/plain; charset=utf-8", body)
}

// EventRow is a simplified event.
type EventRow struct {
	Type      string `json:"type"`
	Reason    string `json:"reason"`
	Message   string `json:"message"`
	Object    string `json:"object"`
	Namespace string `json:"namespace"`
	Last      string `json:"last"`
}

// ListEvents GET /k8s/clusters/:id/events?namespace=
func (h *Handler) ListEvents(c *gin.Context) {
	cli, ok := h.bindClientset(c)
	if !ok {
		return
	}
	ns := nsOrDefault(c.Query("namespace"))
	ctx := c.Request.Context()
	list, err := cli.CoreV1().Events(ns).List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	out := make([]EventRow, 0, len(list.Items))
	for _, e := range list.Items {
		out = append(out, EventRow{
			Type:      e.Type,
			Reason:    e.Reason,
			Message:   e.Message,
			Object:    e.InvolvedObject.Name,
			Namespace: e.Namespace,
			Last:      humanAge(e.LastTimestamp),
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}

// NodeRow is a simplified node.
type NodeRow struct {
	Name       string   `json:"name"`
	Status     string   `json:"status"`
	KubeletVer string   `json:"kubelet_version"`
	Addresses  []string `json:"addresses"`
}

// ListNodes GET /k8s/clusters/:id/nodes
func (h *Handler) ListNodes(c *gin.Context) {
	cli, ok := h.bindClientset(c)
	if !ok {
		return
	}
	ctx := c.Request.Context()
	list, err := cli.CoreV1().Nodes().List(ctx, metav1.ListOptions{})
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	out := make([]NodeRow, 0, len(list.Items))
	for _, n := range list.Items {
		st := string(corev1.ConditionUnknown)
		for _, c := range n.Status.Conditions {
			if c.Type == corev1.NodeReady {
				if c.Status == corev1.ConditionTrue {
					st = "Ready"
				} else {
					st = "NotReady"
				}
				break
			}
		}
		kubelet := n.Status.NodeInfo.KubeletVersion
		addrs := make([]string, 0)
		for _, a := range n.Status.Addresses {
			if a.Type == corev1.NodeInternalIP || a.Type == corev1.NodeHostName {
				addrs = append(addrs, fmt.Sprintf("%s=%s", a.Type, a.Address))
			}
		}
		out = append(out, NodeRow{
			Name:       n.Name,
			Status:     st,
			KubeletVer: kubelet,
			Addresses:  addrs,
		})
	}
	c.JSON(http.StatusOK, gin.H{"data": out})
}
