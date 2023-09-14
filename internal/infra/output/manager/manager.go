package manager

import (
	"context"
	"encoding/json"

	"github.com/1995parham/saf/internal/domain/model/event"
	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/output"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Manager manages output channels.
type Manager struct {
	Plugins []output.Channel

	logger *zap.Logger
	tracer trace.Tracer
	cmq    *cmq.CMQ

	channels []chan output.TracedEvent
}

// Provide create new manager for output channels.
func Provide(lc fx.Lifecycle, logger *zap.Logger, tracer trace.Tracer, cmq *cmq.CMQ) *Manager {
	manager := &Manager{
		Plugins: make([]output.Channel, 0),
		logger:  logger,
		tracer:  tracer,
		cmq:     cmq,
	}

	return manager
}

// Setup registers the given channel. please note that you should add each channel here.
func (m *Manager) Setup(ctx context.Context, enabled []string, cfg map[string]interface{}) {
	for _, p := range channels {
		for _, e := range enabled {
			if p.Name() == e {
				m.Register(ctx, p, cfg[p.Name()])
			}
		}
	}
}

// Register registers the given channel and passes its configration to it.
// Also runs it in new goroutine.
func (m *Manager) Register(ctx context.Context, p output.Channel, cfg interface{}) {
	m.logger.Info("channel started", zap.String("channel", p.Name()))

	m.Plugins = append(m.Plugins, p)

	c := make(chan output.TracedEvent)

	m.cmq.Subscribe(ctx, p.Name(), func(ctx context.Context, data []byte) {
		var ev event.Event

		ctx, span := m.tracer.Start(ctx, "manager.subscriber", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		if err := json.Unmarshal(data, &ev); err != nil {
			m.logger.Error("cannot parse the event", zap.Error(err))
		}

		cev := output.TracedEvent{Event: ev, SpanContext: trace.SpanContextFromContext(ctx)}

		select {
		case c <- cev:
		default:
		}
	})

	p.Init(m.logger.Named(p.Name()), m.tracer, cfg, c)

	go p.Run()
}
