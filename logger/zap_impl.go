package logger

import (
	"context"
	"runtime"
	"strconv"

	"gopkg.in/natefinch/lumberjack.v2"

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

	pc, file, line, ok := runtime.Caller(z.skip + 2)
	if ok {
		// 获取函数名
		funcName := runtime.FuncForPC(pc).Name()
		funcName = trimFuncName(funcName) // 简化函数名

		zapFields = append(zapFields, zap.String("caller", funcName+" "+file+":"+strconv.Itoa(line)))
	}

	// 添加 trace_id 和 caller
	zapFields = append(zapFields, zap.String("trace_id", z.traceID))

	logger := z.logger.WithOptions()
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

func getZapLogWriter(filename string, maxsize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,  // 文件位置
		MaxSize:    maxsize,   // 进行切割之前,日志文件的最大大小(MB为单位)
		MaxAge:     maxAge,    // 保留旧文件的最大天数
		MaxBackups: maxBackup, // 保留旧文件的最大个数
		Compress:   true,      // 是否压缩/归档旧文件
	}
	// AddSync 将 io.Writer 转换为 WriteSyncer。
	// 它试图变得智能：如果 io.Writer 的具体类型实现了 WriteSyncer，我们将使用现有的 Sync 方法。
	// 如果没有，我们将添加一个无操作同步。

	return zapcore.AddSync(lumberJackLogger)
}

// 负责设置 encoding 的日志格式
func getZapEncoder() zapcore.Encoder {
	// 获取一个指定的的EncoderConfig，进行自定义
	encodeConfig := zap.NewProductionEncoderConfig()

	// 序列化时间。eg: 2022-09-01T19:11:35.921+0800
	encodeConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	// "time":"2022-09-01T19:11:35.921+0800"
	encodeConfig.TimeKey = "ts"
	// 将Level序列化为全大写字符串。例如，将info level序列化为INFO。
	encodeConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	// 以 package/file:行 的格式 序列化调用程序，从完整路径中删除除最后一个目录外的所有目录。
	encodeConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encodeConfig)
}
