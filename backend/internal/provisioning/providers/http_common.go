// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go.uber.org/zap"
)

const httpTimeout = 15 * time.Second

type baseConfig struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
}

func parseBase(cfgJSON []byte) (baseConfig, error) {
	var c baseConfig
	if len(cfgJSON) == 0 {
		return c, fmt.Errorf("empty config")
	}
	if err := json.Unmarshal(cfgJSON, &c); err != nil {
		return c, err
	}
	c.BaseURL = strings.TrimSuffix(strings.TrimSpace(c.BaseURL), "/")
	return c, nil
}

func newHTTPClient() *http.Client {
	return &http.Client{Timeout: httpTimeout}
}

// LogSoftError logs network / 401 issues without failing the whole provisioning batch (Phase 1 stubs).
func LogSoftError(log *zap.Logger, kind string, op string, err error, status int) {
	if err != nil {
		if strings.Contains(err.Error(), "timeout") || strings.Contains(err.Error(), "connection refused") {
			log.Warn("provisioning soft network error", zap.String("kind", kind), zap.String("op", op), zap.Error(err))
			return
		}
	}
	if status == 401 || status == 403 {
		log.Warn("provisioning soft auth error", zap.String("kind", kind), zap.String("op", op), zap.Int("status", status))
		return
	}
	if err != nil {
		log.Info("provisioning request completed with error", zap.String("kind", kind), zap.String("op", op), zap.Error(err))
	}
}

// DoJSON performs an HTTP request and returns status; body is optionally logged on failure.
func DoJSON(ctx context.Context, log *zap.Logger, kind string, op string, req *http.Request) (int, error) {
	cli := newHTTPClient()
	resp, err := cli.Do(req.WithContext(ctx))
	if err != nil {
		LogSoftError(log, kind, op, err, 0)
		return 0, err
	}
	defer resp.Body.Close()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode >= 400 {
		LogSoftError(log, kind, op, fmt.Errorf("status %d", resp.StatusCode), resp.StatusCode)
	}
	return resp.StatusCode, nil
}

// BuildGET creates a GET request with optional bearer token.
func BuildGET(ctx context.Context, baseURL, path, token string) (*http.Request, error) {
	u := baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

// BuildPOSTJSON creates POST with JSON body.
// BuildDELETE creates DELETE with optional bearer token.
func BuildDELETE(ctx context.Context, baseURL, path, token string) (*http.Request, error) {
	u := baseURL + path
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, u, nil)
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Accept", "application/json")
	return req, nil
}

func BuildPOSTJSON(ctx context.Context, baseURL, path, token string, payload interface{}) (*http.Request, error) {
	u := baseURL + path
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	return req, nil
}
