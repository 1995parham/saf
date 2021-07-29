package printer

import (
	"runtime"

	"github.com/1995parham/saf/internal/model"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

// Printer is a plugin for the saf app. this plugin consumes event
// and log them.
type Printer struct {
	ch     <-chan model.ChanneledEvent
	logger *zap.Logger
	tracer trace.Tracer
}

func (p *Printer) Init(logger *zap.Logger, tracer trace.Tracer, cfg interface{}) {
	p.logger = logger
	p.tracer = tracer
}

func (p *Printer) Run() {
	for i := 0; i < 10*runtime.GOMAXPROCS(0); i++ {
		go func() {
			for e := range p.ch {
				_, span := p.tracer.Start(e.Context, "channels.printer")

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

func (p *Printer) SetChannel(c <-chan model.ChanneledEvent) {
	p.ch = c
}
