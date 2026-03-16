package redis

import (
	"slices"
	"testing"
)

func TestSetSave(t *testing.T) {
	client, _ := newTestClient(t)

	setClient, err := client.Set()
	if err != nil {
		t.Fatalf("创建 Set 客户端失败: %v", err)
	}

	setClient.Add("set:key", 0, "1")
	setClient.Add("set:key", 0, "2")
	setClient.Add("set:key", 0, "3")
	setClient.Del("set:key", "3")
	setClient.Get("set:key")

	res, err := client.Save()
	if err != nil {
		t.Fatalf("保存 Set 管道失败: %v", err)
	}

	setRes := res.GetSet("set:key")
	if setRes == nil {
		t.Fatal("缺少 set:key 的结果")
	}
	if setRes.AddNum != 3 {
		t.Fatalf("AddNum 不匹配，期望 3，实际 %d", setRes.AddNum)
	}
	if setRes.DelNum != 1 {
		t.Fatalf("DelNum 不匹配，期望 1，实际 %d", setRes.DelNum)
	}
	if setRes.Oper != SMembers {
		t.Fatalf("最后操作不匹配，期望 %s，实际 %s", SMembers, setRes.Oper)
	}
	if !slices.Equal(sortedStrings(setRes.Value), []string{"1", "2"}) {
		t.Fatalf("Set 值不匹配，实际 %v", setRes.Value)
	}
}

func TestSetDirectCommands(t *testing.T) {
	client, _ := newTestClient(t)

	setClient, err := client.Set()
	if err != nil {
		t.Fatalf("创建 Set 客户端失败: %v", err)
	}

	added, err := setClient.AddR("set:key", 0, "1")
	if err != nil {
		t.Fatalf("AddR 失败: %v", err)
	}
	if added != 1 {
		t.Fatalf("首次 AddR 结果不匹配，期望 1，实际 %d", added)
	}

	added, err = setClient.AddR("set:key", 0, "1")
	if err != nil {
		t.Fatalf("重复 AddR 失败: %v", err)
	}
	if added != 0 {
		t.Fatalf("重复 AddR 结果不匹配，期望 0，实际 %d", added)
	}

	values, err := setClient.GetR("set:key")
	if err != nil {
		t.Fatalf("GetR 失败: %v", err)
	}
	if !slices.Equal(sortedStrings(values), []string{"1"}) {
		t.Fatalf("GetR 结果不匹配，实际 %v", values)
	}

	deleted, err := setClient.DelR("set:key", "1")
	if err != nil {
		t.Fatalf("DelR 失败: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("DelR 结果不匹配，期望 1，实际 %d", deleted)
	}

	if _, err := setClient.AddR("set:key2", 0, "2"); err != nil {
		t.Fatalf("为 DelAllR 准备数据失败: %v", err)
	}
	deleted, err = setClient.DelAllR("set:key2")
	if err != nil {
		t.Fatalf("DelAllR 失败: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("DelAllR 结果不匹配，期望 1，实际 %d", deleted)
	}
}

func TestSetWrongType(t *testing.T) {
	client, server := newTestClient(t)
	server.Set("wrong:type", "value")

	setClient, err := client.Set()
	if err != nil {
		t.Fatalf("创建 Set 客户端失败: %v", err)
	}

	_, err = setClient.AddR("wrong:type", 0, "member")
	requireRedisErrCode(t, err, Err_0300010004.Code())
}
