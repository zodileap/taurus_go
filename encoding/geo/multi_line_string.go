package geo

import "strings"

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
func (mls MultiLineString) String() string {
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
	return s.String()
}

// param GeometryType接口的实现。
func (mls MultiLineString) param() string {
	return mls.String()
}

// attrType GeometryType接口的实现。
func (mls MultiLineString) attrType() string {
	return "MultiLineString"
}

// valueType GeometryType接口的实现。
func (mls MultiLineString) valueType() string {
	return "geo.MultiLineString"
}

// Decode 用于将字符串解析为MultiLineString。
func (mls *MultiLineString) Decode(s string) error {
	trimmed := strings.TrimPrefix(s, multiLineStringPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	lineStrings := strings.Split(trimmed, "),(")
	ls := make([]LineString, len(lineStrings))
	for i, lineString := range lineStrings {
		ls[i] = LineString{}
		if err := ls[i].Decode(lineString); err != nil {
			return err
		}
	}
	mls.lineStrings = ls
	return nil
}
