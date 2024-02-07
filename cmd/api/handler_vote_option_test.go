package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_voteOptionHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		ip             string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid vote",
			id:             "1",
			ip:             "0.0.0.0:0",
			expectedStatus: http.StatusOK,
			expectedBody:   "vote successful",
		},
		{
			name:           "ip already voted",
			id:             "1",
			ip:             "0.0.0.1:0",
			expectedStatus: http.StatusForbidden,
			expectedBody:   "you have already voted on this poll",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", "1")
			chiCtx.URLParams.Add("optionID", test.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			req.RemoteAddr = test.ip
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
