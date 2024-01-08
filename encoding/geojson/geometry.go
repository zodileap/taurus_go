package geojson

import (
	"encoding/json"
	"fmt"

	"github.com/yohobala/taurus_go/entity/postgresql"
)

// 关于geometry的类型，需要特殊处理的
const (
	// geojson.Geometry类型，通过taurus_go/enconding/geojson.Geometry设定的类型
	GeometryType                       postgresql.FieldType = "geometry"
	Geometry2163Type                   postgresql.FieldType = "geometry(2163)"
	Geometry3395Type                   postgresql.FieldType = "geometry(3395)"
	Geometry4269Type                   postgresql.FieldType = "geometry(4269)"
	Geometry4326Type                   postgresql.FieldType = "geometry(4326)"
	GeometryPointType                  postgresql.FieldType = "geometry(Point)"
	GeometryPoint2163Type              postgresql.FieldType = "geometry(Point,2163)"
	GeometryPoint3395Type              postgresql.FieldType = "geometry(Point,3395)"
	GeometryPoint4269Type              postgresql.FieldType = "geometry(Point,4269)"
	GeometryPoint4326Type              postgresql.FieldType = "geometry(Point,4326)"
	GeometryLineStringType             postgresql.FieldType = "geometry(LineString)"
	GeometryLineString2163Type         postgresql.FieldType = "geometry(LineString,2163)"
	GeometryLineString3395Type         postgresql.FieldType = "geometry(LineString,3395)"
	GeometryLineString4269Type         postgresql.FieldType = "geometry(LineString,4269)"
	GeometryLineString4326Type         postgresql.FieldType = "geometry(LineString,4326)"
	GeometryPolygonType                postgresql.FieldType = "geometry(Polygon)"
	GeometryPolygon2163Type            postgresql.FieldType = "geometry(Polygon,2163)"
	GeometryPolygon3395Type            postgresql.FieldType = "geometry(Polygon,3395)"
	GeometryPolygon4269Type            postgresql.FieldType = "geometry(Polygon,4269)"
	GeometryPolygon4326Type            postgresql.FieldType = "geometry(Polygon,4326)"
	GeometryMultiPointType             postgresql.FieldType = "geometry(MultiPoint)"
	GeometryMultiPoint2163Type         postgresql.FieldType = "geometry(MultiPoint,2163)"
	GeometryMultiPoint3395Type         postgresql.FieldType = "geometry(MultiPoint,3395)"
	GeometryMultiPoint4269Type         postgresql.FieldType = "geometry(MultiPoint,4269)"
	GeometryMultiPoint4326Type         postgresql.FieldType = "geometry(MultiPoint,4326)"
	GeometryMultiLineStringType        postgresql.FieldType = "geometry(MultiLineString)"
	GeometryMultiLineString2163Type    postgresql.FieldType = "geometry(MultiLineString,2163)"
	GeometryMultiLineString3395Type    postgresql.FieldType = "geometry(MultiLineString,3395)"
	GeometryMultiLineString4269Type    postgresql.FieldType = "geometry(MultiLineString,4269)"
	GeometryMultiLineString4326Type    postgresql.FieldType = "geometry(MultiLineString,4326)"
	GeometryMultiPolygonType           postgresql.FieldType = "geometry(MultiPolygon)"
	GeometryMultiPolygon2163Type       postgresql.FieldType = "geometry(MultiPolygon,2163)"
	GeometryMultiPolygon3395Type       postgresql.FieldType = "geometry(MultiPolygon,3395)"
	GeometryMultiPolygon4269Type       postgresql.FieldType = "geometry(MultiPolygon,4269)"
	GeometryMultiPolygon4326Type       postgresql.FieldType = "geometry(MultiPolygon,4326)"
	GeometryGeometryCollectionType     postgresql.FieldType = "geometry(GeometryCollection)"
	GeometryGeometryCollection2163Type postgresql.FieldType = "geometry(GeometryCollection,2163)"
	GeometryGeometryCollection3395Type postgresql.FieldType = "geometry(GeometryCollection,3395)"
	GeometryGeometryCollection4269Type postgresql.FieldType = "geometry(GeometryCollection,4269)"
	GeometryGeometryCollection4326Type postgresql.FieldType = "geometry(GeometryCollection,4326)"
)

// 用于github.com/jackc/pgx/v5 的scan方法
// 能够实现从数据库数据变成Geometry数据
func (g *Geometry[T]) Scan(value interface{}) error {
	// 将值转换为字符串
	data, ok := value.(string)
	if !ok {
		return Err_geometry_scan_type_error
	}
	// 使用 json.Unmarshal 解析字符串为 Geometry 类型
	err := json.Unmarshal([]byte(data), g)
	if err != nil {
		return Err_geometry_scan_convert_error
	}

	return nil
}

