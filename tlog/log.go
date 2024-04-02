// Description: 日志模块
package tlog

/*
示例:
	tlog.SetOutputPath("api", "log/api/api.log")
	log := tlog.GetLogger("api")
	logger.Info("",
		tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")), // 请求时间                         // 请求的ip              // 请求的方法
		tlog.Int("code", 321321),                        // 状态码
	)

	tlog.Print("api",
		tlog.Int("code", 321321),
	)

*/

import (
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Field 日志字段
type Field = zapcore.Field

// String 生成一个字符串类型的Field。
//
// Params:
//
//   - key: 字段名
//   - val: 字符串值
func String(key string, val string) Field {
	return Field{Key: key, Type: zapcore.StringType, String: val}
}

// Strings 生成一个字符串切片类型的Field。
//
// Params:
//
//   - key: 字段名
//   - ss: 字符串切片
func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

// Int64 生成一个Int64类型的Field。
//
// Params:
//
//   - key: 字段名
//   - val: int64值
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: zapcore.Int64Type, Integer: val}
}

// Int 生成一个Int类型的Field。
//
// Params:
//
//   - key: 字段名
//   - val: int值
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// Duration 生成一个Duration类型的Field。
//
// Params:
//
//   - key: 字段名
//   - val: time.Duration值
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: zapcore.DurationType, Integer: int64(val)}
}

// Time 生成一个Time类型的Field。
//
// Params:
//
//   - key: 字段名
//   - val: time.Time值
func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

// Any 生成一个Any类型的Field。
//
// Params:
//
//   - key: 字段名
//   - value: 任意类型的值
func Any(key string, value interface{}) zapcore.Field {
	return zap.Any(key, value)
}

// Reflect 生成一个Reflect类型的Field。
func Reflect(key string, val interface{}) Field {
	Print(val)
	Print(key)
	// return Field{Key: key, Type: zapcore.ReflectType, Interface: val}
	return zap.Reflect(key, val)
}

// Level 日志级别
type Level = zapcore.Level

const (
	// DebugLevel debug级别
	DebugLevel = zapcore.DebugLevel
	// InfoLevel info级别
	InfoLevel = zapcore.InfoLevel
	// WarnLevel warn级别
	WarnLevel = zapcore.WarnLevel
	// ErrorLevel error级别
	ErrorLevel = zapcore.ErrorLevel
	// DPanicLevel dpanic级别
	DPanicLevel = zapcore.DPanicLevel
	// PanicLevel panic级别
	PanicLevel = zapcore.PanicLevel
	// FatalLevel fatal级别
	FatalLevel = zapcore.FatalLevel
)

// Encoder 日志编码格式
type Encoder = zapcore.Encoder

// WriteSyncer 日志输出路径
type WriteSyncer = zapcore.WriteSyncer

// Logger 日志对象
type Logger struct {
	logger         *zap.Logger
	name           string
	consoleEncoder Encoder
	fileEncoder    Encoder
	writerSyncer   WriteSyncer
	level          Level
}

var loggers = make(map[string]*Logger)

// 默认输出到控制台
var defaultWriterSyncer WriteSyncer = nil

// 默认日志控制台格式
var defaultConsoleEncoder Encoder = getEncoder(true)

// 默认日志文件格式
var defaultFileEncoder Encoder = getEncoder(false)

// 默认日志级别为debug
var defaultLevel zapcore.Level = zapcore.DebugLevel

// GetLogger 获取日志对象，如果获取的对象不存在则新建一个日志对象并返回。
//
// 创建的日志对象默认只输出到控制条，日志级别为Debug。
//
// Params:
//
//   - loggerName: 日志对象名称，用于匹配日志
//
// Returns:
//
//	0: 日志对象。
//
// Example:
//
//	logger := tlog.GetLogger("api")
func GetLogger(loggerName string) *Logger {
	logger, exits := loggers[loggerName]
	if !exits {
		logger = createLogger(loggerName, true)
	}
	return logger
}

