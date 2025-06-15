package tlog

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestRotateByDate(t *testing.T) {
	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "tlog_test")
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test_date.log")

	// 创建按日期切割的LogWriter
	writer, err := newLogWriterWithDays(logFile, 10, 3, 7, 1) // 1天切割一次
	if err != nil {
		t.Fatalf("创建LogWriter失败: %v", err)
	}
	defer writer.file.Close()

	// 写入一些日志
	_, err = writer.Write([]byte("第一条日志消息\n"))
	if err != nil {
		t.Fatalf("写入日志失败: %v", err)
	}

	// 模拟日期变化，设置为1天前
	writer.currentDate = time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	// 再写入日志，应该触发切割
	_, err = writer.Write([]byte("第二条日志消息\n"))
	if err != nil {
		t.Fatalf("写入日志失败: %v", err)
	}

	// 检查是否生成了备份文件
	dir := filepath.Dir(logFile)
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("读取目录失败: %v", err)
	}

	hasBackup := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gz" {
			hasBackup = true
			break
		}
	}

	if !hasBackup {
		t.Error("没有找到备份文件")
	}
}

func TestRotateBoth(t *testing.T) {
	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "tlog_test")
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test_both.log")

	// 创建同时按大小和日期切割的LogWriter
	writer, err := newLogWriterWithDays(logFile, 0, 3, 7, 1) // maxSize = 0MB, 1天切割
	if err != nil {
		t.Fatalf("创建LogWriter失败: %v", err)
	}
	defer writer.file.Close()

	// 设置一个很小的maxSize来测试大小切割
	writer.maxSize = 10

	// 写入一些日志，超过大小限制
	_, err = writer.Write([]byte("这是一条很长的日志消息，应该会触发大小切割\n"))
	if err != nil {
		t.Fatalf("写入日志失败: %v", err)
	}

	// 检查是否生成了备份文件
	dir := filepath.Dir(logFile)
	files, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("读取目录失败: %v", err)
	}

	hasBackup := false
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".gz" {
			hasBackup = true
			break
		}
	}

	if !hasBackup {
		t.Error("没有找到备份文件")
	}
}

func TestLoggerSetOutputPathWithDays(t *testing.T) {
	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "tlog_test")
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test_logger_days.log")

	// 创建logger并设置按日期切割
	logger := Get("test_days")
	logger.SetOutputPathWithDays(logFile, 100, 3, 7, 2) // 100MB, 2天切割一次

	// 写入日志
	logger.Info("测试按日期切割的日志", 2)

	// 检查文件是否创建
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("日志文件未创建")
	}
}

func TestLoggerDefaultDays(t *testing.T) {
	// 创建临时目录
	tmpDir := filepath.Join(os.TempDir(), "tlog_test")
	defer os.RemoveAll(tmpDir)

	logFile := filepath.Join(tmpDir, "test_logger_default.log")

	// 创建logger并使用默认设置(默认1天切割)
	logger := Get("test_default")
	logger.SetOutputPath(logFile, 100, 3, 7) // 100MB, 默认1天切割

	// 写入日志
	logger.Info("测试默认日期切割的日志", 2)

	// 检查文件是否创建
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		t.Error("日志文件未创建")
	}
}
