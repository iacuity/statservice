package writer

import (
	"database/sql"

	"github.com/statservice/data"
	"github.com/statservice/util"
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
	return nil
}
