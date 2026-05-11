// Copyright (c) 2025 kk
//
// This software is released under the MIT License.
// https://opensource.org/licenses/MIT

package providers

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey int

const loggerCtxKey ctxKey = 1

// WithLogger attaches a zap logger for provider stubs.
func WithLogger(ctx context.Context, log *zap.Logger) context.Context {
	if log == nil {
		return ctx
	}
	return context.WithValue(ctx, loggerCtxKey, log)
}

// Logger returns the logger from context or a no-op logger.
func Logger(ctx context.Context) *zap.Logger {
	if v := ctx.Value(loggerCtxKey); v != nil {
		if l, ok := v.(*zap.Logger); ok && l != nil {
			return l
		}
	}
	return zap.NewNop()
}
