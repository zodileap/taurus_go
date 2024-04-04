package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

// Hash 追踪所有对Redis中Hash类型的事务管道操作，单次命令操作不会追踪，会在调用Save方法后执行。
// Redis中Hash类型用于存储键值对的集合
type Hash struct {
	client  *Client
	tracker *HashTracker
}

// Add 在Hash中存入一个字段和它的值（事务管道）。
// HashRes.Oper为HSet。
// 如果增加的字段已经存在，返回的AddNum为0。
//
// Params:
//   - key: Hash类型的键名
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。注意：这是对一整个键的设置。
//   - field: 存入的字段
//   - value: 字段的值
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.Add("hash_key", 0, "filed_1", "value1")
//	h.Add("hash_key", 0, "filed_2", "value2")
//	h.Add("hash_key2", 0, "filed_1", "value1")
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetHash("hash_key")
//	fmt.Println(keyRes.AddNum)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashAdd
func (h *Hash) Add(key string, expiration time.Duration, field string, value string) *Hash {
	h.tracker.Add(key, expiration, field, value)
	return h
}

// AddM 以Map的形式存入多个字段和值（事务管道）。
// HashRes.Oper为HMSet。
//
// Params:
//   - key: Hash类型的键名
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。注意：这是对一整个键的设置。
//   - pairs: 存入的键值对Map
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.AddM("hash_key", 0, map[string]string{"filed_3": "value3", "filed_4": "value4"})
//	h.AddM("hash_key2", 0, map[string]string{"filed_2": "value2"})
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetHash("hash_key")
//	fmt.Println(keyRes.AddNum)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashAddM
func (h *Hash) AddM(key string, expiration time.Duration, pairs map[string]string) *Hash {
	h.tracker.AddM(key, expiration, pairs)
	return h
}

// Get 获取字段的值（事务管道）。
// 获得的值存放在HashRes.Value，HashRes.Oper为HGet。
//
// Params:
//   - key: Hash类型的键名
//   - field: 获取的字段
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.Get("hash_key", "filed_1")
//	h.Get("hash_key", "filed_2")
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//
// keyRes := r.GetHash("hash_key")
// fmt.Println(keyRes.Value)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGet
func (h *Hash) Get(key string, field string) *Hash {
	h.tracker.Get(key, field)
	return h
}

// GetM 获取多个字段的值（事务管道）。
// 获得的值存放在HashRes.Value，HashRes.Oper为HMGet。
//
// Params:
//   - key: Hash类型的键名
//   - fields: 获取的字段
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.GetM("hash_key", []string{"filed_1", "filed_2"})
//	if err != nil {
//		fmt.Println(err)
//	}
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetHash("hash_key")
//	fmt.Println(keyRes.Value)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGetM
func (h *Hash) GetM(key string, fields []string) *Hash {
	h.tracker.GetM(key, fields)
	return h
}

// GetVals 获得Hash中的所有值（事务管道）。
// 获得的值存放在HashRes.Value，HashRes.Oper为HVals。
//
// Params:
//   - key: Hash类型的键名
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.GetVals("hash_key")
//	if err != nil {
//		fmt.Println(err)
//	}
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetHash("hash_key")
//	fmt.Println(keyRes.Value)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGetVals
func (h *Hash) GetVals(key string) *Hash {
	h.tracker.GetVals(key)
	return h
}

// GetAll 获得Hash中的所有字段和值（事务管道）。
// 获得的值存放在HashRes.MapValue，HashRes.Oper为HGetAll。
//
// Params:
//   - key: Hash类型的键名
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.GetAll("hash_key")
//	if err != nil {
//		fmt.Println(err)
//	}
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetHash("hash_key")
//	fmt.Println(keyRes.MapValue)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGetAll
func (h *Hash) GetAll(key string) *Hash {
	h.tracker.GetAll(key)
	return h
}

// Del 删除Hash中的一个字段（事务管道）。
// HashRes.Oper为HDel。与DelAll不同，它不会忽略其他的操作。
//
// Params:
//   - key: Hash类型的键名
//   - fields: 删除的字段
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.Del("hash_key", "filed_1")
//	h.Del("hash_key2", "filed_2")
//	r, err := c.Save()
//
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetHash("hash_key")
//	fmt.Println(keyRes.DelNum)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashDel
func (h *Hash) Del(key string, fields ...string) *Hash {
	h.tracker.Del(key, fields)
	return h
}

