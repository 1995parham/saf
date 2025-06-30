package manager

import (
	"context"
	"encoding/json"

	"github.com/1995parham/saf/internal/domain/model/event"
	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/output"
	"github.com/1995parham/saf/internal/infra/telemetry"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// Manager manages output channels.
type Manager struct {
	Plugins []output.Channel

	logger    *zap.Logger
	tracer    trace.Tracer
	cmq       *cmq.CMQ
	consumers []jetstream.ConsumeContext
}

// Provide create new manager for output channels.
func Provide(lc fx.Lifecycle, cfg output.Config, logger *zap.Logger, _ telemetry.Telemetery, cmq *cmq.CMQ) *Manager {
	manager := &Manager{
		Plugins:   make([]output.Channel, 0),
		logger:    logger.Named("output").Named("manager"),
		tracer:    otel.GetTracerProvider().Tracer("output.manager"),
		cmq:       cmq,
		consumers: make([]jetstream.ConsumeContext, 0),
	}

	enabled := []string{}
	for name := range cfg.Configurations {
		enabled = append(enabled, name)
	}

	lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				manager.Setup(ctx, enabled, cfg.Configurations)

				return nil
			},
			OnStop: func(_ context.Context) error {
				for _, consumer := range manager.consumers {
					consumer.Stop()
				}

				return nil
			},
		},
	)

	return manager
}

// Setup registers the given channel. please note that you should add each channel here.
func (m *Manager) Setup(ctx context.Context, enabled []string, cfg map[string]interface{}) {
	for _, p := range channels() {
		for _, e := range enabled {
			if p.Name() == e {
				m.logger.Info("register new plugin", zap.String("plugin", p.Name()))
				m.Register(ctx, p, cfg[p.Name()])
			}
		}
	}
}

// Register registers the given channel and passes its configuration to it.
// Also runs it in new goroutine.
func (m *Manager) Register(ctx context.Context, p output.Channel, cfg interface{}) {
	m.logger.Info("channel started", zap.String("channel", p.Name()))

	m.Plugins = append(m.Plugins, p)

	c := make(chan output.TracedEvent)

	con, err := m.cmq.Subscribe(ctx, p.Name(), func(ctx context.Context, data []byte) {
		var ev event.Event

		ctx, span := m.tracer.Start(ctx, "manager.subscriber", trace.WithSpanKind(trace.SpanKindConsumer))
		defer span.End()

		err := json.Unmarshal(data, &ev)
		if err != nil {
			m.logger.Error("cannot parse the event", zap.Error(err))
		}

		cev := output.TracedEvent{Event: ev, SpanContext: trace.SpanContextFromContext(ctx)}

		select {
		case c <- cev:
		default:
		}
	})
	if err != nil {
		m.logger.Error("cannot create subscription", zap.Error(err))
	}

	m.consumers = append(m.consumers, con)

	p.Init(m.logger.Named(p.Name()), m.tracer, cfg, c)

	go p.Run()
}
