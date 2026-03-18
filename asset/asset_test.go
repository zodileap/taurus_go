package asset

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	terr "github.com/zodileap/taurus_go/err"
)

func withWorkingDir(t *testing.T, dir string) {
	t.Helper()

	oldDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("切换目录失败: %v", err)
	}
	t.Cleanup(func() {
		if err := os.Chdir(oldDir); err != nil {
			t.Fatalf("恢复目录失败: %v", err)
		}
	})
}

func requireErrCode(t *testing.T, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("期望错误码 %s，实际无错误", want)
	}

	var errCode terr.ErrCode
	if !errors.As(got, &errCode) {
		t.Fatalf("期望 ErrCode，实际为 %T: %v", got, got)
	}
	if errCode.Code() != want {
		t.Fatalf("错误码不匹配，期望 %s，实际 %s", want, errCode.Code())
	}
}

func TestAssetsWrite(t *testing.T) {
	workDir := t.TempDir()
	withWorkingDir(t, workDir)

	var assets Assets
	assets.Add("test.txt", []byte("Hello, World!"))

	if err := assets.Write(); err != nil {
		t.Fatalf("写入文件失败: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(workDir, "test.txt"))
	if err != nil {
		t.Fatalf("读取写入结果失败: %v", err)
	}
	if string(content) != "Hello, World!" {
		t.Fatalf("文件内容不匹配，实际: %q", string(content))
	}
}

func TestAssetsAddDir(t *testing.T) {
	t.Run("创建目录", func(t *testing.T) {
		workDir := t.TempDir()
		withWorkingDir(t, workDir)

		var assets Assets
		if err := assets.AddDir(filepath.Join("nested", "dir")); err != nil {
			t.Fatalf("添加目录失败: %v", err)
		}
		if err := assets.Write(); err != nil {
			t.Fatalf("写入目录失败: %v", err)
		}

		info, err := os.Stat(filepath.Join(workDir, "nested", "dir"))
		if err != nil {
			t.Fatalf("检查目录失败: %v", err)
		}
		if !info.IsDir() {
			t.Fatal("目标路径不是目录")
		}
	})

	t.Run("空目录路径返回错误码", func(t *testing.T) {
		workDir := t.TempDir()
		withWorkingDir(t, workDir)

		var assets Assets
		if err := assets.AddDir(""); err != nil {
			t.Fatalf("AddDir 不应在收集阶段报错: %v", err)
		}

		requireErrCode(t, assets.Write(), Err_0200030001.Code())
	})
}

func TestAssetsFormat(t *testing.T) {
	t.Run("格式化成功", func(t *testing.T) {
		workDir := t.TempDir()
		withWorkingDir(t, workDir)

		var assets Assets
		assets.Add("test.go", []byte(`
package main

func Main() {fmt.Println("Hello, World!")}
`))

		if err := assets.Write(); err != nil {
			t.Fatalf("写入 Go 文件失败: %v", err)
		}
		if err := assets.Format(); err != nil {
			t.Fatalf("格式化 Go 文件失败: %v", err)
		}

		content, err := os.ReadFile(filepath.Join(workDir, "test.go"))
		if err != nil {
			t.Fatalf("读取格式化结果失败: %v", err)
		}
		formatted := string(content)
		if !strings.Contains(formatted, "package main") {
			t.Fatalf("格式化结果缺少 package 声明: %s", formatted)
		}
		if !strings.Contains(formatted, "import \"fmt\"") {
			t.Fatalf("格式化结果缺少 fmt 导入: %s", formatted)
		}
		if !strings.Contains(formatted, "fmt.Println(\"Hello, World!\")") {
			t.Fatalf("格式化结果缺少目标调用: %s", formatted)
		}
	})

	t.Run("缺少 package 时返回错误码", func(t *testing.T) {
		workDir := t.TempDir()
		withWorkingDir(t, workDir)

		var assets Assets
		assets.Add("test.go", []byte(`
func Main() {fmt.Println("Hello, World!")}
`))

		if err := assets.Write(); err != nil {
			t.Fatalf("写入无效 Go 文件失败: %v", err)
		}

		requireErrCode(t, assets.Format(), Err_0200020002.Code())
	})
}

func TestAssetsCopyFile(t *testing.T) {
	workDir := t.TempDir()
	src := filepath.Join(workDir, "source.txt")
	dst := filepath.Join(workDir, "dest.txt")

	if err := os.WriteFile(src, []byte("copy via assets"), 0o644); err != nil {
		t.Fatalf("创建源文件失败: %v", err)
	}

	var assets Assets
	if err := assets.CopyFile(src, dst); err != nil {
		t.Fatalf("Assets.CopyFile 失败: %v", err)
	}
	if err := assets.Write(); err != nil {
		t.Fatalf("写入复制文件失败: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("读取复制结果失败: %v", err)
	}
	if string(content) != "copy via assets" {
		t.Fatalf("复制内容不匹配，实际: %q", string(content))
	}
}

func TestAssetsClear(t *testing.T) {
	var assets Assets
	assets.Add("test.txt", []byte("value"))
	if err := assets.AddDir("nested"); err != nil {
		t.Fatalf("添加目录失败: %v", err)
	}

	assets.Clear()

	if assets.Files != nil || assets.Dirs != nil {
		t.Fatalf("Clear 后应清空全部内容，实际 Files=%v Dirs=%v", assets.Files, assets.Dirs)
	}
}

func TestCopyFile(t *testing.T) {
	workDir := t.TempDir()
	src := filepath.Join(workDir, "source.txt")
	dst := filepath.Join(workDir, "dest.txt")

	if err := os.WriteFile(src, []byte("copy me"), 0o644); err != nil {
		t.Fatalf("创建源文件失败: %v", err)
	}

	if err := CopyFile(src, dst); err != nil {
		t.Fatalf("复制文件失败: %v", err)
	}

	content, err := os.ReadFile(dst)
	if err != nil {
		t.Fatalf("读取目标文件失败: %v", err)
	}
	if string(content) != "copy me" {
		t.Fatalf("复制内容不匹配，实际: %q", string(content))
	}
}

func TestCopyDir(t *testing.T) {
	workDir := t.TempDir()
	src := filepath.Join(workDir, "src")
	dst := filepath.Join(workDir, "dst")

	if err := os.MkdirAll(filepath.Join(src, "nested"), 0o755); err != nil {
		t.Fatalf("创建源目录失败: %v", err)
	}
	if err := os.WriteFile(filepath.Join(src, "nested", "file.txt"), []byte("dir copy"), 0o644); err != nil {
		t.Fatalf("创建源目录文件失败: %v", err)
	}

	if err := CopyDir(src, dst); err != nil {
		t.Fatalf("复制目录失败: %v", err)
	}

	content, err := os.ReadFile(filepath.Join(dst, "nested", "file.txt"))
	if err != nil {
		t.Fatalf("读取复制后的目录文件失败: %v", err)
	}
	if string(content) != "dir copy" {
		t.Fatalf("目录复制内容不匹配，实际: %q", string(content))
	}
}
