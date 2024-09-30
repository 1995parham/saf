package cmq

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type CMQ struct {
	jetstream jetstream.JetStream
	nats      *nats.Conn
	logger    *zap.Logger
}

func Provide(lc fx.Lifecycle, cfg Config, logger *zap.Logger) (*CMQ, error) {
	cmq := &CMQ{
		nats:      nil,
		logger:    logger,
		jetstream: nil,
	}

	nc, err := nats.Connect(cfg.URL)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed %w", err)
	}

	cmq.nats = nc

	logger.Info("nats connection successful",
		zap.String("connected-addr", nc.ConnectedAddr()),
		zap.Strings("discovered-servers", nc.DiscoveredServers()))

	nc.SetDisconnectErrHandler(cmq.disconnectHandler)
	nc.SetClosedHandler(cmq.closeHandler)
	nc.SetReconnectHandler(cmq.reconnectHandler)

	js, err := jetstream.New(nc)
	if err != nil {
		return nil, fmt.Errorf("jetstream creation failed %w", err)
	}

	cmq.jetstream = js

	lc.Append(
		fx.StopHook(func() { nc.Close() }),
	)

	return cmq, nil
}

// close handler is called when the connection is closed which means
// we are not going to retry again for getting a live connection
// to server anymore.
func (c *CMQ) closeHandler(nc *nats.Conn) {
	c.logger.Fatal("connection closed",
		zap.Strings("urls", nc.DiscoveredServers()),
		zap.Error(nc.LastError()),
	)
}

// disconnection handler is called when we lost connection
// and we are going to retry, so we may get connected in the future.
func (c *CMQ) disconnectHandler(nc *nats.Conn, err error) {
	c.logger.Error("got disconnected",
		zap.Strings("urls", nc.DiscoveredServers()),
		zap.Error(err),
	)
}

func (c *CMQ) reconnectHandler(nc *nats.Conn) {
	c.logger.Info("got reconnected", zap.String("url", nc.ConnectedUrl()))
}

// Streams creates required streams on jetstream.
// On production you may want to create streams manually to have
// more control. Stream creation process is like migration.
func (c *CMQ) Streams(ctx context.Context) error {
	info, err := c.jetstream.Stream(ctx, EventsChannel)

	switch {
	case errors.Is(err, jetstream.ErrStreamNotFound):
		// Each stream contains multiple topics, here we use a
		// same name for stream and its topic.
		// nolint: exhaustruct
		stream, err := c.jetstream.CreateStream(ctx, jetstream.StreamConfig{
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

type Handler func(context.Context, []byte)

// Only pull consumers are supported in jetstream package. However, unlike the JetStream API in nats package,
// pull consumers allow for continuous message retrieval (similarly to how nats.Subscribe() works).
// Because of that, push consumers can be easily replace by pull consumers for most of the use cases.
//
// a consumer can also be ephemeral or durable. A consumer is considered durable when an explicit name is set on
// the Durable field when creating the consumer, otherwise it is considered ephemeral.
// Durables and ephemeral behave exactly the same except that an ephemeral will be automatically cleaned up (deleted)
// after a period of inactivity, specifically when there are no subscriptions bound to the consumer.
// By default, durables will remain even when there are periods
// of inactivity (unless InactiveThreshold is set explicitly).
func (c *CMQ) Subscribe(ctx context.Context, name string, handler Handler) (jetstream.ConsumeContext, error) {
	// nolint: exhaustruct
	con, err := c.jetstream.CreateOrUpdateConsumer(ctx, EventsChannel, jetstream.ConsumerConfig{
		Name:              fmt.Sprintf("%s-%s", QueueName, name),
		Durable:           fmt.Sprintf("%s-%s", DurableName, name),
		AckPolicy:         jetstream.AckExplicitPolicy,
		DeliverPolicy:     jetstream.DeliverLastPerSubjectPolicy,
		InactiveThreshold: time.Hour, // remove durable consumer after 1 hour of inactivity
		FilterSubject:     EventsChannel,
	})
	if err != nil {
		return nil, fmt.Errorf("consumer creation failed %w", err)
	}

	conCtx, err := con.Consume(c.handler(handler)) // nolint: contextcheck
	if err != nil {
		return nil, fmt.Errorf("consume failed %w", err)
	}

	return conCtx, nil
}

func (c *CMQ) handler(h Handler) jetstream.MessageHandler {
	return func(msg jetstream.Msg) {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(msg.Headers()))

		metadata, _ := msg.Metadata()

		c.logger.Info("receive new message",
			zap.String("timestamp", metadata.Timestamp.String()),
			zap.ByteString("payload", msg.Data()),
		)

		if err := msg.Ack(); err != nil {
			c.logger.Error("cannot ack message",
				zap.String("timestamp", metadata.Timestamp.String()),
				zap.Error(err),
			)
		}

		h(ctx, msg.Data())
	}
}

func (c *CMQ) Publish(ctx context.Context, id string, data []byte) error {
	msg := new(nats.Msg)

	msg.Subject = EventsChannel
	msg.Data = data
	msg.Header = make(nats.Header)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

	if _, err := c.jetstream.PublishMsg(context.WithoutCancel(ctx), msg, jetstream.WithMsgID(id)); err != nil {
		return fmt.Errorf("jetstream publish message failed %w", err)
	}

	return nil
}
