package logger

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
	error_messages "udp-chat/infra/logger/constants"
	"udp-chat/infra/logger/model"
)

const (
	ErrorLevel = "ERROR"
	WarnLevel  = "WARNING"
	InfoLevel  = "INFO"
)

type Logger struct {
	logFilePath string
	serviceName string
}

func NewLogger(logFilePath, service string) Logger {
	return Logger{
		logFilePath: logFilePath,
		serviceName: service,
	}
}

func (l Logger) Error(err error) {
	errorLog := model.Log{
		Level:   ErrorLevel,
		Service: l.serviceName,
		Error:   err.Error(),
		Time:    time.Now(),
	}
	l.writeError(errorLog)
}

func (l Logger) Warn(msg string) {
	errorLog := model.Log{
		Level:   WarnLevel,
		Service: l.serviceName,
		Message: msg,
		Time:    time.Now(),
	}
	l.writeError(errorLog)
}

func (l Logger) Info(msg string) {
	errorLog := model.Log{
		Level:   InfoLevel,
		Service: l.serviceName,
		Message: msg,
		Time:    time.Now(),
	}
	l.writeError(errorLog)
}

func (l Logger) getLogFile(path string) (f *os.File) {
	filePath, err := filepath.Abs(path)
	if err != nil {
		return
	}

	if _, err = os.Stat(filePath); os.IsNotExist(err) {
		err = os.MkdirAll(filePath, 0755)
		if err != nil {
			log.Println(error_messages.FailedToCreateLogsFolder, err.Error())
			return
		}
	}

	fileName := fmt.Sprintf("%s.log", time.Now().Format("2006-01-02"))
	fullPath := fmt.Sprintf("%s/%s", filePath, fileName)

	f, err = os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Println(error_messages.FailedToCreateLogsFile, err.Error())
		return nil
	}

	return
}

func (l Logger) writeError(errorLog model.Log) {
	f := l.getLogFile(l.logFilePath)
	defer f.Close()

	b, e := json.Marshal(errorLog)
	if e != nil {
		log.Println(e.Error())
	}

	_, e = f.Write(b)
	if e != nil {
		log.Println(e.Error())
	}
}
