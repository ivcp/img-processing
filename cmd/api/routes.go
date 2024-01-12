package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.NotFound(app.notFoundResponse)

	mux.Get("/v1/healthcheck", app.healthcheckHandler)
	mux.Post("/v1/polls", app.createPollHandler)
	mux.Get("/v1/polls/{id}", app.showPollHandler)
	mux.Patch("/v1/polls/{id}", app.updatePollHandler)
	mux.Delete("/v1/polls/{id}", app.deletePollHandler)

	mux.Post("/v1/polls/{id}/options", app.createPollOptionHandler)

	return mux
}
