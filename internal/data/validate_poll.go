package data

import (
	"time"

	"github.com/ivcp/polls/internal/validator"
)

func ValidatePoll(v *validator.Validator, poll *Poll) {
	v.Check(poll.Question != "", "question", "must not be empty")
	v.Check(len(poll.Question) <= 500, "question", "must not be more than 500 bytes long")
	v.Check(len(poll.Description) <= 1000, "description", "must not be more than 1000 bytes long")
	v.Check(poll.Options != nil, "options", "must be provided")
	v.Check(len(poll.Options) >= 2, "options", "must contain at least two options")
	var optValues []string
	var optPositions []int
	for _, opt := range poll.Options {
		optValues = append(optValues, opt.Value)
		optPositions = append(optPositions, opt.Position)
	}
	v.Check(validator.Unique(optValues), "options", "must not contain duplicate values")
	v.Check(validator.Unique(optPositions), "options", "positions must be unique")
	for _, o := range optValues {
		v.Check(o != "", "options", "option values must not be empty")
		v.Check(len(o) <= 500, "options", "option value must not be more than 500 bytes long")
	}
	for _, p := range optPositions {
		v.Check(p >= 0, "options", "position must be greater or equal to 0")
		v.Check(p <= len(poll.Options)-1, "options", "position must not excede the number of options")
	}
	if !poll.ExpiresAt.IsZero() {
		v.Check(poll.ExpiresAt.After(time.Now().Add(time.Minute)), "expires_at", "must be more than a minute in the future")
	}
}
