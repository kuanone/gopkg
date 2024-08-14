package logger

import "context"

type Logger interface {
	Info(ctx context.Context, msg string, fields ...interface{})
	Warn(ctx context.Context, msg string, fields ...interface{})
	Error(ctx context.Context, msg string, fields ...interface{})
	Debug(ctx context.Context, msg string, fields ...interface{})

	WithSkip(skip int) Logger
	WithContext(ctx context.Context) Logger
}
