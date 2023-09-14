package output

type Config struct {
	Enabled        []string               `koanf:"enabled"`
	Configurations map[string]interface{} `koanf:"configurations"`
}
