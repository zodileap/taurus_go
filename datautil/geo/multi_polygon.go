package geo

import (
	"strings"

	"encoding/json"
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
// 实现Coordinate接口。
func (mp *MultiPolygon) String() string {
	s, _ := mp.WKT()
	return s
}

// MarshalJSON 实现json.Marshaler接口。
// 这个和GeoJSON相比相比，输出的JSON只是coordinates的值。
// 实现Coordinate接口和json.Marshaler接口。
func (mp *MultiPolygon) MarshalJSON() ([]byte, error) {
	var jsonString strings.Builder
	jsonString.WriteString("[")
	for i, polygon := range mp.polygons {
		if i > 0 {
			jsonString.WriteString(",")
		}
		j, err := json.Marshal(&polygon)
		if err != nil {
			return nil, err
		}
		jsonString.Write(j)
	}
	jsonString.WriteString("]")
	return []byte(jsonString.String()), nil
}

// UnmarshalJSON MultiPolygon的JSON反序列化。
// 实现json.Unmarshaler接口。
func (mp *MultiPolygon) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "[]")
	polygonStrings := strings.Split(s, "]],[[")
	polygons := make([]Polygon, len(polygonStrings))
	for i, polygonString := range polygonStrings {
		polygon := &Polygon{}
		if err := polygon.UnmarshalJSON([]byte(polygonString)); err != nil {
			return err
		}
		polygons[i] = *polygon
	}
	mp.polygons = polygons
	return mp.valid()
}

// WKT MultiPolygon的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (mp *MultiPolygon) WKT() (string, error) {
	if mp == nil {
		return "", nil
	}
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
	return s.String(), nil
}

// UnmarshalWKT 用于将字符串解析为MultiPolygon。
func (mp *MultiPolygon) UnmarshalWKT(src string) error {
	trimmed := strings.TrimPrefix(src, multiPolygonPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	polygonStrings := strings.Split(trimmed, ")),((")
	polygons := make([]Polygon, len(polygonStrings))
	for i, polygonString := range polygonStrings {
		polygon := &Polygon{}
		if err := polygon.UnmarshalWKT(polygonString); err != nil {
			return err
		}
		polygons[i] = *polygon
	}
	mp.polygons = polygons
	return mp.valid()
}

// GeoJSON 用于将MultiPolygon转换为GeoJSON。
// 这个和MarshalJSON相比，输出的JSON是一个完整的GeoJSON。
// 实现Coordinate接口。
func (mp *MultiPolygon) GeoJSON() (string, error) {
	g := NewGeometry(MultiPolygonType, mp)
	json, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

// UnmarshalGeoJSON 用于将字符串解析为MultiPolygon。
func (mp *MultiPolygon) UnmarshalGeoJSON(s string) error {
	if mp == nil {
		return Err_0200020202.Sprintf(MultiPolygonType, "fromGeoJSON() MultiPolygon is nil")
	}
	g := NewGeometry(MultiPolygonType, mp)
	if err := json.Unmarshal([]byte(s), g); err != nil {
		return Err_0200020202.Sprintf(MultiPolygonType, err)
	}
	return g.Coordinates.valid()
}

// valid 用于验证MultiPolygon是否有效。
func (mp *MultiPolygon) valid() error {
	return nil
}

// attrType 获取这个几何类型在数据库中字段的类型。
// PostGISGeometry接口的实现。
func (mp *MultiPolygon) attrType() string {
	return "MultiPolygon"
}

// valueType 获取这个几何类型在go中的类型。
// PostGISGeometry接口的实现。
func (mp *MultiPolygon) valueType() string {
	return "*geo.MultiPolygon"
}

// decode 用于解析Scan获得的数据，并存储到实体中。
// PostGISGeometry接口的实现。
func (mp *MultiPolygon) decode(src string, geomType string) error {
	switch geomType {
	case "ST_AsText":
		if err := mp.UnmarshalWKT(src); err != nil {
			return err
		}
	case "ST_AsGeoJSON":
		if err := mp.UnmarshalGeoJSON(src); err != nil {
			return err
		}
	default:
		return Err_0200020102.Sprintf(geomType, MultiPolygonType)
	}
	return nil
}
