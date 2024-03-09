package main

import (
	"net/http"
	"strings"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) updateOptionValueHandler(w http.ResponseWriter, r *http.Request) {
	poll := app.pollFromContext(r.Context())

	var input struct {
		Value string `json:"value"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	optionID, err := app.readIDParam(r, "optionID")
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	var optionToUpdate *data.PollOption
	match := false

	for _, opt := range poll.Options {
		if opt.ID == optionID {
			opt.Value = strings.TrimSpace(input.Value)
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
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.PollOptions.UpdateValue(optionToUpdate)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "option updated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
