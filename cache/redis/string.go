package redis

import (
	"time"
)

// 设置一个string类型的键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 键
//   - value: 值
//   - expiration: 过期时间，单位为秒,如果设置为0，则表示不过期
func StringSet(clientName string, key string, value string, expiration time.Duration) error {
	client, ok := pool[clientName]
	if !ok {
		return Err_nil_options
	}

	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		err := validErr(err)
		return err
	}

	return nil
}

// 获取一个string类型的键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 需要得到数据的key
func StringGet(clientName string, key string) (string, error) {
	client, ok := pool[clientName]
	if !ok {
		return "", Err_nil_options
	}

	val, err := client.Get(ctx, key).Result()
	if err != nil {
		err := validErr(err)
		return "", err
	}

	return val, nil
}

// 删除一个string类型的键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 需要删除数据的key
func StringDelete(clientName string, key string) error {
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
