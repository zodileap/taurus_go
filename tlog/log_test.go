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

func testSetLogger(logger *zap.Logger) {
	for _, data := range testLogData {
		logger.Debug(data.msg, data.fields...)
		logger.Info(data.msg, data.fields...)
		logger.Warn(data.msg, data.fields...)
		logger.Error(data.msg, data.fields...)
		logger.DPanic(data.msg, data.fields...)
	}
}

func TestCreateLogger(t *testing.T) {
	logger := CreateLogger("test", true)

	testSetLogger(logger)
}

func TestGetLogger(t *testing.T) {
	logger := GetLogger("test")

	testSetLogger(logger)
}

func TestSetOutputPath(t *testing.T) {
	logger := SetOutputPath("test", "./test.log", true)

	testSetLogger(logger)
}

func TestSetLevel(t *testing.T) {
	logger := SetLevel("test", zapcore.InfoLevel, true)
	testSetLogger(logger)
	logger = SetLevel("test", zapcore.InfoLevel, false)
	testSetLogger(logger)
}
