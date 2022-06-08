package cmd

import (
	"context"
	"fmt"

	execclient "github.com/skillz-blockchain/go-utils/ethereum/execution/client"
	jsonrpchttp "github.com/skillz-blockchain/go-utils/net/jsonrpc/http"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
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
			return execclient.New(EthELConfigFromViper(v).SetDefault())
		}
	}

	v := ViperFromContext(ctx)

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
	EthELAddrFlag(v, cmds.PersistentFlags())

	cmds.AddCommand(newCmdEthELChainID(ethELCtx))
	cmds.AddCommand(newCmdEthELBlockNumber(ethELCtx))

	return cmds
}
func newCmdEthELChainID(ctx *ethELContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "chain-id",
		Short: "Get execution layer chain ID",
		RunE: PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.client.ChainID(ctx)
		}),
	}

	return cmd
}

func newCmdEthELBlockNumber(ctx *ethELContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocknumber",
		Short: "Get execution layer chain's head number",
		RunE: PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.client.BlockNumber(ctx)
		}),
	}

	return cmd
}

const (
	ethELAddrFlag     = "eth-el-addr"
	ethELAddrViperKey = "eth.el-addr"
	ethELAddrEnv      = "ETH_EL_ADDR"
)

// EthELAddrFlag register flag for Eth1 node to connect to
func EthELAddrFlag(v *viper.Viper, f *pflag.FlagSet) {
	desc := fmt.Sprintf(`JSON-RPC address of the Ethereum execution layer node to connect to.
	Environment variable: %q`, ethELAddrEnv)
	f.String(ethELAddrFlag, "", desc)
	_ = v.BindPFlag(ethELAddrViperKey, f.Lookup(ethELAddrFlag))
	_ = v.BindEnv(ethELAddrViperKey, ethELAddrEnv)
}

func GetEthELAddr(v *viper.Viper) string {
	return v.GetString(ethELAddrViperKey)
}

func EthELConfigFromViper(v *viper.Viper) *jsonrpchttp.Config {
	return &jsonrpchttp.Config{
		Address: GetEthELAddr(v),
	}
}
