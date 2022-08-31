package db

import (
	"fmt"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(conf *config.DbConfig) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		conf.Host, conf.User, conf.Password, conf.DbName, conf.Port)
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
