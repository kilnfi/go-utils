package cmd

import (
	"github.com/kilnfi/go-utils/app"
	consclienthttp "github.com/kilnfi/go-utils/ethereum/consensus/client/http"
	execclient "github.com/kilnfi/go-utils/ethereum/execution/client"
	"github.com/kilnfi/go-utils/hashicorp"
	gethkeystore "github.com/kilnfi/go-utils/keystore/geth"
	"github.com/kilnfi/go-utils/sql"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func NewCmdAllFlags() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-flags",
		Short: "Display all flags (this is a dev command)",
		Run:   func(cmd *cobra.Command, args []string) {},
	}

	// Register flags
	v := viper.New()
	app.Flags(v, cmd.Flags())
	consclienthttp.Flags(v, cmd.Flags())
	execclient.Flags(v, cmd.Flags())
	gethkeystore.Flags(v, cmd.Flags())
	hashicorp.Flags(v, cmd.Flags())
	sql.NewFlagPrefixer("mysql", "My Service").Flags(v, cmd.Flags())

	return cmd
}
