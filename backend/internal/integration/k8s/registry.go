// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package k8s provides cached Kubernetes clients built from integration kubeconfig.
package k8s

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// Credentials is the decrypted JSON shape for kind kubernetes.
type Credentials struct {
	Kubeconfig string `json:"kubeconfig"`
}

// ClusterRegistry caches *kubernetes.Clientset per integration ID when kubeconfig bytes match.
type ClusterRegistry struct {
	mu    sync.RWMutex
	cache map[uint]clusterCacheEntry
}

type clusterCacheEntry struct {
	hash string
	cli  *kubernetes.Clientset
}

// NewClusterRegistry builds an empty registry.
func NewClusterRegistry() *ClusterRegistry {
	return &ClusterRegistry{cache: make(map[uint]clusterCacheEntry)}
}

func kubeconfigHash(yaml []byte) string {
	h := sha256.Sum256(yaml)
	return hex.EncodeToString(h[:])
}

// ParseCredentials unmarshals integration JSON and returns kubeconfig YAML bytes.
func ParseCredentials(decryptedJSON []byte) ([]byte, error) {
	var c Credentials
	if len(decryptedJSON) == 0 {
		return nil, fmt.Errorf("empty integration config")
	}
	if err := json.Unmarshal(decryptedJSON, &c); err != nil {
		return nil, fmt.Errorf("parse kubernetes config: %w", err)
	}
	if len([]byte(c.Kubeconfig)) == 0 {
		return nil, fmt.Errorf("kubeconfig is required")
	}
	return []byte(c.Kubeconfig), nil
}

// Clientset returns a cached clientset for the integration, rebuilding when kubeconfig changes.
func (r *ClusterRegistry) Clientset(integrationID uint, kubeYAML []byte) (*kubernetes.Clientset, error) {
	h := kubeconfigHash(kubeYAML)

	r.mu.RLock()
	if ent, ok := r.cache[integrationID]; ok && ent.hash == h && ent.cli != nil {
		cli := ent.cli
		r.mu.RUnlock()
		return cli, nil
	}
	r.mu.RUnlock()

	cfg, err := clientcmd.RESTConfigFromKubeConfig(kubeYAML)
	if err != nil {
		return nil, fmt.Errorf("kubeconfig: %w", err)
	}
	cli, err := kubernetes.NewForConfig(cfg)
	if err != nil {
		return nil, fmt.Errorf("kubernetes client: %w", err)
	}

	r.mu.Lock()
	defer r.mu.Unlock()
	r.cache[integrationID] = clusterCacheEntry{hash: h, cli: cli}
	return cli, nil
}
