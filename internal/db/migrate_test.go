package db

import (
	"context"
	"github.com/rog-golang-buddies/api-hub_storage-and-update-service/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	container := setupPostgres(ctx)
	code := m.Run()
	err := container.Terminate(ctx)
	if err != nil {
		log.Fatalf("error while stopping container: %s", err)
	}
	os.Exit(code)
}

func setupPostgres(ctx context.Context) testcontainers.Container {
	pgUser := "test_pg_user"
	pgPasswd := "test_pg_password"
	pgDBName := "test_db"
	env := map[string]string{
		"POSTGRES_PASSWORD": pgPasswd,
		"POSTGRES_USER":     pgUser,
		"POSTGRES_DB":       pgDBName,
	}
	req := testcontainers.ContainerRequest{
		Image:        "postgres:14.5",
		ExposedPorts: []string{"5432/tcp"},
		Env:          env,
		WaitingFor:   wait.ForListeningPort("5432"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatalf("couldn't create docker container: %s", err)
	}

	ip, err := container.Host(ctx)
	if err != nil {
		log.Fatalf("couldn't get container host: %s", ip)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		log.Fatalf("couldn't get mapped port: %s", mappedPort.Port())
	}

	if err := setDbEnv(pgUser, pgPasswd, pgDBName, ip, mappedPort.Port()); err != nil {
		log.Fatalf("couldn't set environment: %s", err)
	}
	return container
}

func setDbEnv(pgUser string, pgPasswd string, pgName string, pgHost string, pgPort string) error {
	if err := os.Setenv("DB_USER", pgUser); err != nil {
		return err
	}
	if err := os.Setenv("DB_PASSWORD", pgPasswd); err != nil {
		return err
	}
	if err := os.Setenv("DB_NAME", pgName); err != nil {
		return err
	}
	if err := os.Setenv("DB_HOST", pgHost); err != nil {
		return err
	}
	if err := os.Setenv("DB_PORT", pgPort); err != nil {
		return err
	}
	return nil
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
