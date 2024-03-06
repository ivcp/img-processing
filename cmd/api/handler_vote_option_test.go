package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/ivcp/polls/internal/data"
)

func Test_app_voteOptionHandler(t *testing.T) {
	tests := []struct {
		name           string
		pollID         string
		ip             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid vote",
			pollID:         data.ExamplePollIDValid,
			ip:             "0.0.0.0",
			expectedStatus: http.StatusOK,
			expectedBody:   "vote successful",
		},
		{
			name:           "ip already voted",
			pollID:         data.ExamplePollIDValid,
			ip:             "0.0.0.1",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "you have already voted on this poll",
		},
		{
			name:           "expired poll",
			pollID:         data.ExamplePollIDExpiredPoll,
			ip:             "0.0.0.0",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "poll has expired",
		},
		{
			name:           "expired not set",
			pollID:         data.ExamplePollIDExpiredNotSet,
			ip:             "0.0.0.0",
			expectedStatus: http.StatusOK,
			expectedBody:   "vote successful",
		},
		{
			name:           "unexisting poll",
			pollID:         uuid.NewString(),
			ip:             "0.0.0.0",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "the requested resource could not be found",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", test.pollID)
			chiCtx.URLParams.Add("optionID", data.ExampleOptionID1)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			req.Header.Set("X-Forwarded-For", test.ip)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.voteOptionHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, but got %d", test.expectedStatus, rr.Code)
			}
			if !strings.Contains(rr.Body.String(), test.expectedBody) {
				t.Errorf("expected body to contain %q, but got %q", test.expectedBody, rr.Body)
			}
		})
	}
}
