package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
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
			pollID:         "1",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "invalid poll id",
			pollID:         "99",
			expectedStatus: http.StatusNotFound,
		},
		{
			name:           "don't show results bofore voting",
			pollID:         "35",
			ip:             "10.10.10.10:10",
			expectedStatus: http.StatusForbidden,
		},
		{
			name:           "show results after voting",
			pollID:         "35",
			ip:             "0.0.0.1:0",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "don't show results before deadline",
			pollID:         "36",
			ip:             "0.0.0.1:0",
			expectedStatus: http.StatusForbidden,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", test.pollID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			req.RemoteAddr = test.ip
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.showResultsHandler)
			handler.ServeHTTP(rr, req)
			if rr.Code != test.expectedStatus {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, rr.Code)
			}
		})
	}
}