// DelAll 删除整个key（事务管道）。
// HashRes.Oper为Del。当添加这个操作后，对这个key的其他操作都会被忽略。
//
// Params:
//   - key: Hash类型的键名
//
// Returns:
//
//	0: 返回一个*Hash对象
//
// Example:
//
//	...获得Hash客户端。
//
//	h.DelAll("hash_key")
//	r, err := c.Save()
//	if err != nil {
//		fmt.Println(err)
//	}
//	keyRes := r.GetHash("hash_key")
//	fmt.Println(keyRes.DelNum)
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashDelAll
func (h *Hash) DelAll(key string) *Hash {
	h.tracker.DelAll(key)
	return h
}

// AddR 在Hash中存入一个字段和它的值（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。注意：这是对一整个键的设置。
//   - field: 存入的字段
//   - value: 字段的值
//
// Returns:
//
//	0: 成功添加的数量，如果成功添加，则返回1，如果已经存在，则返回0。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r1, err := h.AddR("hash_key", 0, "filed_1", "value1")
//	if err != nil {
//		fmt.Println(err)
//	}
//	r2, err := h.AddR("hash_key", 0, "filed_2", "value2")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashAddR
//
// ErrCodes:
//
//   - Err_030001000x
//   - Err_0300010004
func (h *Hash) AddR(key string, expiration time.Duration, field string, value string) (int64, error) {
	length, err := h.client.Client.HSet(ctx, key, field, value).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return 0, err
	}
	h.client.SetExpire(key, expiration)
	return length, nil
}

// AddMR 添加多个多个字段和值（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//   - expiration: 过期时间, 单位为秒, 如果设置为0, 则表示不过期。注意：这是对一整个键的设置。
//   - pairs: 存入的键值对Map
//
// Returns:
//
//	0: 成功添加的数量，如果成功添加，则返回1，如果已经存在，则返回0。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r1, err := h.AddMR("hash_key", 0, map[string]string{"filed_3": "value3", "filed_4": "value4"})
//	if err != nil {
//		fmt.Println(err)
//	}
//	r2, err := h.AddMR("hash_key2", 0, map[string]string{"filed_2": "value2"})
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashAddMR
//
// ErrCodes:
//
//   - Err_030001000x
//   - Err_0300010004
func (h *Hash) AddMR(key string, expiration time.Duration, pairs map[string]string) (int64, error) {
	length, err := h.client.Client.HSet(ctx, key, pairs).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return 0, err
	}
	h.client.SetExpire(key, expiration)
	return length, nil
}

// GetR 获取一个set类型的键值对（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//   - field: 获取的字段
//
// Returns:
//
//	0: 返回一个string，如果键不存在，则返回一个空的string。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r1, err := h.GetR("hash_key", "filed_1")
//	if err != nil {
//		fmt.Println(err)
//	}
//	r2, err := h.GetR("hash_key", "filed_2")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGetR
//
// ErrCodes:
//
//   - Err_030001000x
func (h *Hash) GetR(key string, field string) (string, error) {
	val, err := h.client.Client.HGet(ctx, key, field).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return "", err
	}
	return val, nil
}

// GetMR 获取多个字段的值（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//   - fields: 需要获取值的字段
//
// Returns:
//
//	0: 返回查询的字段值的切片，如果查询的字段不存在，则该字段对应的值为nil。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r, err := h.GetMR("hash_key", "filed_1", "filed_5")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGetMR
//
// ErrCodes:
//
//   - Err_030001000x
func (h *Hash) GetMR(key string, fields ...string) ([]interface{}, error) {
	val, err := h.client.Client.HMGet(ctx, key, fields...).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return nil, err
	}
	return val, nil
}

// GetValsR 获得Hash中的所有值（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//
// Returns:
//
//	0: 返回一个值切片，如果键不存在，则返回一个空的[]string。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r, err := h.GetValsR("hash_key")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGetValsR
//
// ErrCodes:
//
//   - Err_030001000x
func (h *Hash) GetValsR(key string) ([]string, error) {
	val, err := h.client.Client.HVals(ctx, key).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return nil, err
	}
	return val, nil
}

