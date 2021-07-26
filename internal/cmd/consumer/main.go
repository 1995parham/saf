package consumer

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/config"
	"github.com/1995parham/saf/internal/metric"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
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

	if _, err := c.JConn.QueueSubscribe(cmq.EventsChannel, cmq.QueueName, func(msg *nats.Msg) {
		ctx := otel.GetTextMapPropagator().Extract(context.Background(), propagation.HeaderCarrier(msg.Header))

		_, span := tracer.Start(ctx, "subscriber.events")
		defer span.End()

		metadata, _ := msg.Metadata()

		logger.Named("subscriber").Info("receive new message",
			zap.String("timestamp", metadata.Timestamp.String()),
			zap.ByteString("payload", msg.Data),
		)
	}, nats.AckExplicit(), nats.DeliverAll(), nats.Durable(cmq.DurableName)); err != nil {
		logger.Fatal("subscription failed", zap.Error(err))
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
