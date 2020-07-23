package service

import (
	"NATS/model"
	"NATS/subjects"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/nats-io/go-nats"
)

func Run() {
	e := echo.New()

	e.GET("/", Get)
	e.Logger.Fatal(e.Start(":1373"))
}

func Get(c echo.Context) error {
	Publish()

	return nil
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

	err = c.Publish(subjects.Topic, model.Message {
		Message:   "Hello",
		CreatedAt: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
}