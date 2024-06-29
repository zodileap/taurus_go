package geo

import (
	"fmt"
	"strconv"
	"strings"
)

// 由多个点构成的闭合多边形。多边形的边界由线段组成，第一个点和最后一个点相同。
type Polygon struct {
	// points 多边形的点集合。
	points []Point
}

const polygonPrefix = "POLYGON"

// NewPolygon通过坐标数组创建一个多边形。
//
// Example:
//
//	p, err := NewPolygon([][]float64{{1, 2}, {3, 4}, {5, 6}})
func NewPolygon(polygon [][]float64) (*Polygon, error) {
	if len(polygon) < 4 {
		return nil, Err_0200020201.Sprintf("Polygon", "The number of points is less than 4")
	}
	ps := make([]Point, len(polygon))
	for i, p := range polygon {
		var err error
		np, err := NewPoint(p[0], p[1])
		if err != nil {
			return nil, err
		}
		ps[i] = *np
	}
	if ps[0].lat != ps[len(ps)-1].lat || ps[0].lng != ps[len(ps)-1].lng {
		return nil, Err_0200020201.Sprintf("Polygon", "The first point and the last point are not the same")
	}
	p := Polygon{points: ps}
	return &p, nil
}

// NewPolygonByPoint通过Point数组创建一个多边形。
//
// Example:
//
//	p1, err := NewPoint(1, 2)
//	p2, err := NewPoint(3, 4)
//	p3, err := NewPoint(5, 6)
//	p4, err := NewPoint(1, 2)
//	polygon, err := NewPolygonByPoint([]Point{*p1, *p2, *p3, *p4})
func NewPolygonByPoint(points []Point) (*Polygon, error) {
	if len(points) < 4 {
		return nil, Err_0200020201.Sprintf("Polygon", "The number of points is less than 4")
	}
	if points[0].lat != points[len(points)-1].lat || points[0].lng != points[len(points)-1].lng {
		return nil, Err_0200020201.Sprintf("Polygon", "The first point and the last point are not the same")
	}
	p := &Polygon{points: points}
	return p, nil
}

// String Polygon的字符串表示，它是一个WKT。
func (p Polygon) String() string {
	var s strings.Builder
	s.WriteString(polygonPrefix + "((")
	for i, p := range p.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString("))")
	return s.String()
}

// param GeometryType接口的实现。
func (p Polygon) param() string {
	return p.String()
}

// attrType GeometryType接口的实现。
func (p Polygon) attrType() string {
	return "Polygon"
}

// valueType GeometryType接口的实现。
func (p Polygon) valueType() string {
	return "geo.Polygon"
}

// Decode 用于将字符串解析为Polygon。
func (p *Polygon) Decode(s string) error {
	trimmed := strings.TrimPrefix(s, polygonPrefix)
	trimmed = strings.TrimPrefix(trimmed, "((")
	trimmed = strings.TrimSuffix(trimmed, "))")
	points := strings.Split(trimmed, ",")
	ps := make([]Point, len(points))
	for i, p := range points {
		coors := strings.Split(p, " ")
		if len(coors) != 2 {
			return Err_0200020201.Sprintf("Polygon", fmt.Sprintf("Decode() coordinate length is 2, but got %d", len(coors)))
		}
		lng, err := strconv.ParseFloat(coors[0], 64)
		if err != nil {
			return Err_0200020201.Sprintf("Polygon", fmt.Sprintf("Decode() lng is float, but got %s", coors[0]))
		}
		lat, err := strconv.ParseFloat(coors[1], 64)
		if err != nil {
			return Err_0200020201.Sprintf("Polygon", fmt.Sprintf("Decode() lat is float, but got %s", coors[1]))
		}
		np, err := NewPoint(lng, lat)
		if err != nil {
			return err
		}
		ps[i] = *np
	}
	p.points = ps
	return nil
}
