package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
)

func (app *application) createPollHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Question    string `json:"question"`
		Description string `json:"description"`
		Options     []struct {
			Value    string `json:"value"`
			Position int    `json:"position"`
		} `json:"options"`
		ExpiresAt         data.ExpiresAt `json:"expires_at"`
		ResultsVisibility string         `json:"results_visibility"`
		IsPrivate         bool           `json:"is_private"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, err)
		return
	}

	options := []*data.PollOption{}
	for _, option := range input.Options {
		options = append(
			options,
			&data.PollOption{Value: strings.TrimSpace(option.Value), Position: option.Position},
		)
	}

	if input.ResultsVisibility == "" {
		input.ResultsVisibility = "always"
	}

	poll := &data.Poll{
		Question:          strings.TrimSpace(input.Question),
		Description:       strings.TrimSpace(input.Description),
		Options:           options,
		ExpiresAt:         input.ExpiresAt,
		ResultsVisibility: input.ResultsVisibility,
		IsPrivate:         input.IsPrivate,
	}

	v := validator.New()
	if data.ValidatePoll(v, poll); !v.Valid() {
		app.failedValidationResponse(w, v.Errors)
		return
	}

	token, err := data.GenerateToken()
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}
	poll.Token = token.Plaintext

	err = app.models.Polls.Insert(poll, token.Hash)
	if err != nil {
		app.serverErrorResponse(w, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/polls/%s", poll.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"poll": poll}, headers)
	if err != nil {
		app.serverErrorResponse(w, err)
	}
}
