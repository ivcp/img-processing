package main

import (
	"net/http"
)

func (app *application) logError(err error) {
	app.logger.Print(err)
}

func (app *application) errorJSONResponse(w http.ResponseWriter, status int, message any) {
	env := envelope{"error": message}

	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(err)
		w.WriteHeader(500)
	}
}

func (app *application) serverErrorResponse(w http.ResponseWriter, err error) {
	app.logError(err)
	message := "the server encountered a problem and could not process your request"
	app.errorJSONResponse(w, http.StatusInternalServerError, message)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "the requested resource could not be found"
	app.errorJSONResponse(w, http.StatusNotFound, message)
}

func (app *application) badRequestResponse(w http.ResponseWriter, err error) {
	app.errorJSONResponse(w, http.StatusBadRequest, err.Error())
}

func (app *application) failedValidationResponse(w http.ResponseWriter, errors map[string]string) {
	app.errorJSONResponse(w, http.StatusUnprocessableEntity, errors)
}

func (app *application) rateLimitExcededResponse(w http.ResponseWriter) {
	message := "rate limit exceeded"
	app.errorJSONResponse(w, http.StatusTooManyRequests, message)
}

func (app *application) cannotVoteResponse(w http.ResponseWriter) {
	message := "you have already voted on this poll"
	app.errorJSONResponse(w, http.StatusForbidden, message)
}

func (app *application) cannotEditResponse(w http.ResponseWriter) {
	message := "editing the poll is not permitted once voting has begun"
	app.errorJSONResponse(w, http.StatusForbidden, message)
}

func (app *application) pollExpiredResponse(w http.ResponseWriter) {
	message := "poll has expired"
	app.errorJSONResponse(w, http.StatusForbidden, message)
}

func (app *application) cannotShowResultsResponse(w http.ResponseWriter, msg string) {
	message := "results will be available " + msg
	app.errorJSONResponse(w, http.StatusForbidden, message)
}

func (app *application) invalidTokenResponse(w http.ResponseWriter) {
	w.Header().Set("WWW-Authenticate", "Bearer")
	message := "invalid or missing token"
	app.errorJSONResponse(w, http.StatusUnauthorized, message)
}
