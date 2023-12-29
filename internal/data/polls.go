package data

import (
	"time"

	"github.com/ivcp/polls/internal/validator"
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
	v.Check(validator.Unique(poll.Options), "options", "must not contain duplicate values")
	v.Check(!poll.ExpiresAt.IsZero(), "expires_at", "must be provided")
	v.Check(poll.ExpiresAt.After(time.Now().Add(time.Minute)), "expires_at", "must be more than a minute in the future")
}
