package config

import (
	"github.com/kelseyhightower/envconfig"
)

type ApplicationConfig struct {
	Env    Environment `default:"dev"`
	Logger LoggerConfig
	Queue  QueueConfig
	GRPC   GRPCConfig
}

// ReadConfig reads configuration from the environment and populates the structure with it
func ReadConfig() (*ApplicationConfig, error) {
	var conf ApplicationConfig
	if err := envconfig.Process("", &conf); err != nil {
		return nil, err
	}
	return &conf, nil
}
