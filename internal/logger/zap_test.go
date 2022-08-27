package logger

import (
	"testing"

	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
)

func TestNewZapLogger_checkConfiguration(t *testing.T) {
	conf := config.ApplicationConfig{
		Env: config.Prod,
		Logger: config.LoggerConfig{
			Level: config.ErrorLevel,
		},
	}
	logger, err := newZapLogger(&conf)
	assert.Nil(t, err)
	assert.Equal(t, config.ErrorLevel, logger.level)
	checkRes := logger.log.Desugar().Check(zapcore.InfoLevel, "test")
	assert.Nil(t, checkRes) //Logger won't write to log
	checkRes = logger.log.Desugar().Check(zapcore.ErrorLevel, "test")
	assert.NotNil(t, checkRes) //Logger will write to log
}

func TestNewZapLogger_notFailFromDefaultConfig(t *testing.T) {
	conf := config.ApplicationConfig{}
	logger, err := newZapLogger(&conf)
	assert.Nil(t, err)
	assert.NotNil(t, logger)
	//Check default info level
	checkRes := logger.log.Desugar().Check(zapcore.DebugLevel, "test")
	assert.Nil(t, checkRes) //Logger won't write to log
	checkRes = logger.log.Desugar().Check(zapcore.InfoLevel, "test")
	assert.NotNil(t, checkRes) //Logger will write to log
}
