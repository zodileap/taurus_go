package unit

import (
	"testing"

	terr "github.com/zodileap/taurus_go/err"
)

func TestValidErr(t *testing.T) {
	got := terr.New("0100020015", "msg %s", "").Sprintf("a")
	want := terr.New("0100020015", "msg %s", "").Sprintf("b")
	ValidErr(&got, &want, t)
}

func TestValidRes(t *testing.T) {
	ValidRes(map[string]int{"a": 1}, map[string]int{"a": 1}, t)
}

func TestMust(t *testing.T) {
	Must(t, nil)
}
