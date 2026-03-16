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
