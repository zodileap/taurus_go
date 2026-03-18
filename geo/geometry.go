package geo

type Coordinate interface {
	// String 用于将几何对象转换为字符串。在fmt.Print等函数中会调用此方法。
	String() string
	// MarshalJSON 用于将几何对象转换为json, 用于json.Marshal。
	MarshalJSON() ([]byte, error)
	// WKT 将几何对象转换为WKT。
	WKT() (string, error)
	// GeoJSON 将几何对象转换为GeoJSON。
	GeoJSON() (string, error)
}

type Geometry[C Coordinate] struct {
	// 几何的类型
	Type GeometryType `json:"type"`
	// 几何的坐标
	Coordinates C `json:"coordinates"`
}

// NewGeometry 创建一个几何对象。
//
// Params:
//
//   - t: 几何的类型。
//   - c: 几何的坐标。
//
// Returns:
//
//   - 几何对象。
//
// Example:
//
//	p, err := geo.NewPoint(lng, lat)
//	if err != nil {
//		return nil, err
//	}
//	g := geo.NewGeometry(PointType, p)
func NewGeometry[T PostGISGeometry](t GeometryType, c T) *Geometry[T] {
	return &Geometry[T]{
		Type:        t,
		Coordinates: c,
	}
}

// NewGeometryByPoint 通过经纬度创建一个点。
//
// Params:
//
//   - lng: 经度。
//   - lat: 纬度。
//
// Returns:
//
//	0: 点的几何对象。
//	1: 错误信息。
//
// Example:
//
//	p, err := geo.NewGeometryByPoint(60, 60)
//
// ErrCodes:
//
//   - Err_0200020201
func NewGeometryByPoint(lng float64, lat float64) (*Geometry[*Point], error) {
	p, err := NewPoint(lng, lat)
	if err != nil {
		return nil, err
	}
	g := NewGeometry(PointType, p)
	return g, nil
}

// NewGeometryByLineString 创建一个lineString。
//
// Params:
//
//   - points: 点。
//
// Returns:
//
//	0: 点的几何对象。
//
// Example:
//
//	p, err := geo.NewGeometryByLineString([][]float64{{1, 2}, {3, 4}})
//
// ErrCodes:
//
//   - Err_0200020201
func NewGeometryByLineString(points [][]float64) (*Geometry[*LineString], error) {
	line, err := NewLineString(points)
	if err != nil {
		return nil, err
	}
	g := NewGeometry(LineStringType, line)
	return g, nil
}

// NewGeometryByLineStringP 创建一个lineString, 通过Point创建。
//
// Params:
//
//   - points: 点。
//
// Returns:
//
//		0: 点的几何对象。
//	 1: 错误信息。
//
// Example:
//
//	p, err := geo.NewGeometryByLineStringP([]Point{p1, p2})
func NewGeometryByLineStringP(points []Point) (*Geometry[*LineString], error) {
	line, err := NewLineStringByPoint(points)
	if err != nil {
		return nil, err
	}
	g := NewGeometry(LineStringType, line)
	return g, nil
}

// String 用于将几何对象转换为字符串。在fmt.Print等函数中会调用此方法。
func (g Geometry[T]) String() string {
	return g.Coordinates.String()
}

type GeometryType string

const (
	PointType              GeometryType = "Point"
	LineStringType         GeometryType = "LineString"
	PolygonType            GeometryType = "Polygon"
	MultiPointType         GeometryType = "MultiPoint"
	MultiLineStringType    GeometryType = "MultiLineString"
	MultiPolygonType       GeometryType = "MultiPolygon"
	GeometryCollectionType GeometryType = "GeometryCollection"
	CircularStringType     GeometryType = "CircularString"
	CompoundCurveType      GeometryType = "CompoundCurve"
	CurvePolygonType       GeometryType = "CurvePolygon"
	MultiCurveType         GeometryType = "MultiCurve"
	MultiSurfaceType       GeometryType = "MultiSurface"
)
