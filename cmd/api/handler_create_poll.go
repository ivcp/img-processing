package main

import (
	"fmt"
	"net/http"
	"time"
)

func (app *application) createPollHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Question    string    `json:"question"`
		Description string    `json:"desription"`
		Options     []string  `json:"options"`
		ExpiresAt   time.Time `json:"expires_at"`
	}

	err := app.readJSON(r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Fprintf(w, "%+v\n", input)
}
