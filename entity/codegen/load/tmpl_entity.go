package load

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	stringutil "github.com/yohobala/taurus_go/datautil/string"
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
		// 添加触发器配置
		Triggers []entity.TriggerConfig
	}

	// Entity 表示了从已经编译好的用户的package中加载的entity
	Entity struct {
		// Name entity的名称
		Name string `json:"name,omitempty"`
		// AttrName entity的属性名称
		AttrName string `json:"attr_name,omitempty"`
		// Comment entity的注释
		Comment string `json:"comment,omitempty"`
		// Config entity配置
		Config entity.EntityConfig `json:"config,omitempty"`
		// Fields entity的字段
		Fields []*Field `json:"fields,omitempty"`
		// ImportPkgs 导入的Go package路径
		ImportPkgs []string
		// Sequences entity的关联序列
		Sequences []entity.Sequence
		// Relations entity的关系
		Relations []*Relation
	}

	// Field 表示entity的字段所包含的信息。
	// 继承了entity.Descriptor
	Field struct {
		entity.Descriptor
		// Validators 字段的验证器数量
		Validators int `json:"validators,omitempty"`
		// StoragerType 字段的存储器的类型，这个是字段的作用是关联已经定义的好的存储器。比如field.IntStorage[int16]
		StoragerType string `json:"storager_type,omitempty"`
		// StoragerOrigType 字段的存储器去除泛型后的名字，比如field.IntStorage[int16]变成field.IntStorage
		StoragerOrigType string `json:"storager_orig_type,omitempty"`
		// StoragerPkg 字段的存储器的包路径
		StoragerPkg string `json:"storager_pkg,omitempty"`
		// Templates 字段关联的额外模版
		Templates []string `json:"templates,omitempty"`
		// Tag 字段的标签信息
		Tag string `json:"tag,omitempty"`
	}

	// Relation 表示entity之间的关系
	Relation struct {
		// Desc entity在自定义时的信息
		Desc RelationDesc
		// Principal 主体实体
		Principal RelationEntity
		// Dependent 依赖实体
		Dependent RelationEntity
	}

	// RelationEntity 存储有关系的entity的信息
	RelationEntity struct {
		Name string `json:"name,omitempty"`
		// AttrName entity的属性名称
		AttrName string `json:"attr_name,omitempty"`
		Field    *Field
		Rel      entity.Rel
	}

	// RelationDesc 表示entity之间的关系的描述，
	// 不用entity.RelationshipDescriptor是因为entity.RelationshipDescriptor有一些字段是接口类型，
	// 在序列化时会产生异常。
	RelationDesc struct {
		Has          Entity
		With         Entity
		HasRel       entity.Rel
		WithRel      entity.Rel
		ForeignKey   Field
		ReferenceKey Field
		Constraint   string
		Update       string
		Delete       string
	}

	// entityInfo 这个是用于解析字段的类型中Builder和Storager的信息。
	entityInfo struct {
		field *fieldInfo
	}

	// fieldInfo entiy中字段的信息
	fieldInfo struct {
		Tag      string
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
	// PostgreSQL保留关键字
	postgresKeywords = []string{
		"all", "analyse", "analyze", "and", "any", "array", "as", "asc",
		"asymmetric", "authorization", "binary", "both", "case", "cast",
		"check", "collate", "column", "constraint", "create", "cross",
		"current_date", "current_role", "current_time", "current_timestamp",
		"current_user", "default", "deferrable", "desc", "distinct", "do",
		"else", "end", "except", "false", "for", "foreign", "freeze", "from",
		"full", "grant", "group", "having", "ilike", "in", "initially", "inner",
		"intersect", "into", "is", "isnull", "join", "leading", "left", "like",
		"limit", "localtime", "localtimestamp", "natural", "not", "notnull",
		"null", "offset", "on", "only", "or", "order", "outer", "overlaps",
		"placing", "primary", "references", "right", "select", "session_user",
		"similar", "some", "symmetric", "table", "then", "to", "trailing",
		"true", "union", "unique", "user", "using", "when", "where", "with",
	}
	// Go语言关键字和保留字
	goKeywords = []string{
		"break", "default", "func", "interface", "select", "case", "defer", "go", "map", "struct",
		"chan", "else", "goto", "package", "switch", "const", "fallthrough", "if", "range", "type",
		"continue", "for", "import", "return", "var",
	}
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
		panic("entity must set AttrName in Config() method")
	}
	ent = &Entity{
		Name: indirect(reflect.TypeOf(ei)).Name(),
		// 利用反射获取entity的名称
		AttrName: entityName,
		Comment:  config.Comment,
		Config:   config,
	}
	ImportPkgs = []string{}
	// 加载entity的字段，调用[entity.EntityInterface]的Fields()方法
	if err := ent.loadEntity(ei); err != nil {
		return nil, err
	}

	for _, f := range ent.Fields {
		ImportPkgs = append(ImportPkgs, f.StoragerPkg)
	}
	for _, pkgName := range ImportPkgs {
		ps := strings.Split(pkgName, "/")
		p := ps[len(ps)-1]
		if entityName == p {
			return nil, entity.Err_0100020018.Sprintf(ent.Name, fmt.Sprintf("entity name %q is the same as package name %q", entityName, pkgName))
		}
	}
	ent.ImportPkgs = ImportPkgs
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

