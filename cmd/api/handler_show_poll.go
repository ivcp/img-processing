package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) showPollHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r, "pollID")
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	poll, err := app.models.Polls.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
