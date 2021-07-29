package printer

import (
	"runtime"

	"github.com/1995parham/saf/internal/model"
	"go.uber.org/zap"
)

// Printer is a plugin for the saf app. this plugin consumes event
// and log them.
type Printer struct {
	ch     <-chan model.Event
	logger *zap.Logger
}

func (p *Printer) Init(logger *zap.Logger, cfg interface{}) {
	p.logger = logger
}

func (p *Printer) Run() {
	for i := 0; i < 10*runtime.GOMAXPROCS(0); i++ {
		go func() {
			for e := range p.ch {
				p.logger.Info("receive event",
					zap.Time("created", e.CreatedAt),
					zap.String("subject", e.Subject),
				)
			}
		}()
	}
}

func (p *Printer) Name() string {
	return "printer"
}

func (p *Printer) SetChannel(c <-chan model.Event) {
	p.ch = c
}
