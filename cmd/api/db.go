package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func (app *application) connectToDB() (*pgxpool.Pool, error) {
	connPoll, err := pgxpool.New(context.Background(), app.config.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to the DB: %w", err)
	}

	err = connPoll.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping the DB: %w", err)
	}

	app.logger.Println("Connected to DB!")

	return connPoll, nil
}

func (app *application) runMigrations(db *pgxpool.Pool, dir string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}
	database := stdlib.OpenDBFromPool(db)

	if err := goose.Up(database, dir); err != nil {
		return fmt.Errorf("migrate: %w", err)
	}

	return nil
}
