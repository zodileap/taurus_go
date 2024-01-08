package redis

import (
	"fmt"

	goRedis "github.com/redis/go-redis/v9"
	"github.com/yohobala/taurus_go/err"
)

func validErr(err error) error {
	if err == goRedis.Nil {
		return Err_not_key
	} else if err.Error() == "WRONGPASS invalid username-password pair or user is disabled." {
		return Err_not_auth_failed
	} else if err.Error() == "WRONGTYPE Operation against a key holding the wrong kind of value" {
		return Err_key_wrongtype_operation
	}
	return Err_030001000x.Sprintf(err)
}

var Err_nil_options error = fmt.Errorf("没有设置redis客户端配置")

var Err_not_key error = fmt.Errorf("没有该key存在")

var Err_not_auth_failed error = fmt.Errorf("用户名和密码认证失败")

var Err_key_wrongtype_operation error = fmt.Errorf("对key的类型,与想要进行的类型操作不一致")

// Err_030001000x 未知的错误。
//
// Verbs:
//
//	0: github.com/redis/go-redis产生的错误
// var Err_030001000x err.ErrCode = err.ErrCode{
// 	Code:   "030001000x",
// 	Format: "%s.",
// 	Reason: "",
// }

// var Err_0300010001 err.ErrCode = err.ErrCode{
// 	Code:   "0300010001",
// 	Format: "Redis client %s is not exist.",
// 	Reason: "",
// }

var Err_030001000x err.ErrCode = err.New("030001000x", "%s.", "")
