package output

type Config struct {
	Configurations map[string]any `json:"configurations,omitempty" koanf:"configurations"`
}
