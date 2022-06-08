package main

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/skillz-blockchain/go-utils/cmd"
	"github.com/spf13/cobra"
)

func main() {
	cmds := &cobra.Command{
		Use:   "main SUBCOMMAND",
		Short: "Methods to operate on River protocol",
	}

	ctx := context.Background()
	cmds.AddCommand(cmd.NewCmdEthEL(ctx, nil))
	cmds.AddCommand(cmd.NewCmdEthCL(ctx, nil))
	cmds.AddCommand(cmd.NewCmdKeystore(ctx, nil))

	if err := cmds.Execute(); err != nil {
		log.WithError(err).Fatalf("main: execution failed")
	}
}
