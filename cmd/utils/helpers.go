package utils

import (
	"encoding/json"

	"github.com/spf13/cobra"
)

func PrintJSON(fnc func(cmd *cobra.Command, args []string) (interface{}, error)) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		res, err := fnc(cmd, args)
		if err != nil {
			return err
		}

		if res == nil {
			return nil
		}

		return json.NewEncoder(cmd.OutOrStdout()).Encode(res)
	}
}
