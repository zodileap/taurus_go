package field

import (
	"errors"
	"fmt"
	"time"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type Timestamptz struct {
	TimestamptzBuilder
	TimestampStorage
}

type TimestamptzBuilder struct {
	desc      *entity.Descriptor
	precision int
}

func (t *TimestamptzBuilder) Init(desc *entity.Descriptor) error {
	if t == nil {
		panic("taurus_go/entity Timestamptz init: nil pointer dereference.")
	}
	t.desc = desc
	return nil
}

func (t *TimestamptzBuilder) Descriptor() *entity.Descriptor {
	return t.desc
}

func (t *TimestamptzBuilder) AttrType(dbType dialect.DbDriver) string {
	if t.precision == 0 {
		t.precision = 6
	}
	switch dbType {
	case dialect.PostgreSQL:
		return fmt.Sprintf("timestamptz(%d)", t.precision)
	default:
		return ""
	}
}

func (t *TimestamptzBuilder) Name(name string) *TimestamptzBuilder {
	t.desc.AttrName = name
	return t
}

// 设置字段的最小长度。
func (t *TimestamptzBuilder) MinLen(size int) *TimestamptzBuilder {
	t.desc.Validators = append(t.desc.Validators, func(b []byte) error {
		if len(b) < size {
			return errors.New("value is less than the required length.")
		}
		return nil
	})
	return t
}

func (t *TimestamptzBuilder) Required() *TimestamptzBuilder {
	t.desc.Required = true
	return t.MinLen(1)
}

func (t *TimestamptzBuilder) Primary(index int) *TimestamptzBuilder {
	t.desc.Required = true
	t.desc.Primary = index
	return t
}

// 设置字段的注释。
func (t *TimestamptzBuilder) Comment(comment string) *TimestamptzBuilder {
	t.desc.Comment = comment
	return t
}

// 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
func (t *TimestamptzBuilder) Default(value string) *TimestamptzBuilder {
	t.desc.Default = true
	t.desc.DefaultValue = value
	return t
}

func (t *TimestamptzBuilder) Precision(precision int) *TimestamptzBuilder {
	t.precision = precision
	return t
}

func (t *TimestamptzBuilder) ValueType() string {
	return "time.Time"
}

type TimestampStorage struct {
	value *time.Time
}

func (t *TimestampStorage) Set(value time.Time) error {
	t.value = &value
	return nil
}

func (t *TimestampStorage) Get() *time.Time {
	return t.value
}

func (t *TimestampStorage) Scan(value interface{}) error {
	if value == nil {
		t.value = nil
		return nil
	}
	return convertAssign(&t.value, value)
}

func (t TimestampStorage) String() string {
	if t.value == nil {
		return "nil"
	}
	return t.value.String()
}

func (t *TimestampStorage) Value() entity.FieldValue {
	if t.value == nil {
		return nil
	}
	return *t.value
}
