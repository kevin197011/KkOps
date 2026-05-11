// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package monitoring

import (
	"strings"
)

// NightingaleClient queries Nightingale via its Prometheus-compatible HTTP prefix (/prometheus).
type NightingaleClient struct {
	*PrometheusClient
}

// NewNightingaleClientFromConfig builds a client from decrypted integration JSON.
func NewNightingaleClientFromConfig(cfgJSON []byte) (*NightingaleClient, error) {
	pc, err := NewPrometheusClientFromConfig(cfgJSON)
	if err != nil {
		return nil, err
	}
	pc.BaseURL = strings.TrimSuffix(pc.BaseURL, "/") + "/prometheus"
	return &NightingaleClient{PrometheusClient: pc}, nil
}
