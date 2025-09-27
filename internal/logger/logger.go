package logger

import (
	"log"
	"os"
)

const (
	green  = "\033[32m"
	yellow = "\033[33m"
	red    = "\033[31m"
	reset  = "\033[0m"
)

type Logger struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	WarnLog  *log.Logger
}

func NewLogger() Logger {
	return Logger{
		InfoLog:  log.New(os.Stdout, green+"INFO:\t"+reset, log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, red+"ERROR:\t"+reset, log.Ldate|log.Ltime|log.Lshortfile),
		WarnLog:  log.New(os.Stdout, yellow+"WARN:\t"+reset, log.Ldate|log.Ltime),
	}
}
