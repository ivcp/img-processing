package main

import (
	"net/http"
	"time"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) showPollHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	poll := data.Poll{
		ID:       id,
		Question: "Test question?",
		Options: data.PollOptions{
			{ID: 1, Value: "One", Position: 0},
			{ID: 2, Value: "Two", Position: 1},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		ExpiresAt: time.Now().Add(12 * time.Hour),
		Version:   1,
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"poll": poll}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
