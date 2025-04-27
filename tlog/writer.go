package tlog

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// LogWriter 日志写入器
type LogWriter struct {
	filename   string
	maxSize    int64
	maxBackups int
	maxAge     int
	size       int64
	file       *os.File
}

// 创建新的LogWriter
func newLogWriter(filename string, maxSize, maxBackups, maxAge int) (*LogWriter, error) {
	writer := &LogWriter{
		filename:   filename,
		maxSize:    int64(maxSize) * 1024 * 1024, // 转换为字节
		maxBackups: maxBackups,
		maxAge:     maxAge,
	}

	if err := writer.openFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

// 打开文件
func (w *LogWriter) openFile() error {
	info, err := os.Stat(w.filename)
	if err == nil {
		w.size = info.Size()
	}

	file, err := os.OpenFile(w.filename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	w.file = file
	return nil
}

// Write 实现io.Writer接口
func (w *LogWriter) Write(p []byte) (n int, err error) {
	// 检查文件是否存在
	if _, err := os.Stat(w.filename); os.IsNotExist(err) {
		// 文件不存在，重新打开
		if err := w.openFile(); err != nil {
			return 0, err
		}
		// 重置size为0
		w.size = 0
	}

	n, err = w.file.Write(p)
	w.size += int64(n)

	if w.size > w.maxSize {
		w.rotate()
	}

	return n, err
}

// 执行日志切割
func (w *LogWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	// 获取当前时间戳，修改为毫秒级
	timestamp := time.Now().Format("2006-01-02-150405")

	// 新的备份文件名
	backupName := fmt.Sprintf("%s.%s.gz", w.filename, timestamp)

	// 压缩当前日志文件
	if err := w.compress(w.filename, backupName); err != nil {
		return err
	}

	// 清空当前日志文件
	if err := os.Truncate(w.filename, 0); err != nil {
		return err
	}

	// 重新打开文件
	if err := w.openFile(); err != nil {
		return err
	}

	// 清理旧日志文件
	w.cleanup()

	return nil
}

// 压缩文件
func (w *LogWriter) compress(src, dst string) error {
	f, err := os.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	gzf, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer gzf.Close()

	gz := gzip.NewWriter(gzf)
	defer gz.Close()

	if _, err := io.Copy(gz, f); err != nil {
		return err
	}
	return nil
}

// 清理旧日志文件
func (w *LogWriter) cleanup() {
	dir := filepath.Dir(w.filename)
	base := filepath.Base(w.filename)

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	var backups []string
	for _, f := range files {
		if strings.HasPrefix(f.Name(), base+".") && strings.HasSuffix(f.Name(), ".gz") {
			backups = append(backups, filepath.Join(dir, f.Name()))
		}
	}

	// 按时间排序
	sort.Slice(backups, func(i, j int) bool {
		return backups[i] > backups[j]
	})

	// 删除超过maxBackups的文件
	if w.maxBackups > 0 && len(backups) > w.maxBackups {
		for _, f := range backups[w.maxBackups:] {
			os.Remove(f)
		}
	}

	// 删除超过maxAge天数的文件
	if w.maxAge > 0 {
		cutoff := time.Now().AddDate(0, 0, -w.maxAge)
		for _, f := range backups {
			info, err := os.Stat(f)
			if err != nil {
				continue
			}
			if info.ModTime().Before(cutoff) {
				os.Remove(f)
			}
		}
	}
}
