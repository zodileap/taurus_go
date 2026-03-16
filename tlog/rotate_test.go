package tlog

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestRotateByDate(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test-date.log")

	writer, err := newLogWriterWithDays(logFile, 10, 3, 7, 1)
	if err != nil {
		t.Fatalf("创建 LogWriter 失败: %v", err)
	}
	defer writer.file.Close()

	if _, err := writer.Write([]byte("first line\n")); err != nil {
		t.Fatalf("首次写入失败: %v", err)
	}

	writer.currentDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	if _, err := writer.Write([]byte("second line\n")); err != nil {
		t.Fatalf("日期轮转写入失败: %v", err)
	}

	backups, err := filepath.Glob(logFile + ".*.gz")
	if err != nil {
		t.Fatalf("查找备份文件失败: %v", err)
	}
	if len(backups) == 0 {
		t.Fatal("未生成日期轮转备份文件")
	}
}

func TestRotateBySize(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "test-size.log")

	writer, err := newLogWriterWithDays(logFile, 0, 3, 7, 1)
	if err != nil {
		t.Fatalf("创建 LogWriter 失败: %v", err)
	}
	defer writer.file.Close()

	writer.maxSize = 10

	if _, err := writer.Write([]byte("this line should trigger rotation\n")); err != nil {
		t.Fatalf("大小轮转写入失败: %v", err)
	}

	backups, err := filepath.Glob(logFile + ".*.gz")
	if err != nil {
		t.Fatalf("查找备份文件失败: %v", err)
	}
	if len(backups) == 0 {
		t.Fatal("未生成大小轮转备份文件")
	}
}

func TestLoggerSetOutputPathWithDays(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "logger-days.log")

	logger := Get("test-logger-days").SetCaller(false)
	logger.SetOutputPathWithDays(logFile, 100, 3, 7, 2)
	logger.Info("按日期切割日志", 0)

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}
	if !strings.Contains(string(content), "按日期切割日志") {
		t.Fatalf("日志文件未写入预期内容: %s", string(content))
	}
}

func TestLoggerDefaultDays(t *testing.T) {
	logFile := filepath.Join(t.TempDir(), "logger-default.log")

	logger := Get("test-logger-default").SetCaller(false)
	logger.SetOutputPath(logFile, 100, 3, 7)
	logger.Info("默认日期切割日志", 0)

	content, err := os.ReadFile(logFile)
	if err != nil {
		t.Fatalf("读取日志文件失败: %v", err)
	}
	if !strings.Contains(string(content), "默认日期切割日志") {
		t.Fatalf("日志文件未写入预期内容: %s", string(content))
	}
}
