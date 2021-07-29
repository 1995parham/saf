package handler

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

type Healthz struct {
	Logger *zap.Logger
	Tracer trace.Tracer
}

// Handle shows server is up and running.
// nolint: wrapcheck
func (h Healthz) Handle(c *fiber.Ctx) error {
	_, span := h.Tracer.Start(c.Context(), "handler.healthz")
	defer span.End()

	return c.Status(http.StatusNoContent).Send(nil)
}

// Register registers the routes of healthz handler on given fiber group.
func (h Healthz) Register(g fiber.Router) {
	g.Get("/healthz", h.Handle)
}
