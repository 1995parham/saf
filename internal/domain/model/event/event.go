package event

import (
	"time"

	"go.opentelemetry.io/otel/trace"
)

type Event struct {
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"createdAt"`
	Payload   []byte    `json:"payload"`
}

type ChanneledEvent struct {
	Event

	SpanContext trace.SpanContext
}
