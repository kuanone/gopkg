package logger

import (
	"context"
	"go.uber.org/zap"
	"log"
	"log/slog"
	"testing"
)

func TestGetTraceIDFromContext(t *testing.T) {
	zapLogger, _ := zap.NewProduction()
	zapLogger = zapLogger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.WithCaller(false))
	logger := NewZapLogger(zapLogger)

	ctx := context.WithValue(context.Background(), "trace_id", "12345")
	logger.WithSkip(2).WithContext(ctx).Info(ctx, "This is an info message with trace id")

	slogLogger := slog.New(slog.NewTextHandler(log.Writer(), nil))
	logger = NewSlogLogger(slogLogger)
	logger.WithContext(ctx).Info(ctx, "This is an info message with slog and trace id")
}
