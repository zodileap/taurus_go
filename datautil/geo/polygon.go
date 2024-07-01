package geo

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Polygon 由多个点构成的闭合多边形。多边形的边界由线段组成，第一个点和最后一个点相同。
type Polygon struct {
	// points 多边形的点集合。
	points []Point
}

const polygonPrefix = "POLYGON"

// NewPolygon 通过坐标数组创建一个多边形。
//
// Example:
//
//	p, err := NewPolygon([][]float64{{1, 2}, {3, 4}, {5, 6}})
func NewPolygon(polygon [][]float64) (*Polygon, error) {
	ps := make([]Point, len(polygon))
	for i, p := range polygon {
		np, err := NewPoint(p[0], p[1])
		if err != nil {
			return nil, err
		}
		ps[i] = *np
	}
	p := &Polygon{points: ps}
	if err := p.valid(); err != nil {
		return nil, err
	}
	return p, nil
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
	p := &Polygon{points: points}
	if err := p.valid(); err != nil {
		return nil, err
	}
	return p, nil
}

// String Polygon的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (p *Polygon) String() string {
	s, _ := p.WKT()
	return s
}

// MarshalJSON Polygon的JSON序列化。
// 这个和GeoJSON相比相比，输出的JSON只是coordinates的值。
// 实现Coordinate接口和json.Marshaler接口。
func (p *Polygon) MarshalJSON() ([]byte, error) {
	var jsonString strings.Builder
	jsonString.WriteString("[[")
	for i, p := range p.points {
		if i > 0 {
			jsonString.WriteString(",")
		}
		j, err := json.Marshal(&p)
		if err != nil {
			return nil, err
		}
		jsonString.Write(j)
	}
	jsonString.WriteString("]]")
	return []byte(jsonString.String()), nil
}

// UnmarshalJSON Polygon的JSON反序列化。
// 实现json.Unmarshaler接口。
func (p *Polygon) UnmarshalJSON(data []byte) error {
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
	p.points = ps
	return p.valid()
}

// WKT Polygon的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (p *Polygon) WKT() (string, error) {
	if p == nil {
		return "", nil
	}
	var s strings.Builder
	s.WriteString(polygonPrefix + "((")
	for i, p := range p.points {
		if i > 0 {
			s.WriteString(",")
		}
		s.WriteString(fmt.Sprintf("%g %g", p.lng, p.lat))
	}
	s.WriteString("))")
	return s.String(), nil
}

// UnmarshalWKT 用于将字符串解析为Polygon。
func (p *Polygon) UnmarshalWKT(src string) error {
	trimmed := strings.TrimPrefix(src, polygonPrefix)
	trimmed = strings.TrimPrefix(trimmed, "((")
	trimmed = strings.TrimSuffix(trimmed, "))")
	pointStrings := strings.Split(trimmed, ",")
	ps := make([]Point, len(pointStrings))
	for i, pointString := range pointStrings {
		p := &Point{}
		if err := p.UnmarshalWKT(pointString); err != nil {
			return err
		}
		ps[i] = *p
	}
	p.points = ps
	return p.valid()
}

// GeoJSON 用于将LineString转换为GeoJSON。
// 这个和MarshalJSON相比，输出的JSON是一个完整的GeoJSON。
// 实现Coordinate接口。
func (p *Polygon) GeoJSON() (string, error) {
	g := NewGeometry(PolygonType, p)
	json, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

// UnmarshalGeoJSON 用于将字符串解析为Polygon。
func (p *Polygon) UnmarshalGeoJSON(src string) error {
	if p == nil {
		return Err_0200020202.Sprintf(PolygonType, "fromGeoJSON() Polygon is nil")
	}
	g := NewGeometry(PolygonType, p)
	if err := json.Unmarshal([]byte(src), &g); err != nil {
		return Err_0200020202.Sprintf(PolygonType, err)
	}
	return g.Coordinates.valid()
}

// valid 检查Polygon是否有效。
func (p *Polygon) valid() error {
	if len(p.points) < 4 {
		return Err_0200020201.Sprintf("Polygon", "The number of points is less than 4")
	}
	if p.points[0].lat != p.points[len(p.points)-1].lat || p.points[0].lng != p.points[len(p.points)-1].lng {
		return Err_0200020201.Sprintf("Polygon", "The first point and the last point are not the same")
	}
	return nil
}

// attrType 获取这个几何类型在数据库中字段的类型。
// PostGISGeometry接口的实现, 。
func (p *Polygon) attrType() string {
	return "Polygon"
}

// valueType 获取这个几何类型在go中的类型。
// PostGISGeometry接口的实现。
func (p *Polygon) valueType() string {
	return "*geo.Polygon"
}

// decode 用于解析Scan获得的数据，并存储到实体中。
// PostGISGeometry接口的实现。
func (p *Polygon) decode(src string, geomType string) error {
	switch geomType {
	case "ST_AsText":
		if err := p.UnmarshalWKT(src); err != nil {
			return err
		}
	case "ST_AsGeoJSON":
		if err := p.UnmarshalGeoJSON(src); err != nil {
			return err
		}
	default:
		return Err_0200020102.Sprintf(geomType, PolygonType)
	}
	return nil
}
