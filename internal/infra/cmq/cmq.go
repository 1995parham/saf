package cmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CMQ struct {
	Jetstream jetstream.JetStream
	nats      *nats.Conn
	logger    *zap.Logger
}

func New(lc fx.Lifecycle, cfg Config, logger *zap.Logger) (*CMQ, error) {
	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed %w", err)
	}

	logger.Info("nats connection successful",
		zap.String("connected-addr", nc.ConnectedAddr()),
		zap.Strings("discovered-servers", nc.DiscoveredServers()))

	nc.SetDisconnectErrHandler(func(_ *nats.Conn, err error) {
		logger.Fatal("nats disconnected", zap.Error(err))
	})

	nc.SetReconnectHandler(func(_ *nats.Conn) {
		logger.Warn("nats reconnected")
	})

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream context extraction failed %w", err)
	}

	lc.Append(
		fx.StopHook(func() { nc.Close() }),
	)

	return &CMQ{
		nats:      nc,
		logger:    logger,
		Jetstream: js,
	}, nil
}

// Streams creates required streams on jetstream.
// On production you may want to create streams manually to have
// more control. Stream creation process is like migration.
func (c *CMQ) Streams(ctx context.Context) error {
	info, err := c.Jetstream.Stream(ctx, EventsChannel)

	switch {
	case errors.Is(err, nats.ErrStreamNotFound):
		// Each stream contains multiple topics, here we use a
		// same name for stream and its topic.
		stream, err := c.Jetstream.CreateStream(ctx, jetstream.StreamConfig{
			Name:                 EventsChannel,
			Description:          "Saf's event channel which only contains events topic",
			Subjects:             []string{EventsChannel},
			Retention:            jetstream.LimitsPolicy,
			MaxConsumers:         0,
			MaxMsgs:              0,
			MaxBytes:             0,
			Discard:              jetstream.DiscardOld,
			DiscardNewPerSubject: false,
			MaxAge:               1 * time.Hour,
			MaxMsgsPerSubject:    0,
			MaxMsgSize:           0,
			Storage:              jetstream.MemoryStorage,
			Replicas:             1,
			NoAck:                false,
			Template:             "",
			Duplicates:           0,
			Placement:            nil,
			Mirror:               nil,
			Sources:              nil,
			Sealed:               false,
			DenyDelete:           false,
			DenyPurge:            false,
			AllowRollup:          false,
			RePublish:            nil,
			AllowDirect:          false,
			MirrorDirect:         false,
		})
		if err != nil {
			return fmt.Errorf("cannot create stream %w", err)
		}

		info = stream
	case err != nil:
		return fmt.Errorf("cannot read stream information %w", err)
	}

	c.logger.Info("events stream", zap.Any("stream", info))

	return nil
}
