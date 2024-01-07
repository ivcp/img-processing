package data

import (
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

type Models struct {
	Polls Polls
}

type Polls interface {
	Insert(poll *Poll) error
	Get(id int) (*Poll, error)
	Update(poll *Poll) error
	Delete(id int) error
}

func NewModels(db *pgxpool.Pool) Models {
	return Models{
		Polls: PollModel{DB: db},
	}
}

func NewMockModels() Models {
	return Models{
		Polls: MockPollModel{},
	}
}
