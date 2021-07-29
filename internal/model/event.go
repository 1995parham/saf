package model

import (
	"context"
	"time"
)

type Event struct {
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"createdAt"`
	Payload   []byte    `json:"payload"`
}

type ChanneledEvent struct {
	Event

	Context context.Context
}
