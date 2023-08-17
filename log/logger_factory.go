package log

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	RegisterWriter(OutputConsole, DefaultConsoleWriterFactory)
	//RegisterWriter(OutputFile, DefaultFileWriterFactory)
	Register(defaultLoggerName, NewZapLog(defaultConfig))
}

const (
	pluginType        = "log"
	defaultLoggerName = "default"
)

var (
	// DefaultLogger 默认的logger实现，初始值为console输出，当框架启动后，通过配置文件初始化后覆盖该值
	DefaultLogger Logger

	mu      sync.RWMutex
	loggers = make(map[string]Logger)
)

// Register 注册日志，支持同时多个日志实现
func Register(name string, logger Logger) {
	mu.Lock()
	defer mu.Unlock()
	if logger == nil {
		panic("log: Register logger is nil")
	}

	if _, dup := loggers[name]; dup && name != defaultLoggerName {
		panic("log: Register called repeated for logger name " + name)
	}

	loggers[name] = logger
	if name == defaultLoggerName {
		DefaultLogger = logger
	}
}

// GetDefaultLogger 默认的logger，通过配置文件key=default来设置, 默认使用console输出
func GetDefaultLogger() Logger {
	mu.RLock()
	l := DefaultLogger
	mu.RUnlock()
	return l
}

// SetLogger 设置默认logger
func SetLogger(logger Logger) {
	mu.Lock()
	DefaultLogger = logger
	mu.Unlock()
}

// Get 通过日志名返回具体的实现 log.Debug使用DefaultLogger打日志，也可以使用 log.Get("name").Debug
func Get(name string) Logger {
	mu.RLock()
	l := loggers[name]
	mu.RUnlock()
	return l
}

// Sync 对注册的所有logger执行Sync动作
func Sync() {
	for _, logger := range loggers {
		_ = logger.Sync()
	}
}

// Decoder log
type Decoder struct {
	OutputConfig *OutputConfig
	Core         zapcore.Core
	ZapLevel     zap.AtomicLevel
}

// Decode 解析writer配置 复制一份
func (d *Decoder) Decode(cfg **OutputConfig) error {
	//output, ok := cfg
	//if !ok {
	//	return fmt.Errorf("decoder config type:%T invalid, not **OutputConfig", cfg)
	//}
	*cfg = d.OutputConfig
	return nil
}

// Factory 日志插件工厂，调用该工厂生成具体日志
type Factory interface {
	Type() string
	Setup(name string, dec *Decoder) error
}
