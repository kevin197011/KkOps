// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package provisioning

import "errors"

var (
	// ErrUnknownKind means no factory registered for the provider kind.
	ErrUnknownKind = errors.New("unknown provisioning provider kind")
)
