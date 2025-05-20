package structutil

import (
	"fmt"
	"testing"
)

func TestGetFields(t *testing.T) {
	type Employee struct {
		Id   int
		Name string
	}
	fields := GetFields(Employee{})
	fmt.Print(fields)
	if len(fields) != 2 {
		t.Error("getFields() failed")
	}

	type Employee2 struct {
	}
	fields2 := GetFields(Employee2{})
	fmt.Print(fields2)
	if len(fields2) != 0 {
		t.Error("getFields() failed")
	}
}
