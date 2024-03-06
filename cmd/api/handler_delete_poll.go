package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) deletePollHandler(w http.ResponseWriter, r *http.Request) {
	id := app.pollIDfromContext(r.Context())

	err := app.models.Polls.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "poll successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
