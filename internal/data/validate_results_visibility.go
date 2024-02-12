package data

import "github.com/ivcp/polls/internal/validator"

type ResultsVisibility struct {
	Value         string
	ValueSafelist []string
}

func ValidateResultsVisibility(v *validator.Validator, f ResultsVisibility) {
	v.Check(validator.PermittedValue(
		f.Value, f.ValueSafelist...,
	), "results_visibility", "invalid results_visibility value")
}
