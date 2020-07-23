package subscriber

import (
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
	nc, err := nats.Connect("nats://localhost:4221")
	if err != nil {
		log.Fatal(err)
	}

	defer nc.Close()

	//c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//defer c.Close()

	ch := make(chan *nats.Msg)

	for {
		if _, err := nc.ChanQueueSubscribe("parham", "raha", ch); err != nil {
			log.Fatal(err)
		}

		//if _, err := c.Subscribe("parham", func(m *model.Message) {
		//	ch<- m
		//});err != nil {
		//	log.Fatal(err)
		//}

		m := <-ch

		fmt.Println(string(m.Data))
	}
}


