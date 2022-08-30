package docker

import (
	"context"
	"github.com/pkg/errors"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
)

// PostgresContainer represents a structure with functionality to start and stop the testcontainer with Postgres DB.
// At the start, it populates environment settings with the DB parameters and you can just 'ReadConfig' to retrieve them.
// It is supposed to use in the TestMain function - to start the container before tests and stop after.
type PostgresContainer struct {
	container testcontainers.Container
}

// Start starts container and populates db environment variables.
func (pc *PostgresContainer) Start(ctx context.Context) error {
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
		return errors.Errorf("couldn't create docker container: %s", err)
	}
	pc.container = container

	ip, err := container.Host(ctx)
	if err != nil {
		return errors.Errorf("couldn't get container host: %s", ip)
	}

	mappedPort, err := container.MappedPort(ctx, "5432")
	if err != nil {
		return errors.Errorf("couldn't get mapped port: %s", mappedPort.Port())
	}

	if err := pc.setDbEnv(pgUser, pgPasswd, pgDBName, ip, mappedPort.Port()); err != nil {
		return errors.Errorf("couldn't set environment: %s", err)
	}
	return nil
}

func (pc *PostgresContainer) setDbEnv(pgUser string, pgPasswd string, pgName string, pgHost string, pgPort string) error {
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

func (pc *PostgresContainer) Stop(ctx context.Context) error {
	if pc == nil {
		return errors.New("container hasn't been started")
	}
	return pc.container.Terminate(ctx)
}
