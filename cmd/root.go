/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 26-04-2020
 * |
 * | File Name:     root.go
 * +===============================================
 */

package cmd

import (
	"os"

	"github.com/1995parham/nats101/cmd/producer"
	"github.com/1995parham/nats101/cmd/subscriber"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// ExitFailure status code
const ExitFailure = 1

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var root = &cobra.Command{
		Use:   "nats101",
		Short: "Have fun with NATS on Kubernetes",
	}

	producer.Register(root)
	subscriber.Register(root)

	if err := root.Execute(); err != nil {
		logrus.Errorf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
