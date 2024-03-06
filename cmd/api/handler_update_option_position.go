package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) updateOptionPositionHandler(w http.ResponseWriter, r *http.Request) {
	poll := app.pollFromContext(r.Context())

	var input struct {
		Options []struct {
			Id       string `json:"id"`
			Position int    `json:"position"`
		} `json:"options"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	optMap := make(map[string]int, len(input.Options))

	for _, inputOpt := range input.Options {
		optMap[inputOpt.Id] = inputOpt.Position
	}

	var optionsToUpdate []*data.PollOption

	for _, option := range poll.Options {
		if position, ok := optMap[option.ID]; ok {
			option.Position = position
			optionsToUpdate = append(optionsToUpdate, option)
		}
	}

	if len(optionsToUpdate) != len(input.Options) || len(optionsToUpdate) == 0 {
		app.badRequestResponse(w, errors.New("invalid option id, or no id provided"))
		return
	}

	v := validator.New()

	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	err = app.models.PollOptions.UpdatePosition(optionsToUpdate)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "options updated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
