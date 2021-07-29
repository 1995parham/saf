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
	"github.com/1995parham/saf/internal/metric"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	metric.NewServer(cfg.Monitoring).Start(logger.Named("metrics"))

	cmq, err := cmq.New(cfg.NATS, logger)
	if err != nil {
		logger.Fatal("nats initiation failed", zap.Error(err))
	}

	if err := cmq.Streams(); err != nil {
		logger.Fatal("nats stream creation failed", zap.Error(err))
	}

	// nolint: exhaustivestruct
	app := fiber.New(fiber.Config{
		AppName: "saf",
	})

	handler.Healthz{
		Logger: logger.Named("handler").Named("healthz"),
		Tracer: tracer,
	}.Register(app.Group(""))

	handler.Event{
		CMQ:    cmq,
		Logger: logger.Named("handler").Named("event"),
		Tracer: tracer,
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
func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "producer",
			Short: "gets events from http and produce them into nats",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg, logger, tracer)
			},
		},
	)
}
