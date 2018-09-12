package data

import (
	"encoding/json"
)

type Pair struct {
	Key   string `json:"key,omitempty"`
	Value int64  `json:"val,omitempty"`
}

type Request struct {
	Pairs []Pair `json:"pairs,omitempty"`
}

type Config struct {
	LogConfig       *LogConfig
	ServerConfig    *ServerConfig
	Servlets        []*ServletConfig
	RefreshInterval *int
	DataDirectory   *string
	DBConfig        *DBConfig
}

type LogConfig struct {
	LogFile  *string
	LogLevel *int
}

type ServerConfig struct {
	Host *string
	Port *int
}

type ServletConfig struct {
	Name *string
	Path *string
}

type DBConfig struct {
	Hostname       *string
	Port           *int
	Username       *string
	Password       *string
	Database       *string
	IdleConnection *int
	MaxConnection  *int
}

func (config *Config) IsValid() bool {
	return true
}

func (config *Config) String() string {
	byts, err := json.Marshal(config)
	if nil != err {
		return ""
	}

	return string(byts)
}

func (config *DBConfig) String() string {
	byts, err := json.Marshal(config)
	if nil != err {
		return ""
	}

	return string(byts)
}
