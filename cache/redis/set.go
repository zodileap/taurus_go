package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
	goRedis "github.com/redis/go-redis/v9"
)

// Set 追踪所有对Redis中Set类型的事务管道操作，单次命令操作不会追踪，会在调用Save方法后执行。
// Redis中Set类型用于存储字符串的无序集合
type Set struct {
	client  *Client
	tracker *setTracker
}

// Add 在Set中增加一个值（事务管道）。
// SetRes.Opert为SAdd。
// 如果增加的值已经存在，则不会增加，返回的AddNum为0。
//
// Params:
//   - key: Set类型的键名
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。注意：这是对一整个键的设置。
//   - value: 键的值
//
// Returns:
//
//	0: 返回一个*Set对象
//
// Example:
//
//	...获得Set客户端。
//
//	s.Add("key", 0, "value")
//	s.Add("key2", 0, "value2")
//	c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetSet("key")
//	fmt.Println(keyRes.AddNum)
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetAdd
func (s *Set) Add(key string, expiration time.Duration, value string) *Set {
	s.tracker.Add(key, expiration, value)
	return s
}

// Get 获取一个set类型的键值对（事务管道）。
// SetRes.Opert为SMembers。
// 在调用Save方法后，可以通过SetRes中的Value，获取键的值。
//
// Params:
//   - key: Set类型的键名
//
// Returns:
//
//	0: 返回一个*Set对象
//
// Example:
//
//	...获得Set客户端。
//
//	s.Get("key")
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetSet("key")
//	fmt.Println(keyRes.Value)
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetGet
func (s *Set) Get(key string) *Set {
	s.tracker.Get(key)
	return s
}

// Del 删除一个set类型的键值对（事务管道）。
// SetRes.Opert为SRem。
// 在调用Save方法后，可以通过SetRes中的DelNum，删除的值的数量。
//
// Params:
//   - key: Set类型的键名
//   - value: 键的值
//
// Returns:
//
//	0: 返回一个*Set对象
//
// Example:
//
//	...获得Set客户端。
//
//	s.Del("key", "value")
//	s.Del("key2", "value2")
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetSet("key")
//	fmt.Println(keyRes.DelNum)
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetDel
func (s *Set) Del(key string, value string) *Set {
	s.tracker.Del(key, value)
	return s
}

// DelAll 删除一个set类型的键（事务管道）。
// SetRes.Opert为Del。
// 在调用Save方法后，可以通过SetRes中的DelNum，查看删除的键的数量。如果删除的键不存在，则DelNum为0，否则是1。
// 请注意：这里的DelNum代表的删除的键的数量，而Del方法中的DelNum代表的是删除的值的数量。
//
// Params:
//   - key: Set类型的键名
//
// Returns:
//
//	0: 返回一个*Set对象
//
// Example:
//
//	...获得Set客户端。
//
//	s, err := c.Set()
//	if err != nil {
//		fmt.Println(err)
//	}
//	s.DelAll("key")
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetSet("key")
//	fmt.Println(keyRes.DelNum)
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetDelAll
func (s *Set) DelAll(key string) *Set {
	s.tracker.DelAll(key)
	return s
}

// AddR 在一个set中添加一个字符串（单次命令操作）。
//
// Params:
//   - key: Set类型的键名
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。注意：这是对一整个键的设置。
//   - value: 键的值
//
// Returns:
//
//	0: 成功添加的数量，如果成功添加，则返回1，如果已经存在，则返回0。
//	1: 错误信息。
//
// Example:
//
//	...获得Set客户端。
//
//	r, err := s.AddR("key", 0, "value")
//	if err != nil {
//		fmt.Println(err)
//	}
//	r, err = s.AddR("key4", 20*time.Second, "value4")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetAddR
//
// ErrCodes:
//
//   - Err_030001000x
//   - Err_0300010002
//   - Err_0300010004
func (s *Set) AddR(key string, expiration time.Duration, value string) (int64, error) {
	length, err := s.client.SAdd(ctx, key, value).Result()
	if err != nil {
		m := NewMutation(key)
		err := s.client.validErr(err, &m)
		return 0, err
	}
	s.client.SetExpire(key, expiration)
	return length, nil
}

// GetR 获取一个set的值（单次命令操作）。
//
// Params:
//   - key: Set类型的键名
//
// Returns:
//
//	0: 返回一个[]string，如果键不存在，则返回一个空的[]string。
//	1: 错误信息。
//
// Example:
//
//	...获得Set客户端。
//
//	r, err := s.GetR("key3")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(r)
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetGetR
//
// ErrCodes:
//
//   - Err_030001000x
func (s *Set) GetR(key string) ([]string, error) {
	val, err := s.client.SMembers(ctx, key).Result()
	if err != nil {
		m := NewMutation(key)
		err := s.client.validErr(err, &m)
		return nil, err
	}
	return val, nil
}

