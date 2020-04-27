/*
 *
 * In The Name of God
 *
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 26-04-2020
 * |
 * | File Name:     main.go
 * +===============================================
 */

package producer

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/1995parham/nats101/model"
	"github.com/nats-io/nats.go"
	"github.com/spf13/cobra"
)

func main(server string) {
	nc, err := nats.Connect(server)
	if err != nil {
		log.Fatal(err)
	}

	c, err := nats.NewEncodedConn(nc, nats.GOB_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	defer c.Close()

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")

		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			continue
		}

		line = strings.TrimSuffix(line, "\n")

		splited := strings.SplitN(line, " ", 2)

		var cmd, args string
		if len(splited) > 1 {
			cmd, args = splited[0], splited[1]
		} else {
			cmd = splited[0]
		}

		switch cmd {
		case "send":
			splited := strings.SplitN(args, " ", 2)
			from, message := splited[0], splited[1]

			if err := c.Publish("message", model.Message{
				From:      from,
				Text:      message,
				CreatedAt: time.Now(),
			}); err != nil {
				log.Fatal(err)
			}
		case "exit":
			return
		default:
			fmt.Println("Please enter valid command")
		}
	}
}

// Register producer command.
func Register(root *cobra.Command, server *string) {
	cmd := &cobra.Command{
		Use:   "producer",
		Short: "Produce messages to NATS",
		Run: func(cmd *cobra.Command, args []string) {
			main(*server)
		},
	}

	root.AddCommand(cmd)
}