// SetOutputPath 设置日志输出路径，保留除输出路径外的其他设置。
//
// Params:
//
//   - path: 日志输出路径
//   - maxSize: 日志文件大小，单位MB。如果文件超过这个大小会被切割。默认100MB。
//   - maxBackups: 需要保留的旧日志文件数，默认保留所有旧日志文件，但是maxAge参数还是会导致文件被删除。
//   - maxAge: 日志文件最大保存天数。
//
// Returns:
//
//	0: 日志对象。
//
// Example:
//
//	logger := tlog.GetLogger("api").SetOutputPath("log/api/api.log")
func (l *Logger) SetOutputPath(path string, maxSize int, maxBackups int, maxAge int) *Logger {
	if maxSize == 0 {
		maxSize = 100
	}
	if maxAge == 0 {
		maxAge = 7
	}
	writerSyncer := getLogWriter(path, maxSize, maxBackups, maxAge)
	return setLogger(l.name, writerSyncer, l.consoleEncoder, l.fileEncoder, l.level, true)
}

// SetLevel 设置日志级别，保留除日志级别外的其他设置。
//
// Params:
//
//   - level: 日志级别
//
// Returns:
//
//	0: 日志对象。
//
// Example:
//
//	logger := tlog.GetLogger("api").SetLevel(tlog.InfoLevel)
func (l *Logger) SetLevel(level Level) *Logger {
	return setLogger(l.name, l.writerSyncer, l.consoleEncoder, l.fileEncoder, level, true)
}

// SetEncoder 设置输出格式，保留除输出格式外的其他设置。
//
// Params:
//
//   - consoleEncoder: 控制台输出格式
//   - fileEncoder: 文件输出格式
//
// Returns:
//
//	0: 日志对象。
//
// Example:
//
//	 consoleEncoder := zap.NewProductionEncoderConfig()
//	 consoleEncoder.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
//			enc.AppendString("[\033[36m" + t.Format("2006-01-02 15:04:05.000") + "\033[0m]")
//	 }
//	 fileEncoder := zap.NewProductionEncoderConfig()
//	 fileEncoder.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
//	 logger := tlog.GetLogger("api").SetEncoder(consoleEncoder, fileEncoder)
func (l *Logger) SetEncoder(consoleEncoder Encoder, fileEncoder Encoder) *Logger {
	return setLogger(l.name, l.writerSyncer, consoleEncoder, fileEncoder, l.level, true)
}

// SetCaller 设置是否输出Caller信息，保留除输出Caller信息外的其他设置。
//
// Params:
//
//   - hasCaller: 是否输出Caller信息
//
// Returns:
//
//	0: 日志对象。
//
// Example:
//
//	logger := tlog.GetLogger("api").SetCaller(false)
func (l *Logger) SetCaller(hasCaller bool) *Logger {
	return setLogger(l.name, l.writerSyncer, l.consoleEncoder, l.fileEncoder, l.level, hasCaller)
}

// Debug 记录一个Debug级别的日志。
//
// Params:
//
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
func (l *Logger) Debug(msg string, fields ...Field) {
	l.logger.Debug(msg, fields...)
}

// Info 记录一个Info级别的日志。
//
// Params:
//
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
func (l *Logger) Info(msg string, fields ...Field) {
	l.logger.Info(msg, fields...)
}

// Warn 记录一个Warn级别的日志。
//
// Params:
//
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
func (l *Logger) Warn(msg string, fields ...Field) {
	l.logger.Warn(msg, fields...)
}

// Error 记录一个Error级别的日志。
//
// Params:
//
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
func (l *Logger) Error(msg string, fields ...Field) {
	l.logger.Error(msg, fields...)
}

// DPanic 记录一个DPanic级别的日志。
//
// Params:
//
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
func (l *Logger) DPanic(msg string, fields ...Field) {
	l.logger.DPanic(msg, fields...)
}

// Panic 记录一个Panic级别的日志。
//
// Params:
//
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
func (l *Logger) Panic(msg string, fields ...Field) {
	l.logger.Panic(msg, fields...)
}

// Fatal 记录一个Fatal级别的日志。
//
// Params:
//
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
func (l *Logger) Fatal(msg string, fields ...Field) {
	l.logger.Fatal(msg, fields...)
}

