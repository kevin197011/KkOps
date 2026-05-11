// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package ai

import "errors"

// ErrNoAIIntegration is returned when no enabled AI integration exists for the requested scope.
var ErrNoAIIntegration = errors.New("no enabled ai integration configured")
