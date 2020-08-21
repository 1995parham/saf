package subscriber

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-ir/nats101/model"
	"github.com/spf13/cobra"
)

func main(server string) {
	nc, err := nats.Connect(server)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to %s from %v\n", nc.ConnectedAddr(), nc.DiscoveredServers())

	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	ch := make(chan struct{})

	if _, err := c.Subscribe("message", func(m *model.Message) {
		fmt.Printf("Received a message: %+v\n", m)
		if m.Text == "quit" {
			close(ch)
		}
	}); err != nil {
		log.Fatal(err)
	}

	<-ch
}

// Register subscriber command.
func Register(root *cobra.Command, server *string) {
	cmd := &cobra.Command{
		Use:   "subscriber",
		Short: "Subscribe to messages from NATS",
		Run: func(cmd *cobra.Command, args []string) {
			main(*server)
		},
	}

	root.AddCommand(cmd)
}
