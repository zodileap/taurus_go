package template

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

type mockTemplatePath struct {
	dir string
}

func (m mockTemplatePath) Dir() string {
	return m.dir
}

func TestParseDir(t *testing.T) {
	tempDir := t.TempDir()
	tmplFile := filepath.Join(tempDir, "hello.tmpl")
	goFile := filepath.Join(tempDir, "skip.go")

	if err := os.WriteFile(tmplFile, []byte(`{{ define "hello.tmpl" }}{{ stringJoin "he" "llo" }} {{ .Name }}{{ end }}`), 0o644); err != nil {
		t.Fatalf("写入模板文件失败: %v", err)
	}
	if err := os.WriteFile(goFile, []byte("package ignore"), 0o644); err != nil {
		t.Fatalf("写入 Go 文件失败: %v", err)
	}

	tmpl, err := NewTemplate("root").ParseDir(tempDir)
	if err != nil {
		t.Fatalf("ParseDir 失败: %v", err)
	}

	var buf bytes.Buffer
	if err := tmpl.ExecuteTemplate(&buf, "hello.tmpl", map[string]string{"Name": "world"}); err != nil {
		t.Fatalf("执行模板失败: %v", err)
	}
	if buf.String() != "hello world" {
		t.Fatalf("模板输出不正确: %q", buf.String())
	}
}

func TestMustParsePanics(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Fatal("MustParse 在错误场景下应 panic")
		}
	}()

	MustParse(NewTemplate("panic").ParseFiles(filepath.Join(t.TempDir(), "missing.tmpl")))
}

func TestDir(t *testing.T) {
	format := Dir("%s/output")
	if got := format(mockTemplatePath{dir: "root"}); got != "root/output" {
		t.Fatalf("Dir 格式化结果不正确: %s", got)
	}
}
