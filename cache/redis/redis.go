package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	// Options redis客户端的配置
	Options = redis.Options
	Client  struct {
		*redis.Client
	}
)

var (
	// ctx redis客户端的上下文,[go-redis配置项](https://redis.uptrace.dev/zh/guide/go-redis-option.html#redis-client)
	ctx context.Context = context.Background()

	// pool redis客户端的连接池
	clientOptions map[string]*Options = make(map[string]*Options)
)

// SetClient 用于新建一个redis连接，并添加到连接池中。
//
// Params:
//
//   - clientName: redis客户端的名称。
//   - options: redis客户端的配置。
//
// Example:
//
//	redis.SetClient("test", &redis.Options{
//		Addr:     "localhost:30001",
//		Username: "root",
//		Password: "root",
//		DB:       15,
//	})
//
// ExamplePath:  taurus_go_demo/cache/redis_test.go
func SetClient(clientName string, options *redis.Options) {
	clientOptions[clientName] = options
}

// GetClient 用于从redis连接池中获取一个redis客户端，如果不存在，则返回nil。
//
// Params:
//
//   - clientName: redis客户端的名称，与`SetClient`中的`clientName`相同。
//
// Returns:
//
//	0: redis客户端的指针。
//
// Example:
//
//	redis.SetClient("test", &redis.Options{
//		Addr:     "localhost:30001",
//		Username: "root",
//		Password: "root",
//		DB:       15,
//	})
//
// c := redis.GetClient("test")
//
//	if c == nil {
//		t.Errorf("GetClient() return nil")
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis_test.go
func GetClient(clientName string) (*Client, err) {
	option := clientOptions[clientName]
	if option == nil {
		return nil
	}
	return &Client{redis.NewClient(option)}
}

// ClearClient 用于清空redis客户端的连接池。
//
// Example:
//
//	redis.SetClient("test", &redis.Options{
//		Addr:     "localhost:30001",
//		Username: "root",
//		Password: "root",
//		DB:       15,
//	})
//
// defer redis.ClearClient()
//
// ExamplePath:  taurus_go_demo/cache/redis_test.go
func ClearClient() {
	clientOptions = make(map[string]*Options)
}

func (c *Client) Set(value []string, expiration time.Duration) *Set {
	return &Set{
		client:            c,
		defaultValue:      value,
		defaultExpiration: expiration,
	}
}
