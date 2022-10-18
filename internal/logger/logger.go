package logger

type Logger interface {
	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
}
