package asset

import "github.com/yohobala/taurus_go/err"

/**************** 通用错误 ***************/

// Err_0200010001 计算两个路径的相对路径失败。
//
// Verbs:
//
//	0: 基础路径。
//	1: 目标路径。
//	2: filepath.Rel产生的错误信息。
var Err_0200010001 err.ErrCode = err.New(
	"0200010001",
	"Relative path %s with %s : %+v",
	"",
)

/**************** 文件错误 ***************/

// Err_020002000x 未知的错误。
//
// Verbs:
//
//	0: 错误的信息。
var Err_020002000x err.ErrCode = err.New(
	"020002000x",
	"%s.",
	"",
)

// Err_0200020001 写入文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: os.WriteFile产生的错误信息。
var Err_0200020001 err.ErrCode = err.New(
	"0200020001",
	"Write file %s : %+v",
	"",
)

// Err_0200020002 格式化文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: imports.Process产生的错误信息
var Err_0200020002 err.ErrCode = err.New(
	"0200020002",
	"Format file  %s : %+v",
	"可能得原因是：1.传入的文件内容缺少package。",
)

// Err_0200020003 打开文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: os.Open产生的错误信息。
var Err_0200020003 err.ErrCode = err.New(
	"0200020003",
	"Open file %s : %+v",
	"",
)

// Err_0200020004 创建文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: os.Create产生的错误信息。
var Err_0200020004 err.ErrCode = err.New(
	"0200020004",
	"Create file %s : %+v",
	"可能原因是：1.传入的路径是一个文件而不是文件夹。",
)

// Err_0200020005 复制文件失败。
//
// Verbs:
//
//	0: 源文件路径。
//	1: 目标文件路径。
//	2: io.Copy产生的错误信息。
var Err_0200020005 err.ErrCode = err.New(
	"0200020005",
	"Copy file %s to %s : %+v",
	"",
)

// Err_0200020006 文件不存在。
//
// Verbs:
//
//	0: 文件路径。
var Err_0200020006 err.ErrCode = err.New(
	"0200020006",
	"File not exists %s",
	"",
)

// Err_0200020007 读取文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: bufio.Scanner产生的错误信息。
var Err_0200020007 err.ErrCode = err.New(
	"0200020007",
	"Read file %s : %+v",
	"",
)

// Err_0200020008 在文件中插入内容的位置不应该小于指定值
var Err_0200020008 err.ErrCode = err.New(
	"0200020008",
	"Insert position should not less than %d",
	"",
)

/**************** 文件夹错误 ***************/

// Err_0200030001 创建文件夹失败。
//
// Verbs:
//
//	0: 文件夹路径。
//	1: os.MkdirAll产生的错误信息。
var Err_0200030001 err.ErrCode = err.New(
	"0200030001",
	"Create directory %s : %+v",
	"可能得原因是：1.传入的文件夹路径不正确，比如传入一个空的路径``。 ",
)

// Err_0200030002 打开文件夹失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: filepath.Walk产生的错误信息。
var Err_0200030002 err.ErrCode = err.New(
	"0200030002",
	"Open directory %s : %+v",
	"",
)
