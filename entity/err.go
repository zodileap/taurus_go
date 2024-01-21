package entity

import (
	"github.com/yohobala/taurus_go/err"
)

/**************** 初始化遇到的问题 ***************/

// Err_010001000x 初始化时遇到的未知错误。
//
// Verbs:
//
//	0: 未知的错误。
var Err_010001000x err.ErrCode = err.New(
	"010001000x",
	"%s.",
	"",
)

// Err_0100010001 添加连接时,配置中的tag为空。
var Err_0100010001 err.ErrCode = err.New(
	"0100010001",
	"Add Connection failed, tag is empty.",
	"",
)

// Err_0100010002 添加连接时,配置中的driverName不支持。
//
// Verbs:
//
//	0: 配置中的driverName。
var Err_0100010002 err.ErrCode = err.New(
	"0100010002",
	"Add Connection failed, driverName '%s' not support.",
	"",
)

// Err_0100010003 得到连接错误，因为连接信息标签不存在。
//
// Verbs:
//
//	0: 连接信息标签。
var Err_0100010003 err.ErrCode = err.New(
	"0100010003",
	"Received a connection error because the connection tag '%s' does not exist.",
	"",
)

// Err_0100010004 添加连接错误，因为连接信息标签已经存在。
//
// Verbs:
//
//	0: 连接信息标签。
var Err_0100010004 err.ErrCode = err.New(
	"0100010004",
	"Add Connection failed, the connection tag '%s' already exists.",
	"",
)

// Err_0100010005 创建数据库实例时，遇到未知的驱动，可能是忘记导入驱动包。
var Err_0100010005 err.ErrCode = err.New(
	"0100010005",
	"unknown driver %q.",
	"1. check import driver package.",
)

/**************** codegen遇到的问题 ***************/

/**************** entitySQL遇到的问题 ***************/

// Err_0100030001 在创建语句中，必填但没有默认值的字段的值为空。
//
// Verbs:
//
//	0: 实体表的名字。
//	1: 字段名。
var Err_0100030001 err.ErrCode = err.New(
	"0100030001",
	"entity table %s field %s required,but value is nil.",
	"",
)

// Err_0100030002 在升级语句中，没有需要更新的字段。
//
// Verbs:
//
//	0: 实体表的名字。
var Err_0100030002 err.ErrCode = err.New(
	"0100030002",
	"entity table %s no fields need to update.",
	"",
)

// Err_0100030003 改变实体类的跟踪状态失败。
//
// Verbs:
//
//	0: 需要改变成的状态。
//	1: 当前应该属于的状态。
var Err_0100030003 err.ErrCode = err.New(
	"0100030003",
	"change entity tracking status to %s failed.The current state needs to be %s.",
	"",
)

/**************** dialect遇到的问题 ***************/
