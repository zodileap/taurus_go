package geo

import (
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
func (mp MultiPoint) String() string {
	var s strings.Builder
	s.WriteString(multiPointPrefix + "(")
	for i, p := range mp.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString(")")
	return s.String()
}

// param GeometryType接口的实现。
func (mp MultiPoint) param() string {
	return mp.String()
}

// attrType GeometryType接口的实现。
func (mp MultiPoint) attrType() string {
	return "MultiPoint"
}

// valueType GeometryType接口的实现。
func (mp MultiPoint) valueType() string {
	return "geo.MultiPoint"
}

// Decode 用于将字符串解析为MultiPoint。
func (mp *MultiPoint) Decode(s string) error {
	trimmed := strings.TrimPrefix(s, multiPointPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	pointStrings := strings.Split(trimmed, ",")
	ps := make([]Point, len(pointStrings))
	for i, pointString := range pointStrings {
		p := &Point{}
		err := p.Decode(pointString)
		if err != nil {
			return err
		}
		ps[i] = *p
	}
	mp.points = ps
	return nil
}
