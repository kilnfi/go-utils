package client

import (
	cmdutils "github.com/kilnfi/go-utils/cmd/utils"
	jsonrpchttp "github.com/kilnfi/go-utils/net/jsonrpc/http"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func Flags(v *viper.Viper, f *pflag.FlagSet) {
	EthELAddrFlag(v, f)

}

func ConfigFromViper(v *viper.Viper) *jsonrpchttp.Config {
	return &jsonrpchttp.Config{
		Address: GetEthELAddr(v),
	}
}

const (
	ethELAddrFlag     = "eth-el-addr"
	ethELAddrViperKey = "eth.el-addr"
	ethELAddrEnv      = "ETH_EL_ADDR"
)

// EthELAddrFlag register flag for Eth1 node to connect to
func EthELAddrFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		"JSON-RPC address of the Ethereum execution layer node to connect to",
		ethELAddrEnv,
	)
	f.String(ethELAddrFlag, "", desc)
	_ = v.BindPFlag(ethELAddrViperKey, f.Lookup(ethELAddrFlag))
	_ = v.BindEnv(ethELAddrViperKey, ethELAddrEnv)
}

func GetEthELAddr(v *viper.Viper) string {
	return v.GetString(ethELAddrViperKey)
}
