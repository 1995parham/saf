package channel

import (
	"github.com/1995parham/saf/internal/model"
	"go.uber.org/zap"
)

type Channel interface {
	Init(*zap.Logger, interface{})
	Run()
	SetChannel(<-chan model.Event)
	Name() string
}

// Manager manages Channels.
type Manager struct {
	Plugins []Channel
	logger  *zap.Logger

	channels []chan model.Event
}

// New create new manager.
func New(logger *zap.Logger) *Manager {
	return &Manager{
		Plugins: make([]Channel, 0),
		logger:  logger,

		channels: make([]chan model.Event, 0),
	}
}

// Setup registers the given channel. please note that you should add each channel here.
func (m *Manager) Setup(enabled []string, cfg map[string]interface{}) {
	// list of available channles, please add each channel into this list to make them available.
	channels := []Channel{}

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

	c := make(chan model.Event)

	m.channels = append(m.channels, c)

	p.Init(m.logger.Named(p.Name()), cfg)
	p.SetChannel(c)

	go p.Run()
}

// Channels return the channels to register them on subscriber in the main.
func (m *Manager) Channels() []chan<- model.Event {
	c := make([]chan<- model.Event, len(m.channels))
	for i := range m.channels {
		c[i] = m.channels[i]
	}

	return c
}
