package geo

import (
	"fmt"

	"github.com/zodileap/taurus_go/entity/entitysql"
)

// PostGIS 这个是'taurus_go/entity'字段类型的拓展，用于支持PostGIS的字段类型。
type PostGIS[G PostGISGeometry, S SRID, T GeomType] struct {
	// GeometryBuilder 字段的构造器，用于代码生成阶段。
	GeometryBuilder[G, S, T]
	// GeometryStorage 字段的存储器。
	GeometryStorage[G, S, T]
}

// PostGISGeometry PostGIS的几何类型。
type PostGISGeometry interface {
	Coordinate
	// attrType 获取这个几何类型在数据库中字段的类型。
	attrType() string
	// valueType 获取这个几何类型在go中的类型。
	valueType() string
	// decode 用于解析Scan获得的数据，并存储到实体中。
	decode(src string, geomType string) error
}

// SRID 空间参考标识符。
type SRID interface {
	String() string
}

// SDefault 默认的空间参考标识符。
type SDefault struct{}

// String 返回空间参考标识符的字符串表示。
func (s SDefault) String() string {
	return "0"
}

// S2163 空间参考标识符2163。
type S2163 struct{}

// String 返回空间参考标识符的字符串表示。
func (s S2163) String() string {
	return "2163"
}

// S3395 空间参考标识符3395。
type S3395 struct{}

func (s S3395) String() string {
	return "3395"
}

// S4269 空间参考标识符4269。
type S4269 struct{}

// String 返回空间参考标识符的字符串表示。
func (s S4269) String() string {
	return "4269"
}

// S4326 空间参考标识符4326。
type S4326 struct{}

// String 返回空间参考标识符的字符串表示。
func (s S4326) String() string {
	return "4326"
}

// GeomType 存储的Geometry的类型
type GeomType interface {
	// Param 返回在SQL中VALUE序列化参数的函数名
	Param() string
	// Column 返回在SQL中SELECT参数的函数名
	Column() string
}

// GeomFromText 用于PostGIS的ST_GeomFromText函数，
// 存储的是WKT格式的Geometry。
type GeomFromText struct{}

func (g GeomFromText) Param() string {
	return "ST_GeomFromText"
}

func (g GeomFromText) Column() string {
	return "ST_AsText"
}

// GeomFromGeoJSON 用于PostGIS的ST_GeomFromGeoJSON函数，
// 存储的是GeoJSON格式的Geometry。
type GeomFromGeoJSON struct{}

func (g GeomFromGeoJSON) Param() string {
	return "ST_GeomFromGeoJSON"
}

func (g GeomFromGeoJSON) Column() string {
	return "ST_AsGeoJSON"
}

// PostGISFunc PostGIS函数。函数的实例通过“entity_tmpl/where.tmpl"写入到生成的代码中。
type PostGISFunc struct {
	// Entity 实体名称，如果为空则使用默认实体名称。
	Entity string
	// Column 数据库列名称，如果为空则使用Name。
	Column string
	// Name 函数名称。当Name不为空时，会使用Value和Children。
	Name string
	// Value 函数的值类型参数。
	Value []any
	// Children 函数的子函数。
	Children []PostGISFunc
}

// Pred PostGIS的函数写入到Predicate中。
func (p PostGISFunc) Pred(pred *entitysql.Predicate, entityName string) *entitysql.Predicate {
	return pred.Append(func(b *entitysql.Builder) {
		p.write(b, pred, entityName)
	})
}

// write 写入到Builder中。
func (p PostGISFunc) write(b *entitysql.Builder, pred *entitysql.Predicate, entityName string) {
	var as string
	if p.Entity != "" {
		as = pred.Builder.FindAs(p.Entity)
	} else {
		as = pred.Builder.FindAs(entityName)
	}
	if p.Column != "" {
		if b.IsAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
		b.Ident(p.Column)
	} else if p.Name != "" {
		if b.IsAs && as != "" {
			b.WriteString(b.Quote(as))
			b.WriteByte('.')
		}
		b.WriteString(p.Name)
		b.WriteByte('(')
		for i, child := range p.Children {
			if i > 0 {
				b.Comma()
			}
			child.write(b, pred, entityName)
		}
		for i, v := range p.Value {
			s := fmt.Sprintf("%v", v)
			if s != "" {
				if i > 0 || len(p.Children) > 0 {
					b.Comma()
				}
				b.WriteString(s)
			}
		}
		b.WriteByte(')')
	}
}
