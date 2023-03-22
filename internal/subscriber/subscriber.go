package subscriber

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/model"
	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Subscriber struct {
	CMQ    *cmq.CMQ
	Tracer trace.Tracer
	Logger *zap.Logger

	handlers []chan<- model.ChanneledEvent
}

func New(c *cmq.CMQ, logger *zap.Logger, tracer trace.Tracer) *Subscriber {
	var subscriber Subscriber

	subscriber.handlers = make([]chan<- model.ChanneledEvent, 0)
	subscriber.CMQ = c
	subscriber.Tracer = tracer
	subscriber.Logger = logger.Named("subscriber")

	return &subscriber
}

func (s *Subscriber) Subscribe() error {
	// change between pull and push consumers
	// return s.PushSubscribe()
	return s.PullSubscribe()
}

func (s *Subscriber) PushSubscribe() error {
	// subscribe finds the stream name automatically based on given subject and also creates the consumer.
	// we can create the consumer manually with nats.Bind or set the stream name manually with nats.BindStream.
	if _, err := s.CMQ.JConn.QueueSubscribe(cmq.EventsChannel, cmq.QueueName, s.handler,
		nats.AckExplicit(),
		nats.DeliverAll(),
		nats.Durable(cmq.DurableName),
		nats.InactiveThreshold(time.Hour), // remove durable consumer after 1 hour of inactivity
	); err != nil {
		return fmt.Errorf("queue subscrption failed %w", err)
	}

	return nil
}

func (s *Subscriber) PullSubscribe() error {
	sub, err := s.CMQ.JConn.PullSubscribe(cmq.EventsChannel, cmq.DurableName,
		nats.AckExplicit(),
		nats.DeliverAll(),
		nats.BindStream(cmq.EventsChannel),
		nats.InactiveThreshold(time.Hour), // remove durable consumer after 1 hour of inactivity
	)
	if err != nil {
		return fmt.Errorf("pull subscrption failed %w", err)
	}

	go func() {
		for {
			msg, err := sub.Fetch(1)
			if err != nil {
				if errors.Is(err, nats.ErrTimeout) {
					continue
				}

				s.Logger.Error("fetching messages from a pull consumer failed", zap.Error(err))
			}

			for _, msg := range msg {
				s.handler(msg)
			}
		}
	}()

	return nil
}

func (s *Subscriber) handler(msg *nats.Msg) {
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(msg.Header))

	ctx, span := s.Tracer.Start(ctx, "subscriber.events", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	metadata, _ := msg.Metadata()

	s.Logger.Info("receive new message",
		zap.String("timestamp", metadata.Timestamp.String()),
		zap.ByteString("payload", msg.Data),
	)

	var ev model.Event

	if err := json.Unmarshal(msg.Data, &ev); err != nil {
		s.Logger.Error("cannot parse the event", zap.Error(err))
	}

	cev := model.ChanneledEvent{Event: ev, SpanContext: trace.SpanContextFromContext(ctx)}

	for _, c := range s.handlers {
		c <- cev
	}
}

func (s *Subscriber) RegisterHandler(c chan<- model.ChanneledEvent) {
	s.handlers = append(s.handlers, c)
}
