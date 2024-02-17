package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func Test_app_createPollHandler(t *testing.T) {
	expiresValid := time.Now().Add(2 * time.Minute).Format(time.RFC3339)
	expiresInvalid := time.Now().Format(time.RFC3339)
	questionInvalid := strings.Repeat("a", 501)
	descriptionInvalid := strings.Repeat("a", 1001)

	tests := []struct {
		name           string
		json           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "epmty question",
			json: `{
				"question":"", 
				"options":[{"value":"first","position":0},{"value":"second","position":1}]
				}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"question":"must not be empty"}}`,
		},
		{
			name: "question too long",
			json: fmt.Sprintf(
				`{
				"question":%q, 
				"options":[{"value":"first","position":0},{"value":"second","position":1}]
				}`,
				questionInvalid,
			),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"question":"must not be more than 500 bytes long"}}`,
		},
		{
			name: "description too long",
			json: fmt.Sprintf(
				`{
					"question":"Test?", 
					"description":%q, 
					"options":[{"value":"first","position":0},{"value":"second","position":1}]}`,
				descriptionInvalid,
			),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"description":"must not be more than 1000 bytes long"}}`,
		},
		{
			name: "expires_at invalid",
			json: fmt.Sprintf(
				`{
					"question":"Test?", 
					"options":[{"value":"first","position":0},{"value":"second","position":1}],
					"expires_at":%q
					}`,
				expiresInvalid,
			),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"expires_at":"must be more than a minute in the future"}}`,
		},
		{
			name: "only one option provided",
			json: `{
				"question":"Test?", 
				"options":[{"value":"first","position":0}]
				}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"must contain at least two options"}}`,
		},
		{
			name: "duplicate options",
			json: `{
				"question":"Test?", 
				"options":[{"value":"first","position":0},{"value":"first","position":1}]
				}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"must not contain duplicate values"}}`,
		},
		{
			name: "duplicate option positions",
			json: `{
				"question":"Test?", 
				"options":[{"value":"first","position":0}, {"value":"second","position":0}]
				}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"positions must be unique"}}`,
		},
		{
			name: "invalid option positions",
			json: `{
				"question":"Test?", 
				"options":[{"value":"first","position":2}, {"value":"second","position":0}]
				}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"position must not excede the number of options"}}`,
		},
		{
			name: "invalid option positions",
			json: `{
				"question":"Test?", 
				"options":[{"value":"first","position":-1}, {"value":"second","position":0}]
				}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"position must be greater or equal to 0"}}`,
		},
		{
			name: "invalid empty option",
			json: `{
				"question":"Test?", 
				"options":[{"value":" ","position":0}, {"value":"second","position":1}]
				}`,
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"option values must not be empty"}}`,
		},
		{
			name: "option value too large",
			json: fmt.Sprintf(
				`{
					"question":"Test?", 
					"options":[{"value":%q,"position":0},{"value":"second","position":1}]				
					}`,
				questionInvalid,
			),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"options":"option value must not be more than 500 bytes long"}}`,
		},
		{
			name: "invalid json field type",
			json: `{
				"question":1, 
				"options":[{"value":"first","position":0}]
				}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"body contains incorrect JSON type for field \"question\""}`,
		},
		{
			name: "insert poll valid",
			json: fmt.Sprintf(
				`{
					"question":"Test?", 
					"options":[{"value":"first","position":0}, {"value":"second","position":1}],
					"expires_at":%q
					}`,
				expiresValid,
			),
			expectedStatus: http.StatusCreated,
			expectedBody:   `"question":"Test?"`,
		},
		{
			name: "invalid results_visibility",
			json: fmt.Sprintf(
				`{
					"question":"Test?", 
					"options":[{"value":"first","position":0}, {"value":"second","position":1}],
					"expires_at":%q,
					"results_visibility": "test"
					}`,
				expiresValid,
			),
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"error":{"results_visibility":"invalid results_visibility value"}}`,
		},
		{
			name: "valid results_visibility",
			json: fmt.Sprintf(
				`{
					"question":"Test?", 
					"options":[{"value":"first","position":0}, {"value":"second","position":1}],
					"expires_at":%q,
					"results_visibility": "after_deadline"
					}`,
				expiresValid,
			),
			expectedStatus: http.StatusCreated,
			expectedBody:   `"question":"Test?"`,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodPost, "/", strings.NewReader(test.json))
			rr := httptest.NewRecorder()
			handler := http.HandlerFunc(app.createPollHandler)
			handler.ServeHTTP(rr, req)
			if rr.Code != test.expectedStatus {
				t.Errorf("expected status %d, but got %d", test.expectedStatus, rr.Code)
			}
			if !strings.Contains(rr.Body.String(), test.expectedBody) {
				t.Errorf("expected body %q, but got %q", test.expectedBody, rr.Body)
			}
			if rr.Code == http.StatusCreated && !strings.Contains(
				rr.Header().Get("Location"),
				"/v1/polls/",
			) {
				t.Errorf("Location does not contain link to created poll")
			}
		})
	}
}
