package asset

import "github.com/yohobala/taurus_go/err"

// Err_020001000x 未知的错误。
//
// Verbs:
//
//	0: 错误的信息。
var Err_020001000x err.ErrCode = err.New(
	"020001000x",
	"%s.",
	"",
)

// Err_0200010001 创建文件夹失败。
//
// Verbs:
//
//	0: 文件夹路径。
//	1: os.MkdirAll产生的错误信息。
var Err_0200010001 err.ErrCode = err.New(
	"0200010001",
	"Create directory %s : %+v",
	"可能得原因是：1.传入的文件夹路径不正确，比如传入一个空的路径``。 ",
)

// Err_0200010002 写入文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: os.WriteFile产生的错误信息。
var Err_0200010002 err.ErrCode = err.New(
	"0200010002",
	"Write file %s : %+v",
	"",
)

// Err_0200010003 格式化文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: imports.Process产生的错误信息
var Err_0200010003 err.ErrCode = err.New(
	"0200010003",
	"Format file  %s : %+v",
	"可能得原因是：1.传入的文件内容缺少package。",
)

// Err_0200010004 打开文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: os.Open产生的错误信息。
var Err_0200010004 err.ErrCode = err.New(
	"0200010004",
	"Open file %s : %+v",
	"",
)

// Err_0200010005 创建文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: os.Create产生的错误信息。
var Err_0200010005 err.ErrCode = err.New(
	"0200010005",
	"Create file %s : %+v",
	"可能原因是：1.传入的路径是一个文件而不是文件夹。",
)

// Err_0200010006 复制文件失败。
//
// Verbs:
//
//	0: 源文件路径。
//	1: 目标文件路径。
//	2: io.Copy产生的错误信息。
var Err_0200010006 err.ErrCode = err.New(
	"0200010006",
	"Copy file %s to %s : %+v",
	"",
)

// Err_0200010007 打开文件失败。
//
// Verbs:
//
//	0: 文件路径。
//	1: filepath.Walk产生的错误信息。
var Err_0200010007 err.ErrCode = err.New(
	"0200010007",
	"Open file %s : %+v",
	"",
)

// Err_0200010008 计算两个路径的相对路径失败。
//
// Verbs:
//
//	0: 基础路径。
//	1: 目标路径。
//	2: filepath.Rel产生的错误信息。
var Err_0200010008 err.ErrCode = err.New(
	"0200010008",
	"Relative path %s with %s : %+v",
	"",
)
