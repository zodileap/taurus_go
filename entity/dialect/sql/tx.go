package sql

import "database/sql/driver"

type Tx struct {
	Conn
	driver.Tx
}
