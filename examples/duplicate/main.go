package main

import (
	"github.com/nats-io/nats.go"
	"github.com/pterm/pterm"
)

const (
	StreamName = "rides"
	Subject    = "ride.finish"
)

func main() {
	clientName := "dispatching-js-examples"
	options := []nats.Option{
		nats.Name(clientName),
	}

	s, _ := pterm.DefaultBigText.WithLetters(pterm.NewLettersFromString("Dispatching")).Srender()
	pterm.DefaultCenter.Println(s)

	nc, err := nats.Connect("127.0.0.1:4222", options...)
	if err != nil {
		pterm.Fatal.Printf("nats connection failed %s\n", err)
	}
	defer nc.Close()

	pterm.DefaultCenter.Printf("NATS connection success\nVersion: %s\nServers: %v",
		nc.ConnectedServerVersion(),
		nc.Servers(),
	)

	js, err := nc.JetStream()
	if err != nil {
		pterm.Fatal.Printf("jetstream connection failed %s\n", err)
	}

	info, err := js.StreamInfo(StreamName)
	if err != nil {
		pterm.Fatal.Printf("threre is an issue for reading %s %s\n", StreamName, err)
	}

	pterm.Success.Printf("%+v\n", info)

	pterm.DefaultParagraph.Println(`
we are going to publish multiple messages with <uniqueId> as their 'Nats-Msg-Id'
so you should get only one new message on your stream
		`)

	for i := 0; i < 10; i++ {
		msg := nats.NewMsg(Subject)
		msg.Header.Add(nats.MsgIdHdr, "uniqueId")

		ack, err := js.PublishMsg(msg)
		if err != nil {
			pterm.Fatal.Printf("publishing message failed failed %s\n", err)
		}

		pterm.Info.Printf("%+v\n", ack)
	}
}
