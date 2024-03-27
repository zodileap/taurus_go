package sql

import (
	"database/sql/driver"

	"github.com/yohobala/taurus_go/entity/dialect"
)

// Tx 包装了数据库事务,实现了driver.Tx接口。
type Tx struct {
	Conn
	driver.Tx
	dialect dialect.DbDriver
}

// Dialect 返回数据库类型。
//
// Returns:
//
//	0: 数据库类型。
func (d Tx) Dialect() dialect.DbDriver {
	return d.dialect
}
