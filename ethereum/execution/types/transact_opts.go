package types

import (
	"context"
	"math/big"

	gethbind "github.com/ethereum/go-ethereum/accounts/abi/bind"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

type SignTxFunc func(ctx context.Context, addr gethcommon.Address, tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error)

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

func (opts *TransactOpts) ToOpts(ctx context.Context, chainID *big.Int, signTx SignTxFunc) *gethbind.TransactOpts {
	return &gethbind.TransactOpts{
		Context: ctx,
		Signer: func(addr gethcommon.Address, tx *gethtypes.Transaction) (*gethtypes.Transaction, error) {
			if opts.NoSign {
				return tx, nil
			}

			return signTx(ctx, addr, tx, chainID)
		},
		From:      opts.From,
		Nonce:     opts.Nonce,
		Value:     opts.Value,
		GasPrice:  opts.GasPrice,
		GasFeeCap: opts.GasFeeCap,
		GasTipCap: opts.GasTipCap,
		GasLimit:  opts.GasLimit,
		NoSend:    !opts.Send || opts.NoSign,
	}
}
