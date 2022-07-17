package hashicorp

import (
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	cmdutils "github.com/kilnfi/go-utils/cmd/utils"
	kilntls "github.com/kilnfi/go-utils/crypto/tls"
	kilnhttp "github.com/kilnfi/go-utils/net/http"
)

var tlsFlag kilntls.FlagPrefixer

func init() {
	tlsFlag = kilntls.NewFlagPrefixer(cmdutils.NewFlagPrefixer("vault-tls"), "Vault")
}

func Flags(v *viper.Viper, f *pflag.FlagSet) {
	VaultAddress(v, f)
	VaultToken(v, f)
	VaultAuthGitHubToken(v, f)
	VaultMount(v, f)
	tlsFlag.Flags(v, f)
}

func ClientConfigFromViper(v *viper.Viper) *ClientConfig {
	cfg := &ClientConfig{
		Address: v.GetString(VaultAddrViperKey),
		Path:    v.GetString(vaultPathViperKey),
		Auth: &AuthConfig{
			Token:       v.GetString(VaultTokenViperKey),
			GitHubToken: v.GetString(VaultAuthGithubTokenViperKey),
		},
		HTTP: &kilnhttp.ClientConfig{
			Transport: &kilnhttp.TransportConfig{
				TLS: tlsFlag.ConfigFromViper(v),
			},
		},
	}

	return cfg
}

const (
	vaultAddrFlag     = "vault-addr"
	VaultAddrViperKey = "vault.addr"
	vaultAddrEnv      = "VAULT_ADDRESS"
)

func VaultAddress(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		"Vault Address",
		vaultAddrEnv,
	)

	f.String(vaultAddrFlag, "", desc)
	_ = v.BindPFlag(VaultAddrViperKey, f.Lookup(vaultAddrFlag))
	_ = v.BindEnv(VaultAddrViperKey, vaultAddrEnv)
}

const (
	vaultPathFlag     = "vault-path"
	vaultPathViperKey = "vault.path"
	vaultPathDefault  = "secret"
	vaultPathEnv      = "VAULT_PATH"
)

func VaultMount(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		"Vault mount path",
		vaultPathEnv,
	)

	f.String(vaultPathFlag, vaultPathDefault, desc)
	_ = v.BindPFlag(vaultPathViperKey, f.Lookup(vaultPathFlag))
	_ = v.BindEnv(vaultPathViperKey, vaultPathEnv)
	v.SetDefault(vaultPathViperKey, vaultPathDefault)
}

const (
	vaultTokenFlag     = "vault-token"
	VaultTokenViperKey = "vault.token"
	vaultTokenEnv      = "VAULT_TOKEN"
)

func VaultToken(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		"Vault token",
		vaultTokenEnv,
	)

	f.String(vaultTokenFlag, "", desc)
	_ = v.BindPFlag(VaultTokenViperKey, f.Lookup(vaultTokenFlag))
	_ = v.BindEnv(VaultTokenViperKey, vaultTokenEnv)
}

const (
	vaultAuthGithubTokenFlag     = "vault-auth-github-token"
	VaultAuthGithubTokenViperKey = "vault.auth.github-token"
	vaultAuthGithubTokenEnv      = "VAULT_AUTH_GITHUB_TOKEN"
)

func VaultAuthGitHubToken(v *viper.Viper, f *pflag.FlagSet) {
	desc := cmdutils.FlagDesc(
		"Vault GitHub token",
		vaultAuthGithubTokenEnv,
	)

	f.String(vaultAuthGithubTokenFlag, "", desc)
	_ = v.BindPFlag(VaultAuthGithubTokenViperKey, f.Lookup(vaultAuthGithubTokenFlag))
	_ = v.BindEnv(VaultAuthGithubTokenViperKey, vaultAuthGithubTokenEnv)
}
