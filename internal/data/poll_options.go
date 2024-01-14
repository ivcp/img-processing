package data

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PollOption struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
	// Position of option in the list, starting at 0
	Position  int `json:"position"`
	VoteCount int `json:"vote_count"`
}

type PollOptionModel struct {
	DB *pgxpool.Pool
}

func (p PollOptionModel) Insert(option *PollOption, pollID int) error {
	query := `
		INSERT INTO poll_options (poll_id, value, position, vote_count)
		VALUES ($1, $2, $3, $4);		
	`

	args := []any{pollID, option.Value, option.Position, option.VoteCount}
	_, err := p.DB.Exec(context.Background(), query, args...)
	if err != nil {
		return fmt.Errorf("insert poll option: %w", err)
	}

	queryPoll := `
		UPDATE polls
		SET version = version + 1, updated_at = NOW()
		WHERE id = $1;
	`
	_, err = p.DB.Exec(context.Background(), queryPoll, pollID)
	if err != nil {
		return fmt.Errorf("insert poll option: %w", err)
	}

	return nil
}

func (p PollOptionModel) UpdateValue(option *PollOption) error {
	query := `
		UPDATE poll_options 
		SET value = $1
		WHERE id = $2
		RETURNING poll_id;	
	`

	var pollID int

	err := p.DB.QueryRow(
		context.Background(), query, option.Value, option.ID,
	).Scan(&pollID)
	if err != nil {
		return fmt.Errorf("update poll option: %w", err)
	}

	queryPoll := `
		UPDATE polls
		SET version = version + 1, updated_at = NOW()
		WHERE id = $1;
	`
	_, err = p.DB.Exec(context.Background(), queryPoll, pollID)
	if err != nil {
		return fmt.Errorf("insert poll option: %w", err)
	}

	return nil
}

func (p PollOptionModel) UpdatePosition(options []*PollOption) error {
	query := `
		UPDATE poll_options 
		SET position = $1
		WHERE id = $2
		RETURNING poll_id;	
	`

	var pollID int

	for _, option := range options {
		err := p.DB.QueryRow(
			context.Background(), query, option.Position, option.ID,
		).Scan(&pollID)
		if err != nil {
			return fmt.Errorf("update option position: %w", err)
		}
	}

	queryPoll := `
		UPDATE polls
		SET version = version + 1, updated_at = NOW()
		WHERE id = $1;
	`
	_, err := p.DB.Exec(context.Background(), queryPoll, pollID)
	if err != nil {
		return fmt.Errorf("insert poll option: %w", err)
	}

	return nil
}

// mocks
type MockPollOptionModel struct {
	DB *pgxpool.Pool
}

func (p MockPollOptionModel) Insert(option *PollOption, pollID int) error {
	return nil
}

func (p MockPollOptionModel) UpdateValue(option *PollOption) error {
	return nil
}

func (p MockPollOptionModel) UpdatePosition(options []*PollOption) error {
	return nil
}
