package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) updateOptionPositionHandler(w http.ResponseWriter, r *http.Request) {
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

	var input struct {
		Options []struct {
			Id       int `json:"id"`
			Position int `json:"position"`
		} `json:"options"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	optMap := make(map[int]int, len(input.Options))

	for _, inputOpt := range input.Options {
		optMap[inputOpt.Id] = inputOpt.Position
	}

	var optionsToUpdate []*data.PollOption

	for i, option := range poll.Options {
		if position, ok := optMap[option.ID]; ok {
			poll.Options[i].Position = position
			optionsToUpdate = append(optionsToUpdate, poll.Options[i])
		}
	}

	v := validator.New()

	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.PollOptions.UpdatePosition(optionsToUpdate)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusCreated, envelope{"message": "options updated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
