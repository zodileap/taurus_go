package telegram

import "github.com/zodileap/taurus_go/err"

// Err_0400010001 Telegram 参数校验失败。
var Err_0400010001 err.ErrCode = err.New(
	"0400010001",
	"telegram notify invalid argument: %s",
	"",
)

// Err_0400010002 Telegram 请求构造失败。
var Err_0400010002 err.ErrCode = err.New(
	"0400010002",
	"telegram notify build request failed: %v",
	"",
)

// Err_0400010003 读取本地文件失败。
var Err_0400010003 err.ErrCode = err.New(
	"0400010003",
	"telegram notify read local file %s failed: %v",
	"",
)

// Err_0400010004 发送 HTTP 请求失败。
var Err_0400010004 err.ErrCode = err.New(
	"0400010004",
	"telegram notify send request failed: %v",
	"",
)

// Err_0400010005 Telegram 返回非 2xx 状态码。
var Err_0400010005 err.ErrCode = err.New(
	"0400010005",
	"telegram notify received http status %d: %s",
	"",
)

// Err_0400010006 Telegram 返回 ok=false。
var Err_0400010006 err.ErrCode = err.New(
	"0400010006",
	"telegram notify api error %d: %s",
	"",
)

// Err_0400010007 Telegram 响应解析失败。
var Err_0400010007 err.ErrCode = err.New(
	"0400010007",
	"telegram notify decode response failed: %v",
	"",
)

// Err_0400010008 Telegram 不支持的媒体组组合。
var Err_0400010008 err.ErrCode = err.New(
	"0400010008",
	"telegram notify unsupported media group combination: %s",
	"",
)
