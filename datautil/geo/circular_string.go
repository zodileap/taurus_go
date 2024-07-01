package geo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// CircularString 一系列点连接成的圆弧线段。
type CircularString struct {
	points []Point
}

const circularStringPrefix = "CIRCULARSTRING"

// NewCircularString 通过坐标数组创建一个CircularString。
//
// Example:
//
//	cs, err := NewCircularString([][]float64{{1, 2}, {3, 4}, {5, 6}})
func NewCircularString(circularString [][]float64) (*CircularString, error) {
	ps := make([]Point, len(circularString))
	for i, p := range circularString {
		var err error
		np, err := NewPoint(p[0], p[1])
		if err != nil {
			return nil, err
		}
		ps[i] = *np
	}
	c := CircularString{points: ps}
	return &c, nil
}

// NewCircularStringByPoint 通过Point数组创建一个CircularString。
//
// Example:
//
//	p1, err := geo.NewPoint(1, 2)
//	p2, err := geo.NewPoint(3, 4)
//	p3, err := geo.NewPoint(5, 6)
//	cs, err := geo.NewCircularStringByPoint([]geo.Point{*p1, *p2, *p3})
func NewCircularStringByPoint(points []Point) (*CircularString, error) {
	c := &CircularString{points: points}
	return c, nil
}

// String  CircularString的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (c *CircularString) String() string {
	s, _ := c.WKT()
	return s
}

// MarshalJSON 用于将CircularString转换为JSON。
// 这个和GeoJSON相比相比，输出的JSON只是coordinates的值。
// 实现Coordinate接口和json.Marshaler接口。
func (c *CircularString) MarshalJSON() ([]byte, error) {
	return nil, Err_0200020203.Sprintf(CircularStringType)
}

// UnmarshalJSON 用于将JSON解析为CircularString。
// 实现json.Unmarshaler接口。
func (c *CircularString) UnmarshalJSON(data []byte) error {
	return Err_0200020203.Sprintf(CircularStringType)
}

// WKT CircularString的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (c *CircularString) WKT() (string, error) {
	if c == nil {
		return "", nil
	}
	var s strings.Builder
	s.WriteString(circularStringPrefix + "(")
	for i, p := range c.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString(")")
	return s.String(), nil
}

// UnmarshalWKT 用于将WKT解析为CircularString。
func (c *CircularString) UnmarshalWKT(src string) error {
	trimmed := strings.TrimPrefix(src, circularStringPrefix)
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
	c.points = ps
	return c.valid()
}

// GeoJSON 用于将CircularString转换为GeoJSON。
// 这个和MarshalJSON相比，输出的JSON是一个完整的GeoJSON。
// 实现Coordinate接口。
func (c *CircularString) GeoJSON() (string, error) {
	g := NewGeometry(CircularStringType, c)
	json, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

// UnmarshalGeoJSON 用于将GeoJSON解析为CircularString。
func (c *CircularString) UnmarshalGeoJSON(src string) error {
	if c == nil {
		return Err_0200020202.Sprintf(CircularStringType, "fromGeoJSON() CircularString is nil")
	}
	g := NewGeometry(CircularStringType, c)
	if err := json.Unmarshal([]byte(src), g); err != nil {
		return Err_0200020202.Sprintf(CircularStringType, err)
	}
	return g.Coordinates.valid()
}

// valid 用于验证CircularString。
func (c *CircularString) valid() error {
	return nil
}

// attrType GeometryType接口的实现。
func (c *CircularString) attrType() string {
	return "CircularString"
}

// valueType GeometryType接口的实现。
func (c *CircularString) valueType() string {
	return "*geo.CircularString"
}

// decode 用于将字符串解析为CircularString。
func (c *CircularString) decode(src string, geomType string) error {
	switch geomType {
	case "ST_AsText":
		if err := c.UnmarshalWKT(src); err != nil {
			return err
		}
	case "ST_AsGeoJSON":
		if err := c.UnmarshalGeoJSON(src); err != nil {
			return err
		}
	default:
		return Err_0200020102.Sprintf(geomType, CircularStringType)
	}
	return nil
}
