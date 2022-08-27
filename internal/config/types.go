package config

import "go.uber.org/zap/zapcore"

// Environment defines the application environment to adjust settings to it.
type Environment string

const (
	Dev  Environment = "dev"
	Prod Environment = "prod"
)

// LoggerLevel defines the minimum logging level to process
type LoggerLevel string

const (
	TraceLevel LoggerLevel = "trace"
	DebugLevel LoggerLevel = "debug"
	InfoLevel  LoggerLevel = "info"
	WarnLevel  LoggerLevel = "warn"
	ErrorLevel LoggerLevel = "error"
	PanicLevel LoggerLevel = "panic"
)

func (ll LoggerLevel) ToZapLevel() zapcore.Level {
	switch ll {
	case TraceLevel, DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	default:
		return zapcore.InfoLevel
	}
}
