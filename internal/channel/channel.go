package channel

import (
	"github.com/1995parham/saf/internal/channel/mqtt"
	"github.com/1995parham/saf/internal/channel/printer"
	"github.com/1995parham/saf/internal/model"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Channel interface {
	Init(*zap.Logger, trace.Tracer, interface{})
	Run()
	SetChannel(<-chan model.ChanneledEvent)
	Name() string
}

// Manager manages Channels.
type Manager struct {
	Plugins []Channel
	logger  *zap.Logger
	tracer  trace.Tracer

	channels []chan model.ChanneledEvent
}

// New create new manager.
func New(logger *zap.Logger, tracer trace.Tracer) *Manager {
	return &Manager{
		Plugins: make([]Channel, 0),
		logger:  logger,
		tracer:  tracer,

		channels: make([]chan model.ChanneledEvent, 0),
	}
}

// Setup registers the given channel. please note that you should add each channel here.
func (m *Manager) Setup(enabled []string, cfg map[string]interface{}) {
	// list of available channles, please add each channel into this list to make them available.
	channels := []Channel{
		&printer.Printer{},
		&mqtt.MQTT{},
	}

	for _, p := range channels {
		for _, e := range enabled {
			if p.Name() == e {
				m.Register(p, cfg[p.Name()])
			}
		}
	}
}

// Register registers the given channel and passes its configration to it.
// Also runs it in new goroutine.
func (m *Manager) Register(p Channel, cfg interface{}) {
	m.logger.Info("channel started", zap.String("channel", p.Name()))

	m.Plugins = append(m.Plugins, p)

	c := make(chan model.ChanneledEvent)

	m.channels = append(m.channels, c)

	p.Init(m.logger.Named(p.Name()), m.tracer, cfg)
	p.SetChannel(c)

	go p.Run()
}

// Channels return the channels to register them on subscriber in the main.
func (m *Manager) Channels() []chan<- model.ChanneledEvent {
	c := make([]chan<- model.ChanneledEvent, len(m.channels))
	for i := range m.channels {
		c[i] = m.channels[i]
	}

	return c
}
