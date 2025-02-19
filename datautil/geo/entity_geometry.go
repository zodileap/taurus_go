package geo

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/yohobala/taurus_go/entity"
	"github.com/yohobala/taurus_go/entity/dialect"
)

type GeometryBuilder[G PostGISGeometry, S SRID, T GeomType] struct {
	desc *entity.Descriptor
}

// Init 初始化字段的描述信息，在代码生成阶段初始化时调用。
//
// Params:
//
//   - desc: 字段的描述信息。
func (g *GeometryBuilder[G, S, T]) Init(desc *entity.Descriptor) error {
	if g == nil {
		panic("taurus_go/entity field init: nil pointer dereference.")
	}
	g.desc = desc
	return nil
}

// Descriptor 获取字段的描述信息。
func (g *GeometryBuilder[G, S, T]) Descriptor() *entity.Descriptor {
	return g.desc
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
func (g *GeometryBuilder[G, S, T]) AttrType(dbType dialect.DbDriver) string {
	switch dbType {
	case dialect.PostgreSQL:
		var gg G
		var gs S
		vgStr := gg.attrType()
		vsStr := gs.String()
		if vsStr == "" {
			return fmt.Sprintf("geometry(%s)", vgStr)
		} else {
			return fmt.Sprintf("geometry(%s, %s)", vgStr, vsStr)
		}
	default:
		return ""
	}
}

// ValueType 用于设置字段的值在go中类型名称。例如entity.Int64的ValueType为"int64"。
//
// Returns:
//
//   - 字段的值在go中类型名称。
func (g *GeometryBuilder[G, S, T]) ValueType() string {
	var gt G
	return fmt.Sprintf("%s", gt.valueType())
}

// ExtTemplate 用于在使用字段时，调用外部模版生成代码，
// 这个相比在 go run github.com/yohobala/taurus_go/entity/cmd generate -t <template>，
// `ExtTemplate`是和字段相关联，只要调用字段就会生成代码，避免了每次都要手动调用模版。
//
// Returns:
//
//	0: 模版的路径。
func (g *GeometryBuilder[G, S, T]) ExtTemplate() []string {
	return []string{
		loadTemplate("entity_tmpl/where.tmpl"),
	}
}

// loadTemplate 加载指定模板并返回 *template.Template
//
// Params:
//
//   - fileName: 模板文件名。
func loadTemplate(fileName string) string {
	// 获取当前文件所在的目录
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	dir := filepath.Dir(filename)
	// 拼接模板文件路径
	tmplPath := filepath.Join(dir, fileName)
	return tmplPath
}

// Name 用于设置字段在数据库中的名称。
//
// 如果不设置，会默认采用`snake_case`的方式将字段名转换为数据库字段名，比如示例中的ID字段会被转换为`i_d`。
//
// Params:
//
//   - name: 字段在数据库中的名称。
func (g *GeometryBuilder[G, S, T]) Name(name string) *GeometryBuilder[G, S, T] {
	g.desc.AttrName = name
	return g
}

// Required 是否非空,默认可以为null,如果调用[Required],则字段为非空字段。
func (i *GeometryBuilder[G, S, T]) Required() *GeometryBuilder[G, S, T] {
	i.desc.Required = true
	return i
}

// Primary设置字段为主键。
//
// Params:
//
//   - index: 主键的索引，从1开始，对于多个主键，需要设置不同大小的索引。
func (i *GeometryBuilder[G, S, T]) Primary(index int) *GeometryBuilder[G, S, T] {
	i.desc.Required = true
	i.desc.Primary = index
	return i
}

// Comment 设置字段的注释。
//
// Params:
//
//   - comment: 字段的注释。
func (i *GeometryBuilder[G, S, T]) Comment(comment string) *GeometryBuilder[G, S, T] {
	i.desc.Comment = comment
	return i
}

// Locked 设置字段为只读字段。
func (i *GeometryBuilder[G, S, T]) Locked() *GeometryBuilder[G, S, T] {
	i.desc.Locked = true
	return i
}

// GeometryStorage 字段的存储器。
type GeometryStorage[G PostGISGeometry, S SRID, T GeomType] struct {
	value G
}

// Set 设置字段的值。
func (g *GeometryStorage[G, S, T]) Set(value G) error {
	g.value = value
	return nil
}

// Get 获取字段的值。
func (g *GeometryStorage[G, S, T]) Get() G {
	return g.value
}

// Scan 从数据库中读取字段的值。
func (g *GeometryStorage[G, S, T]) Scan(src interface{}) error {
	if src == nil {
		var v G
		g.value = v
		return nil
	}
	// 将值转换为字符串
	s, ok := src.(string)
	if !ok {
		return Err_0200020101
	}
	return g.convert(&g.value, s)
}

// String 返回字段的字符串表示。
func (g *GeometryStorage[G, S, T]) String() string {
	if reflect.ValueOf(g.value).IsNil() {
		return "nil"
	}
	return fmt.Sprintf("%v", g.value)
}

// SqlParam 用于sql中获取字段参数并赋值。如 INSERT INTO "blog" ( "desc") VALUES ($1)，给$1传递具体的值。
func (g *GeometryStorage[G, S, T]) SqlParam(dpType dialect.DbDriver) (entity.FieldValue, error) {
	switch dpType {
	case dialect.PostgreSQL:
		if reflect.ValueOf(g.value).IsNil() {
			return nil, nil
		}
		var gt T
		var v any
		var err error
		geomType := gt.Column()
		switch geomType {
		case "ST_AsText":
			v, err = g.value.WKT()
		case "ST_AsGeoJSON":
			v, err = g.value.GeoJSON()
		default:
			return nil, Err_0200020104.Sprintf(geomType)
		}
		if err != nil {
			return nil, err
		}
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported database type: %v", dpType)
	}
}

// SqlFormatParam 用于sql中获取字段的值的格式化字符串。如 INSERT INTO "blog" ( "desc" ) VALUES ( ST_GeomFromGeoJSON($1) ) 中添加的ST_GeomFromGeoJSON()。
func (g *GeometryStorage[G, S, T]) SqlFormatParam() func(dbType dialect.DbDriver, param string) string {
	return func(dbType dialect.DbDriver, param string) string {
		var gt T
		var gs S
		if gs.String() == "" {
			return fmt.Sprintf("%s(%s)", gt.Param(), param)
		} else {
			return fmt.Sprintf("ST_SetSRID(%s(%s),%s)", gt.Param(), param, gs.String())
		}
	}
}

// SqlSelectClause 用于sql语句中获取字段的select子句部分，通过这个能够扩展SELECT部分实现复杂的查询，比如 SELECT id, ST_AsText(point)。
func (g *GeometryStorage[G, S, T]) SqlSelectFormat() func(dbType dialect.DbDriver, name string) string {
	var gt T
	return func(dbType dialect.DbDriver, name string) string {
		return fmt.Sprintf("%s(%s)", gt.Column(), name)
	}
}
