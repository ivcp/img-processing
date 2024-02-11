package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) updateOptionValueHandler(w http.ResponseWriter, r *http.Request) {
	pollID := r.Context().Value("pollID").(int)

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
		Value string `json:"value"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	optionID, err := app.readIDParam(r, "optionID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	var optionToUpdate *data.PollOption
	match := false

	for _, opt := range poll.Options {
		if opt.ID == optionID {
			opt.Value = input.Value
			optionToUpdate = opt
			match = true
		}
	}

	if !match {
		app.notFoundResponse(w, r)
		return
	}

	v := validator.New()

	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.PollOptions.UpdateValue(optionToUpdate)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "option updated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
