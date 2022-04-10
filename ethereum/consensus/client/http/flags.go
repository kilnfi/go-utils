package eth2http

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

const (
	addrFlag     = "eth2-addr"
	addrViperKey = "eth2.addr"
	addrEnv      = "ETH2_ADDR"
)

// Eth2Addr register flag for Eth2 node to connect to
func Address(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Beacon address.
Environment variable: %q`, addrEnv)
	f.String(addrFlag, "", desc)
	_ = v.BindPFlag(addrViperKey, f.Lookup(addrFlag))
	_ = v.BindEnv(addrViperKey, addrEnv)
}

func NewConfigFromViper(v *viper.Viper) *Config {
	return &Config{
		Address: v.GetString(addrViperKey),
	}
}
