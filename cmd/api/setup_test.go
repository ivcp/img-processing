package main

import (
	"os"
	"testing"
)

var app application

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
