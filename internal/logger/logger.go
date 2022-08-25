package logger

import "github.com/rog-golang-buddies/internal/config"

// Logger represents common logger interface
//
//go:generate mockgen -source=logger.go -destination=./mocks/logger.go
type Logger interface {
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Trace(args ...interface{})
	Tracef(format string, args ...interface{})
}

func NewLogger(conf *config.ApplicationConfig) (Logger, error) {
	return newZapLogger(conf)
}
