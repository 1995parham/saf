package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/http/request"
	"github.com/1995parham/saf/internal/model"
	"github.com/labstack/echo/v4"
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
func (h Event) Receive(c echo.Context) error {
	ctx, span := h.Tracer.Start(c.Request().Context(), "handler.event")
	defer span.End()

	var rq request.Event

	if err := c.Bind(&rq); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	if err := rq.Validate(); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	ev := model.Event{
		Subject:   rq.Subject,
		CreatedAt: time.Now(),
		Payload:   rq.Payload,
		Span:      nil,
	}

	data, err := json.Marshal(ev)
	if err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	msg := new(nats.Msg)

	msg.Subject = cmq.EventsChannel
	msg.Data = data
	msg.Header = make(nats.Header)
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(msg.Header))

	if _, err := h.CMQ.JConn.PublishMsg(msg); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// Register registers the routes of healthz handler on given echo group.
func (h Event) Register(g *echo.Group) {
	g.POST("/event", h.Receive)
}
