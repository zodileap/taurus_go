package rand

import (
	"fmt"
	"testing"
)

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func TestStringWithCharset(t *testing.T) {
	s, err := StringWithCharset(16, charset)
	fmt.Println(s)
	if err != nil {
		t.Error(err)
	}
}
