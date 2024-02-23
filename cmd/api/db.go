package main

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
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
