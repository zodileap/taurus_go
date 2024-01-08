package redis

import (
	"testing"

	"github.com/yohobala/taurus_go/tlog"
)

func TestDel(t *testing.T) {
	t.Run("测试删除", func(t *testing.T) {
		SetClient(testClientName, testPptions)
		defer ClearClient()
		c, err := GetClient(testClientName)
		if err != nil {
			t.Errorf(err.Error())
		}
		defer c.Close()
		l, err := c.Del("key1", "key2")
		if err != nil {
			t.Errorf(err.Error())
		}
		tlog.Print(l)
	})
}
