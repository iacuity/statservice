package util

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/llog"

	"github.com/statservice/data"
)

const (
	MYSQL_DRIVER_NAME = "mysql"
	CONN_MAX_LIFETIME = 1 * 60 * 60 // 1 day
)

// return routine safe connection pool object
func GetSQLConnection(conf *data.DBConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s", *conf.Username, *conf.Password,
		*conf.Hostname, *conf.Port, *conf.Database)

	llog.Debug("DATABASE CONFIG: %s", conf.String())
	llog.Debug("DATABASE DSN: %s", dsn)

	conn, err := sql.Open(MYSQL_DRIVER_NAME, dsn)

	if nil != err {
		return nil, err
	}

	llog.Debug("Connection success")
	conn.SetMaxIdleConns(1)
	conn.SetMaxOpenConns(10)
	conn.SetConnMaxLifetime(time.Second * CONN_MAX_LIFETIME)

	err = conn.Ping()

	if nil != err {
		conn.Close()
		return nil, err
	}

	llog.Debug("Connection ping success")

	return conn, nil
}
