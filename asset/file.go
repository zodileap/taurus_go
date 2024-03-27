package asset

import (
	"bytes"
	"io"
	"os"
)

// FileExists 检查文件是否存在。
//
// Params:
//
//   - filePath: 文件路径。
//
// Returns:
//
//	0: 文件是否存在。
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// ReadFileToBuffer 读取文件内容到缓冲区。
//
// Params:
//
//   - filename: 文件名。
//
// Returns:
//
//   - 文件内容。
//   - 错误信息。
func ReadFileToBuffer(filePath string) (*bytes.Buffer, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		return nil, err
	}

	return &buffer, nil
}
