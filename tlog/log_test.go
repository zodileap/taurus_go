package tlog

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type testLogStruct struct {
	msg    string
	fields []Field
}

var testLogData = []testLogStruct{
	{
		msg: "This is a debug message",
		fields: []Field{
			String("url", "http://example.com"),
			Int("attempt", 3),
			String("status", "ok"),
		},
	},
}

func testSetLogger(logger *Logger) {
	l := logger
	for _, data := range testLogData {
		l.Debug(data.msg, data.fields...)
		l.Info(data.msg, data.fields...)
		l.Warn(data.msg, data.fields...)
		l.Error(data.msg, data.fields...)
	}
}

// TestGetLogger 测试设置日志级别
func TestSetLevel(t *testing.T) {

	logger := Get("test").SetLevel(InfoLevel)
	testSetLogger(logger)
}

// TestGetLogger 测试文件输出
func TestFileOutput(t *testing.T) {
	tmpDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	logPath := filepath.Join(tmpDir, "test.log")

	logger := Get("test_file")
	logger.SetOutputPath(logPath, 1, 3, 1) // 1MB, 3 backups, 1 day

	for i := 0; i < 10; i++ {
		logger.Info("Test log message", Int("count", i))
	}
	// 验证文件是否创建
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		t.Errorf("Log file was not created: %v", err)
	}

	// 清理测试文件
	// os.Remove(logPath)
}

// TestLogRotation 测试日志切割
func TestLogRotation(t *testing.T) {
	tmpDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current directory: %v", err)
	}
	logPath := filepath.Join(tmpDir, "rotation_test.log")

	logger := Get("test_rotation")
	logger.SetOutputPath(logPath, 1, 2, 1) // 1MB, 2 backups, 1 day

	// 生成一个较大的可读文本
	line := "This is a test log line that will be repeated many times to create a large log file for testing rotation.\n"
	// 写入足够多的日志触发切割
	for i := 0; i < 10000; i++ { // 写入足够多的行以触发日志切割
		logger.Info(fmt.Sprintf("Log line %d: %s", i, line),
			Int("iteration", i),
			String("test", "rotation"),
			String("content", "readable text"))

		// 每1000行暂停一下，确保时间戳不同
		if i%1000 == 0 {
			time.Sleep(time.Millisecond * 100)
		}
	}

	// 检查备份文件
	files, err := filepath.Glob(logPath + ".*")
	if err != nil {
		t.Errorf("Failed to list backup files: %v", err)
	}

	// 验证备份文件数量
	if len(files) == 0 {
		t.Error("No backup files were created")
	}

	if len(files) > 2 {
		t.Errorf("Expected at most 2 backup files, got %d", len(files))
	}

	// 清理测试文件
	os.Remove(logPath)
	for _, f := range files {
		os.Remove(f)
	}
}

// TestFields 测试字段
func TestFields(t *testing.T) {
	logger := Get("test_fields")

	fields := []Field{
		String("str", "value"),
		Int("int", 123),
		Any("map", map[string]string{"key": "value"}),
	}

	logger.Info("Test fields", fields...)
}

// TestGlobalFunctions 测试全局函数
func TestGlobalFunctions(t *testing.T) {
	tests := []struct {
		name     string
		logFunc  func(string, string, ...Field)
		expected Level
	}{
		{"Debug", Debug, DebugLevel},
		{"Info", Info, InfoLevel},
		{"Warn", Warn, WarnLevel},
		{"Error", Error, ErrorLevel},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc("test_global", "Test message")
		})
	}
}

// BenchmarkLogger 测试日志性能
func BenchmarkLoggerInfo(b *testing.B) {
	logger := Get("benchmark")

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Benchmark message",
				String("key1", "value1"),
				Int("key2", 123))
		}
	})
}

// BenchmarkLoggerWithFields 测试带字段的日志性能
func BenchmarkLoggerWithFields(b *testing.B) {
	logger := Get("benchmark_fields")
	fields := []Field{
		String("str", "value"),
		Int("int", 123),
		Any("map", map[string]string{"key": "value"}),
	}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Benchmark message", fields...)
		}
	})
}

// BenchmarkFileOutput 测试文件输出性能
func BenchmarkFileOutput(b *testing.B) {
	tmpDir, err := os.Getwd()
	if err != nil {
		b.Fatalf("Failed to get current directory: %v", err)
	}
	logPath := filepath.Join(tmpDir, "benchmark.log")

	logger := Get("benchmark_file")
	logger.SetOutputPath(logPath, 100, 3, 1)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			logger.Info("Benchmark message to file")
		}
	})

	// 清理测试文件
	os.Remove(logPath)
}
