package field

import (
	"errors"
	"fmt"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type Int16 struct {
	IntStorage[int16]
	IntBuilder[int16]
}

type Int32 struct {
	IntStorage[int32]
	IntBuilder[int32]
}

type Int64 struct {
	IntStorage[int64]
	IntBuilder[int64]
}

type IntBuilder[T any] struct {
	desc *entity.Descriptor
}

func (i *IntBuilder[T]) Init(desc *entity.Descriptor) error {
	i.desc = desc
	return nil
}

func (i *IntBuilder[T]) Descriptor() *entity.Descriptor {
	return i.desc
}

func (i *IntBuilder[T]) AttrType(dbType dialect.DbDriver) string {
	var t T
	return attrType(t)
}

func (i *IntBuilder[T]) Name(name string) *IntBuilder[T] {
	i.desc.AttrName = name
	return i
}

// 设置字段的最小长度。
func (i *IntBuilder[T]) MinLen(size int) *IntBuilder[T] {
	i.desc.Validators = append(i.desc.Validators, func(b []byte) error {
		if len(b) < size {
			return errors.New("value is less than the required length.")
		}
		return nil
	})
	return i
}

// 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (i *IntBuilder[T]) Required() *IntBuilder[T] {
	i.desc.Required = true
	return i.MinLen(1)
}

// 设置字段为主键。
// index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (i *IntBuilder[T]) Primary(index int) *IntBuilder[T] {
	i.desc.Required = true
	i.desc.Primary = index
	return i
}

// 设置字段的注释。
func (i *IntBuilder[T]) Comment(comment string) *IntBuilder[T] {
	i.desc.Comment = comment
	return i
}

// 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
func (i *IntBuilder[T]) Default(value int16) *IntBuilder[T] {
	i.desc.Default = true
	i.desc.DefaultValue = fmt.Sprintf("%d", value)
	return i
}

// 设置字段的序列。
// 如果序列不存在，则会自动创建序列。
// 优先级高于[Default]。
func (i *IntBuilder[T]) Sequence(s entity.Sequence) *IntBuilder[T] {
	i.desc.Default = true
	i.desc.DefaultValue = fmt.Sprintf(`nextval('%s'::regclass)`, *s.Name)
	i.desc.Sequence = s
	return i
}

func (i *IntBuilder[T]) ValueType() string {
	var t T
	return intType(t)
}

type IntStorage[T any] struct {
	value *T
}

func (i *IntStorage[T]) Set(value T) error {
	i.value = &value
	return nil
}

func (i *IntStorage[T]) Get() *T {
	return i.value
}

func (i *IntStorage[T]) Scan(value interface{}) error {
	if value == nil {
		i.value = nil
		return nil
	}
	return convertAssign(&i.value, value)
}

func (i IntStorage[T]) String() string {
	if i.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%d", *i.value)
}

func (i *IntStorage[T]) Value() entity.FieldValue {
	if i.value == nil {
		return nil
	}
	return *i.value
}

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
