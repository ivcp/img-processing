package main

import (
	"errors"
	"net"
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

	optionID, err := app.readIDParam(r, "optionID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if !poll.ExpiresAt.Time.IsZero() && poll.ExpiresAt.Time.Before(time.Now()) {
		app.pollExpiredResponse(w, r)
		return
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	voted, err := app.checkIP(r, pollID, ip)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if voted {
		app.cannotVoteResponse(w, r)
		return
	}

	err = app.models.PollOptions.Vote(optionID, ip)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "vote successful"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