// GetAllR 获得Hash中的所有字段和值（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//
// Returns:
//
//	0: 返回全部的字段和值的对，如果键不存在，则返回一个空的map[string]string。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r, err := h.GetAllR("hash_key")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashGetAllR
//
// ErrCodes:
//
//   - Err_030001000x
func (h *Hash) GetAllR(key string) (map[string]string, error) {
	val, err := h.client.Client.HGetAll(ctx, key).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return nil, err
	}
	return val, nil
}

// DelR 删除Hash中的一个字段（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//   - field: 删除的字段
//
// Returns:
//
//	0: 成功删除的字段数量，如果成功删除，返回1，否则返回0。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r, err := h.DelR("hash_key", "filed_1")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashDelR
//
// ErrCodes:
//
//   - Err_030001000x
func (h *Hash) DelR(key string, field string) (int64, error) {
	length, err := h.client.Client.HDel(ctx, key, field).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return 0, err
	}
	return length, nil
}

// DelAllR 删除整个key（单次命令操作）。
//
// Params:
//   - key: Hash类型的键名
//
// Returns:
//
//	0: 成功删除的键数量，如果成功删除，返回1，否则返回0。
//	1: 错误信息。
//
// Example:
//
//	...获得Hash客户端。
//
//	r, err := h.DelAllR("hash_key")
//	if err != nil {
//		fmt.Println(err)
//	}
//
// ExamplePath:  taurus_go_demo/cache/redis/hash_test.go - TestHashDelAllR
//
// ErrCodes:
//
//   - Err_030001000x
func (h *Hash) DelAllR(key string) (int64, error) {
	length, err := h.client.Client.Del(ctx, key).Result()
	if err != nil {
		m := NewMutation(key)
		err := h.client.validErr(err, &m)
		return 0, err
	}
	return length, nil
}

type (
	HashTracker struct {
		tracking
		mutation []*HashMutation
	}

	HashMutation struct {
		Mutation
		Ops   []HashOp
		Del   bool
		AddOp HashOp
	}

	HashOp struct {
		Op   string
		MVal map[string]string
		Val  Pairs
	}

	Pairs struct {
		Field string
		Value string
	}
)

func NewHashTracker(name string) *HashTracker {
	return &HashTracker{
		tracking: newTracking(name),
		mutation: make([]*HashMutation, 0),
	}
}

func (h *HashTracker) Exec(p *redis.Pipeliner, r *[]cmdRes) error {
	_p := *p
	for _, m := range h.mutation {
		if m.Del {
			nr := newCmdRes(m.Key, HashType, _p.Del(ctx, m.Key))
			*r = append(*r, nr)
			continue
		} else {
			for _, op := range m.Ops {
				if op.Op == HSet {
					nr := newCmdRes(m.Key, HashType, _p.HSet(ctx, m.Key, op.Val.Field, op.Val.Value))
					*r = append(*r, nr)
				} else if op.Op == HMSet {
					nr := newCmdRes(m.Key, HashType, _p.HSet(ctx, m.Key, op.MVal))
					*r = append(*r, nr)
				} else if op.Op == HDel {
					fields := []string{}
					for f, _ := range op.MVal {
						fields = append(fields, f)
					}
					nr := newCmdRes(m.Key, HashType, _p.HDel(ctx, m.Key, fields...))
					*r = append(*r, nr)
				}
			}
			if m.AddOp.Op != "" {
				if m.AddOp.Op == HGet {
					nr := newCmdRes(m.Key, HashType, _p.HGet(ctx, m.Key, m.AddOp.Val.Field))
					*r = append(*r, nr)
				} else if m.AddOp.Op == HMGet {
					fields := []string{}
					for f, _ := range m.AddOp.MVal {
						fields = append(fields, f)
					}
					nr := newCmdRes(m.Key, HashType, _p.HMGet(ctx, m.Key, fields...))
					*r = append(*r, nr)
				} else if m.AddOp.Op == HVals {
					nr := newCmdRes(m.Key, HashType, _p.HVals(ctx, m.Key))
					*r = append(*r, nr)
				} else if m.AddOp.Op == HGetAll {
					nr := newCmdRes(m.Key, HashType, _p.HGetAll(ctx, m.Key))
					*r = append(*r, nr)
				}
			}
			pipeSetExpire(_p, m.Key, m.Exp)
		}

	}
	return nil
}

