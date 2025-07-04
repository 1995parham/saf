package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/1995parham/saf/internal/cmd/consumer"
	"github.com/1995parham/saf/internal/cmd/producer"
	"github.com/urfave/cli/v3"
)

// ExitFailure status code.
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// nolint: exhaustruct
	root := &cli.Command{
		Name:        "saf",
		Description: "Using NATS Jetstream as queue manager to replace RabbitMQ, etc.",
		Authors: []any{
			"Parham Alvani <parham.alvani@gmail.com>",
			"Elahe Dastan <elahe.dstn@gmail.com>",
		},
		Version: func() string {
			revision := ""
			timestamp := ""
			modified := ""

			if info, ok := debug.ReadBuildInfo(); ok {
				for _, setting := range info.Settings {
					switch setting.Key {
					case "vcs.revision":
						revision = setting.Value
					case "vcs.time":
						timestamp = setting.Value
					case "vcs.modified":
						modified = setting.Value
					}
				}
			}

			if revision == "" {
				return ""
			}

			if modified == "true" {
				return fmt.Sprintf("%s (%s) [dirty]", revision, timestamp)
			}

			return fmt.Sprintf("%s (%s)", revision, timestamp)
		}(),
		Commands: []*cli.Command{
			producer.Register(),
			consumer.Register(),
		},
	}

	err := root.Run(context.Background(), os.Args)
	if err != nil {
		log.Printf("failed to execute root command %s", err)

		os.Exit(ExitFailure)
	}
}
