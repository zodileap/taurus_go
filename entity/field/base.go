package field

import (
	"database/sql/driver"
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
	var zero T
	refType := reflect.TypeOf(zero)

	// 处理数组、切片、指针类型
	for refType.Kind() == reflect.Array || refType.Kind() == reflect.Slice || refType.Kind() == reflect.Ptr {
		if refType.Kind() == reflect.Ptr {
			refType = refType.Elem()
		} else {
			refType = refType.Elem()
			depth++
		}
	}

	// 使用 PkgPath() + "." + Name() 获取完整类型名
	if refType.PkgPath() != "" {
		return depth, refType.PkgPath() + "." + refType.Name()
	}

	// 对于内置类型，直接返回 String()
	return depth, refType.String()
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

// Unique 设置字段为唯一字段或参与联合唯一约束。
// 相同的序号表示这些字段组成联合唯一约束。
func (b *BaseBuilder[T]) Unique(index int) *BaseBuilder[T] {
	b.desc.Uniques = append(b.desc.Uniques, index)
	return b
}

// Index 设置字段为索引。
func (b *BaseBuilder[T]) Index(index int) *BaseBuilder[T] {
	// 追加而不是覆盖
	b.desc.Indexes = append(b.desc.Indexes, index)
	return b
}

// IndexName 设置索引名称。
//
// Params:
//   - name: 索引名称。
func (b *BaseBuilder[T]) IndexName(name string) *BaseBuilder[T] {
	b.desc.IndexName = name
	return b
}

// IndexMethod 设置索引方法。
//
// Params:
//   - method: 索引方法,如"btree","hash"等。
func (b *BaseBuilder[T]) IndexMethod(method string) *BaseBuilder[T] {
	b.desc.IndexMethod = method
	return b
}

// CheckBuilder 是检查约束的构建器函数类型
type CheckBuilder func(fieldName string) string

// Check 添加CHECK约束到字段
// 参数 builder 是一个回调函数，接收字段名并返回CHECK约束的内容
func (b *BaseBuilder[T]) Check(builder CheckBuilder) *BaseBuilder[T] {
	if b == nil {
		panic("taurus_go/entity field check: nil pointer dereference.")
	}
	// 获取字段在数据库中的实际名称
	fieldName := b.desc.AttrName
	if fieldName == "" {
		fieldName = b.desc.Name
	}

	// 调用构建器生成约束内容
	constraint := builder(fieldName)
	b.desc.CheckConstraint = constraint

	return b
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
	return i.toValue(dbType)
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
func (b *BaseStorage[T]) toValue(dbType dialect.DbDriver) (entity.FieldValue, error) {
	switch dbType {
	case dialect.PostgreSQL:
		if b.value == nil {
			return nil, nil
		}
		v := reflect.ValueOf(*b.value)
		switch v.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return int64(v.Int()), nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return int64(v.Uint()), nil
		case reflect.Bool:
			return v.Bool(), nil
		case reflect.String:
			return v.String(), nil
		case reflect.Float32, reflect.Float64:
			return v.Float(), nil
		case reflect.Slice:
			return handleSlice(v)
		default:
			return nil, fmt.Errorf("unsupported database type: %v", reflect.TypeOf(*b.value))
		}
	default:
		return nil, fmt.Errorf("unsupported database type: %v", reflect.TypeOf(*b.value))
	}
}

func handleSlice(v reflect.Value) (driver.Value, error) {
	switch v.Type().Elem().Kind() {
	case reflect.Int16, reflect.Int32, reflect.Int64:
		return arrayToPGString(v.Interface(), func(a any) (string, error) {
			return fmt.Sprintf("%d", a), nil
		})
	case reflect.Bool:
		return arrayToPGString(v.Interface(), func(a any) (string, error) {
			if a.(bool) {
				return "true", nil
			}
			return "false", nil
		})
	case reflect.Slice:
		// 处理嵌套数组
		return handleNestedSlice(v)
	default:
		return nil, fmt.Errorf("unsupported slice type: %v", v.Type())
	}
}

func handleNestedSlice(v reflect.Value) (driver.Value, error) {
	depth := getSliceDepth(v.Type())
	if depth > 5 {
		return nil, fmt.Errorf("slice nesting too deep: %d", depth)
	}
	var convFunc func(a any) (string, error)
	elemKind := getDeepestSliceElemKind(v.Type())
	switch elemKind {
	case reflect.Int16, reflect.Int32, reflect.Int64:
		convFunc = func(a any) (string, error) { return fmt.Sprintf("%d", a), nil }
	case reflect.Bool:
		convFunc = func(a any) (string, error) {
			if a.(bool) {
				return "true", nil
			}
			return "false", nil
		}
	default:
		return nil, fmt.Errorf("unsupported nested slice element type: %v", elemKind)
	}
	return arrayToPGString(v.Interface(), convFunc)
}

func getSliceDepth(t reflect.Type) int {
	depth := 0
	for t.Kind() == reflect.Slice {
		depth++
		t = t.Elem()
	}
	return depth
}

func getDeepestSliceElemKind(t reflect.Type) reflect.Kind {
	for t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	return t.Kind()
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
	case []string:
		return "[]string"
	case [][]string:
		return "[][]string"
	case [][][]string:
		return "[][][]string"
	case [][][][]string:
		return "[][][][]string"
	case [][][][][]string:
		return "[][][][][]string"
	default:
		return ""
	}
}
