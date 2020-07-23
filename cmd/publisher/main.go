package publisher

import (
	"NATS/service"

	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	c := cobra.Command{
		Use: "publish",
		Run: func(cmd *cobra.Command, args []string) {
			service.Run()
		},
	}

	root.AddCommand(
		&c,
	)
}
