package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func (app *application) connectToDB() (*pgxpool.Pool, error) {
	cfg, err := dbConfig(app.config.db.dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to create DB config: %w", err)
	}

	connPoll, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the DB: %w", err)
	}

	err = connPoll.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("Failed to ping the DB: %w", err)
	}

	app.logger.Println("Connected to DB!")

	return connPoll, nil
}

func dbConfig(dsn string) (*pgxpool.Config, error) {
	const defaultConnectTimeout = time.Second * 5

	dbConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}
	dbConfig.ConnConfig.ConnectTimeout = defaultConnectTimeout

	return dbConfig, nil
}
