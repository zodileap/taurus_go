package load

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	stringutil "github.com/yohobala/taurus_go/encoding/string"
	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type (
	// Database 用户定义的Schema中已经处理好的database
	Database struct {
		// Name 数据库的名称。
		Name string
		// Tag 数据库的标签，开发者通过tag来实现数据库的连接和数据库实例的匹配。
		Tag string
		// Type 数据库类型。
		Type dialect.DbDriver
		// EntityMap 数据库中的entity的信息。
		EntityMap EntityMap
		// Entities 数据库中的entity的信息。
		Entities map[string]*Entity
	}

	// Entity 表示了从已经编译好的用户的package中加载的entity
	Entity struct {
		// Name entity的名称
		Name string `json:"name,omitempty"`
		// AttrName entity的属性名称
		AttrName string `json:"attr_name,omitempty"`
		// Config entity配置
		Config entity.EntityConfig `json:"config,omitempty"`
		// Fields entity的字段
		Fields []*Field `json:"fields,omitempty"`
		// ImportPkgs 导入的Go package路径
		ImportPkgs []string
		// Sequences entity的关联序列
		Sequences []entity.Sequence
	}

	// Field 表示entity的字段所包含的信息。
	// 继承了entity.Descriptor
	Field struct {
		entity.Descriptor
		// ValueType 字段的值类型，比如"entity.Int64"的ValueType为"int64"。
		ValueType    string `json:"value_type,omitempty"`
		Validators   int    `json:"validators,omitempty"`
		StoragerType string `json:"storager_type,omitempty"`
		// StoragerOrigType 获得去除泛型后的名字
		StoragerOrigType string `json:"storager_orig_type,omitempty"`
		StoragerPkg      string `json:"storager_pkg,omitempty"`
	}

	// fieldInfo 这个是用于解析字段的类型中Builder和Storager的信息。
	fieldInfo struct {
		Builder  entity.FieldBuilder
		Storager fieldInfoStorager
	}

	// fieldInfoStorager 用于解析字段的类型中的Storager的信息。
	fieldInfoStorager struct {
		Pkg      string
		Name     string
		Type     string
		OrigType string
	}
)

var (
	// ImportPkgs 保存了导入的Go package路径。
	ImportPkgs []string = []string{}
	// db 从Schema中加载的database。
	db *Database = &Database{}
)

// MarshalEntity 将entity.EntityInterface序列化为Entity，用于生成代码。
//
// Params:
//
//   - ei: 实现了entity.EntityInterface的entity。
//
// Returns:
//
//	0: 序列化后的Entity。
//
// ErrCodes:
//
//   - Err_0100020018
func MarshalEntity(ei entity.EntityInterface) (ent *Entity, err error) {
	var (
		entityName string
	)
	config := ei.Config()
	if config.AttrName != "" {
		entityName = config.AttrName
	} else {
		// 利用反射获取entity的名称
		entityName = indirect(reflect.TypeOf(ei)).Name()
	}
	ent = &Entity{
		Name: indirect(reflect.TypeOf(ei)).Name(),
		// 利用反射获取entity的名称
		AttrName:   entityName,
		Config:     ei.Config(),
		ImportPkgs: ImportPkgs,
	}
	ImportPkgs = []string{}
	// 加载entity的字段，调用[entity.EntityInterface]的Fields()方法
	if err := ent.loadFields(ei); err != nil {
		return nil, entity.Err_0100020018.Sprintf(ent.Name, err)
	}
	return ent, nil
}

// Unmarshal 实现了[entity.EntityInterface]的entity反序列化。
//
// Params:
//
//   - b: 序列化后的数据库信息。
//
// Returns:
//
//	0: 反序列化后的Entity。
func Unmarshal(b []byte) (*Database, error) {
	s := &Database{}
	if err := json.Unmarshal(b, s); err != nil {
		return nil, err
	}
	return s, nil
}

// loadFields 从entity.EntityInterface中使用Fields()方法加载字段信息。
//
// Params:
//
//   - ei: 实现了entity.EntityInterface的entity。
func (e *Entity) loadFields(ei entity.EntityInterface) error {
	fields, err := initFields(ei)
	if err != nil {
		return err
	}
	_, err = checkFields(ei)
	if err != nil {
		return err
	}
	for _, f := range fields {
		sf, err := newField(f.Builder, f.Builder.Descriptor())
		if err != nil {
			return err
		}
		if sf.Sequence.Name != nil {
			e.Sequences = append(e.Sequences, sf.Sequence)
		}
		sf.StoragerPkg = f.Storager.Pkg
		sf.StoragerType = f.Storager.Type
		sf.StoragerOrigType = f.Storager.OrigType
		e.Fields = append(e.Fields, sf)
	}
	return nil
}

