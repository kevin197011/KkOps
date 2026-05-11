// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

// Package provider defines the integration Provider abstraction and registry for tool connectors.
package provider

import (
	"context"
)

// Provider is a compile-time registered connector for an integration kind (monitoring, CI, etc.).
type Provider interface {
	Kind() string
	// Verify checks connectivity with stored credentials; returns a typed error on failure (never panics).
	Verify(ctx context.Context) error
	// Metadata returns non-secret diagnostic labels (e.g. server version when known).
	Metadata() map[string]string
}

// Factory builds a Provider from decrypted JSON configuration bytes.
type Factory func(cfgJSON []byte) (Provider, error)
