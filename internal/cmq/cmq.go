package cmq

import (
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

type CMQ struct {
	Conn   *nats.Conn
	JConn  nats.JetStreamContext
	Logger *zap.Logger
}

func New(cfg Config, logger *zap.Logger) (*CMQ, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed %w", err)
	}

	logger.Info("nats connection successful",
		zap.String("connected-addr", nc.ConnectedAddr()),
		zap.Strings("discovered-servers", nc.DiscoveredServers()))

	nc.SetDisconnectErrHandler(func(c *nats.Conn, err error) {
		logger.Warn("nats disconnected", zap.Error(err))
	})

	nc.SetReconnectHandler(func(c *nats.Conn) {
		logger.Warn("nats reconnected")
	})

	jsm, err := nc.JetStream()
	if err != nil {
		return nil, fmt.Errorf("jetstream context extraction failed %w", err)
	}

	return &CMQ{
		Conn:   nc,
		JConn:  jsm,
		Logger: logger,
	}, nil
}

func (c *CMQ) Streams() error {
	info, err := c.JConn.StreamInfo("ORDERS")

	switch {
	case errors.Is(err, nats.ErrStreamNotFound):
		// nolint: exhaustivestruct
		stream, err := c.JConn.AddStream(&nats.StreamConfig{
			Name:     "events",
			Subjects: []string{"events.>"},
			MaxAge:   1 * time.Hour,
			Storage:  nats.MemoryStorage,
		})
		if err != nil {
			return fmt.Errorf("cannot create stream %w", err)
		}

		info = stream
	case err != nil:
		return fmt.Errorf("cannot read stream information %w", err)
	}

	c.Logger.Info("events stream", zap.Any("stream", info))

	return nil
}
