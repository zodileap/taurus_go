package geo

import "testing"

func TestPointRoundTrip(t *testing.T) {
	point, err := NewPoint(120.5, 30.25)
	if err != nil {
		t.Fatalf("创建 Point 失败: %v", err)
	}

	wkt, err := point.WKT()
	if err != nil {
		t.Fatalf("生成 WKT 失败: %v", err)
	}
	if wkt != "POINT(120.5 30.25)" {
		t.Fatalf("WKT 不正确: %s", wkt)
	}

	var decoded Point
	if err := decoded.UnmarshalWKT(wkt); err != nil {
		t.Fatalf("解析 WKT 失败: %v", err)
	}
	if decoded.String() != wkt {
		t.Fatalf("WKT 往返结果不一致: %s", decoded.String())
	}

	json, err := point.GeoJSON()
	if err != nil {
		t.Fatalf("生成 GeoJSON 失败: %v", err)
	}

	var fromGeoJSON Point
	if err := fromGeoJSON.UnmarshalGeoJSON(json); err != nil {
		t.Fatalf("解析 GeoJSON 失败: %v", err)
	}
	if fromGeoJSON.String() != wkt {
		t.Fatalf("GeoJSON 往返结果不一致: %s", fromGeoJSON.String())
	}
}
