package field

import (
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type RawBytes []byte

// Varchar 字符串类型的字段。
type Varchar struct {
	VarcharBuilder[string]
	StringStorage[string]
}

// UUID UUID类型的字段。
type UUID struct {
	UUIDBuilder[string]
	StringStorage[string]
}

// VarcharBuilder 字符串类型的字段构造器。
type VarcharBuilder[T any] struct {
	BaseBuilder[T]
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
func (s *VarcharBuilder[T]) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		v := varchar(s.desc.Size)
		return v
	default:
		return ""
	}
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (s *VarcharBuilder[T]) Name(name string) *VarcharBuilder[T] {
	s.desc.AttrName = name
	return s
}

// MaxLen 设置字段的最大长度。
//
// Params:
//
//   - i: 字段的最大长度。
func (s *VarcharBuilder[T]) MaxLen(i int64) *VarcharBuilder[T] {
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
func (s *VarcharBuilder[T]) MinLen(i int) *VarcharBuilder[T] {
	s.desc.Validators = append(s.desc.Validators, func(v string) error {
		if len(v) < i {
			return errors.New("value is less than the required length")
		}
		return nil
	})
	return s
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (s *VarcharBuilder[T]) Required() *VarcharBuilder[T] {
	s.desc.Required = true
	return s.MinLen(1)
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (s *VarcharBuilder[T]) Primary(index int) *VarcharBuilder[T] {
	s.desc.Required = true
	s.desc.Primary = index
	return s
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (s *VarcharBuilder[T]) Comment(comment string) *VarcharBuilder[T] {
	s.desc.Comment = comment
	return s
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (s *VarcharBuilder[T]) Default(value string) *VarcharBuilder[T] {
	s.desc.Default = true
	s.desc.DefaultValue = value
	return s
}

// Locked 设置字段为只读字段。
func (s *VarcharBuilder[T]) Locked() *VarcharBuilder[T] {
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
type UUIDBuilder[T any] struct {
	BaseBuilder[T]
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
func (u *UUIDBuilder[T]) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		return "uuid"
	default:
		return ""
	}
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (u *UUIDBuilder[T]) Name(name string) *UUIDBuilder[T] {
	u.desc.AttrName = name
	return u
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (u *UUIDBuilder[T]) Required() *UUIDBuilder[T] {
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
func (u *UUIDBuilder[T]) Primary(index int) *UUIDBuilder[T] {
	u.desc.Required = true
	u.desc.Primary = index
	return u
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (u *UUIDBuilder[T]) Comment(comment string) *UUIDBuilder[T] {
	u.desc.Comment = comment
	return u
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (u *UUIDBuilder[T]) Default(value string) *UUIDBuilder[T] {
	u.desc.Default = true
	u.desc.DefaultValue = value
	return u
}

// Locked 设置字段为只读字段。
func (u *UUIDBuilder[T]) Locked() *UUIDBuilder[T] {
	u.desc.Locked = true
	return u
}

// StringStorage[T] 字符串类型的字段存储。
type StringStorage[T any] struct {
	BaseStorage[T]
}
