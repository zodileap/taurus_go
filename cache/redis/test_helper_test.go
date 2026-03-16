package redis

import (
	"errors"
	"sort"
	"strings"
	"testing"

	miniredis "github.com/alicebob/miniredis/v2"

	terr "github.com/zodileap/taurus_go/err"
)

func newTestClient(t *testing.T) (*Client, *miniredis.Miniredis) {
	t.Helper()

	server, err := miniredis.Run()
	if err != nil {
		t.Fatalf("启动 miniredis 失败: %v", err)
	}
	t.Cleanup(server.Close)
	t.Cleanup(ClearClient)

	clientName := strings.NewReplacer("/", "-", " ", "-").Replace(t.Name())
	SetClient(clientName, &Options{Addr: server.Addr()})

	client, err := GetClient(clientName)
	if err != nil {
		t.Fatalf("获取测试客户端失败: %v", err)
	}
	t.Cleanup(func() {
		if err := client.Close(); err != nil {
			t.Fatalf("关闭测试客户端失败: %v", err)
		}
	})

	return client, server
}

func requireRedisErrCode(t *testing.T, got error, want string) {
	t.Helper()

	if got == nil {
		t.Fatalf("期望错误码 %s，实际无错误", want)
	}

	var errCode terr.ErrCode
	if !errors.As(got, &errCode) {
		t.Fatalf("期望 ErrCode，实际为 %T: %v", got, got)
	}
	if errCode.Code() != want {
		t.Fatalf("错误码不匹配，期望 %s，实际 %s", want, errCode.Code())
	}
}

func sortedStrings(values []string) []string {
	cloned := append([]string(nil), values...)
	sort.Strings(cloned)
	return cloned
}
