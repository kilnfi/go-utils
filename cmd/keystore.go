package cmd

import (
	"context"

	"github.com/kilnfi/go-utils/cmd/utils"
	"github.com/kilnfi/go-utils/keystore"
	gethkeystore "github.com/kilnfi/go-utils/keystore/geth"
	"github.com/spf13/cobra"
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
			return gethkeystore.New(gethkeystore.ConfigFromViper(v).SetDefault()), nil
		}
	}

	v := utils.ViperFromContext(ctx)

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
	gethkeystore.Flags(v, cmds.PersistentFlags())

	cmds.AddCommand(newCmdGenerateEth1Key(keystoreCtx))
	cmds.AddCommand(newCmdImportEth1Key(keystoreCtx))

	return cmds
}

func newCmdGenerateEth1Key(ctx *keystoreContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate an Ethereum execution layer account",
		RunE: utils.PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.keys.CreateAccount(ctx)
		}),
	}

	return cmd
}

func newCmdImportEth1Key(ctx *keystoreContext) *cobra.Command {
	var (
		pkey string
	)

	cmd := &cobra.Command{
		Use:   "import",
		Short: "Import the given private key",
		RunE: utils.PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.keys.Import(ctx, pkey)
		}),
	}

	cmd.Flags().StringVar(&pkey, "priv-key", "", "secp256k1 private key (in hexadecimal format)")
	_ = cmd.MarkFlagRequired("priv-key")

	return cmd
}
