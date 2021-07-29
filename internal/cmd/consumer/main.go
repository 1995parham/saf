package consumer

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/1995parham/saf/internal/channel"
	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/config"
	"github.com/1995parham/saf/internal/metric"
	"github.com/1995parham/saf/internal/subscriber"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

func main(cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	metric.NewServer(cfg.Monitoring).Start(logger.Named("metrics"))

	c, err := cmq.New(cfg.NATS, logger)
	if err != nil {
		logger.Fatal("nats initiation failed", zap.Error(err))
	}

	if err := c.Streams(); err != nil {
		logger.Fatal("nats stream creation failed", zap.Error(err))
	}

	man := channel.New(logger.Named("channels"), tracer)
	man.Setup(cfg.Channels.Enabled, cfg.Channels.Configurations)

	sub := subscriber.New(c, logger, tracer)

	for _, ch := range man.Channels() {
		sub.RegisterHandler(ch)
	}

	if err := sub.Subcribe(); err != nil {
		logger.Fatal("nats subscription failed", zap.Error(err))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

// Register consumer command.
func Register(root *cobra.Command, cfg config.Config, logger *zap.Logger, tracer trace.Tracer) {
	root.AddCommand(
		// nolint: exhaustivestruct
		&cobra.Command{
			Use:   "consumer",
			Short: "gets events from jetstream",
			Run: func(cmd *cobra.Command, args []string) {
				main(cfg, logger, tracer)
			},
		},
	)
}
