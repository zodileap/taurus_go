package geo

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
type FeatureCollection[P interface{}, C Coordinate] struct {
	// 值为FeatureCollection
	Type string `json:"type"`
	// FeatureCollection的features
	Features []Feature[P, C] `json:"features"`
}

// NewFeatureCollection 用于生成Geojson中的一个FeatureCollection。
//
// Params:
//
//   - features: 一个Feature的数组。
//
// Returns:
//
//	0: 一个FeatureCollection。
//
// Example:
//
// ExamplePath:  taurus_go_demo/asset/asset_test.go
func NewFeatureCollection[P interface{}, C Coordinate](features []Feature[P, C]) FeatureCollection[P, C] {
	featureCollection := FeatureCollection[P, C]{
		Type:     "FeatureCollection",
		Features: features,
	}
	return featureCollection
}
