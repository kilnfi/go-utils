package types

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// TransactOpts is a set of option to fine tune the creation of a valid Ethereum transaction.
type TransactOpts struct {
	From gethcommon.Address // Ethereum account to send the transaction from

	Nonce     *big.Int // Nonce to use for the transaction execution (nil = use pending state)
	Value     *big.Int // Funds to transfer along the transaction (nil = 0 = no funds)
	GasPrice  *big.Int // Gas price to use for the transaction execution (nil = gas price oracle)
	GasFeeCap *big.Int // Gas fee cap to use for the 1559 transaction execution (nil = gas price oracle)
	GasTipCap *big.Int // Gas priority fee cap to use for the 1559 transaction execution (nil = gas price oracle)
	GasLimit  uint64   // Gas limit to set for the transaction execution (0 = estimate)

	NoSign bool // Do all transact steps and stops before signing
	Send   bool // Do all transact steps and send the transaction (can not be true if NoSign is true)
}
