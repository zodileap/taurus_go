package tlog

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// Logger 日志结构体
type Logger struct {
	name      string
	level     Level
	output    io.Writer
	file      *os.File
	mu        sync.Mutex
	fields    []Field
	hasCaller bool
	writer    *LogWriter
}

// 全局logger映射
var (
	loggers = make(map[string]*Logger)
	mu      sync.RWMutex
)

// SetOutputPath 设置日志输出路径
func (l *Logger) SetOutputPath(path string, maxSize int, maxBackups int, maxAge int) *Logger {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 创建目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Printf("创建日志目录失败: %v\n", err)
		return l
	}

	// 创建LogWriter
	writer, err := newLogWriter(path, maxSize, maxBackups, maxAge)
	if err != nil {
		fmt.Printf("创建日志写入器失败: %v\n", err)
		return l
	}

	// 关闭旧的资源
	if l.file != nil {
		l.file.Close()
	}

	l.writer = writer
	l.output = io.MultiWriter(os.Stdout, writer)

	return l
}

// SetLevel 设置日志级别
func (l *Logger) SetLevel(level Level) *Logger {
	l.mu.Lock()
	l.level = level
	l.mu.Unlock()
	return l
}

// SetCaller 设置是否显示调用信息
func (l *Logger) SetCaller(show bool) *Logger {
	l.mu.Lock()
	l.hasCaller = show
	l.mu.Unlock()
	return l
}

// ANSI颜色代码
const (
	colorReset  = "\033[0m"
	colorCyan   = "\033[36m" // 时间戳用天蓝色
	colorGreen  = "\033[32m" // Debug用绿色
	colorBlue   = "\033[34m" // Info用蓝色
	colorYellow = "\033[33m" // Warn用黄色
	colorRed    = "\033[31m" // Error用红色
	colorPurple = "\033[35m" // Fatal用紫色
)

var levelColors = map[Level]string{
	DebugLevel: colorGreen,
	InfoLevel:  colorBlue,
	WarnLevel:  colorYellow,
	ErrorLevel: colorRed,
	FatalLevel: colorPurple,
}

// FormatLog 格式化日志信息并返回带颜色和不带颜色的字符串
func (l *Logger) FormatLog(level Level, skip int, msg string, fields ...Field) (plainContent, colorContent string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")

	// 调用信息
	var caller string
	if l.hasCaller {
		_, file, line, ok := runtime.Caller(skip)
		if ok {
			caller = fmt.Sprintf("%s:%d", file, line)
		}
	}

	// 构建不带颜色的日志内容
	if l.hasCaller {
		plainContent = fmt.Sprintf("[%s] [%s] [%s] %s",
			timestamp,
			levelNames[level],
			caller,
			msg)
	} else {
		plainContent = fmt.Sprintf("[%s] [%s] %s",
			timestamp,
			levelNames[level],
			msg)
	}

	// 构建带颜色的日志内容
	if l.hasCaller {
		colorContent = fmt.Sprintf("[%s%s%s] %s%s%s [%s] %s",
			colorCyan, timestamp, colorReset,
			levelColors[level], levelNames[level], colorReset,
			caller,
			msg)
	} else {
		colorContent = fmt.Sprintf("[%s%s%s] %s%s%s %s",
			colorCyan, timestamp, colorReset,
			levelColors[level], levelNames[level], colorReset,
			msg)
	}

	// 添加字段
	allFields := append(l.fields, fields...)
	fieldStr := ""
	if len(allFields) > 0 {
		for _, field := range allFields {
			fieldStr += fmt.Sprintf(" %s=%v", field.Key, field.Value)
		}
	}

	plainContent += fieldStr
	colorContent += fieldStr

	return plainContent, colorContent
}

func (l *Logger) log(level Level, skip int, msg string, fields ...Field) {
	if level < l.level {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	plainContent, colorContent := l.FormatLog(level, skip, msg, fields...)

	// 分别输出到控制台和文件
	if l.writer != nil {
		fmt.Fprintln(l.writer, plainContent)
	}
	fmt.Fprintln(os.Stdout, colorContent)

	// 如果是Fatal级别，则退出程序
	if level == FatalLevel {
		os.Exit(1)
	}
}

// Debug 输出Debug级别日志
func (l *Logger) Debug(msg string, skip int, fields ...Field) {
	l.log(DebugLevel, skip, msg, fields...)
}

// Info 输出Info级别日志
func (l *Logger) Info(msg string, skip int, fields ...Field) {
	l.log(InfoLevel, skip, msg, fields...)
}

// Warn 输出Warn级别日志
func (l *Logger) Warn(msg string, skip int, fields ...Field) {
	l.log(WarnLevel, skip, msg, fields...)
}

// Error 输出Error级别日志
func (l *Logger) Error(msg string, skip int, fields ...Field) {
	l.log(ErrorLevel, skip, msg, fields...)
}

// Fatal 输出Fatal级别日志
func (l *Logger) Fatal(msg string, skip int, fields ...Field) {
	l.log(FatalLevel, skip, msg, fields...)
}
