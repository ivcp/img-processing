package main

import (
	"fmt"
	"net/http"
)

func (app *application) createPollHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new poll")
}
