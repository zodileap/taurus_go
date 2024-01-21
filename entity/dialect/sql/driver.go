package sql

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/yohobala/taurus_go/entity/dialect"
)

type (
	Driver struct {
		Conn
		dialect dialect.DbDriver
	}

	Conn struct {
		ExecQuerier
	}

	ExecQuerier interface {
		ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
		QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
		QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	}
)

// NewDriver 创建一个新的数据库驱动。
func NewDriver(driver dialect.DbDriver, c Conn) *Driver {
	return &Driver{
		Conn:    c,
		dialect: driver,
	}
}

// Dialect 返回数据库类型。
func (d Driver) Dialect() dialect.DbDriver {
	return d.dialect
}

// DB 返回数据库连接。
func (d Driver) DB() *sql.DB {
	return d.ExecQuerier.(*sql.DB)
}

// Close 关闭数据库连接。
func (d *Driver) Close() error { return d.DB().Close() }

// Tx 返回一个事务。
func (d *Driver) Tx(context.Context) (dialect.Tx, error) {
	return d.BeginTx(context.Background(), nil)
}

// BeginTx 开始一个事务。
func (d *Driver) BeginTx(ctx context.Context, opts *TxOptions) (dialect.Tx, error) {
	tx, err := d.DB().BeginTx(ctx, opts)
	if err != nil {
		return nil, err
	}
	return &Tx{
		Conn: Conn{tx},
		Tx:   tx,
	}, nil
}

// Query dialect.ExecQuerier.Query的实现
func (c Conn) Query(ctx context.Context, query string, args []any, v *dialect.Rows) error {
	rows, err := c.ExecQuerier.QueryContext(ctx, query, args...)
	if err != nil {
		return err
	}
	*v = dialect.Rows{rows}
	return nil
}

// Exec dialect.ExecQuerier.Exec的实现
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
