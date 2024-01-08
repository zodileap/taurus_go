package redis

import (
	"time"

	"github.com/yohobala/taurus_go/tlog"
)

// HashFieldSet 设置一个hash类型的键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 哈希的key
//   - field: 哈希内部的字段
//   - value: 字段的值
func HashFieldSet(clientName string, key string, field string, value string) error {
	client, ok := pool[clientName]
	if !ok {
		return Err_nil_options
	}

	err := client.HSet(ctx, key, field, value).Err()
	if err != nil {
		err := validErr(err)
		return err
	}
	return nil
}

// HashSet 设置整个hash类型的对象
//
// 参数：
//   - clientName: 客户端名称
//   - key: 哈希的key
//   - fields: 表示整个哈希对象的map
//   - expiration: 过期时间，单位为秒, 如果设置为0，则表示不过期
func HashSet(clientName string, key string, fields map[string]interface{}, expiration time.Duration) error {
	client, ok := pool[clientName]
	if !ok {
		return Err_nil_options
	}

	var flatFields []interface{}
	for k, v := range fields {
		flatFields = append(flatFields, k, v)
	}

	// 使用HMSet命令来一次性设置多个字段
	err := client.HSet(ctx, key, flatFields...).Err()
	if err != nil {
		err := validErr(err)
		return err
	}
	tlog.Print(expiration)
	// 设置过期时间
	if expiration > 0 {
		client.Expire(ctx, key, expiration)
	}

	return nil
}

// HashFieldGet 获取hash对象中的一个键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 哈希的key
//   - field: 需要得到数据的字段
func HashFieldGet(clientName string, key string, field string) (string, error) {
	client, ok := pool[clientName]
	if !ok {
		return "", Err_nil_options
	}

	val, err := client.HGet(ctx, key, field).Result()
	if err != nil {
		err := validErr(err)
		return "", err
	}

	return val, nil
}

// HashGet 获取整个hash对象
//
// 参数：
//   - clientName: 客户端名称
//   - key: 哈希的key
func HashGet(clientName string, key string) (map[string]string, error) {
	client, ok := pool[clientName]
	if !ok {
		return nil, Err_nil_options
	}

	// 没有结果是返回空map
	val, err := client.HGetAll(ctx, key).Result()
	if err != nil {
		err := validErr(err)
		return nil, err
	}

	if len(val) == 0 {
		// 处理空 map 的情况
		// 你可以返回一个特定的错误，或者按照你的需求进行处理
		return nil, Err_not_key
	}

	return val, nil
}

// HashFieldDelete 删除hash对象中的一个键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 哈希的key
func HashFieldDelete(clientName string, key string, field string) error {
	client, ok := pool[clientName]
	if !ok {
		return Err_nil_options
	}

	err := client.HDel(ctx, key, field).Err()
	if err != nil {
		err := validErr(err)
		return err
	}

	return nil
}

// HashDelete 删除整个hash对象
//
// 参数：
//   - clientName: 客户端名称
//   - key: 哈希的key
func HashDelete(clientName string, key string) error {
	client, ok := pool[clientName]
	if !ok {
		return Err_nil_options
	}

	err := client.Del(ctx, key).Err()
	if err != nil {
		err := validErr(err)
		return err
	}

	return nil
}