// 用于本库中postgresql包的ConverToFieldType方法
// 实现了把结构体中的类型或者fieldType变成postgresql.FieldType类型
//
// 注意是用Geometry[T]不是 *Geometry[T]
func (g Geometry[T]) ConverToFieldType(fieldType string) (postgresql.FieldType, error) {
	if fieldType == "geojson.Geometry[[]float64]" {
		return GeometryType, nil
	} else if fieldType == "geojson.Geometry[[][]float64]" {
		return GeometryType, nil
	} else if fieldType == "geojson.Geometry[[][][]float64]" {
		return GeometryType, nil
	} else if fieldType == "geometry" {
		return GeometryType, nil
	} else if fieldType == "geometry(2163)" {
		return Geometry2163Type, nil
	} else if fieldType == "geometry(3395)" {
		return Geometry3395Type, nil
	} else if fieldType == "geometry(4269)" {
		return Geometry4269Type, nil
	} else if fieldType == "geometry(4326)" {
		return Geometry4326Type, nil
	} else if fieldType == "geometry(Point)" {
		return GeometryPointType, nil
	} else if fieldType == "geometry(Point,2163)" {
		return GeometryPoint2163Type, nil
	} else if fieldType == "geometry(Point,3395)" {
		return GeometryPoint3395Type, nil
	} else if fieldType == "geometry(Point,4269)" {
		return GeometryPoint4269Type, nil
	} else if fieldType == "geometry(Point,4326)" {
		return GeometryPoint4326Type, nil
	} else if fieldType == "geometry(LineString)" {
		return GeometryLineStringType, nil
	} else if fieldType == "geometry(LineString,2163)" {
		return GeometryLineString2163Type, nil
	} else if fieldType == "geometry(LineString,3395)" {
		return GeometryLineString3395Type, nil
	} else if fieldType == "geometry(LineString,4269)" {
		return GeometryLineString4269Type, nil
	} else if fieldType == "geometry(LineString,4326)" {
		return GeometryLineString4326Type, nil
	} else if fieldType == "geometry(Polygon)" {
		return GeometryPolygonType, nil
	} else if fieldType == "geometry(Polygon,2163)" {
		return GeometryPolygon2163Type, nil
	} else if fieldType == "geometry(Polygon,3395)" {
		return GeometryPolygon3395Type, nil
	} else if fieldType == "geometry(Polygon,4269)" {
		return GeometryPolygon4269Type, nil
	} else if fieldType == "geometry(Polygon,4326)" {
		return GeometryPolygon4326Type, nil
	} else if fieldType == "geometry(MultiPoint)" {
		return GeometryMultiPointType, nil
	} else if fieldType == "geometry(MultiPoint,2163)" {
		return GeometryMultiPoint2163Type, nil
	} else if fieldType == "geometry(MultiPoint,3395)" {
		return GeometryMultiPoint3395Type, nil
	} else if fieldType == "geometry(MultiPoint,4269)" {
		return GeometryMultiPoint4269Type, nil
	} else if fieldType == "geometry(MultiPoint,4326)" {
		return GeometryMultiPoint4326Type, nil
	} else if fieldType == "geometry(MultiLineString)" {
		return GeometryMultiLineStringType, nil
	} else {
		return postgresql.UnknownType, nil
	}
}

func (g Geometry[T]) SelectValue(name string, fieldType postgresql.FieldType) (string, error) {
	var column string
	switch fieldType {
	case GeometryType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case Geometry2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case Geometry3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case Geometry4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case Geometry4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	case GeometryPointType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case GeometryPoint2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case GeometryPoint3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case GeometryPoint4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case GeometryPoint4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	case GeometryLineStringType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case GeometryLineString2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case GeometryLineString3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case GeometryLineString4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case GeometryLineString4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	case GeometryPolygonType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case GeometryPolygon2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case GeometryPolygon3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case GeometryPolygon4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case GeometryPolygon4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	case GeometryMultiPointType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case GeometryMultiPoint2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case GeometryMultiPoint3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case GeometryMultiPoint4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case GeometryMultiPoint4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	case GeometryMultiLineStringType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case GeometryMultiLineString2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case GeometryMultiLineString3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case GeometryMultiLineString4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case GeometryMultiLineString4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	case GeometryMultiPolygonType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case GeometryMultiPolygon2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case GeometryMultiPolygon3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case GeometryMultiPolygon4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case GeometryMultiPolygon4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	case GeometryGeometryCollectionType:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	case GeometryGeometryCollection2163Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,2163))", name)
	case GeometryGeometryCollection3395Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,3395))", name)
	case GeometryGeometryCollection4269Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4269))", name)
	case GeometryGeometryCollection4326Type:
		column = fmt.Sprintf("ST_AsGeoJSON(ST_Transform(%s,4326))", name)
	default:
		column = fmt.Sprintf("ST_AsGeoJSON(%s)", name)
	}
	return column, nil
}

