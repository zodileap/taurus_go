package internal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGenerateCmd(t *testing.T) {
	cmd := GenerateCmd()
	if cmd.Use != "generate [flags] path" {
		t.Fatalf("GenerateCmd use 不正确: %s", cmd.Use)
	}
	if cmd.Flag("template") == nil || cmd.Flag("package") == nil {
		t.Fatal("GenerateCmd 缺少预期 flag")
	}
}

func TestNewCmd(t *testing.T) {
	cmd := NewCmd()
	if cmd.Use != "new [entities]" {
		t.Fatalf("NewCmd use 不正确: %s", cmd.Use)
	}
	if err := cmd.Args(cmd, []string{"User"}); err != nil {
		t.Fatalf("NewCmd 参数校验失败: %v", err)
	}
	if cmd.Flag("target") == nil || cmd.Flag("schema") == nil || cmd.Flag("entities") == nil {
		t.Fatal("NewCmd 缺少预期 flag")
	}
}

func TestHasStructInFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "schema.go")
	content := []byte("package schema\ntype User struct{}\n")
	if err := os.WriteFile(file, content, 0o644); err != nil {
		t.Fatalf("写入测试文件失败: %v", err)
	}

	exists, err := hasStructInFile(file, "User")
	if err != nil {
		t.Fatalf("hasStructInFile 返回错误: %v", err)
	}
	if !exists {
		t.Fatal("未识别到已存在结构体")
	}

	exists, err = hasStructInFile(file, "Role")
	if err != nil {
		t.Fatalf("hasStructInFile 查询不存在结构体时返回错误: %v", err)
	}
	if exists {
		t.Fatal("误识别不存在结构体为存在")
	}
}
