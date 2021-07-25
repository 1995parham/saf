package cmq

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
)

func New(cfg Config, logger *zap.Logger) (*nats.Conn, error) {
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

	return nc, nil
}