func (h *HashTracker) FindKey(key string) *HashMutation {
	for _, m := range h.mutation {
		if m.Key == key {
			return m
		}
	}
	m := &HashMutation{
		Mutation: NewMutation(key),
		Ops:      make([]HashOp, 0),
		Del:      false,
		AddOp:    HashOp{},
	}
	h.mutation = append(h.mutation, m)
	return m
}

func (h *HashTracker) Add(key string, expiration time.Duration, field string, value string) error {
	m := h.FindKey(key)
	m.Ops = append(m.Ops, HashOp{
		Op: HSet,
		Val: Pairs{
			Field: field,
			Value: value,
		},
	})
	m.Exp = expiration
	return nil
}

func (h *HashTracker) AddM(key string, expiration time.Duration, pairs map[string]string) error {
	m := h.FindKey(key)
	m.Ops = append(m.Ops, HashOp{
		Op:   HMSet,
		MVal: pairs,
	})
	m.Exp = expiration
	return nil
}

func (h *HashTracker) Get(key string, field string) error {
	m := h.FindKey(key)
	m.AddOp = HashOp{
		Op: HGet,
		Val: Pairs{
			Field: field,
		},
	}
	return nil
}

func (h *HashTracker) GetM(key string, fields []string) error {
	m := h.FindKey(key)
	o := HashOp{
		Op:   HMGet,
		MVal: make(map[string]string, 0),
	}
	for _, f := range fields {
		o.MVal[f] = ""
	}
	m.AddOp = o
	return nil
}

func (h *HashTracker) GetVals(key string) error {
	m := h.FindKey(key)
	m.AddOp = HashOp{
		Op: HVals,
	}
	return nil
}

func (h *HashTracker) GetAll(key string) error {
	m := h.FindKey(key)
	m.AddOp = HashOp{
		Op: HGetAll,
	}
	return nil
}

func (h *HashTracker) Del(key string, field []string) error {
	m := h.FindKey(key)
	o := HashOp{
		Op:   HDel,
		MVal: make(map[string]string, 0),
	}
	for _, f := range field {
		o.MVal[f] = ""
	}
	m.Ops = append(m.Ops, o)
	return nil
}

func (h *HashTracker) DelAll(key string) error {
	m := h.FindKey(key)
	m.Del = true
	return nil
}

// HashRes Hash类型操作的结果。
// 同一个键只会有一个HashRes。
type HashRes struct {
	// Key 键名
	Key string
	// Value 值，`Get`,`GetM`,`GetVals`的操作结果的值。
	Value []string
	// MapValue 键和值的map结构，它一般用于存储GetAll()的操作结果的值。
	MapValue map[string]string
	// AddNum 成功添加的数量
	AddNum int64
	// DelNum 成功删除的数量
	DelNum int64
	// Oper 最后一次操作的命令
	Oper string
}

func NewHashRes(key string) *HashRes {
	return &HashRes{
		Key:      key,
		Value:    make([]string, 0),
		MapValue: make(map[string]string),
	}
}

func (h *HashRes) encodeCmd(c *cmdRes) {
	switch t := c.Cmd.(type) {
	case *redis.IntCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case HSet:
			h.AddNum += v
			h.Oper = HSet
		case HDel:
			h.DelNum += v
			h.Oper = HDel
		case Del:
			h.DelNum += v
			h.Oper = Del
		}
	case *redis.StringCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case HGet:
			h.Value = append(h.Value, v)
			h.Oper = HGet
		}
	case *redis.StringSliceCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case HVals:
			h.Value = v
			h.Oper = HVals
		}
	case *redis.SliceCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case HMGet:
			stringV := make([]string, 0)
			for _, s := range v {
				if s == nil {
					stringV = append(stringV, "")
				} else {
					if _, ok := s.(string); !ok {
						stringV = append(stringV, "")
					} else {
						stringV = append(stringV, s.(string))
					}
				}
			}
			h.Value = stringV
			h.Oper = HMGet
		}
	case *redis.MapStringStringCmd:
		v := t.Val()
		n := t.Name()
		switch n {
		case HGetAll:
			for k, v := range v {
				h.MapValue[k] = v
			}
			h.Oper = HGetAll
		}
	}
}
