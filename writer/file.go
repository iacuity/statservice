package writer

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"
	"time"

	"github.com/statservice/data"
)

const (
	DATE_FORMAT   = "2006-01-02"
	KEY_SEPARATOR = "="
)

var (
	dataDirectory = "./"
	currentDay    int
	currentFile   string
	fsMap         map[string]int64
)

type FileWriter struct {
}

func getFileName() string {
	return fmt.Sprintf("%s/%s", dataDirectory, time.Now().Local().Format("2006-01-02"))
}

func readFile() error {
	content, err := ioutil.ReadFile(currentFile)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(bytes.NewReader(content))
	for scanner.Scan() {
		line := scanner.Text()
		token := strings.Split(line, KEY_SEPARATOR)
		val, _ := strconv.ParseInt(token[1], 10, 64)
		fsMap[token[0]] = val
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

func writeFile() error {
	var buff bytes.Buffer
	for key, val := range fsMap {
		buff.WriteString(key)
		buff.WriteString(KEY_SEPARATOR)
		buff.WriteString(fmt.Sprintf("%d\n", val))
	}
	if err := ioutil.WriteFile(currentFile, buff.Bytes(), 0644); nil != err {
		return err
	}

	if day := time.Now().Day(); currentDay != day { // handle day change
		currentFile = getFileName()
		currentDay = day
		fsMap = make(map[string]int64)
	}

	return nil
}

func (w *FileWriter) Init(config *data.Config) error {
	dataDirectory = *config.DataDirectory
	currentDay = time.Now().Day()
	currentFile = getFileName()
	fsMap = make(map[string]int64)
	readFile()
	return nil
}

func (w *FileWriter) Write(sMap map[string]int64) error {
	for key, val := range sMap {
		if fval, found := fsMap[key]; found {
			fsMap[key] = fval + val
		} else {
			fsMap[key] = val
		}
	}

	return writeFile()
}
