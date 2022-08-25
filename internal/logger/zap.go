package logger

import (
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"go.uber.org/zap"
)

// ZapLogger represents zap implementation of the Logger interface
type ZapLogger struct {
	log   *zap.SugaredLogger
	level config.LoggerLevel
}

func (l *ZapLogger) Fatal(args ...interface{}) {
	l.log.Fatal(args)
}

func (l *ZapLogger) Fatalf(format string, args ...interface{}) {
	l.log.Fatalf(format, args...)
}

func (l *ZapLogger) Error(args ...interface{}) {
	l.log.Error(args)
}

func (l *ZapLogger) Errorf(format string, args ...interface{}) {
	l.log.Errorf(format, args...)
}

func (l *ZapLogger) Warn(args ...interface{}) {
	l.log.Warn(args)
}

func (l *ZapLogger) Warnf(format string, args ...interface{}) {
	l.log.Warnf(format, args...)
}

func (l *ZapLogger) Info(args ...interface{}) {
	l.log.Info(args)
}

func (l *ZapLogger) Infof(format string, args ...interface{}) {
	l.log.Infof(format, args...)
}

func (l *ZapLogger) Debug(args ...interface{}) {
	l.log.Debug(args)
}

func (l *ZapLogger) Debugf(format string, args ...interface{}) {
	l.log.Debugf(format, args...)
}

func (l *ZapLogger) Trace(args ...interface{}) {
	if l.level != config.TraceLevel {
		l.log.Debug(args)
	}
}

func (l *ZapLogger) Tracef(format string, args ...interface{}) {
	if l.level != config.TraceLevel {
		l.log.Debugf(format, args...)
	}
}

// createZapSugaredLogger creates sugared logger from the application configuration
func createZapSugaredLogger(conf *config.ApplicationConfig) (*zap.SugaredLogger, error) {
	var zConf zap.Config
	switch conf.Env {
	case config.Prod:
		zConf = zap.NewProductionConfig()
	case config.Dev:
		zConf = zap.NewDevelopmentConfig()
	default:
		zConf = zap.NewDevelopmentConfig()
	}
	zConf.Level.SetLevel(conf.Logger.Level.ToZapLevel())

	//We are using the wrapper, so to remove the wrapper from the call trace, we need to add AddCallerSkip option
	logger, err := zConf.Build(zap.AddCallerSkip(1))
	if err != nil {
		return nil, err
	}

	return logger.Sugar(), nil
}

// newZapLogger creates new ZapLogger instance using application configuration
func newZapLogger(conf *config.ApplicationConfig) (*ZapLogger, error) {
	logger, err := createZapSugaredLogger(conf)
	if err != nil {
		return nil, err
	}

	zapLogger := &ZapLogger{
		log:   logger,
		level: conf.Logger.Level,
	}

	return zapLogger, nil
}
