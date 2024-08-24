package asset

import (
	"os"
	"path/filepath"

	"golang.org/x/tools/imports"
)

// Assets 用于存放需要创建的文件夹和文件
type Assets struct {
	// Dirs 用于存放需要创建的文件夹
	Dirs map[string]struct{}
	// Files 用于存放需要创建的文件
	Files map[string][]byte
}

// Add 用于在`Assets`中添加一个新的文件。
//
// Params:
//
//   - path: 文件路径。
//   - b: 文件内容。
//
// Example:
//
//	var assets asset.Assets
//
//	assets.Add("file.txt", []byte("Hellow, World!"))
//	assets.Add("../file_2.go", []byte(`
//	package main
//	import "fmt"
//	func Main() {fmt.Println("Hello, World!")}
//	`))
//	err := assets.Write()
//	if err != nil {
//		fmt.Print(err)
//	}
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go - TestAdd
func (a *Assets) Add(path string, b []byte) {
	if a.Files == nil {
		a.Files = make(map[string][]byte)
	}
	a.Files[path] = b
}

// AddDir 用于在`Assets`中添加一个新的文件夹。
//
// Params:
//
//   - path: 文件夹路径。
//
// Returns:
//
//	0: 成功。
//
// Example:
//
//	var assets asset.Assets
//	assets.AddDir("dir")
//	assets.AddDir("./dir2")
//	err := assets.Write()
//
//	if err != nil {
//	 fmt.Print(err)
//	}
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go - TestAddDir
func (a *Assets) AddDir(path string) error {
	if a.Dirs == nil {
		a.Dirs = make(map[string]struct{})
	}
	a.Dirs[path] = struct{}{}
	return nil
}

// Write 写入全部的Dirs和Files,如果文件已经存在，则会覆盖。
// 执行后会清空Dirs和Files。
//
// Example:
//
// var assets asset.Assets
// assets.Add("file.txt", []byte("Hellow, World!"))
// assets.AddDir("dir")
// err := assets.Write()
//
//	if err != nil {
//		fmt.Print(err)
//	}
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go - TestWrite
//
// ErrCodes:
//   - Err_0200010001
//   - Err_0200010002
func (a Assets) Write() error {
	for dir := range a.Dirs {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return Err_0200030001.Sprintf(dir, err)
		}
	}
	for path, content := range a.Files {
		if err := os.WriteFile(path, content, 0644); err != nil {
			return Err_0200020001.Sprintf(path, err)
		}
	}
	a.Clear()
	return nil
}

// Format 格式化文件，目前只能用来格式化Go源文件。
//
// Example:
//
// var assets asset.Assets
//
//	assets.Add("test.go", []byte(`
//	package main
//	import "fmt"
//	func Main() {fmt.Println("Hello, World!")}
//	`))
//	err := assets.Write()
//	if err != nil {
//		fmt.Print(err)
//	}
//	err = assets.Format()
//	if err != nil {
//		fmt.Print(err)
//	}
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go - TestFormat
//
// ErrCodes:
//   - Err_0200010002
//   - Err_0200010003
func (a Assets) Format() error {
	for path, content := range a.Files {
		// 检查文件是否为 Go 源文件
		if filepath.Ext(path) == ".go" {
			src, err := imports.Process(path, content, nil)
			if err != nil {
				return Err_0200020002.Sprintf(path, err)
			}
			if err := os.WriteFile(path, src, 0644); err != nil {
				return Err_0200020001.Sprintf(path, err)
			}
		}

	}
	return nil
}

// Clear 清空`Assets`中的Dirs和Files。
//
// Example:
//
//	var assets asset.Assets
//	assets.Add("file.txt", []byte("Hellow, World!"))
//	assets.AddDir("dir")
//	assets.Clear()
//	fmt.Print(assets)
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go - TestClear
func (a *Assets) Clear() {
	a.Dirs = nil
	a.Files = nil
}
