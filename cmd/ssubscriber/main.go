/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 27-04-2020
 * |
 * | File Name:     main.go
 * +===============================================
 */

package ssubscriber

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/1995parham/nats101/model"
	stan "github.com/nats-io/stan.go"
	"github.com/spf13/cobra"
)

func main(server string, cid string) {
	rand.Seed(time.Now().UnixNano())
	id := rand.Int63()

	nc, err := stan.Connect(cid, fmt.Sprintf("elahe-%d", id), stan.NatsURL(server))
	if err != nil {
		log.Fatal(err)
	}

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
			main(*server, cmd.Flags().GetString("cluster"))
		},
	}

	cmd.Flags().StringP("cluster", "c", "elahe", "nats streaming cluster-id")

	root.AddCommand(cmd)
}
