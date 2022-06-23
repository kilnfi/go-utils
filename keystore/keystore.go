package keystore

import (
	"context"
	"math/big"

	gethaccounts "github.com/ethereum/go-ethereum/accounts"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
)

type Account struct {
	Addr gethcommon.Address `json:"addr"`
	URL  gethaccounts.URL   `json:"url"`
}

type Store interface {
	CreateAccount(context.Context) (*Account, error)
	SignTx(ctx context.Context, addr gethcommon.Address, tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error)
	Import(ctx context.Context, hexkey string) (*Account, error)
}
