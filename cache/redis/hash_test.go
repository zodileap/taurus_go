package redis

import (
	"slices"
	"testing"
)

func TestHashSave(t *testing.T) {
	client, _ := newTestClient(t)

	hashClient, err := client.Hash()
	if err != nil {
		t.Fatalf("创建 Hash 客户端失败: %v", err)
	}

	hashClient.Add("hash:key", 0, "field1", "1")
	hashClient.AddM("hash:key", 0, map[string]string{
		"field2": "2",
		"field3": "3",
	})
	hashClient.GetAll("hash:key")

	res, err := client.Save()
	if err != nil {
		t.Fatalf("保存 Hash 管道失败: %v", err)
	}

	hashRes := res.GetHash("hash:key")
	if hashRes == nil {
		t.Fatal("缺少 hash:key 的结果")
	}
	if hashRes.AddNum != 3 {
		t.Fatalf("AddNum 不匹配，期望 3，实际 %d", hashRes.AddNum)
	}
	if hashRes.Oper != HGetAll {
		t.Fatalf("最后操作不匹配，期望 %s，实际 %s", HGetAll, hashRes.Oper)
	}
	expected := map[string]string{
		"field1": "1",
		"field2": "2",
		"field3": "3",
	}
	if !slices.Equal(sortedStrings(mapKeys(hashRes.MapValue)), sortedStrings(mapKeys(expected))) {
		t.Fatalf("Hash 字段集合不匹配，实际 %v", hashRes.MapValue)
	}
	for key, value := range expected {
		if hashRes.MapValue[key] != value {
			t.Fatalf("Hash 字段 %s 的值不匹配，期望 %s，实际 %s", key, value, hashRes.MapValue[key])
		}
	}
}

func TestHashDirectCommands(t *testing.T) {
	client, _ := newTestClient(t)

	hashClient, err := client.Hash()
	if err != nil {
		t.Fatalf("创建 Hash 客户端失败: %v", err)
	}

	added, err := hashClient.AddR("hash:key", 0, "field1", "1")
	if err != nil {
		t.Fatalf("AddR 失败: %v", err)
	}
	if added != 1 {
		t.Fatalf("AddR 结果不匹配，期望 1，实际 %d", added)
	}

	added, err = hashClient.AddMR("hash:key", 0, map[string]string{
		"field2": "2",
		"field3": "3",
	})
	if err != nil {
		t.Fatalf("AddMR 失败: %v", err)
	}
	if added != 2 {
		t.Fatalf("AddMR 结果不匹配，期望 2，实际 %d", added)
	}

	value, err := hashClient.GetR("hash:key", "field1")
	if err != nil {
		t.Fatalf("GetR 失败: %v", err)
	}
	if value != "1" {
		t.Fatalf("GetR 结果不匹配，期望 1，实际 %q", value)
	}

	values, err := hashClient.GetMR("hash:key", "field1", "field2", "missing")
	if err != nil {
		t.Fatalf("GetMR 失败: %v", err)
	}
	if len(values) != 3 || values[0] != "1" || values[1] != "2" || values[2] != nil {
		t.Fatalf("GetMR 结果不匹配，实际 %v", values)
	}

	allValues, err := hashClient.GetValsR("hash:key")
	if err != nil {
		t.Fatalf("GetValsR 失败: %v", err)
	}
	if !slices.Equal(sortedStrings(allValues), []string{"1", "2", "3"}) {
		t.Fatalf("GetValsR 结果不匹配，实际 %v", allValues)
	}

	allFields, err := hashClient.GetAllR("hash:key")
	if err != nil {
		t.Fatalf("GetAllR 失败: %v", err)
	}
	if len(allFields) != 3 || allFields["field1"] != "1" || allFields["field2"] != "2" || allFields["field3"] != "3" {
		t.Fatalf("GetAllR 结果不匹配，实际 %v", allFields)
	}

	deleted, err := hashClient.DelR("hash:key", "field1")
	if err != nil {
		t.Fatalf("DelR 失败: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("DelR 结果不匹配，期望 1，实际 %d", deleted)
	}

	deleted, err = hashClient.DelAllR("hash:key")
	if err != nil {
		t.Fatalf("DelAllR 失败: %v", err)
	}
	if deleted != 1 {
		t.Fatalf("DelAllR 结果不匹配，期望 1，实际 %d", deleted)
	}
}

func TestHashMissingAndWrongType(t *testing.T) {
	t.Run("缺失字段返回错误码", func(t *testing.T) {
		client, _ := newTestClient(t)

		hashClient, err := client.Hash()
		if err != nil {
			t.Fatalf("创建 Hash 客户端失败: %v", err)
		}

		_, err = hashClient.GetR("hash:key", "missing")
		requireRedisErrCode(t, err, Err_0300010003.Code())
	})

	t.Run("错误类型返回错误码", func(t *testing.T) {
		client, server := newTestClient(t)
		server.Set("wrong:type", "value")

		hashClient, err := client.Hash()
		if err != nil {
			t.Fatalf("创建 Hash 客户端失败: %v", err)
		}

		_, err = hashClient.AddR("wrong:type", 0, "field", "value")
		requireRedisErrCode(t, err, Err_0300010004.Code())
	})
}

func mapKeys(values map[string]string) []string {
	keys := make([]string, 0, len(values))
	for key := range values {
		keys = append(keys, key)
	}
	return keys
}
