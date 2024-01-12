package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) createPollOptionHandler(w http.ResponseWriter, r *http.Request) {
	pollID, err := app.readIDParam(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	poll, err := app.models.Polls.Get(pollID)
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
		Value    string `json:"value"`
		Position int    `json:"position"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	newOption := &data.PollOption{
		Value:    input.Value,
		Position: input.Position,
	}

	poll.Options = append(poll.Options, newOption)

	v := validator.New()

	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.PollOptions.Insert(newOption, pollID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "option added successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
