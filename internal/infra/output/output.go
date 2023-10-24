package output

import (
	"github.com/1995parham/saf/internal/domain/model/event"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TracedEvent struct {
	event.Event

	SpanContext trace.SpanContext
}

type Channel interface {
	Init(logger *zap.Logger, tracer trace.Tracer, config interface{}, eventChan <-chan TracedEvent)
	Run()
	Name() string
}
