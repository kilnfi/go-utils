package sql

import (
	"fmt"
	"time"

	cmdutils "github.com/kilnfi/go-utils/cmd/utils"
	types "github.com/kilnfi/go-utils/common/types"
	kilntls "github.com/kilnfi/go-utils/crypto/tls"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type FlagPrefixer struct {
	cmdutils.FlagPrefixer

	dialect string

	baseDesc string

	tlsFlag kilntls.FlagPrefixer
}

func NewFlagPrefixer(dialect, name string) *FlagPrefixer {
	if name == "" {
		name = dialect
	}

	fl := cmdutils.NewFlagPrefixer(name)

	baseDesc := fmt.Sprintf("%v %v database", name, dialect)

	return &FlagPrefixer{
		FlagPrefixer: fl,
		dialect:      dialect,
		baseDesc:     baseDesc,
		tlsFlag: kilntls.NewFlagPrefixer(
			cmdutils.NewFlagPrefixer(
				fmt.Sprintf("%v-db-ssl", name),
				cmdutils.SeparatorOpt(""),
			),
			baseDesc,
		),
	}
}

func (fl *FlagPrefixer) Flags(v *viper.Viper, f *pflag.FlagSet) {
	fl.UserFlag(v, f)
	fl.PasswordFlag(v, f)
	fl.HostFlag(v, f)
	fl.PortFlag(v, f)
	fl.NameFlag(v, f)
	fl.SSLModeFlag(v, f)
	fl.tlsFlag.CertFlag(v, f)
	fl.tlsFlag.KeyFlag(v, f)
	fl.tlsFlag.CAFlag(v, f)
	fl.ConnectTimeoutFlag(v, f)
}

func (fl *FlagPrefixer) ConfigFromViper(v *viper.Viper) *Config {
	return &Config{
		Dialect:        fl.dialect,
		User:           fl.GetUser(v),
		Password:       fl.GetPassword(v),
		Host:           fl.GetHost(v),
		Port:           fl.GetPort(v),
		DBName:         fl.GetDBName(v),
		SSLMode:        fl.GetSSLMode(v),
		SSLCert:        fl.tlsFlag.GetCert(v),
		SSLKey:         fl.tlsFlag.GetKey(v),
		SSLCA:          fl.tlsFlag.GetCA(v),
		ConnectTimeout: &types.Duration{Duration: fl.GetConnectTimeout(v)},
	}
}

const (
	userFlag     = "db-user"
	userViperKey = "db.user"
	userEnv      = "DB_USER"
)

func (fl *FlagPrefixer) UserFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's user", fl.baseDesc),
		fl.Env(userEnv),
		fl.dialect,
	)

	f.String(fl.FlagName(userFlag), "", desc)
	_ = v.BindPFlag(fl.ViperKey(userViperKey), f.Lookup(fl.FlagName(userFlag)))
	_ = v.BindEnv(fl.ViperKey(userViperKey), fl.Env(userEnv))
	v.SetDefault(fl.ViperKey(userViperKey), fl.dialect)
}

func (fl *FlagPrefixer) GetUser(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(userViperKey))
}

const (
	passwordFlag     = "db-password"
	passwordViperKey = "db.password"
	passwordEnv      = "DB_PASSWORD"
)

func (fl *FlagPrefixer) PasswordFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's password", fl.baseDesc),
		fl.Env(passwordEnv),
		fl.dialect,
	)

	f.String(fl.FlagName(passwordFlag), "", desc)
	_ = v.BindPFlag(fl.ViperKey(passwordViperKey), f.Lookup(fl.FlagName(passwordFlag)))
	_ = v.BindEnv(fl.ViperKey(passwordViperKey), fl.Env(passwordEnv))
	v.SetDefault(fl.ViperKey(passwordViperKey), fl.dialect)
}

func (fl *FlagPrefixer) GetPassword(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(passwordViperKey))
}

const (
	hostFlag     = "db-host"
	hostViperKey = "db.host"
	hostEnv      = "DB_HOST"
	hostDefault  = "localhost:5432"
)

func (fl *FlagPrefixer) HostFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's host", fl.baseDesc),
		fl.Env(hostEnv),
		hostDefault,
	)

	f.String(fl.FlagName(hostFlag), "", desc)
	_ = v.BindPFlag(fl.ViperKey(hostViperKey), f.Lookup(fl.FlagName(hostFlag)))
	_ = v.BindEnv(fl.ViperKey(hostViperKey), fl.Env(hostEnv))
	v.SetDefault(fl.ViperKey(hostViperKey), hostDefault)
}

func (fl *FlagPrefixer) GetHost(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(hostViperKey))
}

const (
	portFlag     = "db-port"
	portViperKey = "db.port"
	portEnv      = "DB_PORT"
	portDefault  = uint16(5432)
)

func (fl *FlagPrefixer) PortFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's port", fl.baseDesc),
		fl.Env(portEnv),
		portDefault,
	)

	f.Uint16(fl.FlagName(portFlag), portDefault, desc)
	_ = v.BindPFlag(fl.ViperKey(portViperKey), f.Lookup(fl.FlagName(portFlag)))
	_ = v.BindEnv(fl.ViperKey(portViperKey), fl.Env(portEnv))
	v.SetDefault(fl.ViperKey(portViperKey), portDefault)
}

func (fl *FlagPrefixer) GetPort(v *viper.Viper) uint16 {
	return uint16(v.GetUint32(fl.ViperKey(portViperKey)))
}

const (
	nameFlag     = "db-name"
	nameViperKey = "db.name"
	nameEnv      = "DB_NAME"
)

