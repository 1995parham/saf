package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/1995parham/saf/internal/doamin/model/event"
	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/http/request"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Event struct {
	CMQ    *cmq.CMQ
	Logger *zap.Logger
	Tracer trace.Tracer
}

// Receive receives event from api and publish them to jetstream.
// nolint: wrapcheck
func (h Event) Receive(c *fiber.Ctx) error {
	ctx, span := h.Tracer.Start(c.Context(), "handler.event",
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	var rq request.Event

	if err := c.BodyParser(&rq); err != nil {
		span.RecordError(err)

		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		span.RecordError(err)

		return fiber.NewError(http.StatusBadRequest, err.Error())
	}

	ev := event.Event{
		Subject:   rq.Subject,
		CreatedAt: time.Now(),
		Payload:   rq.Payload,
	}

	data, err := json.Marshal(ev)
	if err != nil {
		span.RecordError(err)

		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	{
		ctx, span := h.Tracer.Start(ctx, "handler.event.publish", trace.WithSpanKind(trace.SpanKindProducer))
		msg := new(nats.Msg)

		msg.Subject = cmq.EventsChannel
		msg.Data = data
		msg.Header = make(nats.Header)
		otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

		if _, err := h.CMQ.JConn.PublishMsg(ctx, msg, jetstream.WithMsgID(rq.ID)); err != nil {
			span.RecordError(err)

			return fiber.NewError(http.StatusServiceUnavailable, err.Error())
		}
		span.End()
	}

	return c.Status(http.StatusOK).Send(nil)
}

// Register registers the routes of event handler on given fiber group.
func (h Event) Register(g fiber.Router) {
	g.Post("/event", h.Receive)
}
