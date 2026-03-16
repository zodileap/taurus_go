package codegen

import (
	"errors"
	"testing"

	"github.com/zodileap/taurus_go/entity"
	"github.com/zodileap/taurus_go/entity/codegen/gen"

	terr "github.com/zodileap/taurus_go/err"
)

func TestNormalizePkg(t *testing.T) {
	cfg := &gen.Config{Package: "github.com/example/my-pkg"}
	if err := normalizePkg(cfg); err != nil {
		t.Fatalf("normalizePkg 处理带连字符的包名失败: %v", err)
	}
	if cfg.Package != "github.com/example/my_pkg" {
		t.Fatalf("包名标准化结果不正确: %s", cfg.Package)
	}

	cfg = &gen.Config{Package: "github.com/example/123"}
	err := normalizePkg(cfg)
	if err == nil {
		t.Fatal("非法包名应返回错误")
	}

	var errCode terr.ErrCode
	if !errors.As(err, &errCode) || errCode.Code() != entity.Err_0100020015.Code() {
		t.Fatalf("非法包名错误不正确: %v", err)
	}
}
