package field

import (
	"fmt"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

// Bool 布尔类型的字段。
type Bool struct {
	BoolBuilder
	BoolStorage
}

// BoolStorage 布尔类型的字段存储。
type BoolStorage struct {
	value *bool
}

// Set 设置字段的值。
func (i *BoolStorage) Set(value bool) error {
	i.value = &value
	return nil
}

// Get 获取字段的值。
func (i *BoolStorage) Get() *bool {
	return i.value
}

// Scan 从数据库中读取字段的值。
func (s *BoolStorage) Scan(value any) error {
	if value == nil {
		s.value = nil
		return nil
	}
	return convertAssign(&s.value, value)
}

// String 返回字段的字符串表示。
func (s BoolStorage) String() string {
	if s.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", *s.value)
}

// Value 返回字段的值，和Get方法不同的是，Value方法返回的是接口类型。
func (s *BoolStorage) Value() entity.FieldValue {
	if s.value == nil {
		return nil
	}
	return *s.value
}

// BoolBuilder
type BoolBuilder struct {
	desc *entity.Descriptor
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (i *BoolBuilder) Init(initDesc *entity.Descriptor) error {
	i.desc = initDesc
	return nil
}

// Descriptor 获取字段的描述。
func (i *BoolBuilder) Descriptor() *entity.Descriptor {
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
func (i *BoolBuilder) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		return "boolean"
	default:
		return ""
	}
}

// ValueType 用于设置字段的值在go中类型名称。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (i *BoolBuilder) ValueType() string {
	return "bool"
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (i *BoolBuilder) Name(name string) *BoolBuilder {
	i.desc.Name = name
	return i
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (i *BoolBuilder) Required() *BoolBuilder {
	i.desc.Required = true
	return i
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (i *BoolBuilder) Primary(index int) *BoolBuilder {
	i.desc.Required = true
	i.desc.Primary = index
	return i
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (i *BoolBuilder) Comment(comment string) *BoolBuilder {
	i.desc.Comment = comment
	return i
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (i *BoolBuilder) Default(value bool) *BoolBuilder {
	i.desc.Default = true
	i.desc.DefaultValue = fmt.Sprintf("%v", value)
	return i
}

// Locked 设置字段为只读字段。
func (i *BoolBuilder) Locked() *BoolBuilder {
	i.desc.Locked = true
	return i
}
