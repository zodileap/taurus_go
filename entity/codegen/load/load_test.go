package load

import (
	"strings"
	"testing"
)

func TestFilename(t *testing.T) {
	name := filename("github.com/example/schema")
	if !strings.HasPrefix(name, "gen_github.com_example_schema_") {
		t.Fatalf("filename 前缀不正确: %s", name)
	}
}

func TestContains(t *testing.T) {
	if !contains([]string{"a", "b"}, "a") {
		t.Fatal("contains 未找到已存在元素")
	}
	if contains([]string{"a", "b"}, "c") {
		t.Fatal("contains 误判不存在元素为存在")
	}
}

func TestEntityTemplates(t *testing.T) {
	tmpls, err := entityTemplates()
	if err != nil {
		t.Fatalf("entityTemplates 执行失败: %v", err)
	}
	if len(tmpls) != 2 {
		t.Fatalf("entityTemplates 返回数量不正确: %d", len(tmpls))
	}
	if !strings.Contains(tmpls[0], `define "entity"`) {
		t.Fatalf("实体模板定义缺失: %s", tmpls[0])
	}
	if !strings.Contains(tmpls[1], `define "imports"`) {
		t.Fatalf("导入模板定义缺失: %s", tmpls[1])
	}
}
