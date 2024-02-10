package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_updateOptionValueHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{"valid update", "1", `{"value":"test"}`, http.StatusCreated, "option updated successfully"},
		{"invalid id", "9", `{"value":"test"}`, http.StatusNotFound, "the requested resource could not be found"},
		{"duplicate option values", "1", `{"value":"Two"}`, http.StatusUnprocessableEntity, "must not contain duplicate values"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/", strings.NewReader(test.json))
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("optionID", test.id)
			req = req.WithContext(context.WithValue(req.Context(), "pollID", 1))
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.updateOptionValueHandler)
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
