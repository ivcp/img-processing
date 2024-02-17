package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/ivcp/polls/internal/data"
)

func Test_app_addOptionHandler(t *testing.T) {
	tests := []struct {
		name string

		json           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid add option",
			json:           `{"value":"test", "position":3}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "option added successfully",
		},
		{
			name:           "option already exists",
			json:           `{"value":"Two", "position":2}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "must not contain duplicate values",
		},
		{
			name:           "position not unique",
			json:           `{"value":"test", "position":1}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "positions must be unique",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(test.json))
			poll, _ := app.models.Polls.Get(data.ExamplePollIDValid)
			req = req.WithContext(context.WithValue(req.Context(), ctxPollKey, poll))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.addOptionHandler)
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
