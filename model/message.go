package model

import "time"

// Message represents a message to broadcast.
type Message struct {
	From      string
	CreatedAt time.Time
	Text      string
}
