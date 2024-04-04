package entity

import (
	"database/sql"
	"fmt"
	"strings"
	"sync"

	"github.com/yohobala/taurus_go/entity/dialect"
	dsql "github.com/yohobala/taurus_go/entity/dialect/sql"
)

var (
	clients map[string]ConnectionConfig = make(map[string]ConnectionConfig)
	mu      sync.RWMutex
)

// AddConnection 在添加一个数据库连接配置。
func AddConnection(conn ConnectionConfig) error {
	if conn.Tag == "" {
		return Err_0100010001
	}
	switch conn.Driver {
	case dialect.PostgreSQL, dialect.MySQL:
		if _, ok := clients[conn.Tag]; ok {
			return Err_0100010004.Sprintf(conn.Tag)
		} else {
			clients[conn.Tag] = conn
			return nil
		}
	default:
		return Err_0100010002.Sprintf(conn.Driver)
	}
}

// GetConnection 获取一个数据库连接。
func GetConnection(tag string) (dialect.Driver, error) {
	mu.RLock()
	defer mu.RUnlock()
	conn, ok := clients[tag]
	if !ok {
		return nil, Err_0100010003.Sprintf(tag)
	}

	var dbUrl string
	switch conn.Driver {
	case dialect.PostgreSQL:
		if conn.IsVerifyCa {
			dbUrl = fmt.Sprintf(
				`postgres://%s:%s@%s:%d/%s?sslmode=verify-full&sslrootcert=%s&sslcert=%s&sslkey=%s`,
				conn.User,
				conn.Password,
				conn.Host,
				conn.Port,
				conn.DBName,
				conn.RootCertPath,
				conn.ClientCertPath,
				conn.ClientKeyPath)
		} else {
			dbUrl = fmt.Sprintf(
				`postgres://%s:%s@%s:%d/%s?sslmode=disable`,
				conn.User,
				conn.Password,
				conn.Host,
				conn.Port,
				conn.DBName,
			)
		}
	case dialect.MySQL:
		if conn.IsVerifyCa {
			dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=custom&sslmode=%s&sslrootcert=%s&sslcert=%s&sslkey=%s",
				conn.User,
				conn.Password,
				conn.Host,
				conn.Port,
				conn.DBName,
				"verify-ca",
				conn.RootCertPath,
				conn.ClientCertPath,
				conn.ClientKeyPath)
		} else {
			dbUrl = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?tls=false",
				conn.User,
				conn.Password,
				conn.Host,
				conn.Port,
				conn.DBName)
		}
	}
	db, err := sql.Open(string(conn.Driver), dbUrl)
	if err != nil {
		if strings.Contains(err.Error(), "sql: unknown driver") {
			return nil, Err_0100010005.Sprintf(string(conn.Driver))
		}
		return nil, Err_010001000x.Sprintf(err)
	}
	drv := dsql.NewDriver(conn.Driver, dsql.Conn{ExecQuerier: db})
	return drv, nil
}
