package output

import (
	"github.com/1995parham/saf/internal/domain/model/event"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type TracedEvent struct {
	event.Event
}

type Channel interface {
	Init(*zap.Logger, trace.Tracer, interface{})
	Run()
	SetChannel(<-chan TracedEvent)
	Name() string
}
