package geo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// LineString 一系列点连接成的线段一般是直线，曲线用CircularString。
// 但是如果输出成GeoJSON，无论是直线还是曲线都是LineString。
type LineString struct {
	points []Point
}

const lineStringPrefix = "LINESTRING"

// NewLineString 通过坐标数组创建一个LineString。
//
// Example:
//
//	ls, err := NewLineString([][]float64{{1, 2}, {3, 4}})
func NewLineString(lineString [][]float64) (*LineString, error) {
	ps := make([]Point, len(lineString))
	for i, p := range lineString {
		np, err := NewPoint(p[0], p[1])
		if err != nil {
			return nil, err
		}
		ps[i] = *np
	}
	l := &LineString{points: ps}
	if err := l.valid(); err != nil {
		return nil, err
	}
	return l, nil
}

// NewLineStringByPoint 通过Point数组创建一个LineString。
//
// Example:
//
//	p1, err := geo.NewPoint(1, 2)
//	p2, err := geo.NewPoint(3, 4)
//	line2, err := geo.NewLineStringByPoint([]geo.Point{*p1, *p2})
func NewLineStringByPoint(points []Point) (*LineString, error) {
	l := &LineString{points: points}
	if err := l.valid(); err != nil {
		return nil, err
	}
	return l, nil
}

// String  LineString的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (l *LineString) String() string {
	s, _ := l.WKT()
	return s
}

// MarshalJSON 用于将LineString转换为JSON。
// 这个和GeoJSON相比相比，输出的JSON只是coordinates的值。
// 实现Coordinate接口和json.Marshaler接口。
func (l *LineString) MarshalJSON() ([]byte, error) {
	var jsonString strings.Builder
	jsonString.WriteString("[")
	for i, p := range l.points {
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

// UnmarshalJSON 用于将JSON转换为LineString。
// 实现json.Unmarshaler接口。
func (l *LineString) UnmarshalJSON(data []byte) error {
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
	l.points = ps
	return l.valid()
}

// WKT 用于将LineString转换为WKT。
// 实现Coordinate接口。
func (l *LineString) WKT() (string, error) {
	if l == nil {
		return "", nil
	}
	var s strings.Builder
	s.WriteString(lineStringPrefix + "(")
	for i, p := range l.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString(")")
	return s.String(), nil
}

// UnmarshalWKT 用于将字符串解析为LineString。
// 实现Coordinate接口。
func (l *LineString) UnmarshalWKT(src string) error {
	trimmed := strings.TrimPrefix(src, lineStringPrefix)
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
	l.points = ps
	return l.valid()
}

// GeoJSON 用于将LineString转换为GeoJSON。
// 这个和MarshalJSON相比，输出的JSON是一个完整的GeoJSON。
// 实现Coordinate接口。
func (l *LineString) GeoJSON() (string, error) {
	g := NewGeometry(LineStringType, l)
	json, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

// UnmarshalGeoJSON 用于将字符串解析为LineString。
func (l *LineString) UnmarshalGeoJSON(src string) error {
	if l == nil {
		return Err_0200020202.Sprintf(LineStringType, "fromGeoJSON() LineString is nil")
	}
	g := NewGeometry(LineStringType, l)
	if err := json.Unmarshal([]byte(src), g); err != nil {
		return Err_0200020202.Sprintf(LineStringType, err)
	}
	return g.Coordinates.valid()
}

// valid 用于验证LineString是否有效。
func (l *LineString) valid() error {
	if len(l.points) < 2 {
		return Err_0200020201.Sprintf(LineStringType, "LineString must have at least two points")
	}
	return nil
}

// attrType 获取这个几何类型在数据库中字段的类型。
// PostGISGeometry接口的实现。
func (l *LineString) attrType() string {
	return "LineString"
}

// valueType 获取这个几何类型在go中的类型。
// PostGISGeometry接口的实现。
func (l *LineString) valueType() string {
	return "*geo.LineString"
}

// decode 用于解析Scan获得的数据，并存储到实体中。
// PostGISGeometry接口的实现。
func (l *LineString) decode(src string, geomType string) error {
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
		return Err_0200020102.Sprintf(geomType, LineStringType)
	}
	return nil
}
