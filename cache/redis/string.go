package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// String 追踪所有对Redis中String类型的事务管道操作，单次命令操作不会追踪，会在调用Save方法后执行。
// Redis中String类型用于存储单个字符串。
type String struct {
	client  *Client
	tracker *StringTracker
}

// Add 在String中增加一个值（事务管道）。
// StringRes.Opert为Set。
//
// Params:
//   - key: String类型的键名
//   - value: 键的值
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。注意：这是对一整个键的设置。
//
// Returns:
//
//	0: 返回一个*String对象
//
// Example:
//
//	...获得String客户端。
//
//	s.Add("string_key", 0, "value")
//	r, err := c.Save()
//
// ExamplePath:  taurus_go_demo/cache/redis/string_test.go - TestStringAdd
func (s *String) Add(key string, expiration time.Duration, value string) *String {
	s.tracker.Add(key, expiration, value)
	return s
}

// Get 获取一个String类型的键值对（事务管道）。
// StringRes.Opert为Get。
// 在调用Save方法后，可以通过StringRes中的Value，获取键的值。
//
// Params:
//   - key: String类型的键名
//
// Returns:
//
//	0: 返回一个*String对象
//
// Example:
//
//	...获得String客户端。
//
//	s.Get("string_key")
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetString("string_key")
//	fmt.Println(keyRes.Value)
//
// ExamplePath:  taurus_go_demo/cache/redis/string_test.go - TestStringGet
func (s *String) Get(key string) *String {
	s.tracker.Get(key)
	return s
}

// Del 删除一个String类型的键值对（事务管道）。
// StringRes.Opert为Del。
//
// Params:
//   - key: String类型的键名
//
// Returns:
//
//	0: 返回一个*String对象
//
// Example:
//
//	...获得String客户端。
//
//	s.Del("string_key")
//	_, err = c.Save()
//
// ExamplePath:  taurus_go_demo/cache/redis/string_test.go - TestStringDel
func (s *String) Del(key string) *String {
	s.tracker.Del(key)
	return s
}

// AddR 添加一个字符串（单次命令操作）。
//
// Params:
//   - key: String类型的键名
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。
//   - value: 键的值
//
// Returns:
//
//	0: 错误信息。
//
// Example:
//
//	...获得String客户端。
//
//	err = s.AddR("string_key", 0, "value")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/string_test.go - TestStringAddR
//
// ErrCodes:
//
//   - Err_030001000x
func (s *String) AddR(key string, expiration time.Duration, value string) error {
	err := s.client.Client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		m := NewMutation(key)
		err := s.client.validErr(err, &m)
		return err
	}
	return nil
}

// GetR 获取一个字符串（单次命令操作）。
//
// Params:
//   - key: String类型的键名
//
// Returns:
//
//	0: 返回一个String类型的键值对。
//	1: 错误信息。
//
// Example:
//
//	...获得String客户端。
//
//	r, err := s.GetR("string_key")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(r)
//
// ExamplePath:  taurus_go_demo/cache/redis/string_test.go - TestStringGetR
//
// ErrCodes:
//
//   - Err_030001000x
func (s *String) GetR(key string) (string, error) {
	val, err := s.client.Client.Get(ctx, key).Result()
	if err != nil {
		m := NewMutation(key)
		err := s.client.validErr(err, &m)
		return "", err
	}
	return val, nil
}

// DelR 删除整个键（单次命令操作）。
//
// Params:
//   - key: String类型的键名
//
// Returns:
//
//	0: 错误信息。
//
// Example:
//
//	...获得String客户端。
//
//	err = s.DelR("string_key")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/string_test.go - TestStringDelR
//
// ErrCodes:
//
//   - Err_030001000x
func (s *String) DelR(key string) error {
	err := s.client.Client.Del(ctx, key).Err()
	if err != nil {
		m := NewMutation(key)
		err := s.client.validErr(err, &m)
		return err
	}
	return nil
}

type (
	StringTracker struct {
		baseTracker
		mutation []*StringMutation
	}

	StringMutation struct {
		Mutation
		Ops []StringOp
		Del bool
	}

	StringOp struct {
		Op    string
		Value string
	}
)

func NewStringTracker(name string) *StringTracker {
	b := newBaseTracker(name)
	return &StringTracker{
		baseTracker: b,
		mutation:    make([]*StringMutation, 0),
	}
}

func (s *StringTracker) Exec(p *redis.Pipeliner, r *[]cmdRes) error {
	_p := *p
	for _, m := range s.mutation {
		if m.Del {
			nr := newCmdRes(m.Key, StringType, _p.Del(ctx, m.Key))
			*r = append(*r, nr)
		} else {
			for _, op := range m.Ops {
				if op.Op == Add {
					nr := newCmdRes(m.Key, StringType, _p.Set(ctx, m.Key, op.Value, m.Exp))
					*r = append(*r, nr)
				} else if op.Op == Get {
					nr := newCmdRes(m.Key, StringType, _p.Get(ctx, m.Key))
					*r = append(*r, nr)
				}
			}
			pipeSetExpire(_p, m.Key, m.Exp)
		}
	}
	return nil
}

func (s *StringTracker) FindKey(key string) *StringMutation {
	for _, m := range s.mutation {
		if m.Key == key {
			return m
		}
	}
	m := &StringMutation{
		Mutation: NewMutation(key),
		Ops:      make([]StringOp, 0),
		Del:      false,
	}
	s.mutation = append(s.mutation, m)
	return m
}

func (s *StringTracker) Add(key string, expiration time.Duration, value string) *StringTracker {
	m := s.FindKey(key)
	m.Ops = append(m.Ops, StringOp{
		Op:    Add,
		Value: value,
	})
	m.Exp = expiration
	return s
}

func (s *StringTracker) Get(key string) *StringTracker {
	m := s.FindKey(key)
	m.Ops = append(m.Ops, StringOp{
		Op: Get,
	})
	return s
}

func (s *StringTracker) Del(key string) *StringTracker {
	m := s.FindKey(key)
	m.Del = true
	return s
}

// StringRes String类型的键值对。
// 同一个键只会有一个StringRes。
type StringRes struct {
	// Key 键名
	Key string
	// Value 存入的字符串
	Value string
	// Oper 最后一次操作的命令
	Oper string
}

func NewStringRes(key string) *StringRes {
	return &StringRes{
		Key: key,
	}
}

func (s *StringRes) encodeCmd(c *cmdRes) {
	switch t := c.Cmd.(type) {
	case *redis.StringCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case Get:
			s.Value = v
			s.Oper = Get
		}
	case *redis.IntCmd:
		n := t.Name()
		switch n {
		case Del:
			s.Oper = Del
		}
	case *redis.StatusCmd:
		n := t.Name()
		switch n {
		case Add:
			s.Oper = Add
		}
	}
}
