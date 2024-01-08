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

type Logger = zap.Logger

// 日志字段
type Field = zapcore.Field

type loggerInfo struct {
	logger         *Logger
	consoleEncoder zapcore.Encoder
	fileEncoder    zapcore.Encoder
	writerSyncer   zapcore.WriteSyncer
	level          zapcore.Level
}

/***************  创建、更新logger  ***************/

var loggers = make(map[string]loggerInfo)

// 默认输出到控制台
var defaultWriterSyncer zapcore.WriteSyncer = nil

// 默认日志控制台格式
var defaultConsoleEncoder zapcore.Encoder = getEncoder(true)

// 默认日志文件格式
var defaultFileEncoder zapcore.Encoder = getEncoder(false)

// 默认日志级别为debug
var defaultLevel zapcore.Level = zapcore.DebugLevel

// 获取日志对象
//
// 获取已有的日志对象，如果不存在则新建一个日志对象
//
// 参数:
//   - loggerName: 日志对象名称，用于匹配日志
func GetLogger(loggerName string) *Logger {
	loggerInfo, exits := loggers[loggerName]
	var logger *Logger
	if !exits {
		logger = CreateLogger(loggerName, true)
	} else {
		logger = loggerInfo.logger
	}
	return logger
}

// 创建日志对象
//
//	新建一个日志对象，如果已经存在则返回已有的日志对象
//
// 参数:
//   - loggerName: 日志对象名称，用于匹配日志
func CreateLogger(loggerName string, hasCaller bool) *Logger {
	// core := zapcore.NewCore(defaultEncoder, defaultWriterSyncer, defaultLevel)
	logger := setLogger(loggerName, defaultWriterSyncer, defaultConsoleEncoder, defaultFileEncoder, zapcore.DebugLevel, hasCaller)
	return logger
}

// 设置日志输出路径
//
// 如果日志对象不存在则新建一个日志对象，并对其进行初始化设置
// 如果日志对象已存在，则保留除输出路径外的其他设置
//
// 参数:
//   - loggerName: 日志对象名称，用于匹配日志
//   - path: 日志输出路径
//   - hasCaller: 是否输出调用者信息
func SetOutputPath(loggerName string, path string, hasCaller bool) *Logger {
	loggerInfo, ok := loggers[loggerName]
	var writerSyncer zapcore.WriteSyncer = getLogWriter(path)
	var consoleEncoder zapcore.Encoder
	var fileEncoder zapcore.Encoder
	var level zapcore.Level
	if !ok {
		consoleEncoder = defaultConsoleEncoder
		fileEncoder = defaultFileEncoder
		level = defaultLevel
	} else {
		consoleEncoder = loggerInfo.consoleEncoder
		fileEncoder = loggerInfo.fileEncoder
		level = loggerInfo.level
	}

	logger := setLogger(loggerName, writerSyncer, consoleEncoder, fileEncoder, level, hasCaller)

	return logger
}

// 设置日志级别
//
// 如果日志对象不存在则新建一个日志对象，并对其进行初始化设置，默认为info级别
// 如果日志对象已存在，则保留除日志级别外的其他设置
//
// 参数:
//   - loggerName: 日志对象名称，用于匹配日志
//   - level: 日志级别
//   - hasCaller: 是否输出调用者信息
func SetLevel(loggerName string, level zapcore.Level, hasCaller bool) *Logger {
	loggerInfo, ok := loggers[loggerName]
	var writerSyncer zapcore.WriteSyncer
	var consoleEncoder zapcore.Encoder
	var fileEncoder zapcore.Encoder
	if !ok {
		writerSyncer = defaultWriterSyncer
		consoleEncoder = defaultConsoleEncoder
		fileEncoder = defaultFileEncoder

	} else {
		writerSyncer = loggerInfo.writerSyncer
		consoleEncoder = loggerInfo.consoleEncoder
		fileEncoder = loggerInfo.fileEncoder
	}
	logger := setLogger(loggerName, writerSyncer, consoleEncoder, fileEncoder, level, hasCaller)

	return logger
}

// 设置输出格式
//
// 如果日志对象不存在则新建一个日志对象，并对其进行初始化设置
// 如果日志对象已存在，则保留除输出格式外的其他设置
//
// 参数:
//   - loggerName: 日志对象名称，用于匹配日志
//   - consoleEncoder: 控制台输出格式
//   - fileEncoder: 文件输出格式
//   - hasCaller: 是否输出调用者信息
func SetEncoder(loggerName string, consoleEncoder zapcore.Encoder, fileEncoder zapcore.Encoder, hasCaller bool) *Logger {
	loggerInfo, ok := loggers[loggerName]
	var writerSyncer zapcore.WriteSyncer
	var level zapcore.Level
	if !ok {
		writerSyncer = defaultWriterSyncer
		level = defaultLevel
	} else {
		writerSyncer = loggerInfo.writerSyncer
		level = loggerInfo.level
	}
	logger := setLogger(loggerName, writerSyncer, consoleEncoder, fileEncoder, level, hasCaller)

	return logger
}

