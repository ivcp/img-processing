//go:build integration

package data

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

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
	resource   *dockertest.Resource
	pool       *dockertest.Pool
	testDB     *pgxpool.Pool
	testModels Models
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

	testModels = NewModels(testDB)

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

func TestPollsInsert(t *testing.T) {
	poll := Poll{
		Question: "Test?",
		Options: []*PollOption{
			{Value: "One", Position: 0},
			{Value: "Two", Position: 1},
			{Value: "Three", Position: 2},
		},
	}

	if err := testModels.Polls.Insert(&poll); err != nil {
		t.Errorf("insert poll returned an error: %s", err)
	}

	if poll.ID != 1 {
		t.Errorf("expected id to be 1 but got %d", poll.ID)
	}

	if poll.CreatedAt.IsZero() || poll.UpdatedAt.IsZero() {
		t.Errorf("expected created and updated not to be zero values")
	}

	for _, opt := range poll.Options {
		if opt.ID == 0 {
			t.Errorf("expected option id not to be zero: %s %d", opt.Value, opt.ID)
		}
	}
}

func TestPollsGet(t *testing.T) {
	poll, err := testModels.Polls.Get(1)
	if err != nil {
		t.Errorf("get poll returned an error: %s", err)
	}

	if poll.Question != "Test?" {
		t.Errorf("get poll returned wrong question: expected 'Test?' but got %s", poll.Question)
	}

	_, err = testModels.Polls.Get(9)
	if !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("expected error on non-existent poll")
	}

	_, err = testModels.Polls.Get(0)
	if !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("expected error on bad poll id")
	}
}

func TestPollsUpdate(t *testing.T) {
	newQuestion := "Is this a test?"
	newDescription := "Test description."
	newExpires := ExpiresAt{time.Now().Add(10 * time.Minute)}

	poll, _ := testModels.Polls.Get(1)

	oldUpdatedAt := poll.UpdatedAt

	poll.Question = newQuestion
	poll.Description = newDescription
	poll.ExpiresAt = newExpires

	time.Sleep(1 * time.Second)
	if err := testModels.Polls.Update(poll); err != nil {
		t.Errorf("update poll returned an error: %s", err)
	}

	poll, _ = testModels.Polls.Get(1)

	if poll.Question != newQuestion {
		t.Errorf("expected question to be %s, but got %s", newQuestion, poll.Question)
	}
	if poll.Description != newDescription {
		t.Errorf("expected description to be %s, but got %s", newDescription, poll.Description)
	}
	if poll.ExpiresAt.IsZero() {
		t.Errorf("expected expires at not to be zero value")
	}

	if poll.UpdatedAt.Equal(oldUpdatedAt) {
		t.Errorf("expected updated at to be changed")
	}
}

func TestPollsDelete(t *testing.T) {
	if err := testModels.Polls.Delete(10); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("expected error on non-existent poll")
	}
	if err := testModels.Polls.Delete(0); !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("expected error on bad poll id")
	}

	if err := testModels.Polls.Delete(1); err != nil {
		t.Errorf("delete poll returned an error: %s", err)
	}
	_, err := testModels.Polls.Get(1)
	if !errors.Is(err, ErrRecordNotFound) {
		t.Errorf("expected error on getting deleted poll")
	}
}
