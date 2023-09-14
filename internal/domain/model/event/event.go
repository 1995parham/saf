package event

import (
	"time"
)

type Event struct {
	Subject   string    `json:"subject"`
	CreatedAt time.Time `json:"createdAt"`
	Payload   []byte    `json:"payload"`
}
