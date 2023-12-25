package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (app *application) readIDParam(r *http.Request) (int, error) {
	param := chi.URLParam(r, "id")
	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id")
	}
	return id, nil
}
