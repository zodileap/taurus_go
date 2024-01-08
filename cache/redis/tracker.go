package redis

import (
	"time"

	"github.com/redis/go-redis/v9"
)

type Tracker interface {
	GetName() string
	SetName(name string)
	Exec(*redis.Pipeliner, *[]cmdRes) error
}

type baseTracker struct {
	name string
}

func newBaseTracker(name string) baseTracker {
	return baseTracker{
		name: name,
	}
}

func (b *baseTracker) GetName() string {
	return b.name
}

func (b *baseTracker) SetName(name string) {
	b.name = name
}

type Mutation struct {
	Key string
	Exp time.Duration
}

func NewMutation(key string) Mutation {
	return Mutation{
		Key: key,
		Exp: 0,
	}
}

func (m *Mutation) SetExp(exp time.Duration) {
	m.Exp = exp
}

func (m *Mutation) SetKey(key string) {
	m.Key = key
}

// cmdRes 这是是用于在管线中获取命令执行后的结果。
type cmdRes struct {
	// Key 这是命令的键。
	Key string
	// KeyType 这是键的类型。
	KeyType KeyType
	// Cmd 存入管线中的命令。
	Cmd redis.Cmder
}

func newCmdRes(key string, typ KeyType, cmd redis.Cmder) cmdRes {
	return cmdRes{
		Key:     key,
		KeyType: typ,
		Cmd:     cmd,
	}
}

func (c *cmdRes) classify(r *Res) {
	if c.KeyType == SetType {
		setRes := r.GetSet(c.Key)
		if setRes == nil {
			setRes = NewSetRes(c.Key)
			r.Set = append(r.Set, setRes)
		}
		setRes.encodeCmd(c)
	} else if c.KeyType == HashType {
		hashRes := r.GetHash(c.Key)
		if hashRes == nil {
			hashRes = NewHashRes(c.Key)
			r.Hash = append(r.Hash, hashRes)
		}
		hashRes.encodeCmd(c)
	} else if c.KeyType == StringType {
		stringRes := r.GetString(c.Key)
		if stringRes == nil {
			stringRes = NewStringRes(c.Key)
			r.String = append(r.String, stringRes)
		}
		stringRes.encodeCmd(c)
	}
}

type Res struct {
	Set    []*SetRes
	Hash   []*HashRes
	String []*StringRes
}

func (r *Res) GetSet(key string) *SetRes {
	for _, v := range r.Set {
		if v.Key == key {
			return v
		}
	}
	return nil
}

func (r *Res) GetHash(key string) *HashRes {
	for _, v := range r.Hash {
		if v.Key == key {
			return v
		}
	}
	return nil
}

func (r *Res) GetString(key string) *StringRes {
	for _, v := range r.String {
		if v.Key == key {
			return v
		}
	}
	return nil
}
