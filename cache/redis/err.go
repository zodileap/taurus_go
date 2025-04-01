package redis

import (
	goRedis "github.com/redis/go-redis/v9"
	"github.com/zodileap/taurus_go/err"
)

func (c *Client) validErr(err error, m *Mutation) error {
	var key string = ""
	if m != nil {
		key = m.Key
	}
	if err == goRedis.Nil {
		return Err_0300010003.Sprintf(c.Name, key)
	} else if err.Error() == "WRONGPASS invalid username-password pair or user is disabled." {
		return Err_0300010002.Sprintf(c.Name)
	} else if err.Error() == "WRONGTYPE Operation against a key holding the wrong kind of value" {
		return Err_0300010004.Sprintf(c.Name, key)
	}
	return Err_030001000x.Sprintf(err)
}

// Err_030001000x 未知的错误。
//
// Verbs:
//
//	0: github.com/redis/go-redis产生的错误
var Err_030001000x err.ErrCode = err.New(
	"030001000x",
	"%s.",
	"",
)

// Err_0300010001 redis客户端的配置不存在。
//
// Verbs:
//
//	0: redis客户端的名称。
var Err_0300010001 err.ErrCode = err.New(
	"0300010001",
	"Redis client %s,options does not exist.",
	"可能得原因是：1.是否使用SetClient方法设置了redis客户端配置。 2.传入的redis客户端名称与SetClient中的不匹配。",
)

// Err_0300010002 redis客户端连接失败，用户名或密码认证失败。
//
// Verbs:
//
//	0: redis客户端的名称。
var Err_0300010002 err.ErrCode = err.New(
	"0300010002",
	"Redis client %s connect failed, username or password authentication failed.",
	"传入的用户名或密码不正确。",
)

// Err_0300010003 因为key不存在，所以没有结果返回。这有时候也不能被当做是错误，比如判断一个key是否存在，
// 就不应该把这个当做错误处理。
//
// Verbs:
//
//	0: redis客户端的名称。
//	1: key的名称。
var Err_0300010003 err.ErrCode = err.New(
	"0300010003",
	"Redis client %s,no result is returned because the key %s does not exist.",
	"",
)

// Err_0300010004 对key执行不适合其值类型的操作。如果一个键是作为列表（list）存储的，
// 而你尝试对它执行针对字符串（string）类型的操作，就会出现这个错误
//
// Verbs:
//
//	0: redis客户端的名称。
//	1: key的名称。
var Err_0300010004 err.ErrCode = err.New(
	"0300010004",
	"Redis client %s,operation against key %s holding the wrong kind of value.",
	"可能的原因：1.传入的值,与key要求的类型不一致。",
)

// Err_0300010005 redis客户端的追踪器类型错误,这个错误不应该出现。出现在内部新建tracker时，
// 对tracker类型进行断言出现错误
var Err_0300010005 err.ErrCode = err.New(
	"0300010005",
	"Redis client %s,tracker type error.",
	"",
)
