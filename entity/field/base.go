package field

import (
	"fmt"
	"reflect"
	"time"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type BaseBuilder[T any] struct {
	desc *entity.Descriptor
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (b *BaseBuilder[T]) Init(desc *entity.Descriptor) error {
	if b == nil {
		panic("taurus_go/entity field init: nil pointer dereference.")
	}

	depth, elemType := b.sliceTypeDetails()
	desc.Depth = depth
	desc.BaseType = elemType
	b.desc = desc
	return nil
}

func (b *BaseBuilder[T]) sliceTypeDetails() (int, string) {

	depth := 0
	// Safely get the type of T, checking for nil pointer dereference safety.
	var zero T
	refType := reflect.TypeOf(zero) // Use a zero value of T instead of nil pointer
	// Process the type to determine depth and base type.
	for refType.Kind() == reflect.Array || refType.Kind() == reflect.Slice || refType.Kind() == reflect.Ptr {
		if refType.Kind() == reflect.Ptr {
			refType = refType.Elem() // Dereference pointer type
		} else {
			refType = refType.Elem()
			depth++ // Increment depth for arrays or slices
		}
	}
	return depth, refType.Name()
}

// Descriptor 获取字段的描述信息。
func (b *BaseBuilder[T]) Descriptor() *entity.Descriptor {
	return b.desc
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
func (b *BaseBuilder[T]) AttrType(dbType dialect.DbDriver) string {
	var t T
	return attrType(t, dbType)
}

// ValueType 用于设置字段的值在go中类型名称。例如entity.Int64的ValueType为"int64"。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (b *BaseBuilder[T]) ValueType() string {
	var t T
	return valueType(t)
}

// ExtTemplate 用于在使用字段时，调用外部模版生成代码，
// 这个相比在 go run github.com/yohobala/taurus_go/entity/cmd generate -t <template>，
// `ExtTemplate`是和字段相关联，只要调用字段就会生成代码，避免了每次都要手动调用模版。
func (b *BaseBuilder[T]) ExtTemplate() []string {
	return []string{}
}

type BaseStorage[T any] struct {
	value *T
}

// Set 设置字段的值。
func (b *BaseStorage[T]) Set(value T) error {
	b.value = &value
	return nil
}

// Get 获取字段的值。
func (b *BaseStorage[T]) Get() *T {
	return b.value
}

// Scan 从数据库中读取字段的值。
func (b *BaseStorage[T]) Scan(value interface{}) error {
	if value == nil {
		b.value = nil
		return nil
	}
	return convertAssign(b.value, value)
}

// String 返回字段的字符串表示。
func (b BaseStorage[T]) String() string {
	if b.value == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", *b.value)
}

// SqlParam 用于sql中获取字段参数并赋值。如 INSERT INTO "blog" ( "desc") VALUES ($1)，给$1传递具体的值。
func (i *BaseStorage[T]) SqlParam(dbType dialect.DbDriver) (entity.FieldValue, error) {
	var t T
	return i.toValue(t, dbType)
}

// SqlFormatParam 用于sql中获取字段的值的格式化字符串。如 INSERT INTO "blog" ( "desc" ) VALUES ( ST_GeomFromGeoJSON($1) ) 中添加的ST_GeomFromGeoJSON()。
func (i *BaseStorage[T]) SqlFormatParam() func(dbType dialect.DbDriver, param string) string {
	return func(dbType dialect.DbDriver, param string) string {
		return param
	}
}

// SqlSelectClause 用于sql语句中获取字段的select子句部分，通过这个能够扩展SELECT部分实现复杂的查询，比如 SELECT id, ST_AsText(point)。
func (i *BaseStorage[T]) SqlSelectFormat() func(dbType dialect.DbDriver, name string) string {
	return func(dbType dialect.DbDriver, name string) string {
		return name
	}
}

// toValue 将字段的值转换为数据库中的值。
//
// Params:
//
//   - t: 字段的值。
//   - dbType: 数据库类型。
func (b *BaseStorage[T]) toValue(t any, dbType dialect.DbDriver) (entity.FieldValue, error) {
	switch dbType {
	case dialect.PostgreSQL:
		if b.value == nil {
			return nil, nil
		}
		switch t.(type) {
		case int16:
			return *b.value, nil
		case int32:
			return *b.value, nil
		case int64:
			return *b.value, nil
		case []int16:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case []int32:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case []int64:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][]int16:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][]int32:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][]int64:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][]int16:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][]int32:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][]int64:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][][]int16:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][][]int32:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][][]int64:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][][][]int16:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][][][]int32:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case [][][][][]int64:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				return fmt.Sprintf("%d", a), nil
			})
		case bool:
			return *b.value, nil
		case []bool:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				if a.(bool) {
					return "true", nil
				}
				return "false", nil
			})
		case [][]bool:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				if a.(bool) {
					return "true", nil
				}
				return "false", nil
			})
		case [][][]bool:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				if a.(bool) {
					return "true", nil
				}
				return "false", nil
			})
		case [][][][]bool:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				if a.(bool) {
					return "true", nil
				}
				return "false", nil
			})
		case [][][][][]bool:
			return arrayToPGString(*b.value, func(a any) (string, error) {
				if a.(bool) {
					return "true", nil
				}
				return "false", nil
			})
		case string:
			return *b.value, nil
		default:
			return nil, fmt.Errorf("unsupported database type: %v", reflect.TypeOf(t))
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %v", reflect.TypeOf(t))
	}
}

