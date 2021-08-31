package config

type Config struct {
	Trace    `koanf:"trace"`
	Profiler `koanf:"profiler"`
}

type Trace struct {
	Enabled bool `koanf:"enabled"`
	Agent   `koanf:"agent"`
	Ratio   float64 `koanf:"ratio"`
}

type Agent struct {
	Host string `koanf:"host"`
	Port string `koanf:"port"`
}

type Profiler struct {
	Enabled bool   `koanf:"enabled"`
	Address string `koanf:"address"`
}
