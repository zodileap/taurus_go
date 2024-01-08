package geojson

func GenertaeFeatureCollection[T interface{}, G interface{}](features []Feature[T, G]) FeatureCollection[T, G] {
	featureCollection := FeatureCollection[T, G]{
		Type:     "FeatureCollection",
		Features: features,
	}
	return featureCollection
}

func GenerateFeature[T interface{}, G interface{}](properties T, geometry Geometry[G]) Feature[T, G] {
	feature := Feature[T, G]{
		Type:       "Feature",
		Properties: properties,
		Geometry:   geometry,
	}
	return feature
}

func GeneratePoint(x, y float64) Geometry[Point] {
	g := Geometry[Point]{
		Type:        "Point",
		Coordinates: Point{x, y},
	}
	return g
}

func GenerateLineString(lineString LineString) Geometry[LineString] {
	g := Geometry[LineString]{
		Type:        "LineString",
		Coordinates: lineString,
	}
	return g
}
