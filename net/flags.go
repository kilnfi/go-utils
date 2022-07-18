package net

import (
	"fmt"

	cmdutils "github.com/kilnfi/go-utils/cmd/utils"
	kilntls "github.com/kilnfi/go-utils/crypto/tls"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type EntrypointFlagPrefixer struct {
	cmdutils.FlagPrefixer

	defaultAddr string
	baseDesc    string

	tlsFlag kilntls.FlagPrefixer
}

func NewEntrypointFlagPrefixer(name, defaultAddr string) *EntrypointFlagPrefixer {
	fl := cmdutils.NewFlagPrefixer(name)

	baseDesc := fmt.Sprintf("%v entrypoint", name)

	return &EntrypointFlagPrefixer{
		FlagPrefixer: fl,
		defaultAddr:  defaultAddr,
		baseDesc:     baseDesc,
		tlsFlag: kilntls.NewFlagPrefixer(
			cmdutils.NewFlagPrefixer(
				fmt.Sprintf("%v-ep-tls", name),
			),
			baseDesc,
		),
	}
}

func (fl *EntrypointFlagPrefixer) Flags(v *viper.Viper, f *pflag.FlagSet) {
	fl.AddrFlag(v, f)
	fl.tlsFlag.Flags(v, f)
}

func (fl *EntrypointFlagPrefixer) ConfigFromViper(v *viper.Viper) *EntrypointConfig {
	return &EntrypointConfig{
		Address:   fl.GetAddr(v),
		TLSConfig: fl.tlsFlag.ConfigFromViper(v),
	}
}

const (
	addrFlag     = "ep-addr"
	addrViperKey = "ep.addr"
	addrEnv      = "EP_ADDR"
)

func (fl *EntrypointFlagPrefixer) AddrFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		fmt.Sprintf("%v's address", fl.baseDesc),
		fl.Env(addrEnv),
	)

	f.String(fl.FlagName(addrFlag), fl.defaultAddr, desc)
	_ = v.BindPFlag(fl.ViperKey(addrViperKey), f.Lookup(fl.FlagName(addrFlag)))
	_ = v.BindEnv(fl.ViperKey(addrViperKey), fl.Env(addrEnv))
	v.SetDefault(fl.ViperKey(addrViperKey), fl.defaultAddr)
}

func (fl *EntrypointFlagPrefixer) GetAddr(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(addrViperKey))
}
