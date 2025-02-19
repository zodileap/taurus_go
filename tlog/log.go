package tlog

import (
	"fmt"
	"os"
)

// Level 定义日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

var levelNames = map[Level]string{
	DebugLevel: "DEBUG",
	InfoLevel:  "INFO",
	WarnLevel:  "WARN",
	ErrorLevel: "ERROR",
	FatalLevel: "FATAL",
}

// Field 定义日志字段
type Field struct {
	Key   string
	Value interface{}
}

// 创建Field的辅助函数
func String(key string, val string) Field {
	return Field{Key: key, Value: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Value: val}
}

func Any(key string, val interface{}) Field {
	return Field{Key: key, Value: val}
}

// Get 获取或创建logger
func Get(name string) *Logger {
	mu.RLock()
	logger, exists := loggers[name]
	mu.RUnlock()

	if !exists {
		logger = &Logger{
			name:      name,
			level:     DebugLevel,
			output:    os.Stdout,
			hasCaller: true,
		}
		mu.Lock()
		loggers[name] = logger
		mu.Unlock()
	}
	return logger
}

// 全局函数
func Debug(name string, msg string, fields ...Field) {
	logger := Get(name)
	logger.log(DebugLevel, msg, fields...)
}

func Info(name string, msg string, fields ...Field) {
	looger := Get(name)
	looger.log(InfoLevel, msg, fields...)
}

func Warn(name string, msg string, fields ...Field) {
	logger := Get(name)
	logger.log(WarnLevel, msg, fields...)
}

func Error(name string, msg string, fields ...Field) {
	logger := Get(name)
	logger.log(ErrorLevel, msg, fields...)
}

func Fatal(name string, msg string, fields ...Field) {
	logger := Get(name)
	logger.log(FatalLevel, msg, fields...)
}

// Print 简单打印日志
func Print(v ...interface{}) {
	logger := Get("print")
	msg := fmt.Sprint(v...)
	logger.log(DebugLevel, msg)
}

// Printf 格式化打印日志
func Printf(format string, v ...interface{}) {
	logger := Get("print")
	msg := fmt.Sprintf(format, v...)
	logger.log(DebugLevel, msg)
}
