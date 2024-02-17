package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ivcp/polls/internal/data"
)

func Test_app_deletePollHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		expectedStatus int
		expectedBody   string
	}{
		{"delete a poll", data.ExamplePollIDValid, http.StatusOK, "poll successfully deleted"},
		{"poll not found", uuid.NewString(), http.StatusNotFound, "the requested resource could not be found"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodDelete, "/", nil)
			req = req.WithContext(context.WithValue(req.Context(), ctxPollIDKey, test.id))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.deletePollHandler)
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
