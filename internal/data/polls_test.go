package data

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose"
)

var (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbName   = "polls_test"
	port     = "5434"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable"
)

var (
	resource *dockertest.Resource
	pool     *dockertest.Pool
	testDB   *pgxpool.Pool
)

func TestMain(m *testing.M) {
	endpoint := os.Getenv("DOCKER_TEST")
	p, err := dockertest.NewPool(endpoint)
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	pool = p

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {{HostIP: "0.0.0.0", HostPort: port}},
		},
	}

	resource, err := pool.RunWithOptions(&opts)
	if err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("could not connect to docker: %s", err)
	}

	if err := pool.Retry(func() error {
		var err error
		testDB, err = pgxpool.New(context.Background(), fmt.Sprintf(dsn, host, port, user, password, dbName))
		if err != nil {
			log.Println("Error: ", err)
			return err
		}
		return testDB.Ping(context.Background())
	}); err != nil {
		_ = pool.Purge(resource)
		log.Fatalf("something went wrong: %s", err)
	}

	if err = runMigrations(); err != nil {
		log.Fatalf("something went wrong: %s", err)
	}

	code := m.Run()
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}
	os.Exit(code)
}

func runMigrations() error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("createTables: %w", err)
	}

	db := stdlib.OpenDBFromPool(testDB)

	if err := goose.Up(db, "../../migrations"); err != nil {
		return fmt.Errorf("createTables: %w", err)
	}

	return nil
}

func Test_pingDB(t *testing.T) {
	err := testDB.Ping(context.Background())
	if err != nil {
		t.Error("can't ping DB")
	}
}
