package redis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type (
	// Options redis客户端的配置,[go-redis配置项](https://redis.uptrace.dev/zh/guide/go-redis-option.html#redis-client)。
	Options = redis.Options
	// Client 调用`GetClient`后生成的redis客户端，用于执行redis操作。
	Client struct {
		// Name 客户端名称。
		Name string
		// Client go-redis客户端。
		*redis.Client
		// nodes redis客户端的操作队列。
		nodes []Tracker
	}
	// KeyType 键的类型,目前只有Set,String,Hash。
	KeyType = int16
)

const (
	// SetType Set类型的键。
	SetType KeyType = 1
	// StringType String类型的键。
	StringType KeyType = 2
	// HashType Hash类型的键。
	HashType KeyType = 3
)

const (
	// Del 所有类型的键的删除操作，会删除整个key，在如果对某个键添加了这个操作，则其他全部的操作都会忽略。
	Del string = "del"
	// Add String类型键的添加操作。
	Add string = "set"
	// Get String类型键的获取操作。
	Get string = "get"
	// SAdd Set类型键的添加操作。
	SAdd string = "sadd"
	// SMembers Set类型键的获取操作。
	SMembers string = "smembers"
	// SRem Set类型键的删除操作。
	SRem string = "srem"
	// HSet Hash类型键添加单个键值对的操作。
	HSet string = "hset"
	// HMSet Hash类型键添加多个键值对的操作，但是最后还是会调用HSet操作，所以输出会的Oper是HSet,
	// 只是为了区分Add,和AddM方法
	HMSet string = "hmset"
	// HGet Hash类型键获取单个键值对的操作。
	HGet string = "hget"
	// HMGet Hash类型键获取多个键值对的操作。
	HMGet string = "hmget"
	// HVals Hash类型键获取所有值的操作。
	HVals string = "hvals"
	// HGetAll Hash类型键获取所有键值对的操作。
	HGetAll string = "hgetall"
	// HDel Hash类型键删除单个键值对的操作。
	HDel string = "hdel"
)

var (
	// ctx redis客户端的上下文。
	ctx context.Context = context.Background()

	// pool redis客户端的连接池。
	clientOptions map[string]*Options = make(map[string]*Options)
)

// SetClient 用于添加一个redis客户端的配置。这个操作不会立即连接redis，只有在调用`GetClient`时才会连接。
//
// Params:
//
//   - clientName: redis客户端的名称。
//   - options: redis客户端的配置。
//
// Example:
//
//	redis.SetClient("test", &redis.Options{
//		Addr: "localhost:6379",
//		Username: "",
//		Password: "",
//		DB:       1,
//	})
//
// ExamplePath:  taurus_go_demo/cache/redis/client_test.go - TestSetClient
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
//	0: redis客户端。
//	1: 错误信息。
//
// Example:
//
//	redis.SetClient("test", &redis.Options{
//		Addr:     "localhost:6379",
//		Username: "",
//		Password: "",
//		DB:       1,
//	})
//	// highlight-start
//	c, err := redis.GetClient("test")
//	defer c.Close()
//	// highlight-end
//	if err != nil {
//		fmt.Print(err.Error())
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/client_test.go - TestGetClient
//
// ErrCodes:
//
//   - Err_0300010001
func GetClient(clientName string) (*Client, error) {
	option := clientOptions[clientName]
	if option == nil {
		return nil, Err_0300010001
	}
	return &Client{
		Name:   clientName,
		Client: redis.NewClient(option),
	}, nil
}

// ClearClient 用于清空redis客户端的连接池。
//
// Example:
//
//	redis.SetClient("test", &redis.Options{
//		Addr:     "localhost:30001",
//		Username: "",
//		Password: "",
//		DB:       15,
//	})
//	defer redis.ClearClient()
//
// ExamplePath:  taurus_go_demo/cache/redis/client_test.go - TestClearClient
func ClearClient() {
	clientOptions = make(map[string]*Options)
}

// Close 用于关闭redis客户端。每次调用`GetClient`时，都会生成一个新的redis客户端，
// 在使用完后，需要调用`Close`方法关闭。
//
// Example:
//
//	redis.SetClient("test", &redis.Options{
//		Addr:     "localhost:6379",
//		Username: "",
//		Password: "",
//		DB:       1,
//	})
//
//	c, err := redis.GetClient("test")
//	defer c.Close()
//
// ExamplePath:  taurus_go_demo/cache/redis/client_test.go - TestClose
//
// ErrCodes:
//
//   - Err_030001000x
func (c *Client) Close() error {
	err := c.Client.Close()
	if err != nil {
		return Err_030001000x.Sprintf(err)
	}
	return nil
}

