package field

import (
	"errors"
	"fmt"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

// Int16 用于定义int16类型的字段。
type Int16 struct {
	IntStorage[int16]
	IntBuilder[int16]
}

// Int32 用于定义int32类型的字段。
type Int32 struct {
	IntStorage[int32]
	IntBuilder[int32]
}

// Int64 用于定义int64类型的字段。
type Int64 struct {
	IntStorage[int64]
	IntBuilder[int64]
}

// IntBuilder 用于构建int类型的字段。
type IntBuilder[T any] struct {
	desc *entity.Descriptor
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (i *IntBuilder[T]) Init(desc *entity.Descriptor) error {
	i.desc = desc
	return nil
}

// Descriptor 获取字段的描述信息。
func (i *IntBuilder[T]) Descriptor() *entity.Descriptor {
	return i.desc
}

// AttrType 获取字段的数据库中的类型名，如果返回空字符串，会出现错误。
//
// Params:
//
//   - dbType: 数据库类型。
//
// Returns:
//
//   - 字段的数据库中的类型名。
func (i *IntBuilder[T]) AttrType(dbType dialect.DbDriver) string {
	var t T
	return attrType(t)
}

// ValueType 用于设置字段的值在go中类型名称。例如entity.Int64的ValueType为"int64"。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (i *IntBuilder[T]) ValueType() string {
	var t T
	return intType(t)
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
func (i *IntBuilder[T]) Default(value int) *IntBuilder[T] {
	i.desc.Default = true
	i.desc.DefaultValue = fmt.Sprintf("%d", value)
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
	value *T
}

// Set 设置字段的值。
func (i *IntStorage[T]) Set(value T) error {
	i.value = &value
	return nil
}

// Get 获取字段的值。
func (i *IntStorage[T]) Get() *T {
	return i.value
}

// Scan 从数据库中读取字段的值。
func (i *IntStorage[T]) Scan(value interface{}) error {
	if value == nil {
		i.value = nil
		return nil
	}
	return convertAssign(&i.value, value)
}

// String 返回字段的字符串表示。
func (i IntStorage[T]) String() string {
	if i.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *i.value)
}

// Value 返回字段的值，和Get方法不同的是，Value方法返回的是接口类型。
func (i *IntStorage[T]) Value() entity.FieldValue {
	if i.value == nil {
		return nil
	}
	return *i.value
}

// attrType 返回字段的数据库中的类型名。
func attrType(t any) string {
	switch t.(type) {
	case int16:
		return "int2"
	case int32:
		return "int4"
	case int64:
		return "int8"
	default:
		return "int"
	}
}

// intType 返回字段的值在go中类型名称。
func intType(t any) string {
	switch t.(type) {
	case int16:
		return "int16"
	case int32:
		return "int32"
	case int64:
		return "int64"
	default:
		return "int"
	}
}
