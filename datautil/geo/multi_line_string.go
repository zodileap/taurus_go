package geo

import (
	"encoding/json"
	"strings"
)

// MultiLineString 多个LineString组成的多段线。
type MultiLineString struct {
	lineStrings []LineString
}

const multiLineStringPrefix = "MULTILINESTRING"

// NewMultiLineString 通过LineString数组创建一个MultiLineString。
//
// Example:
//
//	ls1, _ := NewLineString([][]float64{{1, 2}, {3, 4}})
//	ls2, _ := NewLineString([][]float64{{5, 6}, {7, 8}})
//	mls, err := NewMultiLineString([]*LineString{ls1, ls2})
func NewMultiLineString(lineStrings []LineString) (*MultiLineString, error) {
	mls := &MultiLineString{lineStrings: lineStrings}
	return mls, nil
}

// String MultiLineString的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (mls *MultiLineString) String() string {
	s, _ := mls.WKT()
	return s
}

// MarshalJSON 用于将MultiLineString转换为JSON。
// 这个和GeoJSON相比相比，输出的JSON只是coordinates的值。
// 实现Coordinate接口和json.Marshaler接口。
func (mls *MultiLineString) MarshalJSON() ([]byte, error) {
	var jsonString strings.Builder
	jsonString.WriteString("[")
	for i, ls := range mls.lineStrings {
		if i > 0 {
			jsonString.WriteString(",")
		}
		j, err := json.Marshal(&ls)
		if err != nil {
			return nil, err
		}
		jsonString.Write(j)
	}
	jsonString.WriteString("]")
	return []byte(jsonString.String()), nil
}

// UnmarshalJSON 用于将JSON解析为MultiLineString。
// 实现json.Unmarshaler接口。
func (mls *MultiLineString) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "[]")
	lineStringStrings := strings.Split(s, "]],[[")
	ls := make([]LineString, len(lineStringStrings))
	for i, lineStringString := range lineStringStrings {
		l := &LineString{}
		if err := l.UnmarshalJSON([]byte(lineStringString)); err != nil {
			return err
		}
		ls[i] = *l
	}
	mls.lineStrings = ls
	return mls.valid()
}

// WKT Polygon的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (mls *MultiLineString) WKT() (string, error) {
	if mls == nil {
		return "", nil
	}
	var s strings.Builder
	s.WriteString(multiLineStringPrefix + "(")
	for i, ls := range mls.lineStrings {
		if i > 0 {
			s.WriteString(",")
		}
		l := ls.String()
		l = strings.TrimPrefix(l, lineStringPrefix)
		s.WriteString(l)
	}
	s.WriteString(")")
	return s.String(), nil
}

// UnmarshalWKT 用于将WKT解析为MultiLineString。
func (mls *MultiLineString) UnmarshalWKT(src string) error {
	trimmed := strings.TrimPrefix(src, multiLineStringPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	lineStrings := strings.Split(trimmed, "),(")
	ls := make([]LineString, len(lineStrings))
	for i, lineString := range lineStrings {
		l := &LineString{}
		if err := l.UnmarshalWKT(lineString); err != nil {
			return err
		}
		ls[i] = *l
	}
	mls.lineStrings = ls
	return mls.valid()
}

// GeoJSON 用于将MultiLineString转换为GeoJSON。
// 这个和MarshalJSON相比，输出的JSON是一个完整的GeoJSON。
// 实现Coordinate接口。
func (mls *MultiLineString) GeoJSON() (string, error) {
	g := NewGeometry(MultiLineStringType, mls)
	json, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

// UnmarshalGeoJSON 用于将GeoJSON解析为MultiLineString。
func (mls *MultiLineString) UnmarshalGeoJSON(src string) error {
	if mls == nil {
		return Err_0200020202.Sprintf(MultiLineStringType, "fromGeoJSON() MultiLineString is nil")
	}
	g := NewGeometry(MultiLineStringType, mls)
	if err := json.Unmarshal([]byte(src), g); err != nil {
		return Err_0200020202.Sprintf(MultiLineStringType, err)
	}
	return g.Coordinates.valid()
}

// valid 用于验证MultiLineString是否有效。
func (mls *MultiLineString) valid() error {
	return nil
}

// attrType 获取这个几何类型在数据库中字段的类型。
// PostGISGeometry接口的实现。
func (mls *MultiLineString) attrType() string {
	return "MultiLineString"
}

// valueType 获取这个几何类型在go中的类型。
// PostGISGeometry接口的实现。
func (mls *MultiLineString) valueType() string {
	return "*geo.MultiLineString"
}

// decode 用于解析Scan获得的数据，并存储到实体中。
// PostGISGeometry接口的实现。
func (mls *MultiLineString) decode(src string, geomType string) error {
	switch geomType {
	case "ST_AsText":
		if err := mls.UnmarshalWKT(src); err != nil {
			return err
		}
	case "ST_AsGeoJSON":
		if err := mls.UnmarshalGeoJSON(src); err != nil {
			return err
		}
	default:
		return Err_0200020102.Sprintf(geomType, MultiLineStringType)
	}
	return nil
}
