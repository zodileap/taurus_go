package stringutil

import (
	"errors"
	"os"
	"path/filepath"
	"regexp"

	"github.com/google/uuid"
)

// 判断字符串是否是UUID
func IsUUID(u string) error {
	_, err := uuid.Parse(u)
	if err != nil {
		return errors.New(Text_err_not_uuid)
	}
	return nil
}

// 判断字符串是否是Email
func IsEmail(email string) error {
	const emailRegex = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	b := re.MatchString(email)
	if !b {
		return errors.New(Text_err_not_email)
	}
	return nil
}

func IsDir(path string) (bool, error) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false, err // 无法获取文件或目录信息
	}
	return fileInfo.IsDir(), nil
}

// 判断字符串是否有文件的后缀
// 如果没有返回空字符串
func HasSuffix(filename string) string {
	// 获取文件扩展名
	ext := filepath.Ext(filename)
	// 如果扩展名不为空，则表示文件名带有后缀
	return ext
}

func Container(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
