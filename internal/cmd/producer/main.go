package producer

import (
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/config"
	"github.com/1995parham/saf/internal/http/handler"
	"github.com/1995parham/saf/internal/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger) {
	cmq, err := cmq.New(cfg.NATS, logger)
	if err != nil {
		logger.Fatal("nats initiation failed", zap.Error(err))
	}

	if err := cmq.Streams(); err != nil {
		logger.Fatal("nats stream creation failed", zap.Error(err))
	}

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

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	cmq.Conn.Close()
}

// Register producer command.
func Register(cfg config.Config, logger *zap.Logger) *cli.Command {
	tele := telemetry.New(cfg.Telemetry)
	tele.Run()

	// nolint: exhaustruct
	return &cli.Command{
		Name:        "producer",
		Aliases:     []string{"p"},
		Description: "gets events from http and produce them into nats",
		Action: func(_ *cli.Context) error {
			main(cfg, logger)

			return nil
		},
		After: func(ctx *cli.Context) error {
			tele.Shutdown(ctx.Context)

			return nil
		},
	}
}
