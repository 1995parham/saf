package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/model"
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

// Only pull consumers are supported in jetstream package. However, unlike the JetStream API in nats package,
// pull consumers allow for continuous message retrieval (similarly to how nats.Subscribe() works).
// Because of that, push consumers can be easily replace by pull consumers for most of the use cases.
func (s *Subscriber) Subscribe() error {
	// nolint: exhaustruct
	con, err := s.CMQ.JConn.CreateOrUpdateConsumer(context.Background(), cmq.EventsChannel, jetstream.ConsumerConfig{
		Name:              "",
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
