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

// mocks
type MockPollOptionModel struct {
	DB *pgxpool.Pool
}

func (p MockPollOptionModel) Insert(option *PollOption, pollID int) error {
	return nil
}
