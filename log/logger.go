package log

import (
	"io"
)

// Level log level
type Level int

// log level const
const (
	LevelNil Level = iota
	LevelTrace
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelFatal
)

// String turns the LogLevel to string.
func (lv *Level) String() string {
	return LevelStrings[*lv]
}

// LevelStrings log level name map reverse
var LevelStrings = map[Level]string{
	LevelTrace: "trace",
	LevelDebug: "debug",
	LevelInfo:  "info",
	LevelWarn:  "warn",
	LevelError: "error",
	LevelFatal: "fatal",
}

// LevelNames log level name map
var LevelNames = map[string]Level{
	"trace": LevelTrace,
	"debug": LevelDebug,
	"info":  LevelInfo,
	"warn":  LevelWarn,
	"error": LevelError,
	"fatal": LevelFatal,
}

// LoggerOptions log options
type LoggerOptions struct {
	LogLevel Level
	Pattern  string
	Writer   io.Writer
}

// LoggerOption log option
type LoggerOption func(*LoggerOptions)

// Field 日志自定义字段
type Field struct {
	Key   string
	Value interface{}
}

// Logger 日志类接口
type Logger interface {
	// Trace logs to TRACE log. Arguments are handled in the manner of fmt.Print.
	Trace(args ...interface{})

	// Tracef logs to TRACE log. Arguments are handled in the manner of fmt.Printf.
	Tracef(format string, args ...interface{})

	// Debug logs to DEBUG log. Arguments are handled in the manner of fmt.Print.
	Debug(args ...interface{})

	// Debugf logs to DEBUG log. Arguments are handled in the manner of fmt.Printf.
	Debugf(format string, args ...interface{})

	// Info logs to INFO log. Arguments are handled in the manner of fmt.Print.
	Info(args ...interface{})

	// Infof logs to INFO log. Arguments are handled in the manner of fmt.Printf.
	Infof(format string, args ...interface{})

	// Warn logs to WARNING log. Arguments are handled in the manner of fmt.Print.
	Warn(args ...interface{})

	// Warnf logs to WARNING log. Arguments are handled in the manner of fmt.Printf.
	Warnf(format string, args ...interface{})

	// Error logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	Error(args ...interface{})

	// Errorf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Errorf(format string, args ...interface{})

	// Fatal logs to ERROR log. Arguments are handled in the manner of fmt.Print.
	// that all Fatal logs will exit with os.Exit(1).
	// Implementations may also call os.Exit() with a non-zero exit code.
	Fatal(args ...interface{})

	// Fatalf logs to ERROR log. Arguments are handled in the manner of fmt.Printf.
	Fatalf(format string, args ...interface{})

	// Sync calls the underlying Core's Sync method, flushing any buffered log entries.
	// Applications should take care to call Sync before exiting
	Sync() error

	// SetLevel 设置输出端日志级别
	SetLevel(output string, level Level)

	// GetLevel 获取输出端日志级别
	GetLevel(output string) Level

	// With 日志增加自定义字段，支持多种类型 Value
	With(fields ...Field) Logger
}
