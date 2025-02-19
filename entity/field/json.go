package field

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
	"github.com/yohobala/taurus_go/tlog"
)

// JSON 用于定义 json 类型的字段
type JSON[T any] struct {
	JSONBuilder[T]
	JSONStorage[T]
}

// JSONBuilder JSON类型的字段构建器
type JSONBuilder[T any] struct {
	BaseBuilder[T]
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (j *JSONBuilder[T]) Name(name string) *JSONBuilder[T] {
	j.desc.AttrName = name
	return j
}

// Init 初始化字段的描述信息。
//
// Params:
//
//   - desc: 字段的描述信息。
func (j *JSONBuilder[T]) Init(desc *entity.Descriptor) error {
	if j == nil {
		panic("taurus_go/entity JSON init: nil pointer dereference.")
	}
	j.desc = desc
	return nil
}

// AttrType 获取字段的数据库中的类型名。
//
// Params:
//
//   - dbType: 数据库类型。
//
// Returns:
//
//   - 字段的数据库中的类型名。
func (j *JSONBuilder[T]) AttrType(dbType dialect.DbDriver) string {
	var v T
	return j.attrType(v, dbType)
}

func (j *JSONBuilder[T]) attrType(v any, dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		return "json"
	default:
		return ""
	}
}

// ValueType 用于设置字段的值在go中类型名称。JSON默认返回any。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (j *JSONBuilder[T]) ValueType() string {
	var v T
	return reflect.TypeOf(v).String()
}

// MinLen 设置字段的最小长度。
//
// Params:
//
//   - size: 字段的最小长度。
func (j *JSONBuilder[T]) MinLen(size int) *JSONBuilder[T] {
	j.desc.Validators = append(j.desc.Validators, func(b []byte) error {
		if len(b) < size {
			return errors.New("value is less than the required length.")
		}
		return nil
	})
	return j
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (j *JSONBuilder[T]) Required() *JSONBuilder[T] {
	j.desc.Required = true
	return j.MinLen(1)
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (j *JSONBuilder[T]) Primary(index int) *JSONBuilder[T] {
	j.desc.Required = true
	j.desc.Primary = index
	return j
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (j *JSONBuilder[T]) Comment(comment string) *JSONBuilder[T] {
	j.desc.Comment = comment
	return j
}

// Default 设置字段的默认值。
// 如果设置了默认值，则在插入数据时，如果没有设置字段的值，则会使用默认值。
//
// Params:
//
//   - value: 字段的默认值。
func (j *JSONBuilder[T]) Default(value string) *JSONBuilder[T] {
	j.desc.Default = true
	j.desc.DefaultValue = value
	return j
}

// Locked 设置字段为只读字段。
func (j *JSONBuilder[T]) Locked() *JSONBuilder[T] {
	j.desc.Locked = true
	return j
}

// JSONStorage JSON类型的字段存储。
type JSONStorage[T any] struct {
	BaseStorage[T]
}

// SqlParam 用于sql中获取字段参数并赋值。
func (j *JSONStorage[T]) SqlParam(dbType dialect.DbDriver) (entity.FieldValue, error) {
	if j.value == nil {
		return nil, nil
	}
	switch dbType {
	case dialect.PostgreSQL:
		data, err := json.Marshal(*j.value)
		if err != nil {
			return nil, fmt.Errorf("marshal json error: %v", err)
		}
		return string(data), nil
	default:
		return nil, fmt.Errorf("unsupported database type for JSON: %v", reflect.TypeOf(*j.value))
	}
}

// SqlFormatParam 用于sql中获取字段的值的格式化字符串。
func (j *JSONStorage[T]) SqlFormatParam() func(dbType dialect.DbDriver, param string) string {
	return func(dbType dialect.DbDriver, param string) string {
		return param
	}
}

// SqlSelectFormat 用于sql语句中获取字段的select子句部分。
func (j *JSONStorage[T]) SqlSelectFormat() func(dbType dialect.DbDriver, name string) string {
	return func(dbType dialect.DbDriver, name string) string {
		return name
	}
}

func (j *JSONStorage[T]) Scan(value interface{}) error {
	if value == nil {
		j.value = nil
		return nil
	}
	tlog.Print(value)
	return json.Unmarshal([]byte(value.([]byte)), &j.value)
}