// loadRelationship 加载entity的关系。这个用于确定entity之间的关系，并在entity中添加关系。
func (db *Database) loadRelationship(di entity.DbInterface) (err error) {
	rels, err := checkRelationships(di)
	if err != nil {
		return err
	}
	var rs []Relation
	for _, r := range rels {
		var err error
		desc := r.Descriptor()
		err = checkRelationDescriptor(desc)
		if err != nil {
			return err
		}
		rel := desc.WithRel<<2 | desc.HasRel
		var principal entity.EntityInterface
		var dependent entity.EntityInterface
		var principalEntity *Entity
		var dependentEntity *Entity
		var principalField *Field
		var dependentField *Field
		var principalRel entity.Rel
		var dependentRel entity.Rel
		if rel == entity.M2O {
			principal = desc.Has
			principalRel = desc.HasRel
			dependent = desc.With
			dependentRel = desc.WithRel
		} else {
			principal = desc.With
			principalRel = desc.WithRel
			dependent = desc.Has
			dependentRel = desc.HasRel
		}
		principalEntity, err = db.extractEntity(principal)
		if err != nil {
			return err
		}
		dependentEntity, err = db.extractEntity(dependent)
		if err != nil {
			return err
		}
		dependentField, err = db.extractRelField(desc.ForeignKey, dependent)
		if err != nil {
			return err
		}
		principalField, err = db.extractRelField(desc.ReferenceKey, principal)
		if err != nil {
			return err
		}
		if (dependentField.StoragerType != principalField.StoragerType) || (dependentField.StoragerPkg != principalField.StoragerPkg) {
			// 输出两个字段的类型不一致的提示
			return entity.Err_0100020019.Sprintf(dependentField.EntityName, dependentField.Name, dependentField.StoragerType, principalField.EntityName, principalField.Name, principalField.StoragerType)
		}

		r := Relation{
			Principal: RelationEntity{
				Name:     principalEntity.Name,
				AttrName: principalEntity.AttrName,
				Field:    principalField,
				Rel:      principalRel,
			},
			Dependent: RelationEntity{
				Name:     dependentEntity.Name,
				AttrName: dependentEntity.AttrName,
				Field:    dependentField,
				Rel:      dependentRel,
			},
			Desc: RelationDesc{
				Has:          *principalEntity,
				With:         *dependentEntity,
				HasRel:       desc.HasRel,
				WithRel:      desc.WithRel,
				ForeignKey:   *dependentField,
				ReferenceKey: *principalField,
				Constraint:   desc.ConstraintName,
				Update:       desc.Update,
				Delete:       desc.Delete,
			},
		}
		rs = append(rs, r)
	}
	err = db.addRelationship(rs)
	if err != nil {
		return err
	}
	return nil
}

