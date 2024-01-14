package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) deleteOptionHandler(w http.ResponseWriter, r *http.Request) {
	pollID, err := app.readIDParam(r, "pollID")
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

	optionID, err := app.readIDParam(r, "optionID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var newOptions []*data.PollOption
	var optionToDelete *data.PollOption

	match := false
	for _, option := range poll.Options {
		switch option.ID == optionID {
		case true:
			optionToDelete = option
			match = true
		default:
			newOptions = append(newOptions, option)
		}
	}

	if !match {
		app.notFoundResponse(w, r)
		return
	}

	for _, opt := range newOptions {
		if opt.Position > optionToDelete.Position {
			opt.Position -= 1
		}
	}

	poll.Options = newOptions

	v := validator.New()
	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.PollOptions.Delete(optionID, poll.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.PollOptions.UpdatePosition(poll.Options)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "option deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
