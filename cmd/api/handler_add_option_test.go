package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
)

func Test_app_addOptionHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid add option",
			id:             "1",
			json:           `{"value":"test", "position":2}`,
			expectedStatus: http.StatusCreated,
			expectedBody:   "option added successfully",
		},
		{
			name:           "option already exists",
			id:             "1",
			json:           `{"value":"Two", "position":2}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "must not contain duplicate values",
		},
		{
			name:           "position not unique",
			id:             "1",
			json:           `{"value":"test", "position":1}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "positions must be unique",
		},
		{
			name:           "unexisting poll",
			id:             "9",
			json:           `{"value":"test", "position":2}`,
			expectedStatus: http.StatusNotFound,
			expectedBody:   "the requested resource could not be found",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(test.json))
			chiCtx := chi.NewRouteContext()
			chiCtx.URLParams.Add("pollID", test.id)
			req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))
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
