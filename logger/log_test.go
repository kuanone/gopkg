package logger

import (
	"context"
	"log"
	"log/slog"
	"strings"
	"testing"
	"time"

	"go.uber.org/zap"
)

func TestGetTraceIDFromContext(t *testing.T) {
	zapLogger, _ := zap.NewProduction()
	zapLogger = zapLogger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.WithCaller(false))
	logger := NewZapLogger(zapLogger)

	ctx := context.WithValue(context.Background(), "trace_id", "12345")
	logger.WithSkip(2).WithContext(ctx).Info(ctx, "This is an info message with trace id")

	slogLogger := slog.New(slog.NewJSONHandler(log.Writer(), &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				a.Key = "ts"
				ts, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", a.Value.String())
				a.Value = slog.Float64Value(float64(ts.UnixNano()) / 1e9)
			}
			if a.Key == "level" {
				a.Value = slog.StringValue(strings.ToLower(a.Value.String()))
			}
			return a
		},
	}))
	logger = NewSlogLogger(slogLogger)
	logger.WithContext(ctx).Info(ctx, "This is an info message with slog and trace id")
}
