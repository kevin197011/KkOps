// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package monitoring

import (
	"context"
)

// FetchNightingaleAlerts is a stub for Nightingale alert-list APIs.
// TODO: integrate with n9e alert/history endpoints when contract is finalized.
func FetchNightingaleAlerts(_ context.Context, _, _ string) ([]AlertmanagerV2Alert, error) {
	return nil, nil
}
