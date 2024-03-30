package field

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type RawBytes []byte

// Varchar 字符串类型的字段。
type Varchar struct {
	StringBuilder
	StringStorage
}

// UUID UUID类型的字段。
type UUID struct {
	UUIDBuilder
	StringStorage
}

// StringStorage 字符串类型的字段存储。
type StringStorage struct {
	value *string
}

// Set 设置字段的值。
func (i *StringStorage) Set(value string) error {
	i.value = &value
	return nil
}

// Get 获取字段的值。
func (i *StringStorage) Get() *string {
	return i.value
}

// Scan 从数据库中读取字段的值。
func (s *StringStorage) Scan(value any) error {
	if value == nil {
		s.value = nil
		return nil
	}
	return convertAssign(&s.value, value)
}

// String 返回字段的字符串表示。
func (s StringStorage) String() string {
	if s.value == nil {
		return "nil"
	}
	return *s.value
}

// Value 返回字段的值，和Get方法不同的是，Value方法返回的是接口类型。
func (s *StringStorage) Value() entity.FieldValue {
	if s.value == nil {
		return nil
	}
	return *s.value
}

// StringBuilder 字符串类型的字段构造器。
type StringBuilder struct {
	desc *entity.Descriptor
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (s *StringBuilder) Init(desc *entity.Descriptor) error {
	s.desc = desc
	return nil
}

// Descriptor 获取字段的描述信息。
func (s *StringBuilder) Descriptor() *entity.Descriptor {
	return s.desc
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
func (s *StringBuilder) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		v := varchar(s.desc.Size)
		return v
	default:
		return ""
	}
}

// ValueType 用于设置字段的值在go中类型名称。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (s *StringBuilder) ValueType() string {
	return "string"
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (s *StringBuilder) Name(name string) *StringBuilder {
	s.desc.AttrName = name
	return s
}

// MaxLen 设置字段的最大长度。
//
// Params:
//
//   - i: 字段的最大长度。
func (s *StringBuilder) MaxLen(i int64) *StringBuilder {
	s.desc.Size = i
	s.desc.Validators = append(s.desc.Validators, func(v string) error {
		if int64(len(v)) > i {
			return errors.New("value is greater than the required length")
		}
		return nil
	})
	return s
}

// MinLen 设置字段的最小长度。
//
// Params:
//
//   - size: 字段的最小长度。
func (s *StringBuilder) MinLen(i int) *StringBuilder {
	s.desc.Validators = append(s.desc.Validators, func(v string) error {
		if len(v) < i {
			return errors.New("value is less than the required length")
		}
		return nil
	})
	return s
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (s *StringBuilder) Required() *StringBuilder {
	s.desc.Required = true
	return s.MinLen(1)
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (s *StringBuilder) Primary(index int) *StringBuilder {
	s.desc.Required = true
	s.desc.Primary = index
	return s
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (s *StringBuilder) Comment(comment string) *StringBuilder {
	s.desc.Comment = comment
	return s
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (s *StringBuilder) Default(value string) *StringBuilder {
	s.desc.Default = true
	s.desc.DefaultValue = value
	return s
}

// Locked 设置字段为只读字段。
func (s *StringBuilder) Locked() *StringBuilder {
	s.desc.Locked = true
	return s
}

// varchar 返回varchar类型的字段。
func varchar(size int64) string {
	if size <= 0 {
		return "varchar(255)"
	} else {
		return fmt.Sprintf("varchar(%d)", size)
	}
}

// UUIDBuilder UUID类型的字段构造器。
type UUIDBuilder struct {
	desc *entity.Descriptor
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (u *UUIDBuilder) Init(desc *entity.Descriptor) error {
	u.desc = desc
	return nil
}

// Descriptor 获取字段的描述信息。
func (u *UUIDBuilder) Descriptor() *entity.Descriptor {
	return u.desc
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
func (u *UUIDBuilder) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		return "uuid"
	default:
		return "string"
	}
}

// ValueType 用于设置字段的值在go中类型名称。例如entity.Int64的ValueType为"int64"。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (s *UUID) ValueType() string {
	return "string"
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (u *UUID) Name(name string) *UUID {
	u.desc.AttrName = name
	return u
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (u *UUID) Required() *UUID {
	u.desc.Required = true
	u.desc.Validators = append(u.desc.Validators, func(v string) error {
		_, err := uuid.Parse(v)
		if err != nil {
			return errors.New("value is less than the required length")
		}
		return nil
	})
	return u
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (u *UUID) Primary(index int) *UUID {
	u.desc.Primary = index
	return u
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (u *UUID) Comment(comment string) *UUID {
	u.desc.Comment = comment
	return u
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (u *UUID) Default(value string) *UUID {
	u.desc.Default = true
	u.desc.DefaultValue = value
	return u
}

// Locked 设置字段为只读字段。
func (u *UUID) Locked() *UUID {
	u.desc.Locked = true
	return u
}
