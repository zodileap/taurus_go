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
	case *Point:
		if d == nil {
			d = &Point{}
		}
		switch geomType {
		case "ST_AsText":
			if err := d.Decode(src); err != nil {
				return Err_0200020103.Sprintf(src, PointType, err)
			}
		default:
			return Err_0200020102.Sprintf(geomType, PointType) // 默认情况下的错误处理
		}
		return nil
	case *LineString:
		if d == nil {
			d = &LineString{}
		}
		switch geomType {
		case "ST_AsText":
			if err := d.Decode(src); err != nil {
				return Err_0200020103.Sprintf(src, LineStringType, err)
			}
		default:
			return Err_0200020102.Sprintf(geomType, LineStringType) // 默认情况下的错误处理
		}
		return nil
	case *Polygon:
		if d == nil {
			d = &Polygon{}
		}
		switch geomType {
		case "ST_AsText":
			if err := d.Decode(src); err != nil {
				return Err_0200020103.Sprintf(src, PolygonType, err)
			}
		default:
			return Err_0200020102.Sprintf(geomType, PolygonType) // 默认情况下的错误处理
		}
		return nil
	case *MultiPoint:
		if d == nil {
			d = &MultiPoint{}
		}
		switch geomType {
		case "ST_AsText":
			if err := d.Decode(src); err != nil {
				return Err_0200020103.Sprintf(src, MultiPointType, err)
			}
		default:
			return Err_0200020102.Sprintf(geomType, MultiPointType) // 默认情况下的错误处理
		}
		return nil
	case *MultiLineString:
		if d == nil {
			d = &MultiLineString{}
		}
		switch geomType {
		case "ST_AsText":
			if err := d.Decode(src); err != nil {
				return Err_0200020103.Sprintf(src, MultiLineStringType, err)
			}
		default:
			return Err_0200020102.Sprintf(geomType, MultiLineStringType) // 默认情况下的错误处理
		}
		return nil
	case *MultiPolygon:
		if d == nil {
			d = &MultiPolygon{}
		}
		switch geomType {
		case "ST_AsText":
			if err := d.Decode(src); err != nil {
				return Err_0200020103.Sprintf(src, MultiPolygonType, err)
			}
		default:
			return Err_0200020102.Sprintf(geomType, MultiPolygonType) // 默认情况下的错误处理
		}
		return nil
	case *CircularString:
		if d == nil {
			d = &CircularString{}
		}
		switch geomType {
		case "ST_AsText":
			if err := d.Decode(src); err != nil {
				return Err_0200020103.Sprintf(src, CircularStringType, err)
			}
		default:
			return Err_0200020102.Sprintf(geomType, CircularStringType) // 默认情况下的错误处理
		}
		return nil
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
