package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) voteOptionHandler(w http.ResponseWriter, r *http.Request) {
	pollID, err := app.readIDParam(r, "pollID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	optionID, err := app.readIDParam(r, "optionID")
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

	if !poll.ExpiresAt.Time.IsZero() && poll.ExpiresAt.Time.Before(time.Now()) {
		app.pollExpiredResponse(w, r)
		return
	}

	ip := r.Header.Get("X-Forwarded-For")
	if ip == "" {
		app.serverErrorResponse(w, r, errors.New("no ip found"))
		return
	}

	app.mutex.Lock()
	voted, err := app.checkIP(r, poll.ID, ip)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		app.mutex.Unlock()
		return
	}
	if voted {
		app.cannotVoteResponse(w, r)
		app.mutex.Unlock()
		return
	}

	err = app.models.PollOptions.Vote(optionID, poll.ID, ip)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		app.mutex.Unlock()
		return
	}

	app.mutex.Unlock()

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "vote successful"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
