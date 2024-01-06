package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) updatePollHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	poll, err := app.models.Polls.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Question    *string        `json:"question"`
		Description *string        `json:"description"`
		ExpiresAt   data.ExpiresAt `json:"expires_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Question != nil {
		poll.Question = *input.Question
	}

	if input.Description != nil {
		poll.Description = *input.Description
	}

	if !input.ExpiresAt.IsZero() {
		poll.ExpiresAt = input.ExpiresAt
	}

	if input.Question == nil && input.Description == nil && input.ExpiresAt.IsZero() {
		app.badRequestResponse(w, r, errors.New("no fields provided for update"))
		return
	}

	v := validator.New()

	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Polls.Update(poll)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
