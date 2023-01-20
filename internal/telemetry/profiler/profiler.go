package profiler

import (
	"fmt"
	"log"

	"github.com/1995parham/saf/internal/telemetry/config"
	"github.com/pyroscope-io/pyroscope/pkg/agent/profiler"
)

func Start(cfg config.Profiler, component string) {
	if cfg.Enabled {
		// nolint: exhaustruct
		if _, err := profiler.Start(profiler.Config{
			ApplicationName: fmt.Sprintf("1995parham.saf.%s", component),
			ServerAddress:   cfg.Address,
		}); err != nil {
			log.Printf("failed to start the profiler")
		}
	}
}
