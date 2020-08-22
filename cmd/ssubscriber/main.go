package ssubscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/nats-io/stan.go"
	"github.com/nats-ir/nats101/model"
	"github.com/spf13/cobra"
)

func main(server string, cid string) {
	rand.Seed(time.Now().UnixNano())

	// nolint: gosec
	id := rand.Int63()

	nc, err := stan.Connect(cid, fmt.Sprintf("elahe-%d", id), stan.NatsURL(server))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Connected to %s from %v\n", nc.NatsConn().ConnectedAddr(), nc.NatsConn().DiscoveredServers())

	defer nc.Close()

	ch := make(chan struct{})

	if _, err := nc.Subscribe("message", func(msg *stan.Msg) {
		var m model.Message
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			log.Print(err)

			return
		}

		fmt.Printf("Received a message: %+v\n", m)

		if m.Text == "quit" {
			close(ch)
		}
	}, stan.StartAtTimeDelta(1*time.Minute), stan.DurableName("elahe-subscriber")); err != nil {
		log.Fatal(err)
	}

	<-ch
}

// Register subscriber command.
func Register(root *cobra.Command, server *string) {
	cmd := &cobra.Command{
		Use:   "ssubscriber",
		Short: "Subscribe to messages from streaming NATS",
		Run: func(cmd *cobra.Command, args []string) {
			cid, err := cmd.Flags().GetString("cluster")
			if err != nil {
				log.Printf("invalid cluster argument %s", err)
			}

			main(*server, cid)
		},
	}

	cmd.Flags().StringP("cluster", "c", "nats-ir", "nats streaming cluster-id")

	root.AddCommand(cmd)
}
