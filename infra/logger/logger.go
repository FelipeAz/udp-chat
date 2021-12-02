package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	"udp-chat/infra/logger/model"
)

const (
	ErrorLevel = "ERROR"
	WarnLevel  = "WARNING"
	InfoLevel  = "INFO"
)

type Logger struct {
	logFilePath string
}

func NewLogger(logFilePath string) Logger {
	return Logger{
		logFilePath: logFilePath,
	}
}

func (l Logger) Error(err error) {
	errorLog := model.Log{
		Level: ErrorLevel,
		Error: err.Error(),
		Time:  time.Now(),
	}
	l.writeError(errorLog)
}

func (l Logger) Warn(msg string) {
	errorLog := model.Log{
		Level:   WarnLevel,
		Message: msg,
		Time:    time.Now(),
	}
	l.writeError(errorLog)
}

func (l Logger) Info(msg string) {
	errorLog := model.Log{
		Level:   InfoLevel,
		Message: msg,
		Time:    time.Now(),
	}
	l.writeError(errorLog)
}

func (l Logger) getLogFile(path string) (f *os.File) {
	filePath, err := filepath.Abs(path)
	if err != nil {
		log.Println("failed to retrieve log file:", err.Error())
		return
	}

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = os.Mkdir(filePath, 0755)
		if err != nil {
			log.Println("failed to create logs folder", err.Error())
			return
		}
	}

	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
	fullPath := fmt.Sprintf("%s/%s", filePath, fileName)

	f, err = os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println("failed to create log file:", err.Error())
		return nil
	}

	return
}

func (l Logger) writeError(errorLog model.Log) {
	f := l.getLogFile(l.logFilePath)
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Println("failed to close log file", err.Error())
		}
	}(f)

	b, e := json.Marshal(errorLog)
	if e != nil {
		log.Println(e.Error())
	}

	_, e = f.Write(b)
	if e != nil {
		log.Println(e.Error())
	}
}
