package printer

import (
	"context"
	"runtime"

	"github.com/1995parham/saf/internal/infra/output"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Printer is a plugin for the saf app. this plugin consumes event
// and log them.
type Printer struct {
	ch     <-chan output.TracedEvent
	logger *zap.Logger
	tracer trace.Tracer
}

func (p *Printer) Init(logger *zap.Logger, tracer trace.Tracer, _ interface{}, ch <-chan output.TracedEvent) {
	p.logger = logger
	p.tracer = tracer
	p.ch = ch
}

func (p *Printer) Run() {
	for range 10 * runtime.GOMAXPROCS(0) {
		go func() {
			for e := range p.ch {
				ctx := trace.ContextWithSpanContext(context.Background(), e.SpanContext)
				_, span := p.tracer.Start(ctx, "channels.printer")

				p.logger.Info("receive event",
					zap.Time("created", e.CreatedAt),
					zap.String("subject", e.Subject),
				)

				span.End()
			}
		}()
	}
}

func (p *Printer) Name() string {
	return "printer"
}
