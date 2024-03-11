package main

import (
	"expvar"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()

	mux.Use(app.metrics)
	mux.Use(middleware.Recoverer)
	mux.Use(app.enableCORS)
	mux.NotFound(app.notFoundResponse)

	mux.Group(func(mux chi.Router) {
		mux.Use(app.rateLimit)
		mux.Get("/v1/healthcheck", app.healthcheckHandler)
		mux.Post("/v1/polls", app.createPollHandler)
		mux.Get("/v1/polls", app.listPollsHandler)
		mux.Get("/v1/polls/{pollID}", app.showPollHandler)
		mux.Get("/v1/polls/{pollID}/results", app.showResultsHandler)
		mux.Post("/v1/polls/{pollID}/options/{optionID}", app.voteOptionHandler)

		mux.Group(func(mux chi.Router) {
			mux.Use(app.requireToken)
			mux.Delete("/v1/polls/{pollID}", app.deletePollHandler)
			mux.Group(func(mux chi.Router) {
				mux.Use(app.checkPollExpired)
				mux.Patch("/v1/polls/{pollID}", app.updatePollHandler)
				mux.Post("/v1/polls/{pollID}/options", app.addOptionHandler)
				mux.Patch("/v1/polls/{pollID}/options/{optionID}", app.updateOptionValueHandler)
				mux.Patch("/v1/polls/{pollID}/options", app.updateOptionPositionHandler)
				mux.Delete("/v1/polls/{pollID}/options/{optionID}", app.deleteOptionHandler)
			})
		})
	})

	mux.Method(http.MethodGet, "/v1/metrics", expvar.Handler())

	return mux
}
