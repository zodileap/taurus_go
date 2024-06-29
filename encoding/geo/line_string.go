package geo

import (
	"fmt"
	"strconv"
	"strings"
)

// LineString 一系列点连接成的线段，可以是直线或曲线。
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
		var err error
		np, err := NewPoint(p[0], p[1])
		if err != nil {
			return nil, err
		}
		ps[i] = *np
	}
	l := LineString{points: ps}
	return &l, nil
}

// NewLineStringByPoint 通过Point数组创建一个LineString。
//
// Example:
//
//		p1, err := geo.NewPoint(1, 2)
//		p2, err := geo.NewPoint(3, 4)
//	 line2, err := geo.NewLineStringByPoint([]geo.Point{*p1, *p2})
func NewLineStringByPoint(points []Point) (*LineString, error) {
	l := &LineString{points: points}
	return l, nil
}

// String  LineString的字符串表示，它是一个WKT。
func (l LineString) String() string {
	var s strings.Builder
	s.WriteString(lineStringPrefix + "(")
	for i, p := range l.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString(")")
	return s.String()
}

// param GeometryType接口的实现。
func (l LineString) param() string {
	return l.String()
}

// attrType GeometryType接口的实现。
func (l LineString) attrType() string {
	return "LineString"
}

// valueType GeometryType接口的实现。
func (l LineString) valueType() string {
	return "geo.LineString"
}

// Decode 用于将字符串解析为LineString。
func (l *LineString) Decode(s string) error {
	trimmed := strings.TrimPrefix(s, lineStringPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	points := strings.Split(trimmed, ",")
	ps := make([]Point, len(points))
	for i, p := range points {
		coors := strings.Split(p, " ")
		if len(coors) != 2 {
			return Err_0200020201.Sprintf("LineString", fmt.Sprintf("Decode() coordinate length is 2, but got %d", len(coors)))
		}
		lng, err := strconv.ParseFloat(coors[0], 64)
		if err != nil {
			return Err_0200020201.Sprintf("LineString", fmt.Sprintf("Decode() lng is float, but got %s", coors[0]))
		}
		lat, err := strconv.ParseFloat(coors[1], 64)
		if err != nil {
			return Err_0200020201.Sprintf("LineString", fmt.Sprintf("Decode() lat is float, but got %s", coors[1]))
		}
		np, err := NewPoint(lng, lat)
		if err != nil {
			return err
		}
		ps[i] = *np
	}
	l.points = ps
	return nil
}
