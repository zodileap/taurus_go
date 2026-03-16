package redis

import "testing"

func TestGetClientRequiresConfig(t *testing.T) {
	ClearClient()

	_, err := GetClient("missing")
	requireRedisErrCode(t, err, Err_0300010001.Code())
}

func TestClientDel(t *testing.T) {
	client, server := newTestClient(t)
	server.Set("key1", "1")
	server.Set("key2", "2")

	deleted, err := client.Del("key1", "key2")
	if err != nil {
		t.Fatalf("删除键失败: %v", err)
	}
	if deleted != 2 {
		t.Fatalf("删除数量不匹配，期望 2，实际 %d", deleted)
	}
}
