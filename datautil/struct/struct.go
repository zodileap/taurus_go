package structutil

import "reflect"

//	获取结构体的字段名
//
// 示例:
//
//	type Employee struct {
//		ID   int
//	 Name string
//	}
//
// fields := getFields(Employee{})
// 输出：[ID Name]
func GetFields(input interface{}) []string {
	val := reflect.TypeOf(input)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return nil
	}

	var fields []string
	for i := 0; i < val.NumField(); i++ {
		fields = append(fields, val.Field(i).Name)
	}
	return fields
}
