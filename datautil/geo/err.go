package geo

import (
	"github.com/zodileap/taurus_go/err"
)

// Err_020002000x geo: %s.
//
// Verbs:
//
//	0: 未知的错误。
var Err_020002000x err.ErrCode = err.New(
	"020002000x",
	"%s.",
	"",
)

/**************** CRUD遇到的问题 ***************/

// Err_0200020101 Scan()数据类型不是string
var Err_0200020101 err.ErrCode = err.New(
	"0200020001",
	"Scan() data type is not string.",
	"",
)

// Err_0200020102 不支持Scan，将类型%T存储到类型%T中
//
// Verbs:
//
//	0: 存储的数据类型。
//	1: 目标数据类型。
var Err_0200020102 err.ErrCode = err.New(
	"0200020002",
	`unsupported Scan, storing type "%v" into type %v`,
	"",
)

// Err_0200020103 将类型%T存储到类型%T,出现错误
//
// Verbs:
//
//	0: 存储的数据类型。
//	1: 目标数据类型。
//	2: 错误信息。
var Err_0200020103 err.ErrCode = err.New(
	"0200020003",
	`storing type %v into type "%v", error: %v`,
	"",
)

// Err_0200020104 不支持的GeomType类型
//
// Verbs:
//
//	0: GeomType类型。
var Err_0200020104 err.ErrCode = err.New(
	"0200020004",
	"unsupported GeomType type: %q",
	"",
)

/**************** 创建矢量遇到的问题 ***************/

// Err_0200020201 创建几何类型%s错误：%v
//
// Verbs:
//
//	0: 几何类型。
//	1: 错误信息。
var Err_0200020201 err.ErrCode = err.New(
	"0200020201",
	"create geometry type %q error: %v",
	"",
)

// Err_0200020202 从GeoJSON创建%s错误：%v
//
// Verbs:
//
//	0: 几何类型。
//	1: 错误信息。
var Err_0200020202 err.ErrCode = err.New(
	"0200020202",
	"create %q from GeoJSON error: %v",
	"",
)

// Err_0200020203 不是GeoJSON支持的几何类型
//
// Verbs:
//
//	0: 几何类型。
var Err_0200020203 err.ErrCode = err.New(
	"0200020203",
	"unsupported GeoJSON geometry type: %q",
	"",
)
