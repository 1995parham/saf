package model

import "time"

type Event struct {
	Subject   string
	CreatedAt time.Time
	Payload   []byte
}
