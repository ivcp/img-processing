package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) showResultsHandler(w http.ResponseWriter, r *http.Request) {
	pollID, err := app.readIDParam(r, "pollID")
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	poll, err := app.models.Polls.Get(pollID)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, err)
		}
		return
	}

	switch poll.ResultsVisibility {
	case "after_vote":
		if poll.ExpiresAt.Time.Before(time.Now()) {
			ip := r.Header.Get("X-Forwarded-For")
			if ip == "" {
				app.serverErrorResponse(w, errors.New("no ip found"))
				return
			}

			voted, err := app.checkIP(pollID, ip)
			if err != nil {
				app.serverErrorResponse(w, err)
				return
			}
			if !voted {
				app.cannotShowResultsResponse(w, "after voting")
				return
			}
		}

	case "after_deadline":
		if !poll.ExpiresAt.Time.IsZero() && poll.ExpiresAt.Time.After(time.Now()) {
			app.cannotShowResultsResponse(w, "when poll expires")
			return
		}
	}

	options, err := app.models.PollOptions.GetResults(pollID)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	type result struct {
		ID        string `json:"id"`
		Value     string `json:"value"`
		Position  int    `json:"position"`
		VoteCount int    `json:"vote_count"`
	}

	results := make([]result, 0, len(options))

	for _, opt := range options {
		results = append(results, result{
			ID:        opt.ID,
			Value:     opt.Value,
			Position:  opt.Position,
			VoteCount: opt.VoteCount,
		})
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"results": results}, nil)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
