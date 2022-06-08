package gethkeystore

type Config struct {
	Path     string `json:"path"`
	Password string `json:"-"`
}

func (cfg *Config) SetDefault() *Config {
	if cfg.Path == "" {
		cfg.Path = "keystore"
	}

	return cfg
}
