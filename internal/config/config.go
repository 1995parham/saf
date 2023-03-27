package config

import (
	"log"
	"strings"

	"github.com/1995parham/saf/internal/channel"
	"github.com/1995parham/saf/internal/cmq"
	"github.com/1995parham/saf/internal/logger"
	"github.com/1995parham/saf/internal/metric"
	telemetry "github.com/1995parham/saf/internal/telemetry/config"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/structs"
	"github.com/knadh/koanf/v2"
)

const (
	// Prefix indicates environment variables prefix.
	Prefix = "saf_"
)

type (
	// Config holds all configurations.
	Config struct {
		Logger     logger.Config    `koanf:"logger"`
		Telemetry  telemetry.Config `koanf:"telemetry"`
		Monitoring metric.Config    `koanf:"monitoring"`
		NATS       cmq.Config       `koanf:"nats"`
		Channels   channel.Config   `koanf:"channels"`
	}
)

// New reads configuration with viper.
func New() Config {
	var instance Config

	k := koanf.New(".")

	// load default configuration from file
	if err := k.Load(structs.Provider(Default(), "koanf"), nil); err != nil {
		log.Fatalf("error loading default: %s", err)
	}

	// load configuration from file
	if err := k.Load(file.Provider("config.toml"), toml.Parser()); err != nil {
		log.Printf("error loading config.toml: %s", err)
	}

	// load environment variables
	if err := k.Load(env.Provider(Prefix, ".", func(s string) string {
		return strings.ReplaceAll(strings.ToLower(
			strings.TrimPrefix(s, Prefix)), "_", ".")
	}), nil); err != nil {
		log.Printf("error loading environment variables: %s", err)
	}

	if err := k.Unmarshal("", &instance); err != nil {
		log.Fatalf("error unmarshalling config: %s", err)
	}

	log.Printf("following configuration is loaded:\n%+v", instance)

	return instance
}
