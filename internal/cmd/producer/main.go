package producer

import (
	"context"

	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/config"
	"github.com/1995parham/saf/internal/infra/http/server"
	"github.com/1995parham/saf/internal/infra/logger"
	"github.com/1995parham/saf/internal/infra/telemetry"
	"github.com/gofiber/fiber/v2"
	"github.com/urfave/cli/v3"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main(logger *zap.Logger, _ *fiber.App) {
	logger.Info("welcome to producer application")
}

// Register producer command.
func Register() *cli.Command {
	// nolint: exhaustruct
	return &cli.Command{
		Name:        "producer",
		Aliases:     []string{"p"},
		Description: "gets events from http and produce them into the nats jetstream",
		Action: func(_ context.Context, _ *cli.Command) error {
			fx.New(
				fx.Provide(config.Provide),
				fx.Provide(logger.Provide),
				fx.Provide(telemetry.Provide),
				fx.Provide(cmq.Provide),
				fx.Provide(server.Provide),
				fx.Invoke(main),
			).Run()

			return nil
		},
	}
}
