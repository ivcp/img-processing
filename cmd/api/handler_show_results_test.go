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

func Test_app_showResultsHandler(t *testing.T) {
	tests := []struct {
		name           string
		pollID         string
		ip             string
		expectedStatus int
	}{
		{
			name:           "show results valid",
			pollID:         data.ExamplePollIDValid,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid poll id",
			pollID:         uuid.NewString(),
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "don't show results bofore voting",
			pollID:         data.ExamplePollIDAfterVote,
			ip:             "10.10.10.10",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "show results after voting",
			pollID:         data.ExamplePollIDAfterVote,
			ip:             "0.0.0.1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "don't show results before deadline",
			pollID:         data.ExamplePollIDAfterDeadline,
			ip:             "0.0.0.1",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", test.pollID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			req.Header.Set("X-Forwarded-For", test.ip)
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.showResultsHandler)
			handler.ServeHTTP(rr, req)
			if rr.Code != test.expectedStatus {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, rr.Code)
			}
		})
	}
}
