package structutil

import "reflect"

// GetFields 返回结构体的字段名列表。
//
// Example:
//
//	type Employee struct {
//		Id   int
//		Name string
//	}
//
//	fields := structutil.GetFields(Employee{})
//	// [Id Name]
//
// ExamplePath: structutil/struct_test.go - TestGetFields
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
