package cmd

import (
	"os"

	"github.com/1995parham/saf/internal/cmd/consumer"
	"github.com/1995parham/saf/internal/cmd/producer"
	"github.com/1995parham/saf/internal/config"
	"github.com/1995parham/saf/internal/logger"
	"github.com/1995parham/saf/internal/telemetry/trace"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cfg := config.New()

	logger := logger.New(cfg.Logger)

	tracer := trace.New(cfg.Telemetry.Trace)

	// nolint: exhaustruct
	root := &cli.App{
		Name:        "saf",
		Description: "Using NATS Jetstream as queue manager to replace RabbitMQ, etc.",
		Authors: []any{
			"Parham Alvani <parham.alvani@gmail.com>",
			"Elahe Dastan <elahe.dstn@gmail.com>",
		},
		Commands: []*cli.Command{
			producer.Register(cfg, logger, tracer),
			consumer.Register(cfg, logger, tracer),
		},
	}

	if err := root.Run(os.Args); err != nil {
		logger.Error("failed to execute root command", zap.Error(err))

		os.Exit(ExitFailure)
	}
}
