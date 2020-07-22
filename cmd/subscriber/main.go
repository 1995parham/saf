package subscriber

import (
	"NATS/model"
	"fmt"
	"log"

	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	c := cobra.Command{
		Use:"subscribe",
		Run: func(cmd *cobra.Command, args []string) {
			Subscribe()
		},
	}

	root.AddCommand(
		&c,
	)
}

func Subscribe() {
	nc, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		log.Fatal(err)
	}

	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	ch := make(chan *model.Message)

	if _, err := c.Subscribe("parham", func(m *model.Message) {
		ch<- m
	});err != nil {
		log.Fatal(err)
	}

	fmt.Println(<-ch)
}


