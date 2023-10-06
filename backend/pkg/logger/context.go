package logger

import (
	"context"

	"go.uber.org/zap"
)

type ctxLogger struct{}

// FromContext returns logger from context.
func FromContext(ctx context.Context) *zap.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*zap.Logger); ok {
		return l
	}

	return zap.L()
}
