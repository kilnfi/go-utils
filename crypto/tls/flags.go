package tls

import (
	"fmt"

	cmdutils "github.com/kilnfi/go-utils/cmd/utils"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagPrefixer struct {
	cmdutils.FlagPrefixer

	desc string
}

func NewFlagPrefixer(fl cmdutils.FlagPrefixer, desc string) FlagPrefixer {
	return FlagPrefixer{
		FlagPrefixer: fl,
		desc:         desc,
	}
}

func (fl *FlagPrefixer) Flags(v *viper.Viper, f *pflag.FlagSet) {
	fl.CertFlag(v, f)
	fl.KeyFlag(v, f)
	fl.CAFlag(v, f)
	fl.SkipVerify(v, f)
}

func (fl *FlagPrefixer) ConfigFromViper(v *viper.Viper) *Config {
	var cfg *Config

	certPath := fl.GetCert(v)
	keyPath := fl.GetKey(v)
	if certPath != "" || keyPath != "" {
		cfg = new(Config)
		cfg.Certificates = append(
			cfg.Certificates,
			&CertificateFileKeyPair{
				CertPath: certPath,
				KeyPath:  keyPath,
			})
	}

	caPath := fl.GetCA(v)
	if caPath != "" {
		if cfg == nil {
			cfg = new(Config)
		}
		cfg.CAs = append(cfg.CAs, &CertificateFileCA{Path: caPath})
	}

	if cfg != nil {
		cfg.InsecureSkipVerify = fl.GetTLSSkipVerify(v)
	}

	return cfg
}

const (
	certFlag     = "cert"
	certViperKey = "cert"
	certEnv      = "CERT"
)

func (fl *FlagPrefixer) CertFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		fmt.Sprintf("%v's TLS certificate file path (accepts PEM format)", fl.desc),
		fl.Env(certEnv),
	)

	f.String(fl.FlagName(certFlag), "", desc)
	_ = viper.BindPFlag(fl.ViperKey(certViperKey), f.Lookup(fl.FlagName(certFlag)))
	_ = v.BindEnv(fl.ViperKey(certViperKey), fl.Env(certEnv))
}

func (fl *FlagPrefixer) GetCert(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(certViperKey))
}

const (
	keyFlag     = "key"
	keyViperKey = "key"
	keyEnv      = "KEY"
)

func (fl *FlagPrefixer) KeyFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		fmt.Sprintf("%v's TLS private key file path (accepts PEM format)", fl.desc),
		fl.Env(keyEnv),
	)

	f.String(fl.FlagName(keyFlag), "", desc)
	_ = viper.BindPFlag(fl.ViperKey(keyViperKey), f.Lookup(fl.FlagName(keyFlag)))
	_ = v.BindEnv(fl.ViperKey(keyViperKey), fl.Env(keyEnv))
}

func (fl *FlagPrefixer) GetKey(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(keyViperKey))
}

const (
	caFlag     = "ca"
	caViperKey = "ca"
	caEnv      = "CA"
)

func (fl *FlagPrefixer) CAFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		fmt.Sprintf("%v's TLS certificate authority file path (accepts PEM format)", fl.desc),
		fl.Env(caEnv),
	)

	f.String(fl.FlagName(caFlag), "", desc)
	_ = viper.BindPFlag(fl.ViperKey(caViperKey), f.Lookup(fl.FlagName(caFlag)))
	_ = v.BindEnv(fl.ViperKey(caViperKey), fl.Env(caEnv))
}

func (fl *FlagPrefixer) GetCA(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(caViperKey))
}

const (
	skipVerifyFlag     = "skip-verify"
	skipVerifyViperKey = "skip-verify"
	skipVerifyEnv      = "SKIP_VERIFY"
)

func (fl *FlagPrefixer) SkipVerify(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		fmt.Sprintf("Disable TLS certificate verification on %v", fl.desc),
		fl.Env(skipVerifyEnv),
	)

	f.Bool(fl.FlagName(skipVerifyFlag), false, desc)
	_ = v.BindPFlag(fl.ViperKey(skipVerifyViperKey), f.Lookup(fl.FlagName(skipVerifyFlag)))
	_ = v.BindEnv(fl.ViperKey(skipVerifyViperKey), fl.Env(skipVerifyEnv))
}

func (fl *FlagPrefixer) GetTLSSkipVerify(v *viper.Viper) bool {
	return v.GetBool(fl.ViperKey(skipVerifyViperKey))
}
