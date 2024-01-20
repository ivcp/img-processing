package data

import (
	"context"
	"encoding/json"
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

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	err := p.DB.QueryRow(
		ctx, query, args...,
	).Scan(&poll.ID, &poll.CreatedAt, &poll.UpdatedAt)
	if err != nil {
		return fmt.Errorf("insert poll: %w", err)
	}

	rows := make([][]any, 0, len(poll.Options))

	for _, opt := range poll.Options {
		rows = append(rows, []any{opt.Value, poll.ID, opt.Position, opt.VoteCount})
	}

	ctx, cancel = context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()
	_, err = p.DB.CopyFrom(
		ctx,
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

	ctx, cancel = context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rowsOpts, err := p.DB.Query(ctx, queryOptions, poll.ID)
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

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := p.DB.Query(ctx, query, id)
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

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	return p.DB.QueryRow(ctx, queryPoll, args...).Scan(&poll.UpdatedAt)
}

func (p PollModel) Delete(id int) error {
	if id < 1 {
		return ErrRecordNotFound
	}

	query := `
		DELETE FROM polls
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	result, err := p.DB.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("delete poll: %w", err)
	}

	if result.RowsAffected() == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func (p PollModel) GetAll(search string, filters Filters) ([]*Poll, error) {
	query := fmt.Sprintf(`
		SELECT p.id, p.question, p.description, 
		p.created_at, p.updated_at, p.expires_at,
	    jsonb_agg(jsonb_build_object(
			'id', po.id, 'value', po.value, 'position', po.position, 'vote_count', po.vote_count
			)) AS options
		FROM polls p
		JOIN poll_options po ON po.poll_id = p.id 
		WHERE (to_tsvector('simple', question) @@ plainto_tsquery('simple', $1) OR $1 = '') 
		GROUP BY p.id
		ORDER BY %s %s, id ASC;
	`, filters.sortColumn(), filters.sortDirection())

	ctx, cancel := context.WithTimeout(context.Background(), dbTimeout)
	defer cancel()

	rows, err := p.DB.Query(ctx, query, search)
	if err != nil {
		return nil, fmt.Errorf("get all polls: %w", err)
	}
	defer rows.Close()

	polls := []*Poll{}

	for rows.Next() {
		var poll Poll
		var optionsJson string
		err := rows.Scan(
			&poll.ID,
			&poll.Question,
			&poll.Description,
			&poll.CreatedAt,
			&poll.UpdatedAt,
			&poll.ExpiresAt.Time,
			&optionsJson,
		)
		if err != nil {
			return nil, fmt.Errorf("get polls - scan: %w", err)
		}

		if err := json.Unmarshal([]byte(optionsJson), &poll.Options); err != nil {
			return nil, fmt.Errorf("get polls - unmarshal options: %w", err)
		}
		polls = append(polls, &poll)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("get polls: %w", err)
	}

	return polls, nil
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
				{ID: 3, Value: "Three", Position: 2},
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

func (p MockPollModel) GetAll(search string, filters Filters) ([]*Poll, error) {
	return nil, nil
}
