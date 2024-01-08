package geojson

// FeatureCollection 类型,
// 是多个 Feature 对象的集合。
//
// `T`: Feature 的 Properties 类型
// `G`: Feature 的 Geometry 类型,
//
//	例如:
//	"geometry": {
//		"type": "LineString",
//		"coordinates": [[125.6, 10.1], [110.6, 20.1]]
//	},
//	则G为geojson.LineString
type FeatureCollection[T interface{}, G interface{}] struct {
	// 值为FeatureCollection
	Type string `json:"type"`
	// FeatureCollection的features
	Features []Feature[T, G] `json:"features"`
}

// Feature 类型
type Feature[P interface{}, G interface{}] struct {
	// 值为Feature
	Type string `json:"type"`
	// Feature的属性
	Properties P `json:"properties"`
	// Feature的几何
	Geometry Geometry[G] `json:"geometry"`
}

type Geometry[T interface{}] struct {
	// 几何的类型
	Type string `json:"type"`
	// 几何的坐标
	Coordinates T `json:"coordinates"`
}

type Point = []float64
type LineString = [][]float64
