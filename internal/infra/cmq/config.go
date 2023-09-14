package cmq

type Config struct {
	URL string `json:"url,omitempty" koanf:"url"`
}

const (
	EventsChannel = "events"
	QueueName     = "saf"
	DurableName   = "saf"
)
