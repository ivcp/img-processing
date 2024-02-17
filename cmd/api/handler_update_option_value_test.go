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

func Test_app_updateOptionValueHandler(t *testing.T) {
	tests := []struct {
		name           string
		optionID       string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid update",
			optionID:       data.ExampleOptionID1,
			json:           `{"value":"test"}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "option updated successfully",
		},
		{
			name:           "invalid id",
			optionID:       uuid.NewString(),
			json:           `{"value":"test"}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "the requested resource could not be found",
		},
		{
			name:           "duplicate option values",
			optionID:       data.ExampleOptionID1,
			json:           `{"value":"Two"}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "must not contain duplicate values",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/", strings.NewReader(test.json))
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("optionID", test.optionID)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
			poll, _ := app.models.Polls.Get(data.ExamplePollIDValid)
			req = req.WithContext(context.WithValue(req.Context(), ctxPollKey, poll))
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
