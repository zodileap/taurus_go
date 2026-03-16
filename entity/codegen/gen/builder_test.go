package gen

import "testing"

func TestValidSchemaName(t *testing.T) {
	if err := ValidSchemaName("UserSchema"); err != nil {
		t.Fatalf("合法 schema 名称被误判: %v", err)
	}

	if err := ValidSchemaName("Type"); err == nil {
		t.Fatal("与 Go 预定义标识符冲突的 schema 名称应返回错误")
	}

	if err := ValidSchemaName("For"); err == nil {
		t.Fatal("与 Go 关键字冲突的 schema 名称应返回错误")
	}
}
