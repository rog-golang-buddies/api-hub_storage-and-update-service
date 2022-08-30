package db

import (
	"context"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/test/docker"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	pgC := new(docker.PostgresContainer)
	err := pgC.Start(ctx)
	if err != nil {
		log.Fatalf("error while starting container: %s", err)
	}
	code := m.Run()
	err = pgC.Stop(ctx)
	if err != nil {
		log.Fatalf("error while starting container: %s", err)
	}
	os.Exit(code)
}

func TestConnectToDb(t *testing.T) {
	conf, err := config.ReadConfig() //read configuration from env
	if err != nil {
		t.Error("error while reading configuration")
	}
	gormDb, err := Connect(&conf.DB)
	assert.Nil(t, err)
	sqlDb, err := gormDb.DB()
	assert.Nil(t, err)
	assert.NotNil(t, sqlDb)
}

func TestConnectAndMigrate(t *testing.T) {
	conf, err := config.ReadConfig() //read configuration from env
	if err != nil {
		t.Error("error while reading configuration")
	}
	gormDb, err := ConnectAndMigrate(&conf.DB)
	assert.Nil(t, err)
	assert.NotNil(t, gormDb)
}
