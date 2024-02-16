package data

import (
	"errors"
	"net"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrRecordNotFound = errors.New("record not found")

const dbTimeout = time.Second * 3

type Models struct {
	Polls       Polls
	PollOptions PollOptions
}

type Polls interface {
	Insert(poll *Poll, tokenHash []byte) error
	Get(id string) (*Poll, error)
	Update(poll *Poll) error
	Delete(id string) error
	GetAll(search string, filters Filters) ([]*Poll, Metadata, error)
	GetVotedIPs(pollID string) ([]*net.IP, error)
	CheckToken(tokenPlaintext string) (string, error)
}
type PollOptions interface {
	Insert(option *PollOption, pollID string) error
	UpdateValue(option *PollOption) error
	UpdatePosition(options []*PollOption) error
	Vote(optionID string, pollID string, ip string) error
	Delete(optionID string) error
	GetResults(pollID string) ([]*PollOption, error)
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