// Del 用于删除一个键值对。
//
// Params:
//
//   - keys: 需要删除的键。
//
// Returns:
//
//	0: 删除的键的数量。
//	1: 错误信息。
//
// Example:
//
//	...初始化redis客户端
//	l, err := c.Del("key1", "key2")
//	if err != nil {
//		t.Errorf(err.Error())
//	}
//	fmt.Println(l)
//
// ExamplePath:  taurus_go_demo/cache/redis/client_test.go - TestDel
//
// ErrCodes:
//
//   - Err_030001000x
func (c *Client) Del(keys ...string) (int64, error) {
	l, err := c.Client.Del(ctx, keys...).Result()
	if err != nil {
		return 0, Err_030001000x.Sprintf(err)
	}
	return l, nil
}

// Save 用于保存先前调用的redis操作，通过事务管道，一次性执行。
// 在调用`Save`后，如果还需要执行redis操作，请重新调用`Set`、`Hash`、`String`等方法。
//
// Returns:
//
//	0: 执行的结果。
//	1: 错误信息。
//
// Example:
//
//	...初始化redis客户端
//	s, err := c.Set()
//	if err != nil {
//		mt.Println(err.Error())
//	}
//	s.Add("key", 0, "value")
//	s.Add("key2", 0, "value2")
//	s.Add("key", 0, "value3")
//	result, err := c.Save()
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//
//	// 请注意每次调用`Save`后，都需要重新调用`Set`、`Hash`、`String`等方法。
//	s, err = c.Set()
//	s.Add("key", 0, "value3")
//	c.Save()
//
// ExamplePath:  taurus_go_demo/cache/redis/client_test.go - TestSave
//
// ErrCodes:
//
//   - Err_030001000x
//   - Err_0300010002
//   - Err_0300010004
func (c *Client) Save() (*Res, error) {
	pipe := c.Pipeline()
	nr := make([]cmdRes, 0)
	var crs *[]cmdRes = &nr
	for _, node := range c.nodes {
		err := node.Exec(&pipe, crs)
		if err != nil {
			return nil, err
		}
	}
	_, err := pipe.Exec(ctx)
	if err != nil {
		return nil, c.validErr(err, nil)
	}
	res := &Res{}
	for _, cr := range *crs {
		cr.classify(res)

	}
	c.nodes = make([]Tracker, 0)
	return res, nil
}

// SetExpire 用于设置一个键值对的过期时间。最小是1秒，如果为0，则不会设置过期时间。
//
// Params:
//
//   - key: 需要设置过期时间的键。
//   - expiration: 过期时间。
//
// Example:
//
//	...初始化redis客户端
//	c.SetExpire("key", 10*time.Second)
//
// ExamplePath:  taurus_go_demo/cache/redis/client_test.go - TestSetExpire
func (c *Client) SetExpire(key string, expiration time.Duration) {
	if expiration > 0 {
		c.Expire(ctx, key, expiration)
	}
}

// Set 用于开始Set类型的操作
//
// Returns:
//
//	0: Set类型的Client。
//	1: 错误信息。
//
// ErrCodes:
//
//   - Err_0300010005
func (c *Client) Set() (*Set, error) {
	tracker := c.findTracker("Set")
	if tracker == nil {
		tracker = newSetTracker("Set")
		c.nodes = append(c.nodes, tracker)
	}
	st, ok := tracker.(*setTracker)
	if !ok {
		return nil, Err_0300010005.Sprintf(c.Name)
	} else {
		return &Set{
			client:  c,
			tracker: st,
		}, nil
	}
}

// Hash 用于开始Hash类型的操作
//
// Returns:
//
//	0: Hash类型的Client。
//	1: 错误信息。
//
// ErrCodes:
//
//   - Err_0300010005
func (c *Client) Hash() (*Hash, error) {
	tracker := c.findTracker("Hash")
	if tracker == nil {
		tracker = NewHashTracker("Hash")
		c.nodes = append(c.nodes, tracker)
	}
	ht, ok := tracker.(*HashTracker)
	if !ok {
		return nil, Err_0300010005.Sprintf(c.Name)
	} else {
		return &Hash{
			client:  c,
			tracker: ht,
		}, nil
	}
}

// String 用于开始String类型的操作
//
// Returns:
//
//	0: String类型的Client。
//	1: 错误信息。
//
// ErrCodes:
//
//   - Err_0300010005
func (c *Client) String() (*String, error) {
	tracker := c.findTracker("String")
	if tracker == nil {
		tracker = NewStringTracker("String")
		c.nodes = append(c.nodes, tracker)
	}
	st, ok := tracker.(*StringTracker)
	if !ok {
		return nil, Err_0300010005.Sprintf(c.Name)
	} else {
		return &String{
			client:  c,
			tracker: st,
		}, nil
	}
}

// findTracker 用于查找redis客户端的操作队列中是否存在某个操作。
//
// Params:
//
//   - name: 操作的名称。
func (c *Client) findTracker(name string) Tracker {
	for _, node := range c.nodes {
		if node.GetName() == name {
			return node
		}
	}
	return nil
}

// classify 用于将redis操作的结果分类。
//
// Params:
//
//   - res: redis操作的结果。
func pipeSetExpire(pipe redis.Pipeliner, key string, expiration time.Duration) {
	if expiration > 0 {
		pipe.Expire(ctx, key, expiration)
	}
}
