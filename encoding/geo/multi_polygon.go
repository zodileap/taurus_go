package geo

import (
	"strings"
)

// 由多个多边形构成的集合。每个多边形由多个点构成，边界由线段组成，第一个点和最后一个点相同。
type MultiPolygon struct {
	// polygons 多边形的集合。
	polygons []Polygon
}

const multiPolygonPrefix = "MULTIPOLYGON"

// NewMultiPolygon 通过多个多边形数组创建一个多边形集合。
//
// Example:
//
//	p1, err := NewPolygon([][]float64{{1, 2}, {3, 4}, {5, 6}})
//	p2, err := NewPolygon([][]float64{{7, 8}, {9, 10}, {11, 12}})
//	multiPolygon, err := NewMultiPolygon([]Polygon{*p1, *p2})
func NewMultiPolygon(polygons []Polygon) (*MultiPolygon, error) {
	if len(polygons) == 0 {
		return nil, Err_0200020201.Sprintf("MultiPolygon", "The number of polygons is 0")
	}
	mp := MultiPolygon{polygons: polygons}
	return &mp, nil
}

// String MultiPolygon的字符串表示，它是一个WKT。
func (mp MultiPolygon) String() string {
	var s strings.Builder
	s.WriteString(multiPolygonPrefix + "(")
	for i, polygon := range mp.polygons {
		if i > 0 {
			s.WriteString(",")
		}
		p := polygon.String()
		p = strings.TrimPrefix(p, polygonPrefix)
		s.WriteString(p)
	}
	s.WriteString(")")
	return s.String()
}

// param GeometryType接口的实现。
func (mp MultiPolygon) param() string {
	return mp.String()
}

// attrType GeometryType接口的实现。
func (mp MultiPolygon) attrType() string {
	return "MultiPolygon"
}

// valueType GeometryType接口的实现。
func (mp MultiPolygon) valueType() string {
	return "geo.MultiPolygon"
}

// Decode 用于将字符串解析为MultiPolygon。
func (mp *MultiPolygon) Decode(s string) error {
	trimmed := strings.TrimPrefix(s, multiPolygonPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	polygonStrings := strings.Split(trimmed, ",")
	polygons := make([]Polygon, len(polygonStrings))
	for i, polygonString := range polygonStrings {
		polygon := &Polygon{}
		err := polygon.Decode(polygonString)
		if err != nil {
			return err
		}
		polygons[i] = *polygon
	}
	mp.polygons = polygons
	return nil
}
