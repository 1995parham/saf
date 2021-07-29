package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/http/request"
	"github.com/1995parham/saf/internal/model"
	"github.com/gofiber/fiber/v2"
	"github.com/nats-io/nats.go"
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
	ctx, span := h.Tracer.Start(c.Context(), "handler.event")
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

	ev := model.Event{
		Subject:   rq.Subject,
		CreatedAt: time.Now(),
		Payload:   rq.Payload,
	}

	data, err := json.Marshal(ev)
	if err != nil {
		span.RecordError(err)

		return fiber.NewError(http.StatusInternalServerError, err.Error())
	}

	msg := new(nats.Msg)

	msg.Subject = cmq.EventsChannel
	msg.Data = data
	msg.Header = make(nats.Header)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

	{
		_, span := h.Tracer.Start(ctx, "handler.event.publish")
		if _, err := h.CMQ.JConn.PublishMsg(msg); err != nil {
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
