package util

import (
	"log"
	"os"
)

func IsValidFile(file string) bool {
	retVal := false
	for {
		if "" == file {
			log.Println("File Name should not be empty")
			break
		}

		if _, err := os.Stat(file); os.IsNotExist(err) {
			log.Println(file, " does not exists.")
			break
		}

		retVal = true
		break
	}

	return retVal
}
