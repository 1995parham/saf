package cmq

import (
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type Config struct {
	URL             string        `json:"url,omitempty"              koanf:"url"`
	ArtificialSleep time.Duration `json:"artificial_sleep,omitempty" koanf:"artificial_sleep"`
	Events          Stream        `json:"events"                     koanf:"events"`
}

type Stream struct {
	Storage  jetstream.StorageType `json:"storage,omitempty"  koanf:"storage"`
	Replicas int                   `json:"replicas,omitempty" koanf:"replicas"`
}

const (
	EventsChannel = "events"
	QueueName     = "saf"
	DurableName   = "saf"
)
