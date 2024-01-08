package redis

import (
	"time"

	"github.com/yohobala/taurus_go/tlog"
)

type Set struct {
	client            *Client
	defaultValue      []string
	defaultExpiration time.Duration
}

// SetAdd 设置一个set类型的键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 哈希的key
//   - value: 字段的值
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期
func SetAdd(clientName string, key string, value string, expiration time.Duration) (setLength int64, err error) {
	client, ok := pool[clientName]
	if !ok {
		return 0, Err_nil_options
	}
	tlog.Print(key)
	tlog.Print(value)
	tlog.Print(expiration)
	setLength, err = client.SAdd(ctx, key, value).Result()
	if err != nil {
		err := validErr(err)
		return 0, err
	}
	// 设置过期时间
	if expiration > 0 {
		client.Expire(ctx, key, expiration)
	}

	return setLength, nil
}

// SetGet 获取set对象中的一个键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 需要得到数据的key
func SetGet(clientName string, key string) ([]string, error) {
	client, ok := pool[clientName]
	if !ok {
		return nil, Err_nil_options
	}

	val, err := client.SMembers(ctx, key).Result()
	if err != nil {
		err := validErr(err)
		return nil, err
	}

	return val, nil
}

// SetDelete 删除set对象中的一个键值对
//
// 参数：
//   - clientName: 客户端名称
//   - key: 需要删除数据的key
func SetDelete(clientName string, key string, value string) error {
	client, ok := pool[clientName]
	if !ok {
		return Err_nil_options
	}

	err := client.SRem(ctx, key, value).Err()
	if err != nil {
		err := validErr(err)
		return err
	}

	return nil
}
