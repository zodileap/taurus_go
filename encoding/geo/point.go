package geo

import (
	"fmt"
	"strconv"
	"strings"
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
func NewPoint(lng, lat float64) (*Point, error) {
	p := Point{lng: lng, lat: lat}
	if err := p.valid(); err != nil {
		return nil, err
	}
	return &p, nil
}

func (p Point) valid() error {
	if p.lng < -180 || p.lng > 180 {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("lng range [-180, 180], but got %g", p.lng))
	}
	if p.lat < -90 || p.lat > 90 {
		return Err_0200020201.Sprintf("Point", fmt.Sprintf("lat range [-90, 90], but got %g", p.lat))
	}
	return nil
}

// String  Point的字符串表示，它是一个WKT。
func (p Point) String() string {
	return fmt.Sprintf(pointPrefix+"(%g %g)", p.lng, p.lat)
}

// param PostGISGeometry接口的实现。
func (p Point) param() string {
	return p.String()
}

// attrType PostGISGeometry接口的实现。
func (p Point) attrType() string {
	return "Point"
}

// valueType PostGISGeometry接口的实现。
func (p Point) valueType() string {
	return "geo.Point"
}

func (p *Point) Decode(s string) error {
	trimmed := strings.TrimPrefix(s, pointPrefix)
	trimmed = strings.TrimPrefix(trimmed, "(")
	trimmed = strings.TrimSuffix(trimmed, ")")
	coors := strings.Split(trimmed, " ")
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
	err = p.valid()
	return err
}
