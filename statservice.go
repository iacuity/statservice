package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/llog"

	"github.com/statservice/data"
	"github.com/statservice/util"
	"github.com/statservice/writer"
)

var (
	config   data.Config
	msqlConn *sql.DB

	writers = []writer.IWritter{
		&writer.FileWriter{},
		&writer.MySQLWriter{},
	}
)

func init() {
	configFile := flag.String("config", "", "Stat Service Configuration File")
	flag.Parse()

	if !util.IsValidFile(*configFile) {
		log.Println("provide configuration file witth -config=<Config File> option")
		os.Exit(1)
	}

	if nil != util.ReadConfig(*configFile, &config) {
		os.Exit(1)
	}

	llog.Init(*config.LogConfig.LogFile)
	llog.SetLogLevel(llog.LogLevel(*config.LogConfig.LogLevel))
	llog.Info("Read configuration success. configuration is: %s", config.String())
	llog.Debug("Initializing Stat Service...")

	var err error
	if msqlConn, err = util.GetSQLConnection(config.DBConfig); nil != err {
		llog.Error("Failed to connect mysql server:%s", err.Error())
		os.Exit(1)
	}

	llog.Info("Register service access handler...")
	for _, servlet := range config.Servlets {
		switch *servlet.Name {
		case "statservice":
			http.HandleFunc(*servlet.Path, handleStat)
		}
	}

	//	var reqObject interface{}

	//	switch r.Method {
	//	case "GET":
	//		query := r.FormValue("query")
	//		logger.Debug("request data is: %s", query)
	//		if !util.IsBlank(query) {
	//			err := json.Unmarshal([]byte(query), &reqObject)
	//			if nil != err {
	//				logger.Error("Invalid Request: %s:%s", query, err.Error())
	//				w.WriteHeader(http.StatusBadRequest)
	//				return
	//			}
	//		}
	//	case "POST":
	//		decoder := json.NewDecoder(r.Body)
	//		if nil != decoder {
	//			err := decoder.Decode(&reqObject)
	//			if err != nil {
	//				logger.Error("Invalid Request: %s", err.Error())
	//				w.WriteHeader(http.StatusBadRequest)
	//				return
	//			}
	//		}
	//	}
}

func main() {
	llog.Info("Starting adserver...")
	addr := fmt.Sprintf("%s:%d", *config.ServerConfig.Host,
		*config.ServerConfig.Port)
	http.ListenAndServe(addr, nil)
	return
}
