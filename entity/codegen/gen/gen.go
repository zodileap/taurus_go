package gen

import (
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
)

// PrepareEnv 检查是否有runtime.go,如果存在检查导入部分，避免循环导入
// 如果不存在则创建runtime.go文件
// 如果存在导入部分则在文件开头添加 "// +build tools\n"
// 这样在生成代码时，runtime.go文件不会被编译到最终的二进制文件中
//
// Params:
//
//   - c: 代码生成的配置。
//
// Returns:
//
//	0: 无操作函数。
//	1: 错误信息。
func PrepareEnv(c *Config) (undo func() error, err error) {
	var (
		// 无操作函数
		nop = func() error { return nil }
		// 构建路径：使用 filepath.Join 构建 runtime.go 文件的完整路径.
		path = filepath.Join(c.Target, "runtime.go")
	)
	out, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nop, nil
		}
		return nil, err
	}
	fi, err := parser.ParseFile(token.NewFileSet(), path, out, parser.ImportsOnly)
	if err != nil {
		return nil, err
	}
	// Targeted package doesn't import the schema.
	if len(fi.Imports) == 0 {
		return nop, nil
	}
	if err := os.WriteFile(path, append([]byte("// +build tools\n"), out...), 0644); err != nil {
		return nil, err
	}
	return func() error { return os.WriteFile(path, out, 0644) }, nil
}
