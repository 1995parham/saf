package config

import (
	"github.com/1995parham/saf/internal/channel"
	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/logger"
	"github.com/1995parham/saf/internal/metric"
	telemetry "github.com/1995parham/saf/internal/telemetry/config"
)

// Default return default configuration.
func Default() Config {
	return Config{
		Monitoring: metric.Config{
			Address: ":8080",
			Enabled: true,
		},
		Logger: logger.Config{
			Level: "debug",
		},
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: true,
				Ratio:   1.0,
				Agent: telemetry.Agent{
					Host: "127.0.0.1",
					Port: "6831",
				},
			},
			Profiler: telemetry.Profiler{
				Enabled: false,
				Address: "http://127.0.0.1:4040",
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
