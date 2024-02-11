package data

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PollOption struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
	// Position of option in the list, starting at 0
	Position  int `json:"position"`
	VoteCount int `json:"-"`
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
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := p.DB.Exec(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("insert poll option: %w", err)
	}

	return p.setUpdatedAt(pollID)
}

func (p PollOptionModel) UpdateValue(option *PollOption) error {
	query := `
		UPDATE poll_options 
		SET value = $1
		WHERE id = $2
		RETURNING poll_id;	
	`

	var pollID int
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	err := p.DB.QueryRow(
		ctx, query, option.Value, option.ID,
	).Scan(&pollID)
	if err != nil {
		return fmt.Errorf("update poll option: %w", err)
	}

	return p.setUpdatedAt(pollID)
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
		ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
		defer cancel()
		err := p.DB.QueryRow(
			ctx, query, option.Position, option.ID,
		).Scan(&pollID)
		if err != nil {
			return fmt.Errorf("update option position: %w", err)
		}
	}

	return p.setUpdatedAt(pollID)
}

func (p PollOptionModel) Delete(optionID int) error {
	query := `
		DELETE FROM poll_options
		WHERE id = $1
		RETURNING poll_id;	
	`
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var pollID int
	err := p.DB.QueryRow(ctx, query, optionID).Scan(&pollID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("delete option: %w", err)
	}

	return p.setUpdatedAt(pollID)
}

func (p PollOptionModel) Vote(optionID int, ip string) error {
	query := `
		UPDATE poll_options 
		SET vote_count = vote_count + 1
		WHERE id = $1
		RETURNING poll_id;
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	var pollID int
	err := p.DB.QueryRow(ctx, query, optionID).Scan(&pollID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrRecordNotFound
		}
		return fmt.Errorf("vote option: %w", err)
	}

	var paramIP pgtype.Inet
	err = paramIP.Set(ip)
	if err != nil {
		return fmt.Errorf("vote option - set ip: %w", err)
	}
	queryIP := `
		INSERT INTO ips (ip, poll_id)
		VALUES ($1, $2); 		
	`
	_, err = p.DB.Exec(ctx, queryIP, paramIP, pollID)
	if err != nil {
		return fmt.Errorf("vote option - insert ip: %w", err)
	}

	return nil
}

func (p PollOptionModel) GetResults(pollID int) ([]*PollOption, error) {
	query := `
		SELECT id, value, position, vote_count
		FROM poll_options
		WHERE poll_id = $1;
	`

	rows, err := p.DB.Query(context.Background(), query, pollID)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("get votes for poll: %w", err)
	}
	var options []*PollOption

	for rows.Next() {
		var opt PollOption
		err := rows.Scan(
			&opt.ID,
			&opt.Value,
			&opt.Position,
			&opt.VoteCount,
		)
		if err != nil {
			return nil, fmt.Errorf("get votes for poll - scan: %w", err)
		}
		options = append(options, &opt)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get votes for poll: %w", err)
	}

	return options, nil
}

func (p PollOptionModel) setUpdatedAt(pollID int) error {
	query := `
		UPDATE polls
		SET updated_at = NOW()
		WHERE id = $1;
	`
	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err := p.DB.Exec(ctx, query, pollID)
	if err != nil {
		return fmt.Errorf("set updated_at: %w", err)
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

func (p MockPollOptionModel) Delete(optionID int) error {
	return nil
}

func (p MockPollOptionModel) Vote(optionID int, ip string) error {
	return nil
}

func (p MockPollOptionModel) GetResults(pollID int) ([]*PollOption, error) {
	return nil, nil
}
