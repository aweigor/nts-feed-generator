package logger

import (
	"context"
	"log/slog"
)

type Logger struct {
	base      *slog.Logger
	namespace string
}

func NewLogger(base *slog.Logger, namespace string) *Logger {
	return &Logger{base: base, namespace: namespace}
}

func (logger *Logger) Info(ctx *context.Context, msg string, args ...any) {
	logger.base.Info(msg, logger.withCtx(ctx, args)...)
}

func (l *Logger) withCtx(ctx *context.Context, args []any) []any {
	if ctx == nil {
		return args
	}
	if id := RequestIDFrom(*ctx); id != "" {
		args = append(args, "request_id", id)
	}
	return args
}
