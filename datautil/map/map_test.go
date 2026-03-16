package maputil

import "testing"

func TestInterfaceToString(t *testing.T) {
	got := InterfaceToString(map[string]interface{}{
		"int":    1,
		"string": "hello",
		"bool":   true,
	})

	if got["int"] != "1" || got["string"] != "hello" || got["bool"] != "true" {
		t.Fatalf("InterfaceToString 结果不正确: %v", got)
	}
}
