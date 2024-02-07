package main

import (
	"errors"
	"net"
	"net/http"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) voteOptionHandler(w http.ResponseWriter, r *http.Request) {
	pollID, err := app.readIDParam(r, "pollID")
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	_, err = app.models.Polls.Get(pollID)
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

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	ips, err := app.models.Polls.GetVotedIPs(pollID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	for _, storedIP := range ips {
		if storedIP.Equal(net.ParseIP(ip)) {
			app.cannotVoteResponse(w, r)
			return
		}
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
		return
	}
}
