package redis

import (
	"context"
	"testing"
	"time"

	miniredis "github.com/alicebob/miniredis/v2"
)

func TestSetClient(t *testing.T) {
	ClearClient()
	t.Cleanup(ClearClient)

	SetClient("test", &Options{Addr: "127.0.0.1:6379", DB: 1})

	option := clientOptions["test"]
	if option == nil {
		t.Fatal("SetClient 未写入客户端配置")
	}
	if option.Addr != "127.0.0.1:6379" || option.DB != 1 {
		t.Fatalf("SetClient 配置不正确: %+v", option)
	}
}

func TestGetClient(t *testing.T) {
	ClearClient()
	t.Cleanup(ClearClient)

	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 miniredis 失败: %v", err)
	}
	defer server.Close()

	SetClient("test", &Options{Addr: server.Addr()})

	client, err := GetClient("test")
	if err != nil {
		t.Fatalf("GetClient 失败: %v", err)
	}
	defer client.Close()

	if client.Name != "test" {
		t.Fatalf("客户端名称不正确: %s", client.Name)
	}
	if client.Options().Addr != server.Addr() {
		t.Fatalf("客户端地址不正确，期望 %s，实际 %s", server.Addr(), client.Options().Addr)
	}
}

func TestClearClient(t *testing.T) {
	ClearClient()
	SetClient("test", &Options{Addr: "127.0.0.1:6379"})

	ClearClient()

	if len(clientOptions) != 0 {
		t.Fatalf("ClearClient 未清空配置，剩余 %d 项", len(clientOptions))
	}
	_, err := GetClient("test")
	requireRedisErrCode(t, err, Err_0300010001.Code())
}

func TestClose(t *testing.T) {
	ClearClient()
	t.Cleanup(ClearClient)

	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 miniredis 失败: %v", err)
	}
	defer server.Close()

	SetClient("test", &Options{Addr: server.Addr()})

	client, err := GetClient("test")
	if err != nil {
		t.Fatalf("GetClient 失败: %v", err)
	}
	if err := client.Close(); err != nil {
		t.Fatalf("Close 失败: %v", err)
	}
}

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

func TestSave(t *testing.T) {
	client, server := newTestClient(t)
	server.Set("save:key", "value")

	str, err := client.String()
	if err != nil {
		t.Fatalf("创建 String 客户端失败: %v", err)
	}
	str.Get("save:key")

	if _, err := client.Save(); err != nil {
		t.Fatalf("Save 失败: %v", err)
	}
	if len(client.nodes) != 0 {
		t.Fatalf("Save 后应清空 tracker，实际剩余 %d", len(client.nodes))
	}

	str, err = client.String()
	if err != nil {
		t.Fatalf("重新创建 String 客户端失败: %v", err)
	}
	str.Add("save:key-2", 0, "value-2")
	if _, err := client.Save(); err != nil {
		t.Fatalf("Save 后再次执行失败: %v", err)
	}
}

func TestSetExpire(t *testing.T) {
	client, server := newTestClient(t)

	if err := client.Client.Set(context.Background(), "expiring:key", "value", 0).Err(); err != nil {
		t.Fatalf("准备测试数据失败: %v", err)
	}

	client.SetExpire("expiring:key", 5*time.Second)

	ttl, err := client.Client.TTL(context.Background(), "expiring:key").Result()
	if err != nil {
		t.Fatalf("读取 TTL 失败: %v", err)
	}
	if ttl <= 0 {
		t.Fatalf("TTL 未生效: %v", ttl)
	}

	server.FastForward(6 * time.Second)
	if server.Exists("expiring:key") {
		t.Fatal("键在过期后仍然存在")
	}
}
