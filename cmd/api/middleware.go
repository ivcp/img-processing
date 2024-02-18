package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ivcp/polls/internal/data"
	"github.com/ivcp/polls/internal/validator"
	"golang.org/x/time/rate"
)

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	clients := make(map[string]*client)

	go func() {
		for {
			time.Sleep(time.Minute)
			app.mutex.Lock()
			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}
			app.mutex.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.config.limiter.enabled {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			app.mutex.Lock()

			if _, ok := clients[ip]; !ok {
				clients[ip] = &client{
					limiter: rate.NewLimiter(
						rate.Limit(app.config.limiter.rps),
						app.config.limiter.burst,
					),
				}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				app.mutex.Unlock()
				app.rateLimitExcededResponse(w, r)
				return
			}

			app.mutex.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}

func (app *application) requireToken(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("Authorization")
		if authorizationHeader == "" {
			app.invalidTokenResponse(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidTokenResponse(w, r)
			return
		}

		token := headerParts[1]

		v := validator.New()

		if data.ValidateTokenPlaintext(v, token); !v.Valid() {
			app.invalidTokenResponse(w, r)
			return
		}

		pollID, err := app.models.Polls.CheckToken(token)
		if err != nil {
			app.invalidTokenResponse(w, r)
			return
		}

		paramPollID, err := app.readIDParam(r, "pollID")
		if err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		if pollID != paramPollID {
			app.badRequestResponse(w, r, fmt.Errorf("token not valid for this poll"))
			return
		}

		ctx := context.WithValue(r.Context(), ctxPollIDKey, pollID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) checkPollExpired(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.mutex.Lock()
		defer app.mutex.Unlock()

		id := app.pollIDfromContext(r.Context())
		poll, err := app.models.Polls.Get(id)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.notFoundResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return
		}

		if !poll.ExpiresAt.Time.IsZero() && poll.ExpiresAt.Time.Before(time.Now()) {
			app.pollExpiredResponse(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), ctxPollKey, poll)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *application) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Access-Control-Request-Method")

		w.Header().Set("Access-Control-Allow-Origin", "*")

		if r.Method == http.MethodOptions &&
			r.Header.Get("Origin") != "" &&
			r.Header.Get("Access-Control-Request-Method") != "" {

			w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PATCH, DELETE")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
			w.WriteHeader(http.StatusOK)
			return

		}
		next.ServeHTTP(w, r)
	})
}
