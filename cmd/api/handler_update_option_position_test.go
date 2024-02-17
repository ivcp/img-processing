package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/ivcp/polls/internal/data"
)

func Test_app_updateOptionPositionHandler(t *testing.T) {
	tests := []struct {
		name           string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid update positions",
			json: fmt.Sprintf(
				`{"options":[{"id":%q, "position":1}, {"id":%q, "position":0}]}`,
				data.ExampleOptionID1, data.ExampleOptionID2,
			),
			expectedStatus: http.StatusOK,
			expectedBody:   "options updated successfully",
		},
		{
			name: "invalid option id",
			json: fmt.Sprintf(
				`{"options":[{"id":%q, "position":1}, {"id":%q, "position":0}]}`,
				data.ExampleOptionID1, uuid.NewString(),
			),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid option id, or no id provided",
		},
		{
			name:           "invalid position change",
			json:           fmt.Sprintf(`{"options":[{"id":%q, "position":1}]}`, data.ExampleOptionID1),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   "positions must be unique",
		},
		{
			name:           "no options provided",
			json:           `{"options":[]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "invalid option id, or no id provided",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPatch, "/", strings.NewReader(test.json))
			poll, _ := app.models.Polls.Get(data.ExamplePollIDValid)
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
