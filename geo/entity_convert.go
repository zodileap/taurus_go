package geo

import (
	"reflect"
)

// convert 将数据库中的数据转换为实体中的数据。
//
// Params:
//
//   - dest: 目标数据。
//   - src: 源数据。
//
// Returns:
//
//	0: 错误信息。
func (g *GeometryStorage[G, S, T]) convert(dest any, src string) error {
	var gt T
	geomType := gt.Column()
	switch d := dest.(type) {
	case **Point:
		if *d == nil {
			*d = &Point{}
		}
		err := (*d).decode(src, geomType)
		return err
	case **LineString:
		if *d == nil {
			*d = &LineString{}
		}
		return (*d).decode(src, geomType)
	case **Polygon:
		if *d == nil {
			*d = &Polygon{}
		}
		return (*d).decode(src, geomType)
	case **MultiPoint:
		if *d == nil {
			*d = &MultiPoint{}
		}
		return (*d).decode(src, geomType)
	case **MultiLineString:
		if *d == nil {
			*d = &MultiLineString{}
		}
		return (*d).decode(src, geomType)
	case **MultiPolygon:
		if *d == nil {
			*d = &MultiPolygon{}
		}
		return (*d).decode(src, geomType)
	case **CircularString:
		if *d == nil {
			*d = &CircularString{}
		}
		return (*d).decode(src, geomType)
	default:
		dpv := reflect.ValueOf(dest)
		if dpv.Kind() == reflect.Ptr {
			return Err_0200020102.Sprintf(geomType, reflect.TypeOf(dest).Elem().Name())
		} else {
			return Err_0200020102.Sprintf(geomType, reflect.TypeOf(dest).Name())
		}
		// 默认情况下的错误处理
	}
}
