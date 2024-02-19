package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ivcp/polls/internal/data"
)

func Test_app_rateLimit(t *testing.T) {
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	app.config.limiter.rps = 2
	app.config.limiter.burst = 4
	app.config.limiter.enabled = true

	handlerToTest := app.rateLimit(nextHandler)
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.RemoteAddr = "0.0.0.0:0000"
	rr := httptest.NewRecorder()
	for i := 0; i < 6; i++ {
		handlerToTest.ServeHTTP(rr, req)
		if i < 4 && rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
		}
		if i > 4 && rr.Code != http.StatusTooManyRequests {
			t.Errorf("expected status code %d, but got %d", http.StatusTooManyRequests, rr.Code)
		}
	}

	app.config.limiter.enabled = false
	rr = httptest.NewRecorder()
	for i := 0; i < 10; i++ {
		handlerToTest.ServeHTTP(rr, req)
		if rr.Code != http.StatusOK {
			t.Errorf("expected status code %d, but got %d", http.StatusOK, rr.Code)
		}
	}
}

func Test_app_requireToken(t *testing.T) {
	tests := []struct {
		name           string
		authHeader     string
		expectedStatus int
	}{
		{
			name:           "no auth header set",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "bad auth header",
			authHeader:     "invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "invalid token",
			authHeader:     "Bearer invalid",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "missing token",
			authHeader:     "Bearer ",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "valid token",
			authHeader:     "Bearer UBQ2Z7CLB2SJQBNTUCH4IMRI7A",
			expectedStatus: http.StatusOK,
		},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handlerToTest := app.requireToken(nextHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", data.ExamplePollIDValid)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			if test.authHeader != "" {
				req.Header.Set("Authorization", test.authHeader)
			}
			handlerToTest.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, but got %d", test.expectedStatus, rr.Code)
			}
			if test.authHeader == "" && rr.Header().Get("WWW-Authenticate") != "Bearer" {
				t.Errorf("header WWW-Authenticate not set")
			}
		})
	}
}

func Test_app_checkPollExpired(t *testing.T) {
	tests := []struct {
		name           string
		pollID         string
		expectedStatus int
	}{
		{"expired poll", data.ExamplePollIDExpiredPoll, http.StatusForbidden},
		{"valid poll", data.ExamplePollIDValid, http.StatusOK},
		{"unexisting poll", uuid.NewString(), http.StatusNotFound},
		{"expired_at not set", data.ExamplePollIDExpiredNotSet, http.StatusOK},
	}
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handlerToTest := app.checkPollExpired(nextHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxPollIDKey, test.pollID))
			rr := httptest.NewRecorder()
			handlerToTest.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, but got %d", test.expectedStatus, rr.Code)
			}
		})
	}
}

func Test_app_enableCORS(t *testing.T) {
	tests := []struct {
		name      string
		method    string
		origin    string
		reqMethod string
	}{
		{"regular req", http.MethodGet, "", ""},
		{"preflight req", http.MethodOptions, "localhost:8080", http.MethodPost},
	}

	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handlerToTest := app.enableCORS(nextHandler)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(test.method, "/", nil)
			req.Header.Set("Origin", test.origin)
			req.Header.Set("Access-Control-Request-Method", test.reqMethod)
			rr := httptest.NewRecorder()
			handlerToTest.ServeHTTP(rr, req)
			result := rr.Result()
			if result.Header.Get("Access-Control-Allow-Origin") != "*" {
				t.Errorf(
					"Access-Control-Allow-Origin header not set to '*', got %q",
					result.Header.Get("Access-Control-Allow-Origin"),
				)
			}

			if req.Method == http.MethodOptions {
				if result.Header.Get("Access-Control-Allow-Methods") != "OPTIONS, PATCH, DELETE" {
					t.Errorf(
						"Access-Control-Allow-Methods not set to OPTIONS, PATCH, DELETE, got %q",
						result.Header.Get("Access-Control-Allow-Methods"),
					)
				}
				if result.Header.Get("Access-Control-Allow-Headers") != "Authorization, Content-Type" {
					t.Errorf(
						"Access-Control-Allow-Headers not set to 'Authorization, Content-Type', got %q",
						result.Header.Get("Access-Control-Allow-Headers"),
					)
				}
			}
		})
	}
}
