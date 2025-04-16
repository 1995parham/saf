package cmq

import "time"

type Config struct {
	URL             string        `json:"url,omitempty"              koanf:"url"`
	ArtificialSleep time.Duration `json:"artificial_sleep,omitempty" koanf:"artificial_sleep"`
}

const (
	EventsChannel = "events"
	QueueName     = "saf"
	DurableName   = "saf"
)
