package cmq

type Config struct {
	URL string
}

const (
	EventsChannel = "events"
	QueueName     = "saf"
	DurableName   = "saf"
)
