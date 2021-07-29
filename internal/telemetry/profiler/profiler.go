package profiler

import (
	"log"

	"github.com/1995parham/saf/internal/telemetry/config"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
)

func Start(cfg config.Profiler) {
	if cfg.Enabled {
		// nolint: exhaustivestruct
		if _, err := profiler.Start(profiler.Config{
			ApplicationName: "1995parham/saf",
			ServerAddress:   cfg.Address,
		}); err != nil {
			log.Printf("failed to start the profiler")
		}
	}
}
