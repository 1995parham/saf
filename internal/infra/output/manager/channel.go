package manager

import (
	"github.com/1995parham/saf/internal/infra/output"
	"github.com/1995parham/saf/internal/infra/output/mqtt"
	"github.com/1995parham/saf/internal/infra/output/printer"
)

// list of available channles, please add each channel into this list to make them available.
var channels = []output.Channel{
	new(printer.Printer),
	new(mqtt.MQTT),
}