// addRelationship 添加关系。把关系添加到依赖实体中，因为生成的sql语句是在依赖实体中生成的。
func (db *Database) addRelationship(rs []Relation) error {
	entities := db.Entities
	for _, e := range entities {
		for _, r := range rs {
			r := r
			if e.AttrName == r.Dependent.AttrName {
				for _, er := range e.Relations {
					if er.Principal.AttrName == r.Principal.AttrName {
						return fmt.Errorf("relationship already exists principal entity %q,dependent entity %q %v", r.Principal.AttrName, r.Dependent.AttrName, er)
					}
				}
				e.Relations = append(e.Relations, &r)
			} else if e.AttrName == r.Principal.AttrName {
				for _, er := range e.Relations {
					if er.Dependent.AttrName == r.Dependent.AttrName {
						return fmt.Errorf("relationship already exists principal entity %q,dependent entity %q %v", r.Principal.AttrName, r.Dependent.AttrName, er)
					}
				}
				e.Relations = append(e.Relations, &r)
			}
		}
	}
	return nil
}

// extractRelField 从entity中提取关系字段。会判断字段是否为空，如果是空的会寻找实体的主键字段。
func (db *Database) extractRelField(r entity.FieldBuilder, e entity.EntityInterface) (*Field, error) {
	entities := db.Entities
	for _, _e := range entities {
		if _e.AttrName == e.Config().AttrName {
			for _, f := range _e.Fields {
				if r != nil && f.Name == r.Descriptor().Name {
					return f, nil
				} else if r == nil && f.Primary == 1 {
					return f, nil
				}
			}
		}

	}
	return nil, fmt.Errorf("not found field %s", r.Descriptor().Name)
}

func (db *Database) extractEntity(e entity.EntityInterface) (*Entity, error) {
	entities := db.Entities
	for _, _e := range entities {
		if _e.AttrName == e.Config().AttrName {
			return _e, nil
		}
	}
	return nil, fmt.Errorf("not found entity %s", e.Config().AttrName)
}

// loadEntity 从entity.EntityInterface中加载Schema定义的entity。
//
// Params:
//
//   - ei: 实现了entity.EntityInterface的entity。
func (e *Entity) loadEntity(ei entity.EntityInterface) error {
	entityInfos, err := e.initEntity(ei)
	if err != nil {
		return err
	}
	_, err = checkFields(ei)
	if err != nil {
		return err
	}
	for _, f := range entityInfos {
		if f.field != nil {
			sf, err := newField(f.field.Builder, f.field.Builder.Descriptor())
			if err != nil {
				return err
			}
			if sf.Sequence.Name != nil {
				e.Sequences = append(e.Sequences, sf.Sequence)
			}
			sf.StoragerPkg = f.field.Storager.Pkg
			sf.StoragerType = f.field.Storager.Type
			sf.StoragerOrigType = f.field.Storager.OrigType
			sf.Tag = f.field.Tag
			e.Fields = append(e.Fields, sf)
		}
	}

	// 对entity检查
	if err := checkEntity(e); err != nil {
		return err
	}
	// 对字段进行检查
	if err := checkEntityFields(e); err != nil {
		return err
	}
	return nil
}

// initEntity 初始化Shcema中Entity的成员，会生成一个初始的Descriptor，这个Descriptor会有一些默认的配置，并传给字段的Init方法。
// 这个方法保证了调用Fields()等方法时不会nil pointer dereference。
//
// Params:
//
//   - ei: 实现了entity.EntityInterface的entity。
//
// Returns:
//
//	0: entity中的字段信息。
//	1: 错误信息。
func (e *Entity) initEntity(ei entity.EntityInterface) ([]entityInfo, error) {
	infos := make([]entityInfo, 0)
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
		fieldVal := val.Field(i)
		fieldTags := t.Field(i).Tag
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
		if fieldVal.Kind() != reflect.Ptr {
			panic(fmt.Sprintf("field %q must be a non-nil pointer", fieldName))
		}
		fe := fieldVal.Elem()
		if !fe.IsValid() {
			// 处理 fe 是零值的情况
			continue
		}
		fy := fe.Type()
		for j := 0; j < fe.NumField(); j++ {
			storager := analyseField(fe.Field(j), fy.Field(j))
			if storager != nil {
				initDesc := &entity.Descriptor{
					Name:       fieldName,
					AttrName:   stringutil.ToSnakeCase(fieldName),
					Type:       t.Field(i).Type.String(),
					EntityName: e.Name,
				}

				if fieldVal.IsNil() {
					newInstance := reflect.New(fieldVal.Type().Elem()).Interface()
					if ef, ok := newInstance.(entity.FieldBuilder); ok {
						err := e.initField(ei, ef, initDesc)
						if err != nil {
							return infos, err
						}
						// 很重要，将新的实例赋值给���来的字段
						fieldVal.Set(reflect.ValueOf(ef))
						f := entityInfo{
							field: &fieldInfo{
								Tag:      string(fieldTags),
								Builder:  ef,
								Storager: *storager,
							},
						}
						infos = append(infos, f)
					}
				} else {
					if ef, ok := fieldVal.Interface().(entity.FieldBuilder); ok {
						err := e.initField(ei, ef, initDesc)
						if err != nil {
							return infos, err
						}
						f := entityInfo{
							field: &fieldInfo{
								Tag:      string(fieldTags),
								Builder:  ef,
								Storager: *storager,
							},
						}
						infos = append(infos, f)
					}
				}
				continue
			}
		}
	}
	return infos, nil
}

