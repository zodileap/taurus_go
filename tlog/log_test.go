package tlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSetLevel(t *testing.T) {
	logger := Get("test-set-level").SetLevel(InfoLevel)
	if logger.level != InfoLevel {
		t.Fatalf("日志级别设置失败，期望 %d，实际 %d", InfoLevel, logger.level)
	}
}

func TestFormatLog(t *testing.T) {
	logger := Get("test-format-log").SetCaller(false)

	plain, color := logger.FormatLog(InfoLevel, 0, "hello", String("key", "value"), Int("count", 2))
	if !strings.Contains(plain, "[INFO]") || !strings.Contains(plain, "hello") {
		t.Fatalf("普通日志内容不完整: %s", plain)
	}
	if !strings.Contains(plain, "key=value") || !strings.Contains(plain, "count=2") {
		t.Fatalf("普通日志字段缺失: %s", plain)
	}
	if !strings.Contains(color, "hello") || !strings.Contains(color, "key=value") {
		t.Fatalf("彩色日志内容不完整: %s", color)
	}
}

func TestFileOutput(t *testing.T) {
	logPath := filepath.Join(t.TempDir(), "test.log")

	logger := Get("test-file-output").SetCaller(false)
	logger.SetOutputPath(logPath, 1, 3, 1)
	logger.Info("Test log message", 0, Int("count", 1))

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}
	text := string(content)
	if !strings.Contains(text, "Test log message") || !strings.Contains(text, "count=1") {
		t.Fatalf("日志文件内容不正确: %s", text)
	}
}

func TestGlobalFunctions(t *testing.T) {
	tests := []struct {
		name    string
		logFunc func(string, string, ...Field)
	}{
		{"Debug", Debug},
		{"Info", Info},
		{"Warn", Warn},
		{"Error", Error},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.logFunc("test-global-functions", "test message", String("from", tt.name))
		})
	}
}
