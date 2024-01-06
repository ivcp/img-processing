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

func (p PollOptionModel) GetAllByPollID(pollID int) ([]*PollOption, error) {
	query := `
		SELECT id, value, position, vote_count
		FROM poll_options 
		WHERE poll_id = $1;
	`
	var options []*PollOption

	rows, err := p.DB.Query(context.Background(), query, pollID)
	if err != nil {
		return nil, fmt.Errorf("get poll options: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var pollOption PollOption
		err := rows.Scan(
			&pollOption.ID,
			&pollOption.Value,
			&pollOption.Position,
			&pollOption.VoteCount,
		)
		if err != nil {
			return nil, fmt.Errorf("get poll options: %w", err)
		}
		options = append(options, &pollOption)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("get poll options: %w", err)
	}

	return options, nil
}

func (p PollOptionModel) Insert(pollOption *PollOption, pollID int) error {
	query := `
		INSERT INTO poll_options (value, poll_id, position, vote_count)
		VALUES ($1, $2, $3, $4)		
		RETURNING id;
	`
	args := []any{pollOption.Value, pollID, pollOption.Position, pollOption.VoteCount}
	err := p.DB.QueryRow(
		context.Background(), query, args...,
	).Scan(&pollOption.ID)
	if err != nil {
		return fmt.Errorf("insert poll option: %w", err)
	}
	return nil
}

func (p PollOptionModel) Update(pollOption *PollOption) error {
	return nil
}

func (p PollOptionModel) Delete(id int) error {
	return nil
}

// mocks
type MockPollOptionModel struct {
	DB *pgxpool.Pool
}

func (p MockPollOptionModel) GetAllByPollID(pollID int) ([]*PollOption, error) {
	return []*PollOption{
		{ID: 1, Value: "Yes", Position: 0, VoteCount: 0},
		{ID: 2, Value: "No", Position: 1, VoteCount: 0},
	}, nil
}

func (p MockPollOptionModel) Insert(pollOption *PollOption, pollID int) error {
	return nil
}

func (p MockPollOptionModel) Get(id int) (*Poll, error) {
	return nil, nil
}

func (p MockPollOptionModel) Update(pollOption *PollOption) error {
	return nil
}

func (p MockPollOptionModel) Delete(id int) error {
	return nil
}
