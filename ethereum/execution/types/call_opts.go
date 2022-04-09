package types

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
)

// CallOpts is a set of option to fine tune a contract call request
type CallOpts struct {
	Pending     bool               // Whether to operate on the pending state or the last known one
	From        gethcommon.Address // Optional the sender address, otherwise the first account is used
	BlockNumber *big.Int           // Optional the block number on which the call should be performed
}
