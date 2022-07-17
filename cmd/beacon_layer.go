package cmd

import (
	"context"

	"github.com/kilnfi/go-utils/cmd/utils"
	consclient "github.com/kilnfi/go-utils/ethereum/consensus/client"
	consclienthttp "github.com/kilnfi/go-utils/ethereum/consensus/client/http"
	"github.com/kilnfi/go-utils/ethereum/consensus/flag"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ethCLContext struct {
	context.Context
	client consclient.Client
}

// NewCmdEthCL creates the `eth-cl` command
func NewCmdEthCL(
	ctx context.Context,
	newCLClient func(*viper.Viper) (consclient.Client, error),
) *cobra.Command {
	ethCLCtx := &ethCLContext{Context: ctx}

	if newCLClient == nil {
		newCLClient = func(v *viper.Viper) (consclient.Client, error) {
			return consclienthttp.NewClient(consclienthttp.ConfigFromViper(v).SetDefault())
		}
	}

	v := utils.ViperFromContext(ctx)

	cmds := &cobra.Command{
		Use:   "eth-cl SUBCOMMAND",
		Short: "Commands to interact with Ethereum consensus layer node",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			var err error
			ethCLCtx.client, err = newCLClient(v)
			return err
		},
	}

	// Register flags
	consclienthttp.Flags(v, cmds.PersistentFlags())

	cmds.AddCommand(newCmdCLSpec(ethCLCtx))
	cmds.AddCommand(newCmdCLGetValidator(ethCLCtx))

	return cmds
}

func newCmdCLSpec(ctx *ethCLContext) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get-spec",
		Short: "Print validator data",
		RunE: utils.PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.client.GetSpec(ctx)
		}),
	}

	return cmd
}

func newCmdCLGetValidator(ctx *ethCLContext) *cobra.Command {
	var (
		validatorID string
		slot        beaconcommon.Slot
	)
	cmd := &cobra.Command{
		Use:   "get-validator",
		Short: "Print validator data",
		RunE: utils.PrintJSON(func(cmd *cobra.Command, args []string) (res interface{}, err error) {
			return ctx.client.GetValidator(ctx, slot.String(), validatorID)
		}),
	}

	cmd.Flags().SortFlags = false

	cmd.Flags().StringVar(&validatorID, "id", "", "Required validator id on beacon chain")
	_ = cmd.MarkFlagRequired("id")
	flag.SlotVarP(cmd.Flags(), &slot, "slot", "s", 0, "Required beacon chain state-id for which to get validator")
	_ = cmd.MarkFlagRequired("slot")

	return cmd
}