// Debug 记录一个Debug级别的日志。
//
// Params:
//
//   - loggerName: 日志对象名称。
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
//
// Example:
//
//	tlog.Debug("debug",  tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")))
func Debug(loggerName string, msg string, fields ...Field) {
	logger := GetLogger(loggerName)
	logger.Debug(msg, fields...)
}

// Info 记录一个Info级别的日志。
//
// Params:
//
//   - loggerName: 日志对象名称。
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
//
// Example:
//
//	tlog.Info("info",  tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")))
func Info(loggerName string, msg string, fields ...Field) {
	logger := GetLogger(loggerName)
	logger.Info(msg, fields...)
}

// Warn 记录一个Warn级别的日志。
//
// Params:
//
//   - loggerName: 日志对象名称。
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
//
// Example:
//
//	tlog.Warn("warn",  tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")))
func Warn(loggerName string, msg string, fields ...Field) {
	logger := GetLogger(loggerName)
	logger.Warn(msg, fields...)
}

// Error 记录一个Error级别的日志。
//
// Params:
//
//   - loggerName: 日志对象名称。
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
//
// Example:
//
//	tlog.Error("error",  tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")))
func Error(loggerName string, msg string, fields ...Field) {
	logger := GetLogger(loggerName)
	logger.Error(msg, fields...)
}

// DPanic 记录一个DPanic级别的日志。
//
// Params:
//
//   - loggerName: 日志对象名称。
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
//
// Example:
//
//	tlog.DPanic("dpanic",  tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")))
func DPanic(loggerName string, msg string, fields ...Field) {
	logger := GetLogger(loggerName)
	logger.DPanic(msg, fields...)
}

// Panic 记录一个Panic级别的日志。
//
// Params:
//
//   - loggerName: 日志对象名称。
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
//
// Example:
//
//	tlog.Panic("panic",  tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")))
func Panic(loggerName string, msg string, fields ...Field) {
	logger := GetLogger(loggerName)
	logger.Panic(msg, fields...)
}

// Fatal 记录一个Fatal级别的日志。
//
// Params:
//
//   - loggerName: 日志对象名称。
//   - msg: 日志内容。
//   - fields: Field类型的可变参数。
//
// Example:
//
//	tlog.Fatal("fatal",  tlog.String("rquestTime", t.Format("2006-01-02 15:04:05")))
func Fatal(loggerName string, msg string, fields ...Field) {
	logger := GetLogger(loggerName)
	logger.Fatal(msg, fields...)
}

// Print 打印输出日志。使用的loggerName为print,print的日志级别为Debug,默认只输出到控制台。
//
// 如果想要不显示Caller信息，可以使用tlog.GetLogger("print").SetCaller(false)。
//
// 如果需要记录到文件，可以使用tlog.GetLogger("print").SetOutputPath("log/print.log")。
//
// Params:
//
//   - msg: 日志内容
//   - fields: Field类型的可变参数
//
// Example:
//
//	tlog.Print("api",
//		tlog.Int("code", 321321),
//	)
func Print(msg any, fields ...Field) {
	// 获取调用者的信息
	_, file, line, ok := runtime.Caller(1) // 1 代表上一层的调用堆栈
	if !ok {
		file = "???"
		line = 0
	}
	console(fmt.Sprint(msg), file, line, fields...)
}

// Printf 格式化输出日志。使用的loggerName为print,print的日志级别为Debug,默认只输出到控制台。
//
// 如果想要不显示Caller信息，可以使用tlog.GetLogger("print").SetCaller(false)。
//
// 如果需要记录到文件，可以使用tlog.GetLogger("print").SetOutputPath("log/print.log")。
//
// Params:
//
//   - format: 格式化字符串
//   - args: 格式化参数
//
// Example:
//
//	tlog.Printf("format: %s", "test")
func Printf(format string, args ...any) {
	// 获取调用者的信息
	_, file, line, ok := runtime.Caller(1) // 1 代表上一层的调用堆栈
	if !ok {
		file = "???"
		line = 0
	}
	console(fmt.Sprintf(format, args...), file, line)
}

