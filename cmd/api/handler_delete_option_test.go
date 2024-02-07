package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_deleteOptionHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{"valid delete", "1", http.StatusOK, "option deleted successfully"},
		{"invalid id", "9", http.StatusNotFound, "the requested resource could not be found"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/", nil)
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", "1")
			chiCtx.URLParams.Add("optionID", test.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.deleteOptionHandler)
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
