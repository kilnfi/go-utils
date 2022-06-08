package cmd

import (
	"context"
	"fmt"

	"github.com/skillz-blockchain/go-utils/keystore"
	gethkeystore "github.com/skillz-blockchain/go-utils/keystore/geth"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type keystoreContext struct {
	context.Context
	keys keystore.Store
}

// NewCmdKeystore creates the `eth-el` command
func NewCmdKeystore(
	ctx context.Context,
	newKeystore func(*viper.Viper) (keystore.Store, error),
) *cobra.Command {
	keystoreCtx := &keystoreContext{Context: ctx}

	if newKeystore == nil {
		newKeystore = func(v *viper.Viper) (keystore.Store, error) { //nolint
			return gethkeystore.New(KeystoreConfigFromViper(v).SetDefault()), nil
		}
	}

	v := ViperFromContext(ctx)

	cmds := &cobra.Command{
		Use:   "eth1keys SUBCOMMAND",
		Short: "Commands to securely manage Ethereum execution layer keys",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			keystoreCtx.keys, err = newKeystore(v)
			return err
		},
	}

	// Register flags
	KeystoreFlag(v, cmds.PersistentFlags())

	cmds.AddCommand(newCmdGenerateEth1Key(keystoreCtx))

	return cmds
}

func newCmdGenerateEth1Key(ctx *keystoreContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate an Ethereum execution layer account",
		RunE: PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.keys.CreateAccount(ctx)
		}),
	}

	return cmd
}

func KeystoreFlag(v *viper.Viper, f *pflag.FlagSet) {
	KeystorePathFlag(v, f)
	KeystorePasswordFlag(v, f)
}

const (
	keyStorePathFlag     = "keystore-path"
	keyStorePathViperKey = "keystore.path"
	keyStorePathEnv      = "KEYSTORE_PATH"
)

// KeystorePathFlag register flag for the path to the file keystore
func KeystorePathFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Directory where to store keys.
Environment variable: %q`, keyStorePathEnv)
	f.String(keyStorePathFlag, "", desc)
	_ = v.BindPFlag(keyStorePathViperKey, f.Lookup(keyStorePathFlag))
	_ = v.BindEnv(keyStorePathViperKey, keyStorePathEnv)
}

func GetKeystorePath(v *viper.Viper) string {
	return v.GetString(keyStorePathViperKey)
}

const (
	keyStorePasswordFlag     = "keystore-password"
	keyStorePasswordViperKey = "keystore.password"
	keyStorePasswordEnv      = "KEYSTORE_PASSWORD"
)

// KeystorePasswordFlag register flag for the password used to encrypt keys in keystore
func KeystorePasswordFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`Password used to encrypt key files.
Environment variable: %q`, keyStorePasswordEnv)
	f.String(keyStorePasswordFlag, "", desc)
	_ = v.BindPFlag(keyStorePasswordViperKey, f.Lookup(keyStorePasswordFlag))
	_ = v.BindEnv(keyStorePasswordViperKey, keyStorePasswordEnv)
}

func GetKeystorePassword(v *viper.Viper) string {
	return v.GetString(keyStorePasswordViperKey)
}

func KeystoreConfigFromViper(v *viper.Viper) *gethkeystore.Config {
	return &gethkeystore.Config{
		Path:     GetKeystorePath(v),
		Password: GetKeystorePassword(v),
	}
}