// attrType 返回字段的数据库中的类型名。
func attrType(t any, dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		switch t.(type) {
		case int16:
			return "int2"
		case int32:
			return "int4"
		case int64:
			return "int8"
		case []int16:
			return "int2[]"
		case []int32:
			return "int4[]"
		case []int64:
			return "int8[]"
		case [][]int16:
			return "int2[][]"
		case [][]int32:
			return "int4[][]"
		case [][]int64:
			return "int8[][]"
		case [][][]int16:
			return "int2[][][]"
		case [][][]int32:
			return "int4[][][]"
		case [][][]int64:
			return "int8[][][]"
		case [][][][]int16:
			return "int2[][][][]"
		case [][][][]int32:
			return "int4[][][][]"
		case [][][][]int64:
			return "int8[][][][]"
		case [][][][][]int16:
			return "int2[][][][][]"
		case [][][][][]int32:
			return "int4[][][][][]"
		case [][][][][]int64:
			return "int8[][][][][]"
		case bool:
			return "boolean"
		case []bool:
			return "boolean[]"
		case [][]bool:
			return "boolean[][]"
		case [][][]bool:
			return "boolean[][][]"
		case [][][][]bool:
			return "boolean[][][][]"
		case [][][][][]bool:
			return "boolean[][][][][]"
		default:
			return ""
		}
	default:
		return ""
	}

}

// valueType 返回字段的值在go中类型名称。
func valueType(t any) string {
	switch t.(type) {
	case int16:
		return "int16"
	case int32:
		return "int32"
	case int64:
		return "int64"
	case []int16:
		return "[]int16"
	case []int32:
		return "[]int32"
	case []int64:
		return "[]int64"
	case [][]int16:
		return "[][]int16"
	case [][]int32:
		return "[][]int32"
	case [][]int64:
		return "[][]int64"
	case [][][]int16:
		return "[][][]int16"
	case [][][]int32:
		return "[][][]int32"
	case [][][]int64:
		return "[][][]int64"
	case [][][][]int16:
		return "[][][][]int16"
	case [][][][]int32:
		return "[][][][]int32"
	case [][][][]int64:
		return "[][][][]int64"
	case [][][][][]int16:
		return "[][][][][]int16"
	case [][][][][]int32:
		return "[][][][][]int32"
	case [][][][][]int64:
		return "[][][][][]int64"
	case bool:
		return "bool"
	case []bool:
		return "[]bool"
	case [][]bool:
		return "[][]bool"
	case [][][]bool:
		return "[][][]bool"
	case [][][][]bool:
		return "[][][][]bool"
	case [][][][][]bool:
		return "[][][][][]bool"
	case string:
		return "string"
	case time.Time:
		return "time.Time"
	case []time.Time:
		return "[]time.Time"
	case [][]time.Time:
		return "[][]time.Time"
	case [][][]time.Time:
		return "[][][]time.Time"
	case [][][][]time.Time:
		return "[][][][]time.Time"
	case [][][][][]time.Time:
		return "[][][][][]time.Time"
	default:
		return ""
	}
}