func (g Geometry[T]) InsertValue(name string, fieldType postgresql.FieldType) (string, error) {
	var column string
	switch fieldType {
	case GeometryType:
		column = "ST_GeomFromGeoJSON($%d)"
	case Geometry2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case Geometry3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case Geometry4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case Geometry4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	case GeometryPointType:
		column = "ST_GeomFromGeoJSON($%d)"
	case GeometryPoint2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case GeometryPoint3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case GeometryPoint4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case GeometryPoint4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	case GeometryLineStringType:
		column = "ST_GeomFromGeoJSON($%d)"
	case GeometryLineString2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case GeometryLineString3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case GeometryLineString4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case GeometryLineString4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	case GeometryPolygonType:
		column = "ST_GeomFromGeoJSON($%d)"
	case GeometryPolygon2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case GeometryPolygon3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case GeometryPolygon4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case GeometryPolygon4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	case GeometryMultiPointType:
		column = "ST_GeomFromGeoJSON($%d)"
	case GeometryMultiPoint2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case GeometryMultiPoint3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case GeometryMultiPoint4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case GeometryMultiPoint4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	case GeometryMultiLineStringType:
		column = "ST_GeomFromGeoJSON($%d)"
	case GeometryMultiLineString2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case GeometryMultiLineString3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case GeometryMultiLineString4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case GeometryMultiLineString4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	case GeometryMultiPolygonType:
		column = "ST_GeomFromGeoJSON($%d)"
	case GeometryMultiPolygon2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case GeometryMultiPolygon3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case GeometryMultiPolygon4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case GeometryMultiPolygon4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	case GeometryGeometryCollectionType:
		column = "ST_GeomFromGeoJSON($%d)"
	case GeometryGeometryCollection2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163)"
	case GeometryGeometryCollection3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395)"
	case GeometryGeometryCollection4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269)"
	case GeometryGeometryCollection4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326)"
	default:
		column = "ST_GeomFromGeoJSON($%d)"
	}
	return column, nil
}

func (g Geometry[T]) UpdateCaseValue(name string, fieldType postgresql.FieldType) (string, error) {
	var column string
	switch fieldType {
	case GeometryType:
		column = "ST_GeomFromGeoJSON($%d) "
	case Geometry2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case Geometry3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case Geometry4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case Geometry4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	case GeometryPointType:
		column = "ST_GeomFromGeoJSON($%d) "
	case GeometryPoint2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case GeometryPoint3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case GeometryPoint4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case GeometryPoint4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	case GeometryLineStringType:
		column = "ST_GeomFromGeoJSON($%d) "
	case GeometryLineString2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case GeometryLineString3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case GeometryLineString4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case GeometryLineString4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	case GeometryPolygonType:
		column = "ST_GeomFromGeoJSON($%d) "
	case GeometryPolygon2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case GeometryPolygon3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case GeometryPolygon4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case GeometryPolygon4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	case GeometryMultiPointType:
		column = "ST_GeomFromGeoJSON($%d) "
	case GeometryMultiPoint2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case GeometryMultiPoint3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case GeometryMultiPoint4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case GeometryMultiPoint4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	case GeometryMultiLineStringType:
		column = "ST_GeomFromGeoJSON($%d) "
	case GeometryMultiLineString2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case GeometryMultiLineString3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case GeometryMultiLineString4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case GeometryMultiLineString4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	case GeometryMultiPolygonType:
		column = "ST_GeomFromGeoJSON($%d) "
	case GeometryMultiPolygon2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case GeometryMultiPolygon3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case GeometryMultiPolygon4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case GeometryMultiPolygon4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	case GeometryGeometryCollectionType:
		column = "ST_GeomFromGeoJSON($%d) "
	case GeometryGeometryCollection2163Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),2163) "
	case GeometryGeometryCollection3395Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),3395) "
	case GeometryGeometryCollection4269Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4269) "
	case GeometryGeometryCollection4326Type:
		column = "ST_SetSRID(ST_GeomFromGeoJSON($%d),4326) "
	default:
		column = "ST_GeomFromGeoJSON($%d) "
	}
	return column, nil
}

func (g Geometry[T]) ReturningValue(name string, fieldType postgresql.FieldType) (string, error) {
	return g.SelectValue(name, fieldType)
}