// initField 调用字段的Init方法，获得Descriptor。
//
// Params:
//
//   - ei: 实现了entity.EntityInterface的entity。
//   - f: 实现了entity.FieldBuilder的字段。
//   - initDesc: 初始的Descriptor。
func (e *Entity) initField(ei entity.EntityInterface, f entity.FieldBuilder, initDesc *entity.Descriptor) (err error) {
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
		panic(fmt.Sprintf("Unsupported attribute type for entity %q in database %s: attribute %q", ed.EntityName, db.Name, ed.AttrName))
	}
	if ed.AttrType == "" {
		ed.AttrType = attrType
	}
	valueType := f.ValueType()
	if valueType == "" {
		panic(fmt.Sprintf("Unsupported value type for entity %q in database %s: attribute %q", ed.EntityName, db.Name, ed.AttrName))
	}
	tmpls := f.ExtTemplate()

	ef := &Field{}
	// 只复制需要序列化的基本类型字段
	ef.EntityName = ed.EntityName
	ef.Name = ed.Name
	ef.AttrName = ed.AttrName
	ef.Type = ed.Type
	ef.AttrType = ed.AttrType
	ef.Size = int64(ed.Size)
	ef.Required = ed.Required
	ef.Primary = ed.Primary
	ef.Comment = ed.Comment
	ef.Default = ed.Default
	ef.DefaultValue = ed.DefaultValue
	ef.Locked = ed.Locked
	ef.Sequence = ed.Sequence
	ef.Depth = ed.Depth
	ef.BaseType = ed.BaseType
	ef.Uniques = ed.Uniques
	ef.CheckConstraint = ed.CheckConstraint
	ef.Indexes = ed.Indexes
	ef.IndexName = ed.IndexName
	ef.IndexMethod = ed.IndexMethod

	// 设置Field特有的字段
	ef.ValueType = valueType
	ef.Templates = tmpls
	ef.Validators = len(ed.Validators)

	err := checkSequence(ef.Sequence)
	if err != nil {
		return nil, err
	}
	return ef, nil
}

// analyseField 用于分析entity的字段，判断是不是Field类型，提取出里面的builder和storage来。
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
		OrigTypeName := extractOrigTypeName(typeName)
		return &fieldInfoStorager{
			Pkg:      v.Type().PkgPath(),
			Name:     s.Name,
			Type:     extractTypeName(s.Type),
			OrigType: OrigTypeName,
		}

	}
	return nil
}

