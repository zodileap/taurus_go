package redis

import "testing"

func TestStringSave(t *testing.T) {
	client, _ := newTestClient(t)

	str, err := client.String()
	if err != nil {
		t.Fatalf("创建 String 客户端失败: %v", err)
	}

	str.Add("string:key", 0, "value")
	str.Get("string:key")

	res, err := client.Save()
	if err != nil {
		t.Fatalf("保存 String 管道失败: %v", err)
	}

	stringRes := res.GetString("string:key")
	if stringRes == nil {
		t.Fatal("缺少 string:key 的结果")
	}
	if stringRes.Value != "value" {
		t.Fatalf("Get 结果不匹配，期望 value，实际 %q", stringRes.Value)
	}
	if stringRes.Oper != Get {
		t.Fatalf("最后操作不匹配，期望 %s，实际 %s", Get, stringRes.Oper)
	}
}

func TestStringSaveDelOverridesOtherOps(t *testing.T) {
	client, _ := newTestClient(t)

	str, err := client.String()
	if err != nil {
		t.Fatalf("创建 String 客户端失败: %v", err)
	}

	str.Add("string:key", 0, "value")
	str.Get("string:key")
	str.Del("string:key")

	res, err := client.Save()
	if err != nil {
		t.Fatalf("保存 String 管道失败: %v", err)
	}

	stringRes := res.GetString("string:key")
	if stringRes == nil {
		t.Fatal("缺少 string:key 的结果")
	}
	if stringRes.Value != "" {
		t.Fatalf("Del 覆盖前序操作后不应保留值，实际 %q", stringRes.Value)
	}
	if stringRes.Oper != Del {
		t.Fatalf("最后操作不匹配，期望 %s，实际 %s", Del, stringRes.Oper)
	}
}

func TestStringSaveMissingKey(t *testing.T) {
	client, _ := newTestClient(t)

	str, err := client.String()
	if err != nil {
		t.Fatalf("创建 String 客户端失败: %v", err)
	}

	str.Get("missing:key")

	_, err = client.Save()
	requireRedisErrCode(t, err, Err_0300010003.Code())
}

func TestStringDirectCommands(t *testing.T) {
	client, _ := newTestClient(t)

	str, err := client.String()
	if err != nil {
		t.Fatalf("创建 String 客户端失败: %v", err)
	}

	if err := str.AddR("string:key", 0, "value"); err != nil {
		t.Fatalf("AddR 失败: %v", err)
	}

	value, err := str.GetR("string:key")
	if err != nil {
		t.Fatalf("GetR 失败: %v", err)
	}
	if value != "value" {
		t.Fatalf("GetR 结果不匹配，期望 value，实际 %q", value)
	}

	if err := str.DelR("string:key"); err != nil {
		t.Fatalf("DelR 失败: %v", err)
	}

	_, err = str.GetR("string:key")
	requireRedisErrCode(t, err, Err_0300010003.Code())
}
