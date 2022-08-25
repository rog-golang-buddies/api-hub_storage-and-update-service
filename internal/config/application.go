package config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type ApplicationConfig struct {
	Env    Environment `default:"dev"`
	Logger LoggerConfig
	Queue  QueueConfig
	Web    Web
}

// ReadConfig reads configuration from the environment and populates the structure with it
func ReadConfig() (*ApplicationConfig, error) {
	var conf ApplicationConfig
	if err := envconfig.Process("", &conf); err != nil {
		return nil, err
	}
	fmt.Printf("conf: %+v\n", conf)
	return &conf, nil
}
