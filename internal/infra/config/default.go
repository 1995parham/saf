package config

import (
	"github.com/1995parham/saf/internal/infra/cmq"
	"github.com/1995parham/saf/internal/infra/logger"
	"github.com/1995parham/saf/internal/infra/output"
	"github.com/1995parham/saf/internal/infra/telemetry"
)

// Default return default configuration.
func Default() Config {
	// nolint: exhaustruct, mnd
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
				Enabled:  true,
				Endpoint: "127.0.0.1:4317",
				Ratio:    0.1,
			},
		},
		NATS: cmq.Config{
			URL: "nats://127.0.0.1:4222",
			ArtificialSleep: 0,
		},
		Channels: output.Config{
			Configurations: map[string]any{
				"printer": nil,
			},
		},
	}
}