// extractOrigTypeName 提取字段的类型的类型名称。含有泛型参数。
//
// 例如: geo.GeometryStorage[github.com/yohobala/taurus_go/encoding/geo.LineString,github.com/yohobala/taurus_go/encoding/geo.SDefault]
// -> geo.GeometryStorage[geo.LineString,geo.SDefault]
func extractTypeName(t reflect.Type) string {
	typeName := t.String()

	// 自定义函数来处理单个类型名称，保留最后的包名
	trimPackagePath := func(fullTypeName string) string {
		parts := strings.Split(fullTypeName, "/")
		simpleName := parts[len(parts)-1] // 获取最后一部分，可能包含包名和类型名
		return simpleName
	}

	// 检查是否有泛型参数
	if start := strings.Index(typeName, "["); start != -1 {
		// 基本类型部分
		base := trimPackagePath(typeName[:start])

		// 泛型参数部分
		end := strings.LastIndex(typeName, "]")
		params := typeName[start+1 : end]

		// 处理泛型参数，这可能是逗号分隔的列表
		paramTypes := strings.Split(params, ",")
		for i, paramType := range paramTypes {
			trimmedParam := strings.TrimSpace(paramType)
			newParam := trimPackagePath(trimmedParam)
			firstChar := trimmedParam[0]
			if firstChar == '*' {
				// 指针类型
				newParam = "*" + newParam
			}

			paramTypes[i] = trimPackagePath(newParam)
		}

		// 重组类型名称
		return base + "[" + strings.Join(paramTypes, ", ") + "]"
	} else {
		// 非泛型类型，直接处理类型名称
		return trimPackagePath(typeName)
	}
}

// extractOrigTypeName 提取字段的类型的原始类型名称。不含泛型参数和包名。
//
// 例如: field.IntStorage[int16] -> field.IntStorage
func extractOrigTypeName(typeName string) string {
	// 先移除泛型参数
	if genIndex := strings.Index(typeName, "["); genIndex != -1 {
		typeName = typeName[:genIndex]
	}

	// 提取最后一个点号之后的部分，以处理可能的包名
	lastDotIndex := strings.LastIndex(typeName, ".")
	if lastDotIndex != -1 {
		typeName = typeName[lastDotIndex+1:]
	}

	return typeName
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

// checkRelationDescriptor 检查关系描述符是否为空。
func checkRelationDescriptor(desc *entity.RelationshipDescriptor) error {
	if desc == nil {
		return fmt.Errorf("RelationshipDescriptor is nil")
	}
	if desc.ForeignKey == nil {
		return fmt.Errorf("ForeignKey is nil")
	}
	if desc.Has == nil {
		return fmt.Errorf("Has is nil")
	}
	if desc.With == nil {
		return fmt.Errorf("With is nil")
	}
	return nil
}

// checkRelationships 检查entity的Relationships()方法是否有panic，并得到返回值。
func checkRelationships(r interface {
	Relationships() []entity.RelationshipBuilder
}) (rels []entity.RelationshipBuilder, err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%T.Relationships panics: %v", r, v)
		}
	}()
	return r.Relationships(), nil
}

func checkEntity(e *Entity) error {
	if e.AttrName == "" {
		return entity.Err_0100020022
	}

	// 检查是否为PostgreSQL关键字
	lowerAttrName := strings.ToLower(e.AttrName)
	for _, keyword := range postgresKeywords {
		if lowerAttrName == keyword {
			return entity.Err_0100020023.Sprintf(e.AttrName)
		}
	}
	// 检查是否为Go语言关键字和保留字
	for _, keyword := range goKeywords {
		if lowerAttrName == keyword {
			return entity.Err_0100020023.Sprintf(e.AttrName)
		}
	}

	return nil
}

func checkEntityFields(e *Entity) error {
	var hasPrimary bool
	for _, f := range e.Fields {
		if f.Primary >= 1 {
			hasPrimary = true
		}
		lowerAttrName := strings.ToLower(f.AttrName)
		for _, keyword := range postgresKeywords {
			if lowerAttrName == keyword {
				return entity.Err_0100020023.Sprintf(f.AttrName)
			}
		}
		// 检查是否为Go语言关键字和保留字
		for _, keyword := range goKeywords {
			if lowerAttrName == keyword {
				return entity.Err_0100020023.Sprintf(f.AttrName)
			}
		}
	}
	if !hasPrimary {
		return entity.Err_0100020021.Sprintf(e.Name)
	}
	return nil
}

// indirect 穿透指针类型���获取不是指针类型的基础类型
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
		e := v.Elem()
		// 避免零值
		if !e.IsValid() {
			newInstance := reflect.New(v.Type().Elem()).Elem()
			v.Set(newInstance.Addr())
		}
		return v, true
	}
	return v, false
}
