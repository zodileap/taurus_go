package cmd

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSplit(t *testing.T) {
	got := Split("mkdir dir")
	if len(got) != 2 || got[0] != "mkdir" || got[1] != "dir" {
		t.Fatalf("Split 结果不正确: %v", got)
	}
}

func TestSetDirAndRun(t *testing.T) {
	tempDir := t.TempDir()

	output, err := New("pwd").SetDir(tempDir).Run()
	if err != nil {
		t.Fatalf("运行 pwd 失败: %v", err)
	}

	got := strings.TrimSpace(string(output))
	if got != tempDir {
		t.Fatalf("SetDir 未生效，期望 %s，实际 %s", tempDir, got)
	}
}

func TestSetEnv(t *testing.T) {
	output, err := NewSh("sh", "printf %s \"$TEST_VALUE\"").
		SetEnv(append(os.Environ(), "TEST_VALUE=hello")).
		Run()
	if err != nil {
		t.Fatalf("带环境变量执行命令失败: %v", err)
	}
	if string(output) != "hello" {
		t.Fatalf("SetEnv 未生效，实际输出: %q", string(output))
	}
}

func TestRunErrorContainsStdoutAndStderr(t *testing.T) {
	output, err := NewSh("sh", "printf stdout-msg; printf stderr-msg >&2; exit 1").Run()
	if err == nil {
		t.Fatal("期望命令执行失败，实际成功")
	}
	if string(output) != "stdout-msg" {
		t.Fatalf("stdout 返回值不正确: %q", string(output))
	}
	message := err.Error()
	if !strings.Contains(message, "stdout-msg") || !strings.Contains(message, "stderr-msg") {
		t.Fatalf("错误信息未包含 stdout/stderr: %s", message)
	}
}

func TestIsGoModuleInitialized(t *testing.T) {
	tempDir := t.TempDir()
	if IsGoModuleInitialized(tempDir) {
		t.Fatal("空目录不应被识别为已初始化 Go 模块")
	}

	goMod := filepath.Join(tempDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module example.com/test\n"), 0o644); err != nil {
		t.Fatalf("写入 go.mod 失败: %v", err)
	}
	if !IsGoModuleInitialized(tempDir) {
		t.Fatal("存在 go.mod 的目录应被识别为 Go 模块")
	}
}
