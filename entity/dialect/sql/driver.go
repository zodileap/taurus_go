package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/zodileap/taurus_go/entity/dialect"
)

type (
	// ExecQuerier 执行查询需要满足的接口。
	ExecQuerier interface {
		// ExecContext 执行不返回记录的查询。例如，SQL中INSERT或UPDATE, 但是也支持返回一些元数据。
		// 例如，PostgreSQL的INSERT ... RETURNING是使用QueryContext，
		// 但是MySQL的INSERT是使用ExecContext，因为它不支持RETURNING，但是为了兼容PostgreSQL的RETURNING，所以需要支持返回一些元数据。
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		// QueryContext 执行返回记录的查询，通常是SQL中的SELECT，或者有RETURNING子句的INSERT/UPDATE。
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	}

	// Driver 数据库驱动。
	Driver struct {
		Conn
		dialect dialect.DbDriver
	}

	// Conn 当前数据库连接，这个会在Driver.BeginTx()中被初始化，用于存放sql.Tx。
	Conn struct {
		ExecQuerier
	}
)

// NewDriver 创建一个新的数据库驱动。
//
// Params:
//
//   - driver: 数据库类型。
//   - c: 数据库连接。
//
// Returns:
//
//	0: 数据库驱动。
func NewDriver(driver dialect.DbDriver, c Conn) *Driver {
	return &Driver{
		Conn:    c,
		dialect: driver,
	}
}

// Dialect 返回数据库类型。
//
// Returns:
//
//	0: 数据库类型。
func (d Driver) Dialect() dialect.DbDriver {
	return d.dialect
}

// DB 返回数据库连接。
//
// Returns:
//
//	0: 数据库连接。
func (d Driver) DB() *sql.DB {
	return d.ExecQuerier.(*sql.DB)
}

// Close 关闭数据库连接。
func (d *Driver) Close() error { return d.DB().Close() }

// Tx 返回一个事务。
//
// Params:
//
//   - ctx: 上下文。
//
// Returns:
//
//	0: 事务。
//	1: 错误信息。
func (d *Driver) Tx(ctx context.Context) (dialect.Tx, error) {
	return d.BeginTx(ctx, nil)
}

// BeginTx 开始一个事务。
//
// Params:
//
//   - ctx: 上下文。
//   - opts: 事务选项。
//
// Returns:
//
//	0: 事务。
//	1: 错误信息。
func (d *Driver) BeginTx(ctx context.Context, opts *TxOptions) (dialect.Tx, error) {
	tx, err := d.DB().BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		Conn:    Conn{tx},
		Tx:      tx,
		dialect: d.dialect,
	}, nil
}

// Query dialect.ExecQuerier.Query的实现
//
// Params:
//
//   - ctx: 上下文。
//   - query: 查询语句。
//   - args: 查询参数。
//   - v: 查询结果。
func (c Conn) Query(ctx context.Context, query string, args []any, v *dialect.Rows) error {
	rows, err := c.ExecQuerier.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	*v = dialect.Rows{rows}
	return nil
}

// Exec dialect.ExecQuerier.Exec的实现
//
// Params:
//
//   - ctx: 上下文。
//   - query: 查询语句。
//   - args: 查询参数。
//   - v: 查询结果。
func (c Conn) Exec(ctx context.Context, query string, args []any, v any) error {

	switch v := v.(type) {
	case nil:
		if _, err := c.ExecContext(ctx, query, args...); err != nil {
			return err
		}
	case *sql.Result:
		res, err := c.ExecContext(ctx, query, args...)
		if err != nil {
			return err
		}
		*v = res
	default:
		return fmt.Errorf("dialect/sql: invalid type %T. expect *sql.Result", v)
	}
	return nil
}

type (
	TxOptions = sql.TxOptions
)
