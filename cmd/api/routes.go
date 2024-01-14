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
	mux.Get("/v1/polls/{pollID}", app.showPollHandler)
	mux.Patch("/v1/polls/{pollID}", app.updatePollHandler)
	mux.Delete("/v1/polls/{pollID}", app.deletePollHandler)

	mux.Post("/v1/polls/{pollID}/options", app.addOptionHandler)
	mux.Patch("/v1/polls/{pollID}/options/{optionID}", app.updateOptionValueHandler)
	mux.Patch("/v1/polls/{pollID}/options", app.updateOptionPositionHandler)

	return mux
}
