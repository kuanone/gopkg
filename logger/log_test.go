package logger

import (
	"context"
	"log/slog"
	"testing"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func TestGetTraceIDFromContext(t *testing.T) {
	logWriter := &lumberjack.Logger{
		Filename:   "./t.log",
		MaxSize:    500, // megabytes
		MaxBackups: 3,
		MaxAge:     28,   // days
		Compress:   true, // disabled by default
	}

	// zapLogger, _ := zap.NewProduction()
	// zapLogger = zapLogger.WithOptions(zap.AddStacktrace(zap.ErrorLevel), zap.WithCaller(false))

	// 获取日志写入位置
	writeSyncer := getZapLogWriter("./t.log", 500, 3, 28)
	// 获取日志编码格式
	encoder := getZapEncoder()

	// 获取日志最低等级，即>=该等级，才会被写入。
	l := new(zapcore.Level)
	err := l.UnmarshalText([]byte("info"))
	if err != nil {
		return
	}

	// 创建一个将日志写入 WriteSyncer 的核心。
	core := zapcore.NewCore(encoder, writeSyncer, l)
	zapLogger := zap.New(core)

	// 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	zap.ReplaceGlobals(zapLogger)

	logger := NewZapLogger(zapLogger)

	ctx := context.WithValue(context.Background(), "trace_id", "12345")
	logger.WithContext(ctx).Info(ctx, "This is an info message with trace id")

	slogLogger := slog.New(slog.NewJSONHandler(logWriter, &slog.HandlerOptions{
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			if a.Key == "time" {
				a.Key = "ts"
				ts, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", a.Value.String())
				// a.Value = slog.Float64Value(float64(ts.UnixNano()) / 1e9)
				a.Value = slog.StringValue(ts.Format("2006-01-02T15:04:05.000Z0700"))
			}
			return a
		},
	}))
	logger = NewSlogLogger(slogLogger)
	logger.WithContext(ctx).Info(ctx, "This is an info message with trace id")
}
