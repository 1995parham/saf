package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/pterm/pterm"
)

const (
	StreamName = "rides"
	Subject    = "ride.finished"
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
we are going subcribe on given stream with given subject, you must see your messages
on the console.
by using deliver last we are going to receive the very last message in the stream.
		`)

	sub, err := js.QueueSubscribe(
		Subject,
		fmt.Sprintf("dispatching-js-examples-%s", strings.ReplaceAll(Subject, ".", "__")),
		func(msg *nats.Msg) {
			pterm.Info.Printf("\nmessage from %s\n", msg.Subject)

			meta, err := msg.Metadata()
			if err != nil {
				return
			}

			pterm.Info.Printf("%+v\n", meta)
		},
		nats.AckWait(30*time.Second),
		nats.ManualAck(),
		nats.MaxDeliver(3),
		nats.DeliverLast(),
	)
	if err != nil {
		pterm.Fatal.Printf("threre is an issue for subscribing on %s %s\n", Subject, err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit

	sub.Unsubscribe()
}
