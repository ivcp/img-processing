package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) updatePollHandler(w http.ResponseWriter, r *http.Request) {
	poll := app.pollFromContext(r.Context())

	var input struct {
		Question    *string        `json:"question"`
		Description *string        `json:"description"`
		ExpiresAt   data.ExpiresAt `json:"expires_at"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	if input.Question != nil {
		poll.Question = strings.TrimSpace(*input.Question)
	}

	if input.Description != nil {
		poll.Description = strings.TrimSpace(*input.Description)
	}

	if !input.ExpiresAt.IsZero() {
		poll.ExpiresAt = input.ExpiresAt
	}

	if input.Question == nil && input.Description == nil && input.ExpiresAt.IsZero() {
		app.badRequestResponse(w, errors.New("no fields provided for update"))
		return
	}

	v := validator.New()

	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.Polls.Update(poll)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
