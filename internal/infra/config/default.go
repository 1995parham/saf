package config

import (
	"github.com/1995parham/saf/internal/channel"
	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/logger"
	"github.com/1995parham/saf/internal/telemetry"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Logger: logger.Config{
			Level: "debug",
		},
		Telemetry: telemetry.Config{
			Namespace:   "1995parham.me",
			ServiceName: "saf",
			Meter: telemetry.Meter{
				Address: ":8080",
				Enabled: true,
			},
			Trace: telemetry.Trace{
				Enabled:  false,
				Ratio:    1.0,
				Endpoint: "127.0.0.1:4317",
			},
		},
		NATS: cmq.Config{
			URL: "nats://127.0.0.1:4222",
		},
		Channels: channel.Config{
			Enabled: []string{
				"printer",
			},
			Configurations: map[string]interface{}{},
		},
	}
}
