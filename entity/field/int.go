package field

import (
	"errors"
	"fmt"

	"github.com/yohobala/taurus_go/entity"
)

// Int16 用于定义int16类型的字段。
type Int16 struct {
	IntBuilder[int16]
	IntStorage[int16]
}

// Int16A1 用于定义int16类型的数组字段。1维数组。
type Int16A1 struct {
	IntBuilder[[]int16]
	IntStorage[[]int16]
}

// Int16A2 用于定义int16类型的数组字段。2维数组。
type Int16A2 struct {
	IntBuilder[[][]int16]
	IntStorage[[][]int16]
}

// Int16A3 用于定义int16类型的数组字段。3维数组。
type Int16A3 struct {
	IntBuilder[[][][]int16]
	IntStorage[[][][]int16]
}

// Int16A4 用于定义int16类型的数组字段。4维数组。
type Int16A4 struct {
	IntBuilder[[][][][]int16]
	IntStorage[[][][][]int16]
}

// Int16A5 用于定义int16类型的数组字段。5维数组。
type Int16A5 struct {
	IntBuilder[[][][][][]int16]
	IntStorage[[][][][][]int16]
}

// Int32 用于定义int32类型的字段。
type Int32 struct {
	IntBuilder[int32]
	IntStorage[int32]
}

// Int32A1 用于定义int32类型的数组字段。1维数组。
type Int32A1 struct {
	IntBuilder[[]int32]
	IntStorage[[]int32]
}

// Int32A2 用于定义int32类型的数组字段。2维数组。
type Int32A2 struct {
	IntBuilder[[][]int32]
	IntStorage[[][]int32]
}

// Int32A3 用于定义int32类型的数组字段。3维数组。
type Int32A3 struct {
	IntBuilder[[][][]int32]
	IntStorage[[][][]int32]
}

// Int32A4 用于定义int32类型的数组字段。4维数组。
type Int32A4 struct {
	IntBuilder[[][][][]int32]
	IntStorage[[][][][]int32]
}

// Int32A5 用于定义int32类型的数组字段。5维数组。
type Int32A5 struct {
	IntBuilder[[][][][][]int32]
	IntStorage[[][][][][]int32]
}

// Int64 用于定义int64类型的字段。
type Int64 struct {
	IntBuilder[int64]
	IntStorage[int64]
}

// Int64A1 用于定义int64类型的数组字段。1维数组。
type Int64A1 struct {
	IntBuilder[[]int64]
	IntStorage[[]int64]
}

// Int64A2 用于定义int64类型的数组字段。2维数组。
type Int64A2 struct {
	IntBuilder[[][]int64]
	IntStorage[[][]int64]
}

// Int64A3 用于定义int64类型的数组字段。3维数组。
type Int64A3 struct {
	IntBuilder[[][][]int64]
	IntStorage[[][][]int64]
}

// Int64A4 用于定义int64类型的数组字段。4维数组。
type Int64A4 struct {
	IntBuilder[[][][][]int64]
	IntStorage[[][][][]int64]
}

// Int64A5 用于定义int64类型的数组字段。5维数组。
type Int64A5 struct {
	IntBuilder[[][][][][]int64]
	IntStorage[[][][][][]int64]
}

// IntBuilder 用于构建int类型的字段。
type IntBuilder[T any] struct {
	BaseBuilder[T]
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (i *IntBuilder[T]) Name(name string) *IntBuilder[T] {
	i.desc.AttrName = name
	return i
}

// MinLen 设置字段的最小长度。
//
// Params:
//
//   - size: 字段的最小长度。
func (i *IntBuilder[T]) MinLen(size int) *IntBuilder[T] {
	i.desc.Validators = append(i.desc.Validators, func(b []byte) error {
		if len(b) < size {
			return errors.New("value is less than the required length.")
		}
		return nil
	})
	return i
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (i *IntBuilder[T]) Required() *IntBuilder[T] {
	i.desc.Required = true
	return i.MinLen(1)
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (i *IntBuilder[T]) Primary(index int) *IntBuilder[T] {
	i.desc.Required = true
	i.desc.Primary = index
	return i
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (i *IntBuilder[T]) Comment(comment string) *IntBuilder[T] {
	i.desc.Comment = comment
	return i
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (i *IntBuilder[T]) Default(value T) *IntBuilder[T] {
	i.desc.Default = true
	i.desc.DefaultValue = fmt.Sprintf("%v", value)
	return i
}

// Sequence 设置字段的序列。
// 如果序列不存在，则会自动创建序列。
// 优先级高于[Default]。
//
// Params:
//
//   - s: 序列。
func (i *IntBuilder[T]) Sequence(s entity.Sequence) *IntBuilder[T] {
	i.desc.Default = true
	i.desc.DefaultValue = fmt.Sprintf(`nextval('%s'::regclass)`, *s.Name)
	i.desc.Sequence = s
	return i
}

// Locked 设置字段为只读字段。
func (i *IntBuilder[T]) Locked() *IntBuilder[T] {
	i.desc.Locked = true
	return i
}

type IntStorage[T any] struct {
	BaseStorage[T]
}
