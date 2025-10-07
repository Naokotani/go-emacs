package logger

import (
	"log"
	"os"
)

const (
	green     = "\033[32m"
	yellow    = "\033[33m"
	red       = "\033[31m"
	reset     = "\033[0m"
	infoLevel = "info"
	warnLevel = "warn"
	errLevel  = "err"
)

type Logger struct {
	Level    string
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	WarnLog  *log.Logger
}

func NewLogger(level string) Logger {
	logger := Logger{
		Level:    level,
		InfoLog:  log.New(os.Stdout, green+"INFO:\t"+reset, log.Ldate|log.Ltime),
		ErrorLog: log.New(os.Stderr, red+"ERROR:\t"+reset, log.Ldate|log.Ltime|log.Lshortfile),
		WarnLog:  log.New(os.Stdout, yellow+"WARN:\t"+reset, log.Ldate|log.Ltime),
	}
	if logger.Level == "" {
		logger.Level = errLevel
	}
	logger.validateLogLevel()
	return logger
}

func (l *Logger) validateLogLevel() {
	if l.Level == infoLevel || l.Level == warnLevel || l.Level == errLevel {
		return
	}
	l.WarnLog.Printf("%s is not a valid log level. Log level must be '%s', '%s', or '%s'. Falling back to 'error' level",
		l.Level, infoLevel, warnLevel, errLevel)
	l.Level = errLevel
}

func (l *Logger) Info(msg string) {
	if l.Level == infoLevel {
		l.InfoLog.Print(msg)
	}
}

func (l *Logger) Warn(msg string) {
	if l.Level == warnLevel || l.Level == infoLevel {
		l.WarnLog.Print(msg)
	}
}

func (l *Logger) Error(msg string) {
	l.ErrorLog.Print(msg)
}
