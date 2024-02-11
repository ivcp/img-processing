package main

import (
	"errors"
	"net/http"

	"github.com/ivcp/polls/internal/data"
)

func (app *application) showResultsHandler(w http.ResponseWriter, r *http.Request) {
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

	options, err := app.models.PollOptions.GetResults(pollID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	type result struct {
		ID        int    `json:"id"`
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
		app.serverErrorResponse(w, r, err)
	}
}
