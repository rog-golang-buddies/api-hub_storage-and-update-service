package apispecdoc

import (
	"context"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/db"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/test/docker"
	"gorm.io/gorm"
	"log"
	"os"
	"testing"
)

var gDb *gorm.DB

func TestMain(m *testing.M) {
	ctx := context.Background()
	pgC := new(docker.PostgresContainer)
	err := pgC.Start(ctx)
	if err != nil {
		log.Fatalf("error while starting container: %s", err)
	}
	gDb, err = applyMigrations()
	if err != nil {
		log.Fatalf("error while connecting to db: %s", err)
	}
	code := m.Run()
	err = pgC.Stop(ctx)
	if err != nil {
		log.Fatalf("error while starting container: %s", err)
	}
	os.Exit(code)
}

func applyMigrations() (*gorm.DB, error) {
	conf, err := config.ReadConfig()
	gormDb, err := db.Connect(&conf.DB)
	if err != nil {
		return nil, err
	}
	sqlDb, err := gormDb.DB()
	if err != nil {
		return nil, err
	}
	fsDriver, err := iofs.New(db.FS, "migrations")
	if err != nil {
		return nil, err
	}
	err = db.Migrate(sqlDb, fsDriver, &conf.DB)
	if err != nil {
		return nil, err
	}
	return gormDb, nil
}