// initFields 初始化entity的字段，会生成一个初始的Descriptor，这个Descriptor会有一些默认的配置，并传给字段的Init方法。
// 这个方法保证了调用Fields()方法时不会nil pointer dereference。
//
// Params:
//
//   - ei: 实现了entity.EntityInterface的entity。
//
// Returns:
//
//		0: entity中的字段信息。
//	 1: 错误信息。
func initFields(ei entity.EntityInterface) ([]fieldInfo, error) {
	fields := make([]fieldInfo, 0)
	val := reflect.ValueOf(ei)
	// 如果是指针，则获取其指向的元素
	if val.Kind() != reflect.Ptr || val.IsNil() {
		panic("entity must be a non-nil pointer")
	}
	val = val.Elem()
	// 确保指针指向的是结构体
	if val.Kind() != reflect.Struct {
		panic("entity must be a pointer to a struct")
	}
	t := val.Type()
	// 遍历结构体的字段
	for i := 0; i < val.NumField(); i++ {
		// 判断字段是否实现了 entity.Field 接口
		// 同时还要判断字段的指针类型实现了 entity.Field 接口
		// 因为Entity字段可以是指针类型或者是值类型
		fieldVal := val.Field(i)
		fieldName := t.Field(i).Name
		ImportPkgs = append(ImportPkgs, fieldVal.Type().PkgPath())
		fieldVal, ok := assertFieldBuilder(fieldVal)
		if !ok {
			continue
		}
		_, ok = assertFieldStrager(fieldVal)
		if !ok {
			continue
		}

		fe := fieldVal.Elem()
		fy := fe.Type()
		var storager *fieldInfoStorager
		for i := 0; i < fe.NumField(); i++ {
			storager = analyseField(fe.Field(i), fy.Field(i))
			if storager != nil {
				break
			}
		}
		if storager == nil {
			continue
		}

		// 初始化field并传入初始Descriptor
		// 在这里设置了AttrName默认值
		initDesc := &entity.Descriptor{
			Name:     fieldName,
			AttrName: stringutil.ToSnakeCase(fieldName),
			Type:     t.Field(i).Type.String(),
		}
		if fieldVal.IsNil() {
			// 如果字段是 nil，则创建一个新实例
			newInstance := reflect.New(fieldVal.Type().Elem()).Interface()
			if ef, ok := newInstance.(entity.FieldBuilder); ok {
				err := initField(ei, ef, initDesc)
				if err != nil {
					return fields, err
				}
				f := fieldInfo{
					Builder:  ef,
					Storager: *storager,
				}
				fields = append(fields, f)
				// 将新实例赋值给字段
				// fieldVal.Set(reflect.ValueOf(newInstance))
			}

		} else {
			// 如果字段已经被初始化，则只调用 Init 方法
			if ef, ok := fieldVal.Interface().(entity.FieldBuilder); ok {
				err := initField(ei, ef, initDesc)
				if err != nil {
					return fields, err
				}
				f := fieldInfo{
					Builder:  ef,
					Storager: *storager,
				}
				fields = append(fields, f)
			}
		}
	}
	return fields, nil
}

// initField 调用字段的Init方法，获得Descriptor。
//
// Params:
//
//   - ei: 实现了entity.EntityInterface的entity。
//   - f: 实现了entity.FieldBuilder的字段。
//   - initDesc: 初始的Descriptor。
func initField(ei entity.EntityInterface, f entity.FieldBuilder, initDesc *entity.Descriptor) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%T.Init panics: %v", f, v)
		}
	}()
	// 调用字段的Init方法
	f.Init(initDesc)
	desc := f.Descriptor()
	if desc == nil {
		t := reflect.TypeOf(ei).Elem()
		tf := reflect.TypeOf(f).Elem()
		return fmt.Errorf("in the Entity '%s',%s.Descriptor(): nil pointer dereference. Try initialize desc in %s.Init()", t.Name(), tf, tf)
	}
	return nil
}

// checkFields 检查entity的Fields()方法是否有panic，并得到返回值。
//
// Params:
//
//   - fd: 实现了entity.EntityInterface的entity。
//
// Returns:
//
//	0: entity中的字段信息。
//	1: 错误信息。
func checkFields(fd interface {
	Fields() []entity.FieldBuilder
}) (fields []entity.FieldBuilder, err error) {
	defer func() {
		// 如果不是panic那recover为nil
		if v := recover(); v != nil {
			err = fmt.Errorf("%T.Fields panics: %v", fd, v)
			fields = nil
		}
	}()
	return fd.Fields(), nil
}

