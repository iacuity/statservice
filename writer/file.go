package writer

import (
	"fmt"
	"time"
)

const (
	DATE_FORMAT = "2006-01-02"
)

type FileWriter struct {
}

func getFileName() string {
	return time.Now().Local().Format("2006-01-02")
}

func init() {
	fmt.Println("Initializing file writter")
	fmt.Println("The Current time is ", getFileName())
}

func (w *FileWriter) Write(sMap map[string]int64) error {
	return nil
}