func (fl *FlagPrefixer) NameFlag(v *viper.Viper, f *pflag.FlagSet) {
	nameDefault := fl.Prefix()
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's name", fl.baseDesc),
		fl.Env(nameEnv),
		nameDefault,
	)

	f.String(fl.FlagName(nameFlag), "", desc)
	_ = v.BindPFlag(fl.ViperKey(nameViperKey), f.Lookup(fl.FlagName(nameFlag)))
	_ = v.BindEnv(fl.ViperKey(nameViperKey), fl.Env(nameEnv))
	v.SetDefault(fl.ViperKey(nameViperKey), nameDefault)
}

func (fl *FlagPrefixer) GetDBName(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(nameViperKey))
}

const (
	poolSizeFlag     = "db-poolsize"
	poolSizeViperKey = "db.poolsize"
	poolSizeDefault  = 0
	poolSizeEnv      = "DB_POOLSIZE"
)

func (fl *FlagPrefixer) PoolSizeFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's maximum number of connections", fl.baseDesc),
		fl.Env(poolSizeEnv),
		poolSizeDefault,
	)

	f.Int(fl.FlagName(poolSizeFlag), poolSizeDefault, desc)
	_ = viper.BindPFlag(fl.ViperKey(poolSizeViperKey), f.Lookup(fl.FlagName(poolSizeFlag)))
	_ = v.BindEnv(fl.ViperKey(poolSizeViperKey), fl.Env(poolSizeEnv))
	v.SetDefault(fl.ViperKey(poolSizeViperKey), poolSizeDefault)
}

func (fl *FlagPrefixer) GetPoolSize(v *viper.Viper) int {
	return v.GetInt(fl.ViperKey(poolSizeViperKey))
}

const (
	connectTimeoutFlag     = "db-connect-timeout"
	connectTimeoutViperKey = "db.connect-timeout"
	connectTimeoutDefault  = time.Second * 30
	connectTimeoutEnv      = "DB_CONNECT_TIMEOUT"
)

func (fl *FlagPrefixer) ConnectTimeoutFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's time client wait for a free connection", fl.baseDesc),
		fl.Env(connectTimeoutEnv),
		connectTimeoutDefault,
	)

	f.Duration(fl.FlagName(connectTimeoutFlag), connectTimeoutDefault, desc)
	_ = viper.BindPFlag(fl.ViperKey(connectTimeoutViperKey), f.Lookup(fl.FlagName(connectTimeoutFlag)))
	_ = v.BindEnv(fl.ViperKey(connectTimeoutViperKey), fl.Env(connectTimeoutEnv))
	v.SetDefault(fl.ViperKey(connectTimeoutViperKey), connectTimeoutDefault)
}

func (fl *FlagPrefixer) GetConnectTimeout(v *viper.Viper) time.Duration {
	return v.GetDuration(fl.ViperKey(connectTimeoutViperKey))
}

const (
	keepAliveFlag    = "db-keepalive"
	keepAliveKey     = "db.keepalive"
	keepAliveDefault = time.Minute
	keepAliveEnv     = "DB_KEEPALIVE"
)

func (fl *FlagPrefixer) KeepAliveFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("%v's keepalive client connection", fl.baseDesc),
		fl.Env(keepAliveEnv),
		keepAliveDefault,
	)

	f.Duration(fl.FlagName(keepAliveFlag), keepAliveDefault, desc)
	_ = viper.BindPFlag(fl.ViperKey(keepAliveKey), f.Lookup(fl.FlagName(keepAliveFlag)))
	_ = v.BindEnv(fl.ViperKey(keepAliveKey), fl.Env(keepAliveEnv))
	v.SetDefault(fl.ViperKey(keepAliveKey), keepAliveDefault)
}

func (fl *FlagPrefixer) GetKeepAlive(v *viper.Viper) time.Duration {
	return v.GetDuration(fl.ViperKey(keepAliveKey))
}

const (
	disableSSLMode    = "disable"
	requireSSLMode    = "require"
	allowSSLMode      = "allow"
	preferSSLMode     = "prefer"
	verifyCASSLMode   = "verify-ca"
	verifyFullSSLMode = "verify-full"
)

var availableSSLModes = []string{
	disableSSLMode,
	allowSSLMode,
	preferSSLMode,
	requireSSLMode,
	verifyCASSLMode,
	verifyFullSSLMode,
}

const (
	sslModeFlag     = "db-sslmode"
	sslModeViperKey = "db.sslmode"
	sslModeDefault  = disableSSLMode
	sslModeEnv      = "DB_SSLMODE"
)

func (fl *FlagPrefixer) SSLModeFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDescWithDefault(
		fmt.Sprintf("SSL mode to connect to %v\n  Must be one of %q (see https://www.postgresql.org/docs/current/libpq-ssl.html for more information)", fl.baseDesc, availableSSLModes),
		fl.Env(sslModeEnv),
		sslModeDefault,
	)

	f.String(fl.FlagName(sslModeFlag), sslModeDefault, desc)
	_ = viper.BindPFlag(fl.ViperKey(sslModeViperKey), f.Lookup(fl.FlagName(sslModeFlag)))
	_ = v.BindEnv(fl.ViperKey(sslModeViperKey), fl.Env(sslModeEnv))
	v.SetDefault(fl.ViperKey(sslModeViperKey), sslModeDefault)
}

func (fl *FlagPrefixer) GetSSLMode(v *viper.Viper) string {
	return v.GetString(fl.ViperKey(sslModeViperKey))
}
