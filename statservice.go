package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/llog"

	"github.com/statservice/data"
	"github.com/statservice/util"
	"github.com/statservice/writer"
)

const (
	MAX_CHANNEL_BUFFER = 100000
)

var (
	config  data.Config
	msgChan chan []data.Pair

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
	for _, wrtr := range writers {
		if err = wrtr.Init(&config); nil != err {
			llog.Error("Failed to initialize writter:%s", err.Error())
			os.Exit(1)
		}
	}

	llog.Info("Register service access handler...")
	for _, servlet := range config.Servlets {
		switch *servlet.Name {
		case "statservice":
			http.HandleFunc(*servlet.Path, handleRequest)
		}
	}

	msgChan = make(chan []data.Pair, MAX_CHANNEL_BUFFER)
	go updateStat()
}

func updateStat() {
	ticker := time.NewTicker(time.Second * (time.Duration)(*config.RefreshInterval))
	var sMap [2]map[string]int64
	sMap[0] = make(map[string]int64)
	var sMapIdx uint8 = 0
	for {
		select {
		case pairs := <-msgChan:
			for _, pair := range pairs {
				if val, found := sMap[sMapIdx][pair.Key]; !found {
					sMap[sMapIdx][pair.Key] = pair.Value
				} else {
					sMap[sMapIdx][pair.Key] = val + pair.Value
				}
			}
		case <-ticker.C:
			for _, wrtr := range writers {
				go wrtr.Write(sMap[sMapIdx])
			}
			sMap[sMapIdx^1] = make(map[string]int64)
			sMapIdx ^= 1
		}
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
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
		w.WriteHeader(http.StatusBadRequest)
	case "POST":
		req := data.Request{}
		decoder := json.NewDecoder(r.Body)
		if nil != decoder {
			err := decoder.Decode(&req)
			if err != nil {
				llog.Error("Invalid Request: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}

			msgChan <- req.Pairs
		}
	}
}

func main() {
	llog.Info("Starting stat service...")
	addr := fmt.Sprintf("%s:%d", *config.ServerConfig.Host,
		*config.ServerConfig.Port)
	http.ListenAndServe(addr, nil)
	return
}
