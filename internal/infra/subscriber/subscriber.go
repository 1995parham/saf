package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/1995parham/saf/internal/domain/model/event"
	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Subscriber struct {
	CMQ    *cmq.CMQ
	Tracer trace.Tracer
	Logger *zap.Logger

	handlers []chan<- event.ChanneledEvent
}

func New(c *cmq.CMQ, logger *zap.Logger, tracer trace.Tracer) *Subscriber {
	var subscriber Subscriber

	subscriber.handlers = make([]chan<- event.ChanneledEvent, 0)
	subscriber.CMQ = c
	subscriber.Tracer = tracer
	subscriber.Logger = logger.Named("subscriber")

	return &subscriber
}

// Only pull consumers are supported in jetstream package. However, unlike the JetStream API in nats package,
// pull consumers allow for continuous message retrieval (similarly to how nats.Subscribe() works).
// Because of that, push consumers can be easily replace by pull consumers for most of the use cases.
//
// a consumer can also be ephemeral or durable. A consumer is considered durable when an explicit name is set on
// the Durable field when creating the consumer, otherwise it is considered ephemeral.
// Durables and ephemeral behave exactly the same except that an ephemeral will be automatically cleaned up (deleted)
// after a period of inactivity, specifically when there are no subscriptions bound to the consumer.
// By default, durables will remain even when there are periods
// of inactivity (unless InactiveThreshold is set explicitly)
func (s *Subscriber) Subscribe() error {
	// nolint: exhaustruct
	con, err := s.CMQ.Jetstream.CreateOrUpdateConsumer(context.Background(), cmq.EventsChannel, jetstream.ConsumerConfig{
		Name:              "",
		Durable:           "",
		AckPolicy:         jetstream.AckExplicitPolicy,
		DeliverPolicy:     jetstream.DeliverLastPerSubjectPolicy,
		InactiveThreshold: time.Hour, // remove durable consumer after 1 hour of inactivity
		FilterSubject:     cmq.EventsChannel,
	})
	if err != nil {
		return fmt.Errorf("consumer creation failed %w", err)
	}

	if _, err := con.Consume(s.handler); err != nil {
		return fmt.Errorf("consume failed %w", err)
	}

	return nil
}

func (s *Subscriber) handler(msg jetstream.Msg) {
	ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(msg.Headers()))

	ctx, span := s.Tracer.Start(ctx, "subscriber.events", trace.WithSpanKind(trace.SpanKindConsumer))
	defer span.End()

	metadata, _ := msg.Metadata()

	s.Logger.Info("receive new message",
		zap.String("timestamp", metadata.Timestamp.String()),
		zap.ByteString("payload", msg.Data()),
	)

	var ev model.Event

	if err := json.Unmarshal(msg.Data(), &ev); err != nil {
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
