package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_app_createPollHandler(t *testing.T) {
	expiresValid := time.Now().Add(2 * time.Minute)
	expiresInvalid := time.Now()
	questionInvalid := strings.Repeat("a", 501)
	descriptionInvalid := strings.Repeat("a", 1001)

	type opts struct {
		Value    any `json:"value"`
		Position any `json:"position"`
	}
	type input struct {
		Question    any    `json:"question"`
		Description any    `json:"description"`
		Options     []opts `json:"options"`
		Expires_at  any    `json:"expires_at"`
	}

	tests := []struct {
		name           string
		json           any
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "epmty question",
			json: input{
				Question: "",
				Options: []opts{
					{"test", 0},
				},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"question":"must not be empty"}}` + "\n",
		},
		{
			name: "question too long",
			json: input{
				Question:   questionInvalid,
				Options:    []opts{{"test", 0}},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"question":"must not be more than 500 bytes long"}}` + "\n",
		},
		{
			name: "description too long",
			json: input{
				Question:    "test?",
				Description: descriptionInvalid,
				Options:     []opts{{"test", 0}},
				Expires_at:  expiresValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"description":"must not be more than 1000 bytes long"}}` + "\n",
		},
		{
			name: "expires_at missing",
			json: input{
				Question: "test?",
				Options:  []opts{{"test", 0}},
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"expires_at":"must be provided"}}` + "\n",
		},
		{
			name: "expires_at invalid",
			json: input{
				Question:   "test?",
				Options:    []opts{{"test", 0}},
				Expires_at: expiresInvalid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"expires_at":"must be more than a minute in the future"}}` + "\n",
		},
		{
			name: "duplicate options",
			json: input{
				Question:   "test?",
				Options:    []opts{{"test", 0}, {"test", 0}},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"must not contain duplicate values"}}` + "\n",
		},
		{
			name: "duplicate option positions",
			json: input{
				Question:   "test?",
				Options:    []opts{{"test", 0}, {"test2", 0}},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"positions must be unique"}}` + "\n",
		},
		{
			name: "invalid option positions",
			json: input{
				Question:   "test?",
				Options:    []opts{{"test", 2}, {"test2", 0}},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"position must not excede the number of options"}}` + "\n",
		},
		{
			name: "invalid option positions",
			json: input{
				Question:   "test?",
				Options:    []opts{{"test", -1}, {"test2", -2}},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"position must be greater or equal to 0"}}` + "\n",
		},
		{
			name: "invalid json field",
			json: input{
				Question:   1,
				Options:    []opts{{"test", 0}},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"body contains incorrect JSON type for field \"question\""}` + "\n",
		},
		{
			name: "insert poll valid",
			json: input{
				Question:   "test?",
				Options:    []opts{{"test", 0}},
				Expires_at: expiresValid,
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/", bytes.NewReader(createTestJSON(t, test.json)))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.createPollHandler)
			handler.ServeHTTP(rr, req)
			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, but got %d", test.expectedStatus, rr.Code)
			}
			if test.expectedBody != "" && rr.Body.String() != test.expectedBody {
				t.Errorf("expected body %q, but got %q", test.expectedBody, rr.Body)
			}
			t.Log(rr.Body)
		})
	}
}

func createTestJSON(t *testing.T, data any) []byte {
	t.Helper()
	j, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	return j
}
