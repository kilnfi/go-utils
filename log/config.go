package log

type Config struct {
	Level, Format string
}

func (cfg *Config) SetDefault() *Config {
	if cfg.Level == "" {
		cfg.Level = "info"
	}

	if cfg.Format == "" {
		cfg.Format = "json"
	}

	return cfg
}
