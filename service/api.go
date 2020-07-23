package service

import (
	"NATS/model"
	"NATS/subjects"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/go-nats"
)

type Conn struct {
	NatsConn *nats.Conn
}

func (c *Conn) Run() {
	e := echo.New()

	e.GET("/", c.Get)
	e.Logger.Fatal(e.Start(":1373"))
}

func (c *Conn) Get(e echo.Context) error {
	c.Publish()

	return nil
}

func (c *Conn) Publish() {
	ec, err := nats.NewEncodedConn(c.NatsConn, nats.GOB_ENCODER)
	if err != nil {
		log.Fatal(err)
	}

	err = ec.Publish(subjects.Topic, model.Message {
		Message:   "Hello",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
}