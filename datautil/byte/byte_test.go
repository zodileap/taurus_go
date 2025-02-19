package byteutil

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/yohobala/taurus_go/tlog"
)

func TestUpdate(t *testing.T) {

	t.Run("int-1", func(t *testing.T) {
		n := 1
		b := IntToBytes(n)
		tlog.Printf("IntToBytes: %v", b)
		on := BytesToInt(b)
		tlog.Printf("BytesToInt: %v", on)
	})

	t.Run("int-2", func(t *testing.T) {
		n := 69000000
		b := IntToBytes(n)
		tlog.Printf("IntToBytes: %v", b)
	})

	t.Run("intSToBytes", func(t *testing.T) {
		// [123 49 44 50 44 51 125]
		by := []byte{44, 51, 125}
		tlog.Print(string(by))
		n := []int{1, 2, 3}
		b := IntSToBytes(n, "")
		tlog.Printf("IntSliceToBytes: %v", b)
	})

	t.Run("postgreSQL Array To Slice", func(t *testing.T) {
		// [123 49 44 50 44 51 125]
		b := StringToBytes("{}")
		tlog.Printf("StringToBytes: %v", b)
		b = StringToBytes("[]")
		tlog.Printf("StringToBytes: %v", b)
		b = StringToBytes("true")
		tlog.Printf("StringToBytes: %v", b)
		b = StringToBytes("false")
		tlog.Printf("StringToBytes: %v", b)
		pb := []byte{123, 49, 44, 50, 44, 51, 125}
		pb = ReplaceAll(pb, 123, []byte{91})
		pb = ReplaceAll(pb, 125, []byte{93})
		ps := string(pb)
		tlog.Print(ps)
		var oneDSlice []int
		if err := json.Unmarshal(pb, &oneDSlice); err != nil {
			fmt.Println("一维数组解析错误:", err)
			return
		}
		tlog.Print(oneDSlice)
	})
}
