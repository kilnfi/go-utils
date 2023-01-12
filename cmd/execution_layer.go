package cmd

import (
	"context"

	"github.com/kilnfi/go-utils/cmd/utils"
	execclient "github.com/kilnfi/go-utils/ethereum/execution/client/jsonrpc"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ethELContext struct {
	context.Context
	client *execclient.Client
}

// NewCmdEthEL creates the `eth-el` command
func NewCmdEthEL(
	ctx context.Context,
	newELClient func(*viper.Viper) (*execclient.Client, error),
) *cobra.Command {
	ethELCtx := &ethELContext{Context: ctx}

	if newELClient == nil {
		newELClient = func(v *viper.Viper) (*execclient.Client, error) {
			return execclient.New(execclient.ConfigFromViper(v).SetDefault())
		}
	}

	v := utils.ViperFromContext(ctx)

	cmds := &cobra.Command{
		Use:   "eth-el SUBCOMMAND",
		Short: "Commands to interact with Ethereum execution layer node",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			ethELCtx.client, err = newELClient(v)
			return err
		},
	}

	// Register flags
	execclient.EthELAddrFlag(v, cmds.PersistentFlags())

	cmds.AddCommand(newCmdEthELChainID(ethELCtx))
	cmds.AddCommand(newCmdEthELBlockNumber(ethELCtx))

	return cmds
}
func newCmdEthELChainID(ctx *ethELContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-id",
		Short: "Get execution layer chain ID",
		RunE: utils.PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.client.ChainID(ctx)
		}),
	}

	return cmd
}

func newCmdEthELBlockNumber(ctx *ethELContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocknumber",
		Short: "Get execution layer chain's head number",
		RunE: utils.PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.client.BlockNumber(ctx)
		}),
	}

	return cmd
}
