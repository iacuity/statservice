package writer

import (
	"bytes"
	"database/sql"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/llog"

	"github.com/statservice/data"
	"github.com/statservice/util"
)

const (
	TABLE_STAT_METRIC  = "stat_metric"
	INSERT_STAT_METRIC = "INSERT INTO " + TABLE_STAT_METRIC +
		"(timestamp, metric, value) VALUES"
)

var (
	msqlConn *sql.DB
)

type MySQLWriter struct {
}

func (w *MySQLWriter) Init(config *data.Config) error {
	var err error
	msqlConn, err = util.GetSQLConnection(config.DBConfig)
	return err
}

func (w *MySQLWriter) Write(sMap map[string]int64) error {
	var params []interface{}
	var buffer bytes.Buffer
	buffer.WriteString(INSERT_STAT_METRIC)
	var values []string
	timestamp := time.Now().Local().Format("20060102150405")
	for key, val := range sMap {
		values = append(values, "(?, ?, ?)")
		params = append(params, timestamp)
		params = append(params, key)
		params = append(params, val)
	}
	buffer.WriteString(strings.Join(values, ","))
	stmt, err := msqlConn.Prepare(buffer.String())
	if nil != err {
		llog.Error("Insert Syntax Error: %s\n\tError Query: %s", err.Error(), buffer.String())
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(params...)

	if nil != err {
		llog.Error("Insert Execute Error: %s\nError Query: %s", err.Error(), buffer.String())
		return err
	}

	return err
}
