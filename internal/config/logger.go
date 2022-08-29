package config

// LoggerConfig represents configuration of the logger
type LoggerConfig struct {
	Level LoggerLevel `default:"info"` //Level of minimum logging
}
