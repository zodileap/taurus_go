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

// Err_010002000x codegen遇到的未知错误。
var Err_010002000x err.ErrCode = err.New(
	"010002000x",
	"taurus_go/entity: %s.",
	"",
)

// Err_0100020001 调用New命令创建Schema,在创建数据库时验证数据库名出现错误。
//
// Verbs:
//
//	0: 数据库名。
//	1: 错误信息。
var Err_0100020001 err.ErrCode = err.New(
	"0100020001",
	"new database %q : %v",
	"1.传入的数据库名与Go中的关键字冲突。",
)

// Err_0100020002 调用New命令创建Schema,输入的数据库名不是大写字母开头。
//
// Verbs:
//
//	0: 数据库名。
var Err_0100020002 err.ErrCode = err.New(
	"0100020002",
	"database %q' must begin with uppercase",
	"",
)

// Err_0100020003 执行模板时出现错误。
//
// Verbs:
//
//	0: 模板文件名。
//	1: 错误信息。
var Err_0100020003 err.ErrCode = err.New(
	"0100020003",
	"execute template %q: %v",
	"",
)

// Err_0100020004 在加载Shema时出现错误。
//
// Verbs:
//
//	0: 错误信息。
var Err_0100020004 err.ErrCode = err.New(
	"0100020004",
	"load schema: %v",
	"",
)

// Err_0100020005 在加载Shema时，发现没有找到entity。
//
// Verbs:
//
//	0: Schema的路径
var Err_0100020005 err.ErrCode = err.New(
	"0100020005",
	"not entity found  in %s",
	"",
)

// Err_0100020006 格式化模版时出现错误。
var Err_0100020006 err.ErrCode = err.New(
	"0100020006",
	"format template: %v",
	"",
)

// Err_0100020007 创建'.gen'目录时出现错误。
var Err_0100020007 err.ErrCode = err.New(
	"0100020007",
	"create '.gen' directory: %v",
	"",
)

// Err_0100020008 通过模版写入文件时出现错误。
//
// Verbs:
//
//	0: 文件路径。
//	1: 错误信息
var Err_0100020008 err.ErrCode = err.New(
	"0100020008",
	"write file %s: %v",
	"",
)

// Err_0100020009 反序列化Shema配置文件时出现错误。
//
// Verbs:
//
//	0: 出现问题的内容。
//	1: 错误信息。
var Err_0100020009 err.ErrCode = err.New(
	"0100020009",
	"unmarshal schema config %s: %v",
	"",
)

// Err_0100020010 加载Schema和entity.Interface中的Go package时出现错误。
//
// Verbs:
//
//	0: 错误信息。
var Err_0100020010 err.ErrCode = err.New(
	"0100020010",
	"load package: %v",
	"",
)

// Err_0100020011 加载Schema时没有发现package信息。
//
// Verbs:
//
//	0: 加载的Schema的路径。
var Err_0100020011 err.ErrCode = err.New(
	"0100020011",
	"missing package information for: %s",
	"",
)

// Err_0100020012 在断言为 *ast.TypeSpec 类型时出现错误。
//
// Verbs:
//
//	0: 断言的类型。
//	1: 断言的字段名。
var Err_0100020012 err.ErrCode = err.New(
	"0100020012",
	"invalid declaration %T for %s",
	"",
)

// Err_0100020013 在断言为 *ast.StructType 类型时出现错误。
//
// Verbs:
//
//	0: 断言的类型。
//	1: 断言的字段名。
var Err_0100020013 err.ErrCode = err.New(
	"0100020013",
	"invalid spec type %T for %s",
	"",
)

// Err_0100020014 在断言为 *ast.Ident 类型时出现错误。
//
// Verbs:
//
//	0: 断言的类型。
//	1: 断言的字段名。
var Err_0100020014 err.ErrCode = err.New(
	"0100020014",
	"invalid field type %T for %s",
	"",
)

// Err_0100020015 无效的包名。
//
// Verbs:
//
//	0: 包名。
var Err_0100020015 err.ErrCode = err.New(
	"0100020015",
	"invalid package identifier: %s",
	"",
)

// Err_0100020016 调用parser.ParseFile解析go代码出现错误
//
// Verbs:
//
//	0: 错误信息。
var Err_0100020016 err.ErrCode = err.New(
	"0100020016",
	"parse entity file: %v",
	"",
)

// Err_0100020017 把解析的代码格式化，并添加到模板中时出现错误。
//
// Verbs:
//
//	0: 错误信息。
var Err_0100020017 err.ErrCode = err.New(
	"0100020017",
	"format node: %v",
	"",
)

// Err_0100020018 序列化实体时出现错误。
//
// Verbs:
//
//	0: 实体名。
//	1: 错误信息。
var Err_0100020018 err.ErrCode = err.New(
	"0100020018",
	"marshal entity %q : %v",
	"",
)

/**************** CRUD遇到的问题 ***************/

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

// Err_0100030004 子句中参数超过了最大值。
var Err_0100030004 err.ErrCode = err.New(
	"0100030004",
	"Clause Params is too many.",
	"",
)

// Err_0100030005 在update中set和predicate，参数数量不一致。
var Err_0100030005 err.ErrCode = err.New(
	"0100030005",
	"update set and predicate params count not equal.",
	"",
)

/**************** dialect遇到的问题 ***************/
