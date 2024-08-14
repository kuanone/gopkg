package logger

import (
	"context"
	"log/slog"
	"runtime"
	"strconv"
	"strings"
)

type SlogLogger struct {
	logger  *slog.Logger
	skip    int
	traceID string
}

func NewSlogLogger(logger *slog.Logger) Logger {
	return &SlogLogger{logger: logger}
}

func (s *SlogLogger) log(ctx context.Context, level slog.Level, msg string, fields ...interface{}) {
	if ctx != nil {
		traceID := GetTraceIDFromContext(ctx)
		if traceID != "" {
			s.traceID = traceID
		}
	}
	pc, file, line, ok := runtime.Caller(s.skip + 2)
	if ok {
		// 获取函数名
		funcName := runtime.FuncForPC(pc).Name()
		funcName = trimFuncName(funcName) // 简化函数名

		s.logger.With("caller", funcName+" "+file+":"+strconv.Itoa(line), "trace_id", s.traceID).Log(ctx, level, msg, fields...)
	} else {
		s.logger.With("trace_id", s.traceID).Log(ctx, level, msg, fields...)
	}
}

func (s *SlogLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	s.log(ctx, slog.LevelInfo, msg, fields...)
}

func (s *SlogLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	s.log(ctx, slog.LevelWarn, msg, fields...)
}

func (s *SlogLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	s.log(ctx, slog.LevelError, msg, fields...)
}

func (s *SlogLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	s.log(ctx, slog.LevelDebug, msg, fields...)
}

func (s *SlogLogger) WithSkip(skip int) Logger {
	return &SlogLogger{
		logger:  s.logger,
		skip:    s.skip + skip,
		traceID: s.traceID,
	}
}

func (s *SlogLogger) WithContext(ctx context.Context) Logger {
	traceID := GetTraceIDFromContext(ctx)
	return &SlogLogger{
		logger:  s.logger,
		skip:    s.skip,
		traceID: traceID,
	}
}

// 辅助函数: 修剪函数名，保留最重要的部分
func trimFuncName(funcName string) string {
	parts := strings.Split(funcName, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return funcName
}
