package tlog

import (
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"time"
)

type testLogStruct struct {
	msg    string
	fields []zapcore.Field
}

var testLogData = []testLogStruct{
	{
		msg: "This text will appear red.",
		fields: []zapcore.Field{
			zap.String("url", "http://example.com"),
			zap.Int("attempt", 3),
			zap.Duration("backoff", time.Second),
		},
	},
}

func testSetLogger(logger *Logger) {
	l := logger.logger
	for _, data := range testLogData {
		l.Debug(data.msg, data.fields...)
		l.Info(data.msg, data.fields...)
		l.Warn(data.msg, data.fields...)
		l.Error(data.msg, data.fields...)
		l.DPanic(data.msg, data.fields...)
	}
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger("test")

	testSetLogger(logger)
}

func TestSetOutputPath(t *testing.T) {
	logger := GetLogger("test").SetOutputPath("./test.log", 100, 3, 14)

	testSetLogger(logger)
}

func TestSetLevel(t *testing.T) {
	logger := GetLogger("test").SetLevel(InfoLevel)
	testSetLogger(logger)
	logger = GetLogger("test").SetLevel(InfoLevel)
	testSetLogger(logger)
}
