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

// LogWriter 日志写入器，支持文件大小切割、按日期切割、自动压缩和备份清理
type LogWriter struct {
	filename    string    // 日志文件完整路径
	maxSize     int64     // 单个日志文件最大大小，超过此大小将触发切割(单位:字节)
	maxBackups  int       // 保留的备份文件最大数量，超过将删除最旧的备份
	maxAge      int       // 备份文件最大保留天数，超过此天数的文件将被删除
	maxDays     int       // 按日期切割的天数间隔，默认1天切割一次
	size        int64     // 当前日志文件已写入的字节数
	file        *os.File  // 当前打开的日志文件句柄
	lastRotate  time.Time // 上次切割时间
	currentDate string    // 当前日期(YYYY-MM-DD格式)
}

// newLogWriter 创建新的LogWriter
func newLogWriter(filename string, maxSize, maxBackups, maxAge int) (*LogWriter, error) {
	return newLogWriterWithDays(filename, maxSize, maxBackups, maxAge, 1)
}

// newLogWriterWithDays 创建指定日期间隔的LogWriter
func newLogWriterWithDays(filename string, maxSize, maxBackups, maxAge, maxDays int) (*LogWriter, error) {
	now := time.Now()
	writer := &LogWriter{
		filename:    filename,
		maxSize:     int64(maxSize) * 1024 * 1024, // 转换为字节
		maxBackups:  maxBackups,
		maxAge:      maxAge,
		maxDays:     maxDays,
		lastRotate:  now,
		currentDate: now.Format("2006-01-02"),
	}

	if err := writer.openFile(); err != nil {
		return nil, err
	}

	return writer, nil
}

// openFile 打开文件
func (w *LogWriter) openFile() error {
	info, err := os.Stat(w.filename)
	if err == nil {
		w.size = info.Size()
		// 根据文件修改时间更新currentDate
		w.currentDate = info.ModTime().Format("2006-01-02")
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
		// 创建目录
		dir := filepath.Dir(w.filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("创建日志目录失败: %v\n", err)
			return 0, err
		}
		// 文件不存在，重新打开
		if err := w.openFile(); err != nil {
			return 0, err
		}
		// 重置size为0
		w.size = 0
	}

	// 检查是否需要按日期切割
	if w.shouldRotateByDate() {
		w.rotate()
	}

	n, err = w.file.Write(p)
	w.size += int64(n)

	// 检查是否需要按大小切割
	if w.shouldRotateBySize() {
		w.rotate()
	}

	return n, err
}

// shouldRotateBySize 检查是否需要按大小切割
func (w *LogWriter) shouldRotateBySize() bool {
	return w.maxSize > 0 && w.size > w.maxSize
}

// shouldRotateByDate 检查是否需要按日期切割
func (w *LogWriter) shouldRotateByDate() bool {
	if w.maxDays <= 0 {
		return false
	}
	now := time.Now()
	currentFileDate, err := time.Parse("2006-01-02", w.currentDate)
	if err != nil {
		return false
	}

	// 只比较日期，忽略具体时间
	nowDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	fileDate := time.Date(currentFileDate.Year(), currentFileDate.Month(), currentFileDate.Day(), 0, 0, 0, 0, currentFileDate.Location())

	// 计算日期差异的天数
	daysDiff := int(nowDate.Sub(fileDate).Hours() / 24)
	return daysDiff >= w.maxDays
}

// rotate 执行日志切割
func (w *LogWriter) rotate() error {
	if err := w.file.Close(); err != nil {
		return err
	}

	// 获取当前日期
	timestamp := time.Now().Format("2006-01-02")

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

	// 更新状态
	now := time.Now()
	w.lastRotate = now
	w.currentDate = now.Format("2006-01-02")
	w.size = 0

	// 清理旧日志文件
	w.cleanup()

	return nil
}

// compress 压缩文件
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

// cleanup 清理旧日志文件
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
