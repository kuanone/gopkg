package logger

import (
	"context"
	"runtime"
	"strconv"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLogger struct {
	logger  *zap.Logger
	skip    int
	traceID string
}

func NewZapLogger(logger *zap.Logger) Logger {
	return &ZapLogger{logger: logger}
}

func (z *ZapLogger) log(ctx context.Context, level zapcore.Level, msg string, fields ...interface{}) {
	if ctx != nil {
		traceID := GetTraceIDFromContext(ctx)
		if traceID != "" {
			z.traceID = traceID
		}
	}

	// 转换 fields 为 zap.Fields
	zapFields := make([]zap.Field, len(fields)/2)

	for i := 0; i < len(fields); i += 2 {
		key, ok := fields[i].(string)
		if !ok || i+1 >= len(fields) {
			continue
		}
		zapFields[i/2] = zap.Any(key, fields[i+1])
	}

	pc, file, line, ok := runtime.Caller(z.skip)
	if ok {
		// 获取函数名
		funcName := runtime.FuncForPC(pc).Name()
		funcName = trimFuncName(funcName) // 简化函数名

		zapFields = append(zapFields, zap.String("caller", funcName+" "+file+":"+strconv.Itoa(line)))
	}

	// 添加 trace_id 和 caller
	zapFields = append(zapFields, zap.String("trace_id", z.traceID))

	// 设置 caller 的 skip，并打印调用者信息
	logger := z.logger.WithOptions(zap.AddCallerSkip(z.skip))
	logger.Check(level, msg).Write(zapFields...)
}

func (z *ZapLogger) Info(ctx context.Context, msg string, fields ...interface{}) {
	z.log(ctx, zap.InfoLevel, msg, fields...)
}

func (z *ZapLogger) Warn(ctx context.Context, msg string, fields ...interface{}) {
	z.log(ctx, zap.WarnLevel, msg, fields...)
}

func (z *ZapLogger) Error(ctx context.Context, msg string, fields ...interface{}) {
	z.log(ctx, zap.ErrorLevel, msg, fields...)
}

func (z *ZapLogger) Debug(ctx context.Context, msg string, fields ...interface{}) {
	z.log(ctx, zap.DebugLevel, msg, fields...)
}

func (z *ZapLogger) WithSkip(skip int) Logger {
	return &ZapLogger{
		logger:  z.logger,
		skip:    z.skip + skip,
		traceID: z.traceID,
	}
}

func (z *ZapLogger) WithContext(ctx context.Context) Logger {
	traceID := GetTraceIDFromContext(ctx) // 假设你有一个方法从context中获取traceID
	return &ZapLogger{
		logger:  z.logger,
		skip:    z.skip,
		traceID: traceID,
	}
}