// DelR 删除set中的一个值（单次命令操作）。
//
// Params:
//   - key: Set类型的键名
//   - value: 键的值
//
// Returns:
//
//	0: 成功删除的数量，如果成功删除，返回1，否则返回0。
//	1: 错误信息。
//
// Example:
//
//	...获得Set客户端。
//
//	r, err := s.DelR("key3", "value3")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(r)
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetDelR
//
// ErrCodes:
//
//   - Err_030001000x
func (s *Set) DelR(key string, value string) (int64, error) {
	length, err := s.client.SRem(ctx, key, value).Result()
	if err != nil {
		m := NewMutation(key)
		err := s.client.validErr(err, &m)
		return 0, err
	}
	return length, nil
}

// DelAllR 删除一个set类型的键（单次命令操作）。
//
// Params:
//   - key: Set类型的键名
//
// Returns:
//
//	0: 成功删除的数量，如果成功删除，返回1，否则返回0。
//	1: 错误信息。
//
// Example:
//
//	...获得Set客户端。
//
//	r, err := s.DelAllR("key")
//	if err != nil {
//		fmt.Println(err)
//	}
//	fmt.Println(r)
//
// ExamplePath:  taurus_go_demo/cache/redis/set_test.go - TestSetDelAllR
//
// ErrCodes:
//
//   - Err_030001000x
func (s *Set) DelAllR(key string) (int64, error) {
	length, err := s.client.Client.Del(ctx, key).Result()
	if err != nil {
		m := NewMutation(key)
		err := s.client.validErr(err, &m)
		return 0, err
	}
	return length, nil
}

type (
	// setTracker 用于追踪所有对set类型的操作。
	setTracker struct {
		tracking
		mutation []*setMutation
	}

	// setMutation 用于记录一个set类型的键的操作
	setMutation struct {
		Mutation
		Ops []setOp
		Del bool
	}

	// setOp set类型的操作
	setOp struct {
		Op  string
		Val string
	}
)

// newSetTracker 新建一个setTracker
func newSetTracker(name string) *setTracker {
	b := newTracking(name)
	return &setTracker{
		tracking: b,
		mutation: make([]*setMutation, 0),
	}
}

func (s *setTracker) Exec(p *redis.Pipeliner, r *[]cmdRes) error {
	_p := *p
	for _, m := range s.mutation {
		if m.Del {
			nr := newCmdRes(m.Key, SetType, _p.Del(ctx, m.Key))
			*r = append(*r, nr)
		} else {
			for _, op := range m.Ops {
				if op.Op == SAdd {
					nr := newCmdRes(m.Key, SetType, _p.SAdd(ctx, m.Key, op.Val))
					*r = append(*r, nr)
				}
				if op.Op == SRem {
					nr := newCmdRes(m.Key, SetType, _p.SRem(ctx, m.Key, op.Val))
					*r = append(*r, nr)
				}
				if op.Op == SMembers {
					nr := newCmdRes(m.Key, SetType, _p.SMembers(ctx, m.Key))
					*r = append(*r, nr)
				}

			}
			pipeSetExpire(_p, m.Key, m.Exp)
		}

	}
	return nil
}

func (s *setTracker) FindKey(key string) *setMutation {
	for _, m := range s.mutation {
		if m.Key == key {
			return m
		}
	}
	m := &setMutation{
		Mutation: NewMutation(key),
		Ops:      make([]setOp, 0),
	}
	s.mutation = append(s.mutation, m)
	return m
}

func (s *setTracker) Add(key string, expiration time.Duration, value string) error {
	m := s.FindKey(key)
	m.Ops = append(m.Ops, setOp{
		Op:  SAdd,
		Val: value,
	})
	m.Exp = expiration
	return nil
}

func (s *setTracker) Get(key string) error {
	m := s.FindKey(key)
	m.Ops = append(m.Ops, setOp{
		Op: SMembers,
	})
	return nil
}

func (s *setTracker) Del(key string, value string) error {
	m := s.FindKey(key)
	m.Ops = append(m.Ops, setOp{
		Op:  SRem,
		Val: value,
	})
	return nil
}

func (s *setTracker) DelAll(key string) error {
	m := s.FindKey(key)
	m.Del = true
	return nil
}

// SetRes set类型操作的结果。
// 同一个键只会有一个SetRes。
type SetRes struct {
	// Key 键的名字
	Key string
	// Value 键的值，只在Get命令中有值，只会是最后一次Get的值。
	Value []string
	// AddNum 成功添加的数量，如果一次事务中多次执行Add操作，会累加。
	// 但是执行了DelAll,不管调用多少次Add，AddNum都为0，因为Add操作被忽略。
	AddNum int64
	// DelNum 成功删除的数量，如果一次事务中多次执行Del操作，会累加。
	// 但是执行了DelAll,不管调用多少次Del，DelNum都为1，因为Del操作被忽略。
	DelNum int64
	// Oper 最后一次操作的命令
	Oper string
}

func NewSetRes(key string) *SetRes {
	return &SetRes{
		Key:   key,
		Value: make([]string, 0),
	}
}

func (s *SetRes) encodeCmd(c *cmdRes) {
	switch t := c.Cmd.(type) {
	case *goRedis.IntCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case SAdd:
			s.AddNum += v
			s.Oper = SAdd
		case SRem:
			s.DelNum += v
			s.Oper = SRem
		case Del:
			s.DelNum += v
			s.Oper = Del
		}

	case *goRedis.StringSliceCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case SMembers:
			s.Value = v
			s.Oper = SMembers
		}
	}
}
