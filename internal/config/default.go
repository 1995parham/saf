package config

import (
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
			Syslog: logger.Syslog{
				Enabled: false,
				Network: "",
				Address: "",
				Tag:     "",
			},
		},
		Telemetry: telemetry.Config{
			Trace: telemetry.Trace{
				Enabled: false,
				Agent: telemetry.Agent{
					Host: "127.0.0.1",
					Port: "6831",
				},
			},
		},
	}
}
