package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func Test_app_updateOptionPositionHandler(t *testing.T) {
	tests := []struct {
		name           string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{
			"valid update positions",
			`{"options":[{"id":1, "position":1}, {"id":2, "position":0}]}`,
			http.StatusOK,
			"options updated successfully",
		},
		{
			"invalid option id",

			`{"options":[{"id":1, "position":1}, {"id":9, "position":0}]}`,
			http.StatusBadRequest,
			"invalid option id, or no id provided",
		},
		{
			"invalid position change",
			`{"options":[{"id":1, "position":1}]}`,
			http.StatusUnprocessableEntity,
			"positions must be unique",
		},
		{
			"no options provided",
			`{"options":[]}`,
			http.StatusBadRequest,
			"invalid option id, or no id provided",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/", strings.NewReader(test.json))
			poll, _ := app.models.Polls.Get(1)
			req = req.WithContext(context.WithValue(req.Context(), ctxPollKey, poll))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.updateOptionPositionHandler)
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
