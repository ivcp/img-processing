package main

import (
	"os"
	"testing"

	"github.com/ivcp/polls/internal/data"
)

var app application

func TestMain(m *testing.M) {
	app.models = data.NewMockModels()
	os.Exit(m.Run())
}
