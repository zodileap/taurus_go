package geojson

import "fmt"

var Err_geometry_scan_type_error = fmt.Errorf("Geometry: 数据类型错误")

var Err_geometry_scan_convert_error = fmt.Errorf("Geometry: 解析成Geometry数据失败")
