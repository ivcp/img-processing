package data

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/ivcp/polls/internal/validator"
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

type PollOption struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
	// Position of option in the list, starting at 0
	Position  int `json:"position"`
	VoteCount int `json:"vote_count"`
}

func ValidatePoll(v *validator.Validator, poll *Poll) {
	v.Check(poll.Question != "", "question", "must not be empty")
	v.Check(len(poll.Question) <= 500, "question", "must not be more than 500 bytes long")
	v.Check(len(poll.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	v.Check(poll.Options != nil, "options", "must be provided")
	v.Check(len(poll.Options) >= 1, "options", "must contain at least one option")
	var optValues []string
	var optPositions []int
	for _, opt := range poll.Options {
		optValues = append(optValues, opt.Value)
		optPositions = append(optPositions, opt.Position)
	}
	v.Check(validator.Unique(optValues), "options", "must not contain duplicate values")
	v.Check(validator.Unique(optPositions), "options", "positions must be unique")
	for _, p := range optPositions {
		v.Check(p >= 0, "options", "position must be greater or equal to 0")
		v.Check(p <= len(poll.Options)-1, "options", "position must not excede the number of options")
	}
	if !poll.ExpiresAt.IsZero() {
		v.Check(poll.ExpiresAt.After(time.Now().Add(time.Minute)), "expires_at", "must be more than a minute in the future")
	}
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

	queryOption := `
		INSERT INTO poll_options (value, poll_id, position, vote_count)
		VALUES ($1, $2, $3, $4)		
		RETURNING id;
	`

	args := []any{poll.Question, poll.Description, poll.ExpiresAt.Time}

	err := p.DB.QueryRow(
		context.Background(), queryPoll, args...,
	).Scan(&poll.ID, &poll.CreatedAt, &poll.UpdatedAt, &poll.Version)

	for i, option := range poll.Options {
		args := []any{option.Value, poll.ID, option.Position, option.VoteCount}
		err := p.DB.QueryRow(
			context.Background(), queryOption, args...,
		).Scan(&poll.Options[i].ID)
		if err != nil {
			return err
		}
	}

	return err
}

func (p PollModel) Get(id int) (*Poll, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	queryPoll := `
		SELECT id, question, description, created_at, 
		updated_at, expires_at, version		
		FROM polls 	
		WHERE id = $1;
	`
	queryOption := `
		SELECT id, value, position, vote_count
		FROM poll_options 
		WHERE poll_id = $1;
	`
	var poll Poll

	err := p.DB.QueryRow(context.Background(), queryPoll, id).Scan(
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

	var options []*PollOption

	rows, err := p.DB.Query(context.Background(), queryOption, id)
	if err != nil {
		return nil, err
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
			return nil, fmt.Errorf("get poll option: %w", err)
		}
		options = append(options, &pollOption)
	}

	poll.Options = options

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
			Options:   []*PollOption{{ID: 1, Value: "Yes", Position: 0, VoteCount: 0}},
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
	return nil
}
