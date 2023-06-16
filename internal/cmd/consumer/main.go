package consumer

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/1995parham/saf/internal/channel"
	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/config"
	"github.com/1995parham/saf/internal/subscriber"
	"github.com/1995parham/saf/internal/telemetry"
	"github.com/urfave/cli/v3"
	"go.opentelemetry.io/otel"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger) {
	c, err := cmq.New(cfg.NATS, logger)
	if err != nil {
		logger.Fatal("nats initiation failed", zap.Error(err))
	}

	if err := c.Streams(); err != nil {
		logger.Fatal("nats stream creation failed", zap.Error(err))
	}

	man := channel.New(logger.Named("channels"), otel.GetTracerProvider().Tracer("channels"))
	man.Setup(cfg.Channels.Enabled, cfg.Channels.Configurations)

	sub := subscriber.New(c, logger, otel.GetTracerProvider().Tracer("subscriber"))

	for _, ch := range man.Channels() {
		sub.RegisterHandler(ch)
	}

	if err := sub.Subscribe(); err != nil {
		logger.Fatal("nats subscription failed", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Register consumer command.
func Register(cfg config.Config, logger *zap.Logger) *cli.Command {
	tele := telemetry.New(cfg.Telemetry)
	tele.Run()

	// nolint: exhaustruct
	return &cli.Command{
		Name:        "consumer",
		Aliases:     []string{"c"},
		Description: "gets events from jetstream",
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
