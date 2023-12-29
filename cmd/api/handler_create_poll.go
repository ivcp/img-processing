package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) createPollHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Question    string    `json:"question"`
		Description string    `json:"description"`
		Options     []string  `json:"options"`
		ExpiresAt   time.Time `json:"expires_at"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	options := []data.PollOption{}
	for _, option := range input.Options {
		options = append(options, data.PollOption{Value: option})
	}

	poll := &data.Poll{
		Question:    input.Question,
		Description: input.Description,
		Options:     options,
		ExpiresAt:   input.ExpiresAt,
	}

	v := validator.New()
	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
