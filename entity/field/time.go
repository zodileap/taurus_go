package field

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

// Timestamp 时间戳类型的字段。
type Timestamptz struct {
	TimestamptzBuilder[time.Time]
	TimestampStorage[time.Time]
}

type TimestamptzA1 struct {
	TimestamptzBuilder[[]time.Time]
	TimestampStorage[[]time.Time]
}

type TimestamptzA2 struct {
	TimestamptzBuilder[[][]time.Time]
	TimestampStorage[[][]time.Time]
}

type TimestamptzA3 struct {
	TimestamptzBuilder[[][][]time.Time]
	TimestampStorage[[][][]time.Time]
}

type TimestamptzA4 struct {
	TimestamptzBuilder[[][][][]time.Time]
	TimestampStorage[[][][][]time.Time]
}

type TimestamptzA5 struct {
	TimestamptzBuilder[[][][][][]time.Time]
	TimestampStorage[[][][][][]time.Time]
}

// TimestamptzBuilder 时间戳类型的字段构建器。
type TimestamptzBuilder[T any] struct {
	BaseBuilder[T]
	// 精度
	precision int
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (t *TimestamptzBuilder[T]) Init(desc *entity.Descriptor) error {
	if t == nil {
		panic("taurus_go/entity Timestamptz init: nil pointer dereference.")
	}
	t.desc = desc
	return nil
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
func (t *TimestamptzBuilder[T]) AttrType(dbType dialect.DbDriver) string {
	var v T
	return t.attrType(v, dbType)
}

func (t *TimestamptzBuilder[T]) attrType(v any, dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		switch v.(type) {
		case time.Time:
			if t.precision == 0 {
				t.precision = 6
			}
			return fmt.Sprintf("timestamptz(%d)", t.precision)
		case []time.Time:
			return "timestamptz[]"
		case [][]time.Time:
			return "timestamptz[][]"
		case [][][]time.Time:
			return "timestamptz[][][]"
		case [][][][]time.Time:
			return "timestamptz[][][][]"
		case [][][][][]time.Time:
			return "timestamptz[][][][][]"
		default:
			return ""
		}
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
func (t *TimestamptzBuilder[T]) Name(name string) *TimestamptzBuilder[T] {
	t.desc.AttrName = name
	return t
}

// MinLen 设置字段的最小长度。
//
// Params:
//
//   - size: 字段的最小长度。
func (t *TimestamptzBuilder[T]) MinLen(size int) *TimestamptzBuilder[T] {
	t.desc.Validators = append(t.desc.Validators, func(b []byte) error {
		if len(b) < size {
			return errors.New("value is less than the required length.")
		}
		return nil
	})
	return t
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (t *TimestamptzBuilder[T]) Required() *TimestamptzBuilder[T] {
	t.desc.Required = true
	return t.MinLen(1)
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (t *TimestamptzBuilder[T]) Primary(index int) *TimestamptzBuilder[T] {
	t.desc.Required = true
	t.desc.Primary = index
	return t
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (t *TimestamptzBuilder[T]) Comment(comment string) *TimestamptzBuilder[T] {
	t.desc.Comment = comment
	return t
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (t *TimestamptzBuilder[T]) Default(value string) *TimestamptzBuilder[T] {
	t.desc.Default = true
	t.desc.DefaultValue = value
	return t
}

// Precision 设置时间精度。
func (t *TimestamptzBuilder[T]) Precision(precision int) *TimestamptzBuilder[T] {
	t.precision = precision
	return t
}

// Locked 设置字段是否为只读。
func (t *TimestamptzBuilder[T]) Locked() *TimestamptzBuilder[T] {
	t.desc.Locked = true
	return t
}

// TimestampStorage 时间戳类型的字段存储。
type TimestampStorage[T any] struct {
	BaseStorage[T]
}

// SqlParam 用于sql中获取字段参数并赋值。如 INSERT INTO "blog" ( "desc") VALUES ($1)，给$1传递具体的值。
func (i *TimestampStorage[T]) SqlParam(dbType dialect.DbDriver) (entity.FieldValue, error) {
	if i.value == nil {
		return nil, nil
	}
	return i.toValue(*i.value, dbType)
}

func (b *TimestampStorage[T]) toValue(v any, dbType dialect.DbDriver) (entity.FieldValue, error) {
	switch dbType {
	case dialect.PostgreSQL:
		if b.value == nil {
			return nil, nil
		}
		switch val := v.(type) {
		case time.Time:
			return val.Format("2006-01-02 15:04:05.999999-07:00"), nil
		case []time.Time:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return a.(time.Time).Format("2006-01-02 15:04:05.999999-07:00"), nil
			})
		case [][]time.Time:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return a.(time.Time).Format("2006-01-02 15:04:05.999999-07:00"), nil
			})
		case [][][]time.Time:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return a.(time.Time).Format("2006-01-02 15:04:05.999999-07:00"), nil
			})
		case [][][][]time.Time:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return a.(time.Time).Format("2006-01-02 15:04:05.999999-07:00"), nil
			})
		case [][][][][]time.Time:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return a.(time.Time).Format("2006-01-02 15:04:05.999999-07:00"), nil
			})
		default:
			return nil, fmt.Errorf("unsupported database type: %v", reflect.TypeOf(v))
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %v", reflect.TypeOf(v))
	}
}