// newField 根据Descriptor创建新的Field
//
// Params:
//
//   - f: 实现了entity.FieldBuilder的字段。
//   - ed: 字段的Descriptor。
//
// Returns:
//
//	0: 新的Field。
func newField(f entity.FieldBuilder, ed *entity.Descriptor) (*Field, error) {
	attrType := f.AttrType(db.Type)
	if attrType == "" {
		panic(fmt.Sprintf("%T.AttrType in Database %s: empty string", db.Type, f))
	}
	if ed.AttrType == "" {
		ed.AttrType = attrType
	}

	ef := &Field{}
	ef.Name = ed.Name
	ef.AttrName = ed.AttrName
	ef.Type = ed.Type
	ef.AttrType = ed.AttrType
	if size := int64(ed.Size); size != 0 {
		ef.Size = size
	}
	ef.Required = ed.Required
	ef.Primary = ed.Primary
	ef.Comment = ed.Comment
	ef.Default = ed.Default
	ef.DefaultValue = ed.DefaultValue
	ef.Locked = ed.Locked
	ef.Sequence = ed.Sequence
	ef.Validators = len(ed.Validators)
	ef.ValueType = f.ValueType()

	err := checkSequence(ef.Sequence)
	if err != nil {
		return nil, err
	}
	return ef, nil
}

// checkSequence 检查序列的值。
//
// Params:
//
//   - seq: 序列。
func checkSequence(seq entity.Sequence) (err error) {
	if seq.Name != nil && *seq.Name == "" {
		return fmt.Errorf("sequence name is empty")
	}
	if seq.Increament == nil {
		i := int64(1)
		seq.Increament = &i
	}
	if seq.Min == nil {
		i := int64(1)
		seq.Min = &i
	}
	if seq.Max == nil {
		i := int64(9223372036854775807)
		seq.Max = &i
	}
	if seq.Start == nil {
		i := int64(1)
		seq.Start = &i
	}
	if seq.Cache == nil {
		i := int64(1)
		seq.Cache = &i
	}
	return nil
}

// indirect 穿透指针类型，获取不是指针类型的基础类型
//
// Params:
//
//   - t: 反射类型。
//
// Returns:
//
//	0: 不是指针类型的基础类型。
func indirect(t reflect.Type) reflect.Type {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t
}

// analyseField 用于分析entity的字段，提取出里面的builder和storage来。
//
// Params:
//
//   - v: 字段的值。
//   - s: 字段的类型。
//
// Returns:
//
//	0: 字段的Storager的信息。
func analyseField(v reflect.Value, s reflect.StructField) *fieldInfoStorager {
	_, ok := assertFieldBuilder(v)
	if ok {
		return nil
	}
	_, ok = assertFieldStrager(v)
	if ok {
		typeName := s.Type.String()
		OrigTypeName := typeName
		// 提取非泛型部分
		split := strings.Split(typeName, ".")
		if len(split) == 1 {
			OrigTypeName = typeName
		} else {
			OrigTypeName = split[1]
		}
		if strings.Contains(OrigTypeName, "[") && strings.Contains(OrigTypeName, "]") {
			OrigTypeName = OrigTypeName[:strings.Index(OrigTypeName, "[")]
		}
		return &fieldInfoStorager{
			Pkg:      v.Type().PkgPath(),
			Name:     s.Name,
			Type:     typeName,
			OrigType: OrigTypeName,
		}
	}
	return nil
}

// assertFieldBuilder 断言字段是否实现了entity.FieldBuilder接口。
//
// Params:
//
//   - v: 字段的值。
//
// Returns:
//
//	0: 字段的值。
//	1: 是否实现了entity.FieldBuilder接口。
func assertFieldBuilder(v reflect.Value) (reflect.Value, bool) {
	_, ok := v.Interface().(entity.FieldBuilder)
	if !ok {
		if v.CanAddr() {
			v = v.Addr()
			_, ok := v.Interface().(entity.FieldBuilder)
			if ok {
				return v, true
			}
		}
	} else {
		return v, true
	}
	return v, false
}

// assertFieldStrager 断言字段是否实现了entity.FieldStorager接口。
//
// Params:
//
//   - v: 字段的值。
//
// Returns:
//
//	0: 字段的值。
//	1: 是否实现了entity.FieldStorager接口。
func assertFieldStrager(v reflect.Value) (reflect.Value, bool) {
	_, ok := v.Interface().(entity.FieldStorager)
	if !ok {
		if v.CanAddr() {
			v = v.Addr()
			_, ok := v.Interface().(entity.FieldStorager)
			if ok {
				return v, true
			}
		}
	} else {
		return v, true
	}
	return v, false
}
