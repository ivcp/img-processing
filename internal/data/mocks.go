package data

import (
	"net"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MockPollModel struct {
	DB *pgxpool.Pool
}

var (
	ExamplePollIDValid         = "e9da0ad7-6065-40de-8398-2514ce9c566f"
	ExamplePollIDExpiredPoll   = "7a818efb-b94d-49ea-af0e-5f1c8999c1b5"
	ExamplePollIDExpiredNotSet = "e4dd6db9-fa83-45d2-81dd-1f93019a25a2"
	ExamplePollIDAfterVote     = "6e3e617f-b5e6-4627-a2db-c72e29ec1729"
	ExamplePollIDAfterDeadline = "0d5edfad-ba7f-4ddc-a455-4f25ca09bfdd"
	ExampleOptionID1           = "65d7c012-f3f9-43f5-a62c-12ab516c6124"
	ExampleOptionID2           = "b85b14b5-7da6-47d0-8518-07033e199a50"
	ExampleOptionID3           = "b8168cce-4044-4c23-9506-b41915784166"
)

func (p MockPollModel) Insert(poll *Poll, tokenHash []byte) error {
	poll.ID = uuid.NewString()
	return nil
}

func (p MockPollModel) Get(id string) (*Poll, error) {
	if id == ExamplePollIDValid {
		poll := Poll{
			ID:                ExamplePollIDValid,
			Question:          "Test?",
			CreatedAt:         time.Now(),
			UpdatedAt:         time.Now(),
			ExpiresAt:         ExpiresAt{time.Now().Add(2 * time.Minute)},
			ResultsVisibility: "always",
			Options: []*PollOption{
				{ID: ExampleOptionID1, Value: "One", Position: 0},
				{ID: ExampleOptionID2, Value: "Two", Position: 1},
				{ID: ExampleOptionID3, Value: "Three", Position: 2},
			},
		}
		return &poll, nil
	}
	// expired poll
	if id == ExamplePollIDExpiredPoll {
		poll := Poll{
			ExpiresAt: ExpiresAt{time.Now().Add(-1 * time.Minute)},
		}
		return &poll, nil
	}
	// expired not set
	if id == ExamplePollIDExpiredNotSet {
		return &Poll{}, nil
	}
	// results after vote
	if id == ExamplePollIDAfterVote {
		return &Poll{
			ResultsVisibility: "after_vote",
		}, nil
	}
	// results after deadline
	if id == ExamplePollIDAfterDeadline {
		return &Poll{
			ExpiresAt:         ExpiresAt{time.Now().Add(1 * time.Minute)},
			ResultsVisibility: "after_deadline",
		}, nil
	}
	return nil, ErrRecordNotFound
}

func (p MockPollModel) Update(poll *Poll) error {
	if poll.ID == ExamplePollIDValid {
		return nil
	}
	return ErrRecordNotFound
}

func (p MockPollModel) Delete(id string) error {
	if id == ExamplePollIDValid {
		return nil
	}
	return ErrRecordNotFound
}

func (p MockPollModel) GetAll(search string, filters Filters) ([]*Poll, Metadata, error) {
	return nil, Metadata{}, nil
}

func (p MockPollModel) GetVotedIPs(pollID string) ([]*net.IP, error) {
	var ips []*net.IP
	i := net.IPv4(0, 0, 0, 1)
	ips = append(ips, &i)
	return ips, nil
}

func (p MockPollModel) CheckToken(tokenPlaintext string) (string, error) {
	return ExamplePollIDValid, nil
}

// PollOption

type MockPollOptionModel struct {
	DB *pgxpool.Pool
}

func (p MockPollOptionModel) Insert(option *PollOption, pollID string) error {
	return nil
}

func (p MockPollOptionModel) UpdateValue(option *PollOption) error {
	return nil
}

func (p MockPollOptionModel) UpdatePosition(options []*PollOption) error {
	return nil
}

func (p MockPollOptionModel) Delete(optionID string) error {
	return nil
}

func (p MockPollOptionModel) Vote(optionID string, pollID string, ip string) error {
	return nil
}

func (p MockPollOptionModel) GetResults(pollID string) ([]*PollOption, error) {
	return nil, nil
}
