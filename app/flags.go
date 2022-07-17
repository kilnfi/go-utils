package app

import (
	kilnlog "github.com/kilnfi/go-utils/log"
	kilnhttp "github.com/kilnfi/go-utils/net/http"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var serverFlags, healthFlags kilnhttp.ServerFlagPrefixer

func init() {
	serverFlags = kilnhttp.NewServerFlagPrefixer("Main", ":8080")
	healthFlags = kilnhttp.NewServerFlagPrefixer("Health", ":8081")
}

// Flags register viper compatible pflags for app
func Flags(v *viper.Viper, f *pflag.FlagSet) {
	kilnlog.Flags(v, f)
	serverFlags.Flags(v, f)
	healthFlags.Flags(v, f)
}

// ConfigFromViper construct app Config from viper
func ConfigFromViper(v *viper.Viper) *Config {
	return &Config{
		Logger:  kilnlog.ConfigFromViper(v),
		Server:  serverFlags.ConfigFromViper(v),
		Healthz: healthFlags.ConfigFromViper(v),
	}
}
