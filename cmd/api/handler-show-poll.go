package main

import (
	"fmt"
	"net/http"
)

func (app *application) showPollHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	fmt.Fprintf(w, "show the details of poll %d\n", id)
}
