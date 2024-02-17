package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/ivcp/polls/internal/data"
)

func Test_app_updatePollHandler(t *testing.T) {
	tests := []struct {
		name           string
		id             string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid update partial",
			id:             data.ExamplePollIDValid,
			json:           `{"question":"changed"}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `"question":"changed"`,
		},
		{
			name: "valid update complete",
			id:   data.ExamplePollIDValid,
			json: fmt.Sprintf(
				`{"question":"changed", "description":"added description", "expires_at":%q}`,
				time.Now().Add(2*time.Minute).Format(time.RFC3339),
			),
			expectedStatus: http.StatusOK,
			expectedBody:   `"question":"changed","description":"added description"`,
		},
		{
			name:           "empty json",
			id:             data.ExamplePollIDValid,
			json:           `{}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "no fields provided for update",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/", strings.NewReader(test.json))
			poll, _ := app.models.Polls.Get(test.id)
			t.Log(poll.ID)
			req = req.WithContext(context.WithValue(req.Context(), ctxPollKey, poll))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.updatePollHandler)
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
