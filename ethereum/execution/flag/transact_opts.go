package flag

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"

	"github.com/kilnfi/go-utils/ethereum/execution/types"
)

// TransactOptsVar registers a set of custom flags for eth1.TransactOpts
func TransactOptsVar(f *pflag.FlagSet, txOpts *types.TransactOpts) {
	AddressVar(
		f,
		&txOpts.From,
		"from",
		gethcommon.Address{},
		`Optional account used to sign the transaction.
Expects an 0x prefixed Ethereum address`,
	)
	BigIntVarP(
		f,
		&txOpts.Nonce,
		"nonce",
		"n",
		nil,
		`Optional nonce to use for the transaction, if not set then use pending nonce.
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.Value,
		"value",
		"v",
		big.NewInt(0),
		`Optional funds to transfer along the transaction in Wei, if not set then no funds are transferred
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.GasPrice,
		"gas-price",
		"p",
		nil,
		`Optional gas price to use for the transaction execution in Wei. If not set then uses gas price oracle
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.GasFeeCap,
		"gas-fee-cap",
		"f",
		nil,
		`Optional gas fee cap to use for the EIP-1559 transaction execution in Wei. If not set then uses gas price oracle
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.GasTipCap,
		"gas-tip-cap",
		"t",
		nil,
		`Optional gas priority tip fee cap to use for the EIP-1559 transaction execution in Wei. If not set then uses gas price oracle
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	f.Uint64VarP(
		&txOpts.GasLimit,
		"gas-limit",
		"l",
		0,
		`Optional gas limit to set for the transaction execution. If not set then estimate gas
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	f.BoolVar(
		&txOpts.Send,
		"send",
		false,
		"If set then performs all transaction steps and send the transaction to the network",
	)
	f.BoolVar(
		&txOpts.NoSign,
		"no-sign",
		false,
		"If set then performs all transaction steps and stop before signing the transaction",
	)
}
