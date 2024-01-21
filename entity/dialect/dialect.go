package dialect

import (
	"context"
	"database/sql"
	"database/sql/driver"
)

type (
	// 数据库类型。如"postgres"。
	DbDriver string
	// Rows 包装了sql.Rows，以避免locks copy。
	// sql.Rows 对象包含对数据库连接的引用以及用于迭代结果集的内部状态。
	// 因此，复制 sql.Rows 对象（例如，将其作为值传递给函数）可能会导致不可预测的行为，
	// 因为复制可能会导致对内部状态和数据库连接的多个引用。
	Rows struct{ RowsScanner }
)

const (
	PostgreSQL DbDriver = "postgres"
	MySQL      DbDriver = "mysql"
)

type ExecQuerier interface {
	// Exec 执行不返回记录的查询。例如，SQL中INSERT或UPDATE。
	// 它将结果扫描到指针v中，对于SQL驱动程序，它是sql.Rows,
	// v 可以是nil，或者是*sql.Result，如果是nil，将忽略结果。
	Exec(ctx context.Context, query string, args []any, v any) error
	// Query 执行返回记录的查询，通常是SQL中的SELECT，
	// 或者有RETURNING子句的INSERT/UPDATE。
	// 它将结果扫描到指针v中，对于SQL驱动程序，它是sql.Result。
	Query(ctx context.Context, query string, args []any, v *Rows) error
}

type Driver interface {
	ExecQuerier
	// Tx 启动并返回一个新事务。
	// 在事务提交或回滚之前，将使用所提供的上下文
	Tx(context.Context) (Tx, error)
	// Close 关闭数据库连接。
	Close() error
	// Dialect 返回数据库的驱动
	Dialect() DbDriver
}

type Tx interface {
	ExecQuerier
	driver.Tx
}

// RowScanner 封装了sql.Row的标准方法，用于扫描数据库行。
type RowsScanner interface {
	Close() error
	ColumnTypes() ([]*sql.ColumnType, error)
	Columns() ([]string, error)
	Err() error
	Next() bool
	NextResultSet() bool
	Scan(dest ...any) error
}
