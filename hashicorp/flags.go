package hashicorp

import (
	"fmt"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	kilntls "github.com/skillz-blockchain/go-utils/crypto/tls"
	kilnhttp "github.com/skillz-blockchain/go-utils/net/http"
)

func VaultFlags(v *viper.Viper, f *pflag.FlagSet) {
	VaultAddress(v, f)
	VaultToken(v, f)
	VaultAuthGitHubToken(v, f)
	VaultMount(v, f)
	VaultTLSSkipVerify(v, f)
	VaultTLSCert(v, f)
	VaultTLSKey(v, f)
	VaultTLSCA(v, f)
}

const (
	vaultAddrFlag     = "vault-addr"
	VaultAddrViperKey = "vault.addr"
	vaultAddrEnv      = "VAULT_ADDRESS"
)

func VaultAddress(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Vault Address.
Environment variable: %q`, vaultAddrEnv)
	f.String(vaultAddrFlag, "", desc)
	_ = v.BindPFlag(VaultAddrViperKey, f.Lookup(vaultAddrFlag))
	_ = v.BindEnv(VaultAddrViperKey, vaultAddrEnv)

}

const (
	vaultMountFlag     = "vault-mount"
	vaultMountViperKey = "vault.mount"
	vaultMountDefault  = "secret"
	vaultMountEnv      = "VAULT_MOUNT"
)

func VaultMount(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Vault mount path.
Environment variable: %q`, vaultMountEnv)
	f.String(vaultMountFlag, vaultMountDefault, desc)
	_ = v.BindPFlag(vaultMountViperKey, f.Lookup(vaultMountFlag))
	_ = v.BindEnv(vaultMountViperKey, vaultMountEnv)
	v.SetDefault(vaultMountViperKey, vaultMountDefault)
}

const (
	vaultTokenFlag     = "vault-token"
	VaultTokenViperKey = "vault.token"
	vaultTokenEnv      = "VAULT_TOKEN"
)

func VaultToken(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Vault token.
Environment variable: %q`, vaultTokenEnv)
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
	desc := fmt.Sprintf(`Vault GitHub token.
Environment variable: %q`, vaultAuthGithubTokenEnv)
	f.String(vaultAuthGithubTokenFlag, "", desc)
	_ = v.BindPFlag(VaultAuthGithubTokenViperKey, f.Lookup(vaultAuthGithubTokenFlag))
	_ = v.BindEnv(VaultAuthGithubTokenViperKey, vaultAuthGithubTokenEnv)
}

const (
	vaultTLSSkipVerifyFlag     = "vault-tls-skip-verify"
	vaultTLSSkipVerifyViperKey = "vault.tls.skip.verify"
	vaultTLSSkipVerifyEnv      = "VAULT_TLS_SKIP_VERIFY"
)

func VaultTLSSkipVerify(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Key Manager, disables SSL certificate verification.
Environment variable: %q`, vaultTLSSkipVerifyEnv)
	f.Bool(vaultTLSSkipVerifyFlag, false, desc)
	_ = v.BindPFlag(vaultTLSSkipVerifyViperKey, f.Lookup(vaultTLSSkipVerifyFlag))
	_ = v.BindEnv(vaultTLSSkipVerifyViperKey, vaultTLSSkipVerifyEnv)
}

const (
	vaultTLSCertFlag     = "vault-tls-cert"
	vaultTLSCertViperKey = "vault.tls.cert"
	vaultTLSCertEnv      = "VAULT_TLS_CERT"
)

const (
	vaultTLSKeyFlag     = "vault-tls-key"
	vaultTLSKeyViperKey = "vault.tls.key"
	vaultTLSKeyEnv      = "VAULT_TLS_KEY"
)

const (
	vaultTLSCAFlag     = "vault-tls-ca"
	vaultTLSCAViperKey = "vault.tls.ca"
	vaultTLSCAEnv      = "VAULT_TLS_CA"
)

func VaultTLSCert(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Vault TLS certificate file.
Environment variable: %q`, vaultTLSCertEnv)
	f.String(vaultTLSCertFlag, "", desc)
	_ = v.BindPFlag(vaultTLSCertViperKey, f.Lookup(vaultTLSCertFlag))
	_ = v.BindEnv(vaultTLSCertViperKey, vaultTLSCertEnv)
}

func VaultTLSKey(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Vault TLS key file.
Environment variable: %q`, vaultTLSKeyEnv)
	f.String(vaultTLSKeyFlag, "", desc)
	_ = v.BindPFlag(vaultTLSKeyViperKey, f.Lookup(vaultTLSKeyFlag))
	_ = v.BindEnv(vaultTLSKeyViperKey, vaultTLSKeyEnv)
}

func VaultTLSCA(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Vault TLS certificate CA file.
Environment variable: %q`, vaultTLSCAEnv)
	f.String(vaultTLSCAFlag, "", desc)
	_ = v.BindPFlag(vaultTLSCAViperKey, f.Lookup(vaultTLSCAFlag))
	_ = v.BindEnv(vaultTLSCAViperKey, vaultTLSCAEnv)
}

func NewClientConfigFromViper(v *viper.Viper) *ClientConfig {
	cfg := &ClientConfig{
		Address: v.GetString(VaultAddrViperKey),
		Mount:   v.GetString(vaultMountViperKey),
		Auth: &AuthConfig{
			Token:       v.GetString(VaultTokenViperKey),
			GitHubToken: v.GetString(VaultAuthGithubTokenViperKey),
		},
		HTTP: &kilnhttp.ClientConfig{
			Transport: &kilnhttp.TransportConfig{
				TLS: &kilntls.Config{},
			},
		},
	}

	caPath := v.GetString(vaultTLSCAViperKey)
	if caPath != "" {
		cfg.HTTP.Transport.TLS.CAs = append(cfg.HTTP.Transport.TLS.CAs, &kilntls.CertificateFileCA{Path: caPath})
	}

	certPath := v.GetString(vaultTLSCertViperKey)
	keyPath := v.GetString(vaultTLSKeyViperKey)
	if certPath != "" || keyPath != "" {
		cfg.HTTP.Transport.TLS.Certificates = append(
			cfg.HTTP.Transport.TLS.Certificates,
			&kilntls.CertificateFileKeyPair{CertPath: certPath, KeyPath: keyPath},
		)
	}

	return cfg
}
