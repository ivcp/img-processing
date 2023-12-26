package data

import "time"

type Poll struct {
	ID          int         `json:"id"`
	Question    string      `json:"question"`
	Description string      `json:"description"`
	Options     PollOptions `json:"options"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
	ExpiresAt   time.Time   `json:"expires_at"`
	Version     int         `json:"version"`
}

type PollOptions []struct {
	ID    int    `json:"id"`
	Value string `json:"value"`
	// Position of option in the list, starting at 0
	Position  int `json:"position"`
	VoteCount int `json:"vote_count"`
}
