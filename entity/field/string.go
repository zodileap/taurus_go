package field

import (
	"errors"
	"fmt"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type RawBytes []byte

type String struct {
	StringBuilder
	StringStorage
}

type StringBuilder struct {
	desc *entity.Descriptor
}

func (s *StringBuilder) Init(desc *entity.Descriptor) error {
	s.desc = desc
	return nil
}

func (s *StringBuilder) Descriptor() *entity.Descriptor {
	return s.desc
}

func (s *StringBuilder) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		v := varchar(s.desc.Size)
		return v
	default:
		return ""
	}
}

func (s *StringBuilder) Name(name string) *StringBuilder {
	s.desc.AttrName = name
	return s
}

// 设置字段的最大长度。
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

// 设置字段的最小长度。
func (s *StringBuilder) MinLen(i int) *StringBuilder {
	s.desc.Validators = append(s.desc.Validators, func(v string) error {
		if len(v) < i {
			return errors.New("value is less than the required length")
		}
		return nil
	})
	return s
}

// 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (s *StringBuilder) Required() *StringBuilder {
	s.desc.Required = true
	return s.MinLen(1)
}

// 设置字段为主键。
// index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (s *StringBuilder) Primary(index int) *StringBuilder {
	s.desc.Required = true
	s.desc.Primary = index
	return s
}

// 设置字段的注释。
func (s *StringBuilder) Comment(comment string) *StringBuilder {
	s.desc.Comment = comment
	return s
}

// 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
func (s *StringBuilder) Default(value string) *StringBuilder {
	s.desc.Default = true
	s.desc.DefaultValue = value
	return s
}

func (s *StringBuilder) ValueType() string {
	return "string"
}

func varchar(size int64) string {
	if size <= 0 {
		return "varchar(255)"
	} else {
		return fmt.Sprintf("varchar(%d)", size)
	}
}

type UUID struct {
	UUIDBuilder
	StringStorage
}

type UUIDBuilder struct {
	desc *entity.Descriptor
}

func (u *UUIDBuilder) Init(desc *entity.Descriptor) error {
	u.desc = desc
	return nil
}

func (u *UUIDBuilder) Descriptor() *entity.Descriptor {
	return u.desc
}

func (u *UUIDBuilder) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		return "uuid"
	default:
		return "string"
	}
}

func (u *UUID) Name(name string) *UUID {
	u.desc.AttrName = name
	return u
}

func (u *UUID) MinLen(i int) *UUID {
	u.desc.Validators = append(u.desc.Validators, func(v string) error {
		if len(v) < i {
			return errors.New("value is less than the required length")
		}
		return nil
	})
	return u
}

func (u *UUID) Required() *UUID {
	u.desc.Required = true
	return u.MinLen(1)
}

func (u *UUID) Primary(index int) *UUID {
	u.desc.Primary = index
	return u
}

func (u *UUID) Comment(comment string) *UUID {
	u.desc.Comment = comment
	return u
}

func (u *UUID) Default(value string) *UUID {
	u.desc.Default = true
	u.desc.DefaultValue = value
	return u
}

func (s *UUID) ValueType() string {
	return "string"
}

type StringStorage struct {
	value *string
}

func (i *StringStorage) Set(value string) error {
	i.value = &value
	return nil
}

func (i *StringStorage) Get() *string {
	return i.value
}

func (s *StringStorage) Scan(value any) error {
	if value == nil {
		s.value = nil
		return nil
	}
	return convertAssign(&s.value, value)
}

func (s StringStorage) String() string {
	if s.value == nil {
		return "nil"
	}
	return *s.value
}

func (s *StringStorage) Value() entity.FieldValue {
	if s.value == nil {
		return nil
	}
	return *s.value
}
