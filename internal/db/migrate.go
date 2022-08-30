package db

import (
	"database/sql"
	"embed"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"gorm.io/gorm"
)

//go:embed migrations/*.sql
var fs embed.FS

func Migrate(db *sql.DB, fsDriver source.Driver, conf *config.DbConfig) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithInstance("iofs", fsDriver, conf.DbName, driver)
	if err != nil {
		return err
	}
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

func ConnectAndMigrate(conf *config.DbConfig) (*gorm.DB, error) {
	gormDb, err := Connect(conf)
	if err != nil {
		return nil, err
	}
	sqlDb, err := gormDb.DB()
	if err != nil {
		return nil, err
	}
	fsDriver, err := iofs.New(fs, "migrations")
	if err != nil {
		return nil, err
	}
	err = Migrate(sqlDb, fsDriver, conf)
	if err != nil {
		return nil, err
	}
	return gormDb, nil
}
