package geo

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/yohobala/taurus_go/tlog"
)

// Point 一个点，由经度和纬度组成。
type Point struct {
	// lng 经度(longitude)
	lng float64
	// lat 纬度(latitude)
	lat float64
}

const pointPrefix = "POINT"

// NewPoint 通过经度和纬度创建一个点。
//
// Example:
//
//	p, err := NewPoint(1, 2)
func NewPoint(lng, lat float64) (*Point, error) {
	p := Point{lng: lng, lat: lat}
	if err := p.valid(); err != nil {
		return nil, err
	}
	return &p, nil
}

// String  Point的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (p *Point) String() string {
	s, _ := p.WKT()
	return s
}

// MarshalJSON  用于将Point转换为JSON。
// 这个和GeoJSON相比相比，输出的JSON只是coordinates的值。
// 实现Coordinate接口和json.Marshaler接口。
func (p *Point) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf("[%g,%g]", p.lng, p.lat)
	return []byte(json), nil
}

// UnmarshalJSON 用于将JSON解析为Point。
// 实现json.Unmarshaler接口。
func (p *Point) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), "[]")
	coors := strings.Split(s, ",")
	if len(coors) != 2 {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("UnmarshalJSON() coordinate length is 2, but got %d", len(coors)))
	}
	var err error
	p.lng, err = strconv.ParseFloat(coors[0], 64)
	if err != nil {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("UnmarshalJSON() lng is float, but got %s", coors[0]))
	}
	p.lat, err = strconv.ParseFloat(coors[1], 64)
	if err != nil {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("UnmarshalJSON() lat is float, but got %s", coors[1]))
	}
	return p.valid()
}

// WKT  Point的字符串表示，它是一个WKT。
// 实现Coordinate接口。
func (p *Point) WKT() (string, error) {
	if p == nil {
		return "", nil
	}
	return fmt.Sprintf(pointPrefix+"(%g %g)", p.lng, p.lat), nil
}

// UnmarshalWKT 用于将字符串解析为Point。
func (p *Point) UnmarshalWKT(src string) error {
	tlog.Print(src)
	trimmed := strings.TrimPrefix(src, pointPrefix)
	trimmed = strings.Trim(trimmed, " ")
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	coors := strings.Split(trimmed, " ")
	tlog.Print(coors)
	if len(coors) != 2 {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("Decode() coordinate length is 2, but got %d", len(coors)))
	}
	var err error
	p.lng, err = strconv.ParseFloat(coors[0], 64)
	if err != nil {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("Decode() lng is float, but got %s", coors[0]))
	}
	p.lat, err = strconv.ParseFloat(coors[1], 64)
	if err != nil {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("Decode() lat is float, but got %s", coors[1]))
	}
	return p.valid()
}

// GeoJSON 用于将Point转换为GeoJSON。
// 这个和MarshalJSON相比，输出的JSON是一个完整的GeoJSON。
// 实现Coordinate接口。
func (p *Point) GeoJSON() (string, error) {
	g := NewGeometry(PointType, p)
	json, err := json.Marshal(g)
	if err != nil {
		return "", err
	}
	return string(json), nil
}

// UnmarshalGeoJSON 用于将字符串解析为Point。
func (p *Point) UnmarshalGeoJSON(src string) error {
	if p == nil {
		return Err_0200020202.Sprintf(PointType, "fromGeoJSON() Point is nil")
	}
	g := NewGeometry(PointType, p)
	if err := json.Unmarshal([]byte(src), g); err != nil {
		return Err_0200020202.Sprintf(PointType, err)
	}
	return g.Coordinates.valid()
}

// valid 用于验证Point是否合法。
func (p *Point) valid() error {
	if p.lng < -180 || p.lng > 180 {
		return Err_0200020201.Sprintf(PointType, fmt.Sprintf("lng range [-180, 180], but got %g", p.lng))
	}
	if p.lat < -90 || p.lat > 90 {
		return Err_0200020201.Sprintf(PointType, fmt.Sprintf("lat range [-90, 90], but got %g", p.lat))
	}
	return nil
}

// attrType 获取这个几何类型在数据库中字段的类型。
// PostGISGeometry接口的实现, 。
func (p *Point) attrType() string {
	return "Point"
}

// valueType 获取这个几何类型在go中的类型。
// PostGISGeometry接口的实现。
func (p *Point) valueType() string {
	return "*geo.Point"
}

// decode 用于解析Scan获得的数据，并存储到实体中。
// PostGISGeometry接口的实现。
func (p *Point) decode(src string, geomType string) error {
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
		return Err_0200020102.Sprintf(geomType, PointType)
	}
	if err := p.valid(); err != nil {
		return Err_0200020103.Sprintf(src, PointType, err)
	}
	return nil
}
