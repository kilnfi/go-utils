package flag

import (
	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"

	"github.com/skillz-blockchain/go-utils/ethereum/execution/types"
)

// CallOptsVar registers a set of custom flags for eth1.CallOpts
func CallOptsVar(f *pflag.FlagSet, callOpts *types.CallOpts) {
	BlockNumberVarP(
		f,
		&callOpts.BlockNumber,
		"block",
		"b",
		nil,
		"Optional the block number on which the call should be performed",
	)
	f.BoolVar(
		&callOpts.Pending,
		"pending",
		false,
		"Optional whether to operate on the pending state or the last known one",
	)
	AddressVar(
		f,
		&callOpts.From,
		"from",
		gethcommon.Address{},
		"Optional call's sender address in hex format with 0x prefix",
	)
}
