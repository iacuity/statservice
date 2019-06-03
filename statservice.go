package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/llog"

	"github.com/statservice/data"
	"github.com/statservice/util"
	"github.com/statservice/writer"
)

const (
	MAX_CHANNEL_BUFFER = 1000000
)

var (
	config         data.Config
	maxMsgChannels = runtime.NumCPU()
	msgChan        = make([]chan []data.Pair, maxMsgChannels)
	writers        = []writer.IWritter{
		&writer.FileWriter{},
		&writer.MySQLWriter{},
	}
)

func init() {
	rand.Seed(time.Now().UnixNano())
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

	for i := 0; i < maxMsgChannels; i++ {
		msgChan[i] = make(chan []data.Pair, MAX_CHANNEL_BUFFER)
	}
	go updateStat()
}

func updateStat() {
	sMapChan := make(chan map[string]int64)
	for i := 0; i < maxMsgChannels; i++ {
		go func(i int) {
			ticker := time.NewTicker(time.Second * (time.Duration)(*config.RefreshInterval))
			var sMap [2]map[string]int64
			sMap[0] = make(map[string]int64)
			var sMapIdx uint8 = 0
			for {
				select {
				case pairs := <-msgChan[i]:
					for _, pair := range pairs {
						if val, found := sMap[sMapIdx][pair.Key]; !found {
							sMap[sMapIdx][pair.Key] = pair.Value
						} else {
							sMap[sMapIdx][pair.Key] = val + pair.Value
						}
					}
				case <-ticker.C:
					sMapChan <- sMap[sMapIdx]
					sMap[sMapIdx^1] = make(map[string]int64)
					sMapIdx ^= 1
				}
			}
		}(i)
	}
	sMapArray := make([]map[string]int64, 0)

	for {
		select {
		case sMap := <-sMapChan:
			sMapArray = append(sMapArray, sMap)
			if maxMsgChannels == len(sMapArray) {
				aggSMap := make(map[string]int64)
				for _, sMapEle := range sMapArray {
					for key, value := range sMapEle {
						if val, found := aggSMap[key]; !found {
							aggSMap[key] = value
						} else {
							aggSMap[key] = val + value
						}
					}
				}
				for _, wrtr := range writers {
					go wrtr.Write(aggSMap)
				}
				sMapArray = make([]map[string]int64, 0)
			}
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

			rNum := rand.Intn(maxMsgChannels) % maxMsgChannels
			msgChan[rNum] <- req.Pairs
			llog.Debug("Data Come :%v", msgChan[0])
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
