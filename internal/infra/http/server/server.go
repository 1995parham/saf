package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/http/handler"
	"github.com/1995parham/saf/internal/infra/telemetry"
	"github.com/gofiber/fiber/v2"
	"go.opentelemetry.io/otel"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Provide(lc fx.Lifecycle, cmq *cmq.CMQ, logger *zap.Logger, _ telemetry.Telemetery) *fiber.App {
	// nolint: exhaustruct
	app := fiber.New(fiber.Config{
		AppName: "saf",
	})

	handler.Healthz{
		Logger: logger.Named("handler").Named("healthz"),
		Tracer: otel.GetTracerProvider().Tracer("handler.healthz"),
	}.Register(app.Group(""))

	handler.Event{
		CMQ:    cmq,
		Logger: logger.Named("handler").Named("event"),
		Tracer: otel.GetTracerProvider().Tracer("handler.event"),
	}.Register(app.Group("api"))

	if err := app.Listen(":1378"); !errors.Is(err, http.ErrServerClosed) {
		logger.Fatal("echo initiation failed", zap.Error(err))
	}

	lc.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				go func() {
					if err := app.Listen(":1378"); !errors.Is(err, http.ErrServerClosed) {
						logger.Fatal("fiber initiation failed", zap.Error(err))
					}
				}()

				return nil
			},
			OnStop: func(_ context.Context) error {
				if err := app.Shutdown(); err != nil {
					return fmt.Errorf("fiber shutdown failed %w", err)
				}

				return nil
			},
		},
	)

	return app
}
