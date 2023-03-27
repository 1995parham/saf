package publisher

import (
	"NATS/conn"
	"NATS/service"

	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	c := cobra.Command{
		Use: "publish",
		Run: func(cmd *cobra.Command, args []string) {
			nc := conn.Conn()
			api := service.Conn{NatsConn: nc}
			api.Run()
		},
	}

	root.AddCommand(
		&c,
	)
}
