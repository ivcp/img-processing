package data

import (
	"context"
	"time"

	"github.com/ivcp/polls/internal/validator"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Poll struct {
	ID          int          `json:"id"`
	Question    string       `json:"question"`
	Description string       `json:"description"`
	Options     []PollOption `json:"options"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	ExpiresAt   time.Time    `json:"expires_at"`
	Version     int          `json:"version"`
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
	var options []string
	for _, opt := range poll.Options {
		options = append(options, opt.Value)
	}
	v.Check(validator.Unique(options), "options", "must not contain duplicate values")
	v.Check(!poll.ExpiresAt.IsZero(), "expires_at", "must be provided")
	v.Check(poll.ExpiresAt.After(time.Now().Add(time.Minute)), "expires_at", "must be more than a minute in the future")
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

	args := []any{poll.Question, poll.Description, poll.ExpiresAt}

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
	return nil, nil
}

func (p PollModel) Update(poll *Poll) error {
	return nil
}

func (p PollModel) Delete(id int) error {
	return nil
}

// mocks
type MockPollModel struct {
	DB *pgxpool.Pool
}

func (p MockPollModel) Insert(poll *Poll) error {
	return nil
}

func (p MockPollModel) Get(id int) (*Poll, error) {
	return nil, nil
}

func (p MockPollModel) Update(poll *Poll) error {
	return nil
}

func (p MockPollModel) Delete(id int) error {
	return nil
}
