package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Poll struct {
	ID          int           `json:"id"`
	Question    string        `json:"question"`
	Description string        `json:"description"`
	Options     []*PollOption `json:"options"`
	CreatedAt   time.Time     `json:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at"`
	ExpiresAt   ExpiresAt     `json:"expires_at"`
	Version     int           `json:"version"`
}

type PollModel struct {
	DB *pgxpool.Pool
}

func (p PollModel) Insert(poll *Poll) error {
	queryPoll := `
		INSERT INTO polls (question, description, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at, version;				
		`

	args := []any{poll.Question, poll.Description, poll.ExpiresAt.Time}

	return p.DB.QueryRow(
		context.Background(), queryPoll, args...,
	).Scan(&poll.ID, &poll.CreatedAt, &poll.UpdatedAt, &poll.Version)
}

func (p PollModel) Get(id int) (*Poll, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT id, question, description, created_at, 
		updated_at, expires_at, version		
		FROM polls 	
		WHERE id = $1;
	`
	var poll Poll

	err := p.DB.QueryRow(context.Background(), query, id).Scan(
		&poll.ID,
		&poll.Question,
		&poll.Description,
		&poll.CreatedAt,
		&poll.UpdatedAt,
		&poll.ExpiresAt.Time,
		&poll.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, fmt.Errorf("get poll: %w", err)
		}
	}

	return &poll, nil
}

func (p PollModel) Update(poll *Poll) error {
	queryPoll := `
		UPDATE polls
		SET question = $1, description = $2, expires_at = $3, version = version + 1
		WHERE id = $4
		RETURNING version;
	`

	args := []any{
		poll.Question,
		poll.Description,
		poll.ExpiresAt.Time,
		poll.ID,
	}
	return p.DB.QueryRow(context.Background(), queryPoll, args...).Scan(&poll.Version)
}

func (p PollModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM polls
		WHERE id = $1;
	`

	result, err := p.DB.Exec(context.Background(), query, id)
	if err != nil {
		return fmt.Errorf("delete poll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

// mocks
type MockPollModel struct {
	DB *pgxpool.Pool
}

func (p MockPollModel) Insert(poll *Poll) error {
	poll.ID = 1
	return nil
}

func (p MockPollModel) Get(id int) (*Poll, error) {
	if id == 1 {
		poll := Poll{
			ID:        1,
			Question:  "Test?",
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
			ExpiresAt: ExpiresAt{time.Now().Add(2 * time.Minute)},
			Version:   1,
		}
		return &poll, nil
	}
	return nil, ErrRecordNotFound
}

func (p MockPollModel) Update(poll *Poll) error {
	if poll.ID == 1 {
		return nil
	}
	return ErrRecordNotFound
}

func (p MockPollModel) Delete(id int) error {
	if id == 1 {
		return nil
	}
	return ErrRecordNotFound
}
