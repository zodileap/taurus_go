package geo

// Feature 类型
type Feature[P interface{}, C Coordinate] struct {
	// 值为Feature
	Type string `json:"type"`
	// Feature的属性
	Properties P `json:"properties"`
	// Feature的几何
	Geometry Geometry[C] `json:"geometry"`
}

func NewFeature[T interface{}, C Coordinate](properties T, geometry Geometry[C]) Feature[T, C] {
	feature := Feature[T, C]{
		Type:       "Feature",
		Properties: properties,
		Geometry:   geometry,
	}
	return feature
}
