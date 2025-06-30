package consumer

import (
	"context"

	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/config"
	"github.com/1995parham/saf/internal/infra/logger"
	"github.com/1995parham/saf/internal/infra/output/manager"
	"github.com/1995parham/saf/internal/infra/telemetry"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main(cmq *cmq.CMQ, logger *zap.Logger, _ *manager.Manager) {
	logger.Info("welcome to consumer application")

	err := cmq.Streams(context.Background())
	if err != nil {
		logger.Fatal("stream creation failed", zap.Error(err))
	}
}

// Register consumer command.
func Register() *cli.Command {
	// nolint: exhaustruct
	return &cli.Command{
		Name:        "consumer",
		Aliases:     []string{"c"},
		Description: "gets events from jetstream",
		Action: func(_ context.Context, _ *cli.Command) error {
			fx.New(
				fx.Provide(config.Provide),
				fx.Provide(logger.Provide),
				fx.Provide(telemetry.Provide),
				fx.Provide(cmq.Provide),
				fx.Provide(manager.Provide),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
