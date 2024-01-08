package asset

import (
	"io"
	"os"
	"path/filepath"
)

// CopyFile 复制文件到目标路径。
//
// Params:
//   - src: 源文件的路径。
//   - dst: 目标文件的路径。
//
// Example:
//
// err := asset.CopyFile("file.go", "file2.go")
//
//	if err != nil {
//		fmt.Print(err)
//	}
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go
//
// ErrCodes:
//   - Err_0200010004
//   - Err_0200010005
//   - Err_0200010006
func CopyFile(src string, dst string) error {
	_, err := os.Stat(src)
	if os.IsNotExist(err) {
		return nil
	}
	sourceFile, err := os.Open(src)
	if err != nil {
		return Err_0200010004.Sprintf(src, err)
	}
	defer sourceFile.Close()
	destFile, err := os.Create(dst)
	if err != nil {
		return Err_0200010005.Sprintf(dst, err)
	}
	defer destFile.Close()
	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return Err_0200010006.Sprintf(src, dst, err)
	}
	return nil
}

// CopyDir 复制整个目录
//
// Params:
//   - src: 源文件夹的路径.
//     default: 111
//   - dst: 目标文件夹的路径
//
// Example:
//
// err := asset.CopyDir("dir", "dir2")
//
//	if err != nil {
//		fmt.Print(err)
//	}
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go
//
// ErrCodes:
//   - Err_0200010001
//   - Err_0200010004
//   - Err_0200010005
//   - Err_0200010006
//   - Err_0200010007
//   - Err_0200010008
func CopyDir(src string, dst string) error {
	// 检查源目录是否存在
	srcInfo, err := os.Stat(src)
	if os.IsNotExist(err) {
		return nil // 源目录不存在，跳过
	}

	err = filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return Err_0200010007.Sprintf(path, err)
		}

		// 计算目标文件或目录的路径
		relPath, err := filepath.Rel(src, path)
		if err != nil {
			return Err_0200010008.Sprintf(src, path, err)
		}
		targetPath := filepath.Join(dst, relPath)

		// 如果是目录，则创建目录
		if info.IsDir() {
			err = os.MkdirAll(targetPath, srcInfo.Mode())
			if err != nil {
				return Err_0200010001.Sprintf(targetPath, err)
			}
		} else {
			// 如果是文件，则复制文件
			return CopyFile(path, targetPath)
		}
		return nil
	})

	return err
}