// PrintString 打印字符串。使用的loggerName为print,print的日志级别为Debug,默认只输出到控制台。
//
// 如果想要不显示Caller信息，可以使用tlog.GetLogger("print").SetCaller(false)。
//
// 如果需要记录到文件，可以使用tlog.GetLogger("print").SetOutputPath("log/print.log")。
//
// Params:
//
//   - s: 字符串
//
// Example:
//
//	tlog.PrintString("print", "a", "message")
func PrintStrings(s ...string) {
	msg := strings.Join(s, " ")
	// 获取调用者的信息
	_, file, line, ok := runtime.Caller(1) // 1 代表上一层的调用堆栈
	if !ok {
		file = "???"
		line = 0
	}
	console(msg, file, line)
}

// console 输出到控制台
//
// Params:
//
//   - msg: 日志内容
//   - file: 文件名
//   - line: 行号
//   - fs: Field类型的可变参数
func console(msg string, file string, line int, fs ...Field) {
	logger := GetLogger("print")
	allFields := make([]Field, 0)
	allFields = append(allFields, String("file", fmt.Sprintf("%s:%d", file, line)))

	// 将其他 fields 添加到 allFields 切片中
	allFields = append(allFields, fs...)
	logger.logger.Debug(msg, allFields...)
}

// createLogger 创建日志对象，如果已经存在则返回已有的日志对象
//
// Params:
//   - loggerName: 日志对象名称，用于匹配日志
func createLogger(loggerName string, hasCaller bool) *Logger {
	logger := setLogger(loggerName, defaultWriterSyncer, defaultConsoleEncoder, defaultFileEncoder, zapcore.DebugLevel, hasCaller)
	return logger
}

// setLogger 设置日志对象
//
// Params:
//
//   - loggerName: 日志对象名称，用于匹配日志
//   - writerSyncer: 日志输出路径
//   - consoleEncoder: 控制台输出格式
//   - fileEncoder: 文件输出格式
//   - level: 日志级别
//   - hasCaller: 是否输出调用者信息
func setLogger(loggerName string, writerSyncer WriteSyncer, consoleEncoder Encoder, fileEncoder Encoder, level Level, hasCaller bool) *Logger {
	// NOTE:
	// 不论writerSyncer的设置。默认输出到控制台
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
	var l *zap.Logger
	// 如果writerSyncer为nil，则只输出到控制台
	// 否则同时输出到控制台和文件
	if writerSyncer == nil {
		if hasCaller {
			l = zap.New(zapcore.NewTee(consoleCore), zap.AddCaller())
		} else {
			l = zap.New(zapcore.NewTee(consoleCore))
		}
	} else {
		fileCore := zapcore.NewCore(fileEncoder, writerSyncer, level)
		if hasCaller {
			l = zap.New(zapcore.NewTee(consoleCore, fileCore), zap.AddCaller())
		} else {
			l = zap.New(zapcore.NewTee(consoleCore, fileCore))
		}

	}

	logger := &Logger{
		name:           loggerName,
		logger:         l,
		consoleEncoder: consoleEncoder,
		fileEncoder:    fileEncoder,
		writerSyncer:   writerSyncer,
		level:          level,
	}
	loggers[loggerName] = logger
	return logger

}

// getEncoder 获取日志编码格式
//
// 如果isConsole为true，则返回控制台输出格式;则返回文件输出格式.
//
// Params:
//
//   - isConsole: 是否为控制台输出
func getEncoder(isConsole bool) Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	encoderConfig.CallerKey = "caller"
	if isConsole {
		// 控制台输出
		encoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString("[\033[36m" + t.Format("2006-01-02 15:04:05.000") + "\033[0m]")
		}
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

		return zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// 文件输出
		encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
		return zapcore.NewJSONEncoder(encoderConfig)
	}

}

// getLogWriter 把传入的路径实现为一个WriteSyncer。
//
// Params:
//
//   - path: 日志输出路径
//   - maxSize: 日志文件大小，单位MB。如果文件超过这个大小会被切割。默认100MB。
//   - maxBackups: 需要保留的旧日志文件数，默认保留所有旧日志文件，但是maxAge参数还是会导致文件被删除。
//   - maxAge: 日志文件最大保存天数。
func getLogWriter(path string, maxSize int, maxBackups int, maxAge int) WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    maxSize,    // MB
		MaxBackups: maxBackups, // backups
		MaxAge:     maxAge,     // days
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
