package geo

import (
	"fmt"
	"strconv"
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
func (c CircularString) String() string {
	var s strings.Builder
	s.WriteString(circularStringPrefix + "(")
	for i, p := range c.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString(")")
	return s.String()
}

// param GeometryType接口的实现。
func (c CircularString) param() string {
	return c.String()
}

// attrType GeometryType接口的实现。
func (c CircularString) attrType() string {
	return "CircularString"
}

// valueType GeometryType接口的实现。
func (c CircularString) valueType() string {
	return "geo.CircularString"
}

// Decode 用于将字符串解析为CircularString。
func (c *CircularString) Decode(s string) error {
	trimmed := strings.TrimPrefix(s, circularStringPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	points := strings.Split(trimmed, ",")
	ps := make([]Point, len(points))
	for i, p := range points {
		coors := strings.Split(p, " ")
		if len(coors) != 2 {
			return Err_0200020201.Sprintf("CircularString", fmt.Sprintf("Decode() coordinate length is 2, but got %d", len(coors)))
		}
		lng, err := strconv.ParseFloat(coors[0], 64)
		if err != nil {
			return Err_0200020201.Sprintf("CircularString", fmt.Sprintf("Decode() lng is float, but got %s", coors[0]))
		}
		lat, err := strconv.ParseFloat(coors[1], 64)
		if err != nil {
			return Err_0200020201.Sprintf("CircularString", fmt.Sprintf("Decode() lat is float, but got %s", coors[1]))
		}
		np, err := NewPoint(lng, lat)
		if err != nil {
			return err
		}
		ps[i] = *np
	}
	c.points = ps
	return nil
}
