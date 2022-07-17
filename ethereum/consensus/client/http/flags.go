package eth2http

import (
	cmdutils "github.com/kilnfi/go-utils/cmd/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Flags(v *viper.Viper, f *pflag.FlagSet) {
	EthCLAddrFlag(v, f)
}

func ConfigFromViper(v *viper.Viper) *Config {
	return &Config{
		Address: GetCLAddr(v),
	}
}

const (
	ethCLAddrFlag     = "eth-cl-addr"
	ethCLAddrViperKey = "eth.cl-addr"
	ethCLAddrEnv      = "ETH_CL_ADDR"
)

// EthCLAddrFlag register flag for Eth1 node to connect to
func EthCLAddrFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		"Address of the Ethereum consensus layer node to connect to",
		ethCLAddrEnv,
	)

	f.String(ethCLAddrFlag, "", desc)
	_ = v.BindPFlag(ethCLAddrViperKey, f.Lookup(ethCLAddrFlag))
	_ = v.BindEnv(ethCLAddrViperKey, ethCLAddrEnv)
}

func GetCLAddr(v *viper.Viper) string {
	return v.GetString(ethCLAddrViperKey)
}
