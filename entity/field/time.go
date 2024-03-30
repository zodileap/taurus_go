package field

import (
	"errors"
	"fmt"
	"time"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

// Timestamp 时间戳类型的字段。
type Timestamptz struct {
	TimestamptzBuilder
	TimestampStorage
}

// TimestampStorage 时间戳类型的字段存储。
type TimestampStorage struct {
	value *time.Time
}

// Set 设置字段的值。
func (t *TimestampStorage) Set(value time.Time) error {
	t.value = &value
	return nil
}

// Get 获取字段的值。
func (t *TimestampStorage) Get() *time.Time {
	return t.value
}

// Scan 从数据库中读取字段的值。
func (t *TimestampStorage) Scan(value interface{}) error {
	if value == nil {
		t.value = nil
		return nil
	}
	return convertAssign(&t.value, value)
}

// String 返回字段的字符串表示。
func (t TimestampStorage) String() string {
	if t.value == nil {
		return "nil"
	}
	return t.value.String()
}

// Value 返回字段的值，和Get方法不同的是，Value方法返回的是接口类型。
func (t *TimestampStorage) Value() entity.FieldValue {
	if t.value == nil {
		return nil
	}
	return *t.value
}

// TimestamptzBuilder 时间戳类型的字段构建器。
type TimestamptzBuilder struct {
	desc *entity.Descriptor
	// 精度
	precision int
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (t *TimestamptzBuilder) Init(desc *entity.Descriptor) error {
	if t == nil {
		panic("taurus_go/entity Timestamptz init: nil pointer dereference.")
	}
	t.desc = desc
	return nil
}

// Descriptor 获取字段的描述信息。
func (t *TimestamptzBuilder) Descriptor() *entity.Descriptor {
	return t.desc
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

// ValueType 用于设置字段的值在go中类型名称。例如entity.Int64的ValueType为"int64"。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (t *TimestamptzBuilder) ValueType() string {
	return "time.Time"
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (t *TimestamptzBuilder) Name(name string) *TimestamptzBuilder {
	t.desc.AttrName = name
	return t
}

// MinLen 设置字段的最小长度。
//
// Params:
//
//   - size: 字段的最小长度。
func (t *TimestamptzBuilder) MinLen(size int) *TimestamptzBuilder {
	t.desc.Validators = append(t.desc.Validators, func(b []byte) error {
		if len(b) < size {
			return errors.New("value is less than the required length.")
		}
		return nil
	})
	return t
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (t *TimestamptzBuilder) Required() *TimestamptzBuilder {
	t.desc.Required = true
	return t.MinLen(1)
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (t *TimestamptzBuilder) Primary(index int) *TimestamptzBuilder {
	t.desc.Required = true
	t.desc.Primary = index
	return t
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (t *TimestamptzBuilder) Comment(comment string) *TimestamptzBuilder {
	t.desc.Comment = comment
	return t
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (t *TimestamptzBuilder) Default(value string) *TimestamptzBuilder {
	t.desc.Default = true
	t.desc.DefaultValue = value
	return t
}

// Precision 设置时间精度。
func (t *TimestamptzBuilder) Precision(precision int) *TimestamptzBuilder {
	t.precision = precision
	return t
}

// Locked 设置字段是否为只读。
func (t *TimestamptzBuilder) Locked() *TimestamptzBuilder {
	t.desc.Locked = true
	return t
}
