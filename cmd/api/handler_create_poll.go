package main

import (
	"fmt"
	"net/http"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) createPollHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Question    string `json:"question"`
		Description string `json:"description"`
		Options     []struct {
			Value    string `json:"value"`
			Position int    `json:"position"`
		} `json:"options"`
		ExpiresAt data.ExpiresAt `json:"expires_at"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	options := []*data.PollOption{}
	for _, option := range input.Options {
		options = append(options, &data.PollOption{Value: option.Value, Position: option.Position})
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

	err = app.models.Polls.Insert(poll)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/polls/%d", poll.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"poll": poll}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