// 设置Caller
func SetCaller(loggerName string, hasCaller bool) *Logger {
	loggerInfo, ok := loggers[loggerName]
	var writerSyncer zapcore.WriteSyncer
	var consoleEncoder zapcore.Encoder
	var fileEncoder zapcore.Encoder
	var level zapcore.Level
	if !ok {
		writerSyncer = defaultWriterSyncer
		consoleEncoder = defaultConsoleEncoder
		fileEncoder = defaultFileEncoder
		level = defaultLevel
	} else {
		writerSyncer = loggerInfo.writerSyncer
		consoleEncoder = loggerInfo.consoleEncoder
		fileEncoder = loggerInfo.fileEncoder
		level = loggerInfo.level
	}
	logger := setLogger(loggerName, writerSyncer, consoleEncoder, fileEncoder, level, hasCaller)
	return logger
}

/***************  定制化消息 ***************/

// 输出err日志
//
// 参数:
//   - loggerName: 日志对象名称，用于匹配日志
//   - tracers: 调用堆栈的文件位置信息
//   - err: 错误信息
//   - params: 调用参数
func LogErr(loggerName string, tracers []string, code int, err string, reason string, params ...interface{}) {
	logger := GetLogger(loggerName)
	tracersString := strings.Join(tracers, " ===>>> ")

	fields := []Field{
		String("tracers", tracersString),
		Int("code", code),
		String("error", err),
	}

	if reason != "" {
		fields = append(fields, String("reason", reason))
	}

	// 迭代不定参数并添加到日志字段中
	for i, param := range params {
		fields = append(fields, Any(fmt.Sprintf("params-%d", i), param))
	}

	logger.Error(loggerName, fields...)
}

// 打印输出日志
//
// 使用的loggerName为print
//
// 如果想要不显示Caller信息，可以使用SetCaller("print", false)
//
// 参数:
//   - msg: 日志内容
//
// - fields: Field类型的可变参数
//
// 示例:
//
//	tlog.Print("api",
//
//		tlog.Int("code", 321321),
//
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

func Printf(format string, args ...any) {
	// 获取调用者的信息
	_, file, line, ok := runtime.Caller(1) // 1 代表上一层的调用堆栈
	if !ok {
		file = "???"
		line = 0
	}
	console(fmt.Sprintf(format, args...), file, line)
}

// 打印输出日志
//
// 使用的loggerName为print
//
// 如果想要不显示Caller信息，可以使用SetCaller("print", false)
func PrintString(s ...string) {
	msg := strings.Join(s, " ")
	// 获取调用者的信息
	_, file, line, ok := runtime.Caller(1) // 1 代表上一层的调用堆栈
	if !ok {
		file = "???"
		line = 0
	}
	console(msg, file, line)
}

/***************  Field设置 ***************/

// 获得String类型的Field
func String(key string, val string) Field {
	return Field{Key: key, Type: zapcore.StringType, String: val}
}

func Strings(key string, ss []string) Field {
	return zap.Strings(key, ss)
}

// 获得Int64类型的Field
func Int64(key string, val int64) Field {
	return Field{Key: key, Type: zapcore.Int64Type, Integer: val}
}

// 获得Int类型的Field
func Int(key string, val int) Field {
	return Int64(key, int64(val))
}

// 获得Duration类型的Field
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Type: zapcore.DurationType, Integer: int64(val)}
}

func Time(key string, val time.Time) Field {
	return zap.Time(key, val)
}

func Any(key string, value interface{}) zapcore.Field {
	return zap.Any(key, value)
}

func Reflect(key string, val interface{}) Field {
	Print(val)
	Print(key)
	// return Field{Key: key, Type: zapcore.ReflectType, Interface: val}
	return zap.Reflect(key, val)
}

/***************  私有函数和方法  ***************/

func console(msg string, file string, line int, fs ...Field) {
	logger := GetLogger("print")
	allFields := make([]Field, 0)
	allFields = append(allFields, String("file", fmt.Sprintf("%s:%d", file, line)))

	// 将其他 fields 添加到 allFields 切片中
	allFields = append(allFields, fs...)
	logger.Debug(msg,
		allFields...)
}

// 设置日志对象
//
// 参数:
//   - loggerName: 日志对象名称，用于匹配日志
//   - writerSyncer: 日志输出路径
//   - consoleEncoder: 控制台输出格式
//   - fileEncoder: 文件输出格式
func setLogger(loggerName string, writerSyncer zapcore.WriteSyncer, consoleEncoder zapcore.Encoder, fileEncoder zapcore.Encoder, level zapcore.Level, hasCaller bool) *Logger {
	// NOTE:
	// 不论writerSyncer的设置。默认输出到控制台
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.Lock(os.Stdout), level)
	var l *zap.Logger
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
	loggers[loggerName] = loggerInfo{
		logger:         l,
		consoleEncoder: consoleEncoder,
		fileEncoder:    fileEncoder,
		writerSyncer:   writerSyncer,
		level:          level,
	}
	return l

}

// 获取日志编码格式
//
// 如果isConsole为true，则返回控制台输出格式;则返回文件输出格式.
//
// 参数:
//   - isConsole: 是否为控制台输出
func getEncoder(isConsole bool) zapcore.Encoder {
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

// 获取日志输出路径
//
// 参数:
//   - path: 日志输出路径
func getLogWriter(path string) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   path,
		MaxSize:    500, // MB
		MaxBackups: 3,   // backups
		MaxAge:     28,  // days
		Compress:   true,
	}
	return zapcore.AddSync(lumberJackLogger)
}
