package data

import (
	"encoding/json"
	"time"
)

type ExpiresAt struct{ time.Time }

func (e ExpiresAt) MarshalJSON() ([]byte, error) {
	if e.IsZero() {
		return []byte(`""`), nil
	}

	return json.Marshal(e.Time)
}
