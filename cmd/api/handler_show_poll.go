package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) showPollHandler(w http.ResponseWriter, r *http.Request) {
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

	pollOptions, err := app.models.PollOptions.GetAllByPollID(poll.ID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	poll.Options = pollOptions

	err = app.writeJSON(w, http.StatusOK, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
