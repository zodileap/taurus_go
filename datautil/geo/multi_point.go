package geo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// MultiPoint 一系列点的集合。
type MultiPoint struct {
	points []Point
}

const multiPointPrefix = "MULTIPOINT"

// NewMultiPoint 通过Point数组创建一个MultiPoint。
//
// Example:
//
//	p1, err := geo.NewPoint(1, 2)
//	p2, err := geo.NewPoint(3, 4)
//	mp, err := geo.NewMultiPoint([]geo.Point{*p1, *p2})
func NewMultiPoint(points []Point) (*MultiPoint, error) {
	mp := &MultiPoint{points: points}
	return mp, nil
}

// String  MultiPoint的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (mp *MultiPoint) String() string {
	s, _ := mp.WKT()
	return s
}

// MarshalJSON 用于将MultiPoint转换为JSON。
// 这个和GeoJSON相比相比，输出的JSON只是coordinates的值。
// 实现Coordinate接口和json.Marshaler接口。
func (mp *MultiPoint) MarshalJSON() ([]byte, error) {
	var jsonString strings.Builder
	jsonString.WriteString("[")
	for i, p := range mp.points {
		if i > 0 {
			jsonString.WriteString(",")
		}
		j, err := json.Marshal(&p)
		if err != nil {
			return nil, err
		}
		jsonString.Write(j)
	}
	jsonString.WriteString("]")
	return []byte(jsonString.String()), nil
}

// UnmarshalJSON 用于将JSON解析为MultiPoint。
// 实现json.Unmarshaler接口。
func (mp *MultiPoint) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "[]")
	pointstrings := strings.Split(s, "],[")
	ps := make([]Point, len(pointstrings))
	for i, pointString := range pointstrings {
		p := &Point{}
		if err := p.UnmarshalJSON([]byte(pointString)); err != nil {
			return err
		}
		ps[i] = *p
	}
	mp.points = ps
	return mp.valid()
}

// WKT MultiPoint的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (mp *MultiPoint) WKT() (string, error) {
	if mp == nil {
		return "", nil
	}
	var s strings.Builder
	s.WriteString(multiPointPrefix + "(")
	for i, p := range mp.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString(")")
	return s.String(), nil
}

// UnmarshalWKT 用于将字符串解析为MultiPoint。
func (mp *MultiPoint) UnmarshalWKT(s string) error {
	trimmed := strings.TrimPrefix(s, multiPointPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	pointStrings := strings.Split(trimmed, ",")
	ps := make([]Point, len(pointStrings))
	for i, pointString := range pointStrings {
		p := &Point{}
		if err := p.UnmarshalWKT(pointString); err != nil {
			return err
		}
		ps[i] = *p
	}
	mp.points = ps
	return mp.valid()
}

// GeoJSON 用于将MultiPoint转换为GeoJSON。
// 这个和MarshalJSON相比，输出的JSON是一个完整的GeoJSON。
// 实现Coordinate接口。
func (mp *MultiPoint) GeoJSON() (string, error) {
	g := NewGeometry(MultiPointType, mp)
	json, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

// UnmarshalGeoJSON 用于将字符串解析为MultiPoint。
func (mp *MultiPoint) UnmarshalGeoJSON(s string) error {
	if mp == nil {
		return Err_0200020202.Sprintf(MultiPointType, "fromGeoJSON() MultiPoint is nil")
	}
	g := NewGeometry(MultiPointType, mp)
	if err := json.Unmarshal([]byte(s), g); err != nil {
		return Err_0200020202.Sprintf(MultiPointType, err)
	}
	return g.Coordinates.valid()
}

// valid 用于验证MultiPoint是否有效。
func (mp *MultiPoint) valid() error {
	return nil
}

// attrType 获取这个几何类型在数据库中字段的类型。
// PostGISGeometry接口的实现。
func (mp *MultiPoint) attrType() string {
	return "MultiPoint"
}

// valueType 获取这个几何类型在go中的类型。
// PostGISGeometry接口的实现。
func (mp *MultiPoint) valueType() string {
	return "*geo.MultiPoint"
}

// decode 用于解析Scan获得的数据，并存储到实体中。
// PostGISGeometry接口的实现。
func (l *MultiPoint) decode(src string, geomType string) error {
	switch geomType {
	case "ST_AsText":
		if err := l.UnmarshalWKT(src); err != nil {
			return err
		}
	case "ST_AsGeoJSON":
		if err := l.UnmarshalGeoJSON(src); err != nil {
			return err
		}
	default:
		return Err_0200020102.Sprintf(geomType, MultiPointType)
	}
	return nil
}
