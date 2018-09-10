package util

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type IConfig interface {
	IsValid() bool
	String() string
}

func ReadConfig(configFile string, config IConfig) error {
	file, err := ioutil.ReadFile(configFile)

	if nil != err {
		log.Printf("Error while reading config file:%s \n\tError: %s", configFile, err.Error())
		return err
	}

	err = json.Unmarshal(file, config)

	if nil != err {
		log.Printf("Error while reading config file:%s \n\tError: %s", configFile, err.Error())
		return err
	}

	if !config.IsValid() {
		log.Printf("Invalid Configuration file: %s", configFile)
		return err
	}

	return nil
}
