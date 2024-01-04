package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_showPollHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{"valid id", "1", http.StatusOK, `"question":"Test?"`},
		{"invalid id", "-2", http.StatusBadRequest, `invalid id`},
		{"invalid id", "a", http.StatusBadRequest, `invalid id`},
		{"no record found", "44", http.StatusNotFound, `the requested resource could not be found`},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("id", test.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.showPollHandler)
			handler.ServeHTTP(rr, req)

			if rr.Code != test.expectedStatus {
				t.Errorf("expected status code %d, but got %d", test.expectedStatus, rr.Code)
			}

			if !strings.Contains(rr.Body.String(), test.expectedBody) {
				t.Errorf("expected body to contain %q, but got %q", test.expectedBody, rr.Body)
			}
		})
	}
}
