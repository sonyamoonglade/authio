package authio

import (
	"log"
	"os"
)

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
}

type Level int

const (
	ErrorLevel Level = iota + 1
	WarnLevel
	InfoLevel
	DebugLevel
)

const (
	errorPrefix = "[ERROR]: "
	warnPrefix  = "[WARN]: "
	infoPrefix  = "[INFO]: "
	debugPrefix = "[DEBUG]: "
)

type DefaultLogger struct {
	logger *log.Logger
	level  Level
}

func NewDefaultLogger(level Level) *DefaultLogger {
	return &DefaultLogger{
		logger: log.New(os.Stdout, "", 0),
		level:  level,
	}
}

func (l *DefaultLogger) Debugf(msg string, args ...interface{}) {
	if l.level < 4 {
		return
	}
	addCaret(&msg)
	l.printf(debugPrefix, msg, args)
}

func (l *DefaultLogger) Infof(msg string, args ...interface{}) {
	if l.level < 3 {
		return
	}
	addCaret(&msg)
	l.printf(infoPrefix, msg, args)
}

func (l *DefaultLogger) Warnf(msg string, args ...interface{}) {
	if l.level < 2 {
		return
	}
	addCaret(&msg)
	l.printf(warnPrefix, msg, args)
}

func (l *DefaultLogger) Errorf(msg string, args ...interface{}) {
	addCaret(&msg)
	l.printf(errorPrefix, msg, args)
}

func addCaret(msg *string) *string {
	if (*msg)[len(*msg)-1] != '\n' {
		*msg += "\n"
	}
	return msg
}

func (l *DefaultLogger) printf(prefix string, msg string, args ...interface{}) {
	msg = prefix + msg
	l.logger.Printf(msg, args)
}
