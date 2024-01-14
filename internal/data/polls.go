package data

import (
	"context"
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
}

type PollModel struct {
	DB *pgxpool.Pool
}

func (p PollModel) Insert(poll *Poll) error {
	query := `
		INSERT INTO polls (question, description, expires_at)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;				
		`

	args := []any{poll.Question, poll.Description, poll.ExpiresAt.Time}

	err := p.DB.QueryRow(
		context.Background(), query, args...,
	).Scan(&poll.ID, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert poll: %w", err)
	}

	rows := make([][]any, 0, len(poll.Options))

	for _, opt := range poll.Options {
		rows = append(rows, []any{opt.Value, poll.ID, opt.Position, opt.VoteCount})
	}

	_, err = p.DB.CopyFrom(
		context.Background(),
		pgx.Identifier{"poll_options"},
		[]string{"value", "poll_id", "position", "vote_count"},
		pgx.CopyFromRows(rows),
	)

	if err != nil {
		return fmt.Errorf("insert poll options: %w", err)
	}

	queryOptions := `
			SELECT id, value, position, vote_count
			FROM poll_options
			WHERE poll_id = $1;
		`

	options := make([]*PollOption, 0, len(poll.Options))

	rowsOpts, err := p.DB.Query(context.Background(), queryOptions, poll.ID)
	if err != nil {
		return fmt.Errorf("insert poll - get poll options: %w", err)
	}
	defer rowsOpts.Close()

	for rowsOpts.Next() {
		var pollOption PollOption
		err := rowsOpts.Scan(
			&pollOption.ID,
			&pollOption.Value,
			&pollOption.Position,
			&pollOption.VoteCount,
		)
		if err != nil {
			return fmt.Errorf("insert poll - get poll options: %w", err)
		}
		options = append(options, &pollOption)
	}

	if err := rowsOpts.Err(); err != nil {
		return fmt.Errorf("insert poll - get poll options: %w", err)
	}

	poll.Options = options

	return nil
}

func (p PollModel) Get(id int) (*Poll, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT p.id, p. question, p.description, p.created_at, 
		p.updated_at, p.expires_at, po.id, po.value, 
		po.position, po.vote_count
		FROM polls p
		JOIN poll_options po ON po.poll_id = p.id 
		WHERE p.id = $1;
	`

	rows, err := p.DB.Query(context.Background(), query, id)
	defer rows.Close()
	if err != nil {
		return nil, fmt.Errorf("get poll: %w", err)
	}

	var poll Poll
	var options []*PollOption

	first := true
	for rows.Next() {

		var option PollOption

		switch {
		case first:
			err = rows.Scan(
				&poll.ID,
				&poll.Question,
				&poll.Description,
				&poll.CreatedAt,
				&poll.UpdatedAt,
				&poll.ExpiresAt.Time,
				&option.ID,
				&option.Value,
				&option.Position,
				&option.VoteCount,
			)
		default:
			err = rows.Scan(
				nil,
				nil,
				nil,
				nil,
				nil,
				nil,
				&option.ID,
				&option.Value,
				&option.Position,
				&option.VoteCount,
			)
		}

		if err != nil {
			return nil, fmt.Errorf("get poll - scan: %w", err)
		}

		options = append(options, &option)
		first = false
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("get poll: %w", err)
	}

	if len(options) == 0 {
		return nil, ErrRecordNotFound
	}

	poll.Options = options

	return &poll, nil
}

func (p PollModel) Update(poll *Poll) error {
	queryPoll := `
		UPDATE polls
		SET question = $1, description = $2, 
		expires_at = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at;
	`

	args := []any{
		poll.Question,
		poll.Description,
		poll.ExpiresAt.Time,
		poll.ID,
	}
	return p.DB.QueryRow(context.Background(), queryPoll, args...).Scan(&poll.UpdatedAt)
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
			Options: []*PollOption{
				{ID: 1, Value: "One", Position: 0},
				{ID: 2, Value: "Two", Position: 1},
			},
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
