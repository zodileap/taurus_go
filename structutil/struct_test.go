package structutil

import "testing"

func TestGetFields(t *testing.T) {
	type Employee struct {
		Id   int
		Name string
	}
	fields := GetFields(Employee{})
	if len(fields) != 2 || fields[0] != "Id" || fields[1] != "Name" {
		t.Fatalf("GetFields(Employee{}) 结果不正确: %v", fields)
	}

	fieldsByPtr := GetFields(&Employee{})
	if len(fieldsByPtr) != 2 || fieldsByPtr[0] != "Id" || fieldsByPtr[1] != "Name" {
		t.Fatalf("GetFields(&Employee{}) 结果不正确: %v", fieldsByPtr)
	}

	type Empty struct{}
	fields2 := GetFields(Empty{})
	if len(fields2) != 0 {
		t.Fatalf("GetFields(Empty{}) 结果不正确: %v", fields2)
	}

	if got := GetFields(123); got != nil {
		t.Fatalf("非结构体输入应返回 nil，实际: %v", got)
	}
}
