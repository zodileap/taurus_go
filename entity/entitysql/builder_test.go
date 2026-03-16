package entitysql

import (
	"testing"

	"github.com/zodileap/taurus_go/entity/dialect"
)

func TestOpString(t *testing.T) {
	if OpEQ.String() != "=" || OpContains.String() != "@>" {
		t.Fatalf("操作符字符串不正确: %s %s", OpEQ.String(), OpContains.String())
	}
}

func TestRawQuery(t *testing.T) {
	spec, err := Raw("SELECT 1").Query()
	if err != nil {
		t.Fatalf("Raw.Query 失败: %v", err)
	}
	if spec.Query != "SELECT 1" || len(spec.Args) != 0 {
		t.Fatalf("Raw.Query 结果不正确: %+v", spec)
	}
}

func TestBuilderQuoteAndArg(t *testing.T) {
	var builder Builder
	builder.SetDialect(dialect.PostgreSQL)

	builder.WriteString("SELECT ").
		Ident("users").
		Blank().
		WriteString("WHERE ").
		Ident("name").
		WriteOp(OpEQ).
		Arg("alice")

	if builder.String() != `SELECT "users" WHERE "name" = $1` {
		t.Fatalf("Builder 生成 SQL 不正确: %s", builder.String())
	}
	if builder.total != 1 || len(builder.args) != 1 || builder.args[0] != "alice" {
		t.Fatalf("Builder 参数状态不正确: total=%d args=%v", builder.total, builder.args)
	}
}
