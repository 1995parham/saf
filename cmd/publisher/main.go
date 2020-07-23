package publisher

import (
	"NATS/model"
	"log"
	"time"

	"github.com/nats-io/go-nats"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	c := cobra.Command{
		Use:"publish",
		Run: func(cmd *cobra.Command, args []string) {
			Publish()
		},
	}

	root.AddCommand(
		&c,
	)
}

func Publish() {
	nc, err := nats.Connect("nats://localhost:4221")
	if err != nil {
		log.Fatal(err)
	}

	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	err = c.Publish("parham", model.Message{
		Message:   "Hello",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
}

