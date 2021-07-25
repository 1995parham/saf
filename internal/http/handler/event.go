package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/http/request"
	"github.com/1995parham/saf/internal/model"
	"github.com/labstack/echo/v4"
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
	_, span := h.Tracer.Start(c.Request().Context(), "handler.event")
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
	}

	msg, err := json.Marshal(ev)
	if err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if _, err := h.CMQ.JConn.Publish("event", msg); err != nil {
		span.RecordError(err)

		return echo.NewHTTPError(http.StatusServiceUnavailable, err.Error())
	}

	return c.NoContent(http.StatusOK)
}

// Register registers the routes of healthz handler on given echo group.
func (h Event) Register(g *echo.Group) {
	g.POST("/event", h.Receive)
}
