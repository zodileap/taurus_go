package geo

import "testing"

type testFeatureProps struct {
	Name string
}

func TestNewFeatureCollection(t *testing.T) {
	point, err := NewPoint(120.5, 30.25)
	if err != nil {
		t.Fatalf("创建 Point 失败: %v", err)
	}

	feature := NewFeature(testFeatureProps{Name: "demo"}, *NewGeometry(PointType, point))
	collection := NewFeatureCollection([]Feature[testFeatureProps, *Point]{feature})

	if collection.Type != "FeatureCollection" {
		t.Fatalf("FeatureCollection 类型不正确: %s", collection.Type)
	}
	if len(collection.Features) != 1 {
		t.Fatalf("FeatureCollection 数量不正确: %d", len(collection.Features))
	}
	if collection.Features[0].Properties.Name != "demo" {
		t.Fatalf("FeatureCollection 属性不正确: %+v", collection.Features[0].Properties)
	}
	if collection.Features[0].Geometry.Type != PointType {
		t.Fatalf("FeatureCollection 几何类型不正确: %s", collection.Features[0].Geometry.Type)
	}
}
