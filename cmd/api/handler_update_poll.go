package main

import (
	"errors"
	"net/http"
	"time"

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

	type option struct {
		Value    string `json:"value"`
		Position int    `json:"position"`
	}

	var input struct {
		Question    string         `json:"question"`
		Description string         `json:"description"`
		EditOptions map[int]option `json:"edit_options"`
		AddOptions  []option       `json:"add_options"`
		ExpiresAt   time.Time      `json:"expires_at"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.EditOptions != nil {
		if len(input.EditOptions) != len(poll.Options) {
			app.badRequestResponse(w, r, errors.New("wrong num of options provided"))
			return
		}
		for i, option := range poll.Options {
			inputOpt, ok := input.EditOptions[option.ID]
			if ok {
				poll.Options[i].Value = inputOpt.Value
				poll.Options[i].Position = inputOpt.Position
			}
			if !ok {
				app.badRequestResponse(w, r, errors.New("unexisting option ID provided"))
				return
			}
		}
	}

	for _, option := range input.AddOptions {
		poll.Options = append(poll.Options, &data.PollOption{
			Value:    option.Value,
			Position: option.Position,
		})
	}

	poll.Question = input.Question
	poll.Description = input.Description
	poll.ExpiresAt = input.ExpiresAt

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
