package cmq

type Config struct {
	URL string `koanf:"url"`
}

const (
	EventsChannel = "events"
	QueueName     = "saf"
	DurableName   = "saf"
)
