package field

import (
	"fmt"
)

// Bool 布尔类型的字段。
type Bool struct {
	BoolBuilder[bool]
	BoolStorage[bool]
}

type BoolA1 struct {
	BoolBuilder[[]bool]
	BoolStorage[[]bool]
}

type BoolA2 struct {
	BoolBuilder[[][]bool]
	BoolStorage[[][]bool]
}

type BoolA3 struct {
	BoolBuilder[[][][]bool]
	BoolStorage[[][][]bool]
}

type BoolA4 struct {
	BoolBuilder[[][][][]bool]
	BoolStorage[[][][][]bool]
}

type BoolA5 struct {
	BoolBuilder[[][][][][]bool]
	BoolStorage[[][][][][]bool]
}

// BoolBuilder 布尔类型的字段构建器。
type BoolBuilder[T any] struct {
	BaseBuilder[T]
}

// BoolStorage 布尔类型的字段存储。
type BoolStorage[T any] struct {
	BaseStorage[T]
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (i *BoolBuilder[T]) Name(name string) *BoolBuilder[T] {
	i.desc.AttrName = name
	return i
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (i *BoolBuilder[T]) Required() *BoolBuilder[T] {
	i.desc.Required = true
	return i
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (i *BoolBuilder[T]) Primary(index int) *BoolBuilder[T] {
	i.desc.Required = true
	i.desc.Primary = index
	return i
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (i *BoolBuilder[T]) Comment(comment string) *BoolBuilder[T] {
	i.desc.Comment = comment
	return i
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (i *BoolBuilder[T]) Default(value bool) *BoolBuilder[T] {
	i.desc.Default = true
	i.desc.DefaultValue = fmt.Sprintf("%v", value)
	return i
}

// Locked 设置字段为只读字段。
func (i *BoolBuilder[T]) Locked() *BoolBuilder[T] {
	i.desc.Locked = true
	return i
}

// Unique 设置字段为唯一字段或参与联合唯一约束。
// 相同的序号表示这些字段组成联合唯一约束。
func (i *BoolBuilder[T]) Unique(index int) *BoolBuilder[T] {
	i.desc.Uniques = append(i.desc.Uniques, index)
	return i
}
