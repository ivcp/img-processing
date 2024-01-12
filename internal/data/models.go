package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Polls       Polls
	PollOptions PollOptions
}

type Polls interface {
	Insert(poll *Poll) error
	Get(id int) (*Poll, error)
	Update(poll *Poll) error
	Delete(id int) error
}
type PollOptions interface {
	Insert(option *PollOption, pollID int) error
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Polls:       PollModel{DB: db},
		PollOptions: PollOptionModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Polls:       MockPollModel{},
		PollOptions: MockPollOptionModel{},
	}
}
