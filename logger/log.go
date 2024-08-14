package logger

import (
	"context"
	"fmt"
	"path"
	"runtime"
)

func GetTraceIDFromContext(ctx context.Context) string {
	if traceID, ok := ctx.Value("trace_id").(string); ok {
		return traceID
	}
	return ""
}

func getCallerInfo(skip int) (info string) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		info = "runtime.Caller() failed"
		return
	}
	funcName := runtime.FuncForPC(pc).Name()
	fileName := path.Base(file) // Base函数返回路径的最后一个元素
	return fmt.Sprintf("FuncName:%s, file:%s, line:%d ", funcName, fileName, lineNo)
}
