package field

import (
	"testing"

	"github.com/zodileap/taurus_go/entity"
)

func TestBaseBuilderAndStorage(t *testing.T) {
	builder := &BaseBuilder[[]int]{}
	desc := &entity.Descriptor{Name: "ages"}
	if err := builder.Init(desc); err != nil {
		t.Fatalf("Init 失败: %v", err)
	}

	if desc.Depth != 1 || desc.BaseType != "int" {
		t.Fatalf("sliceTypeDetails 结果不正确: depth=%d baseType=%s", desc.Depth, desc.BaseType)
	}

	builder.SetValueType("[]int").Unique(1).Index(2).IndexName("idx_ages").IndexMethod("btree")
	if builder.ValueType() != "[]int" {
		t.Fatalf("ValueType 不正确: %s", builder.ValueType())
	}
	if len(desc.Uniques) != 1 || len(desc.Indexes) != 1 || desc.IndexName != "idx_ages" || desc.IndexMethod != "btree" {
		t.Fatalf("Builder 描述信息不正确: %+v", desc)
	}

	storage := &BaseStorage[int]{}
	if err := storage.Set(42); err != nil {
		t.Fatalf("Set 失败: %v", err)
	}
	if storage.Get() == nil || *storage.Get() != 42 {
		t.Fatalf("Get 结果不正确: %+v", storage.Get())
	}
	if storage.String() != "42" {
		t.Fatalf("String 结果不正确: %s", storage.String())
	}

	data, err := storage.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON 失败: %v", err)
	}
	if string(data) != "42" {
		t.Fatalf("MarshalJSON 结果不正确: %s", string(data))
	}
}
