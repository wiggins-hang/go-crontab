package log

import (
	"context"
	"os"

	"google.golang.org/grpc/metadata"
)

var traceEnabled = traceEnableFromEnv()

// 读取环境变量,判断是否开启Trace
// 默认关闭
// 为空或者为0，关闭Trace
// 非空且非0，开启Trace
func traceEnableFromEnv() bool {
	switch os.Getenv("LogTrace") {
	case "":
		fallthrough
	case "0":
		return false
	default:
		return true
	}
}

// EnableTrace 开启trace级别日志
func EnableTrace() {
	traceEnabled = true
}

// SetLevel 设置不同的输出对应的日志级别, output为输出数组下标 "0" "1" "2"
func SetLevel(output string, level Level) {
	GetDefaultLogger().SetLevel(output, level)
}

// GetLevel 获取不同输出对应的日志级别
func GetLevel(output string) Level {
	return GetDefaultLogger().GetLevel(output)
}

// With 日志增加自定义字段，支持多种类型的 Value
func With(fields ...Field) Logger {
	return GetDefaultLogger().With(fields...)
}

// Trace logs to TRACE log. Arguments are handled in the manner of fmt.Print.
func Trace(args ...interface{}) {
	if traceEnabled {
		GetDefaultLogger().Trace(args...)
	}
}

// Tracef logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func Tracef(format string, args ...interface{}) {
	if traceEnabled {
		GetDefaultLogger().Tracef(format, args...)
	}
}

// TraceContextf logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
func TraceContextf(ctx context.Context, format string, args ...interface{}) {
	if traceEnabled {
		GetDefaultLogger().Tracef(format, args...)
	}
}

// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Print.
func Debug(args ...interface{}) {
	GetDefaultLogger().Debug(args...)
}

// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
func Debugf(format string, args ...interface{}) {
	GetDefaultLogger().Debugf(format, args...)
}

// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
func Info(args ...interface{}) {
	GetDefaultLogger().Info(args...)
}

// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func Infof(format string, args ...interface{}) {
	GetDefaultLogger().Infof(format, args...)
}

// InfoContext logs to INFO log. Arguments are handled in the manner of fmt.Print.
func InfoContext(ctx context.Context, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}
	GetDefaultLogger().Info(args...)
}

// InfoContextf logs to INFO log. Arguments are handled in the manner of fmt.Printf.
func InfoContextf(ctx context.Context, format string, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}
	GetDefaultLogger().Infof(format, args...)
}

// Warn logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func Warn(args ...interface{}) {
	GetDefaultLogger().Warn(args...)
}

// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func Warnf(format string, args ...interface{}) {
	GetDefaultLogger().Warnf(format, args...)
}

// WarnContext logs to WARNING log. Arguments are handled in the manner of fmt.Print.
func WarnContext(ctx context.Context, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}
	GetDefaultLogger().Warn(args...)
}

// WarnContextf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
func WarnContextf(ctx context.Context, format string, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}
	GetDefaultLogger().Warnf(format, args...)
}

// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func Error(args ...interface{}) {
	GetDefaultLogger().Error(args...)
}

// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func Errorf(format string, args ...interface{}) {
	GetDefaultLogger().Errorf(format, args...)
}

// ErrorContext logs to ERROR log. Arguments are handled in the manner of fmt.Print.
func ErrorContext(ctx context.Context, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}
	GetDefaultLogger().Error(args...)
}

// ErrorContextf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func ErrorContextf(ctx context.Context, format string, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}

	GetDefaultLogger().Errorf(format, args...)
}

// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func Fatal(args ...interface{}) {
	GetDefaultLogger().Fatal(args...)
}

// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func Fatalf(format string, args ...interface{}) {
	GetDefaultLogger().Fatalf(format, args...)
}

// FatalContext logs to ERROR log. Arguments are handled in the manner of fmt.Print.
// that all Fatal logs will exit with os.Exit(1).
// Implementations may also call os.Exit() with a non-zero exit code.
func FatalContext(ctx context.Context, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}
	GetDefaultLogger().Fatal(args...)
}

// FatalContextf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
func FatalContextf(ctx context.Context, format string, args ...interface{}) {
	if traceId := getTraceIdFromCtx(ctx); traceId != "" {
		args = append(args, Field{"trace_id", traceId})
	}
	GetDefaultLogger().Fatalf(format, args...)
}

func getTraceIdFromCtx(ctx context.Context) string {
	if id, ok := ctx.Value("trace_id").(string); ok {
		return id
	}

	if md, b := metadata.FromIncomingContext(ctx); b {
		if id, ok := md["trace_id"]; ok {
			return id[0]
		}
	}

	return ""
}
