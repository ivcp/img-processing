package main

import (
	"net/http"
)

func (app *application) logError(r *http.Request, err error) {
	app.logger.Print(err)
}

func (app *application) errorJSONResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "the server encountered a problem and could not process your request"
	app.errorJSONResponse(w, r, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorJSONResponse(w, r, http.StatusNotFound, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorJSONResponse(w, r, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorJSONResponse(w, r, http.StatusUnprocessableEntity, errors)
}

func (app *application) rateLimitExcededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorJSONResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *application) cannotVoteResponse(w http.ResponseWriter, r *http.Request) {
	message := "you have already voted on this poll"
	app.errorJSONResponse(w, r, http.StatusForbidden, message)
}
