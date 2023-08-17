package log

import (
	"errors"
)

var (
	// DefaultConsoleWriterFactory 默认的console输出流实现
	DefaultConsoleWriterFactory = &ConsoleWriterFactory{}
	// DefaultFileWriterFactory 默认的file输出流实现
	//DefaultFileWriterFactory = &FileWriterFactory{}

	writers = make(map[string]Factory)
)

// RegisterWriter 注册日志输出writer，支持同时多个日志实现
func RegisterWriter(name string, writer Factory) {
	writers[name] = writer
}

// GetWriter 获取日志输出writer，不存在返回nil
func GetWriter(name string) Factory {
	return writers[name]
}

// ConsoleWriterFactory  new console writer instance
type ConsoleWriterFactory struct {
}

// Type 日志插件类型
func (f *ConsoleWriterFactory) Type() string {
	return pluginType
}

// Setup 启动加载配置 并注册console output writer
func (f *ConsoleWriterFactory) Setup(name string, decoder *Decoder) error {
	if decoder == nil {
		return errors.New("console writer decoder empty")
	}

	cfg := &OutputConfig{}
	if err := decoder.Decode(&cfg); err != nil {
		return err
	}

	decoder.Core, decoder.ZapLevel = newConsoleCore(cfg)
	return nil
}
