package config

type DbConfig struct {
	Host     string `default:"localhost"`
	Port     string `default:"5432"`
	User     string `default:"postgres"`
	Password string `default:"postgres"`
	DbName   string `default:"storage-and-update" envconfig:"DB_NAME"`
}
