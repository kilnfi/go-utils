package gethkeystore

import (
	"context"
	"math/big"

	gethaccounts "github.com/ethereum/go-ethereum/accounts"
	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	keystore "github.com/kilnfi/go-utils/keystore"
)

type KeyStore struct {
	cfg  *Config
	keys *gethkeystore.KeyStore
}

func New(cfg *Config) *KeyStore {
	return &KeyStore{
		cfg:  cfg,
		keys: gethkeystore.NewKeyStore(cfg.Path, gethkeystore.StandardScryptN, gethkeystore.StandardScryptP),
	}
}

func (s *KeyStore) CreateAccount(_ context.Context) (*keystore.Account, error) {
	acc, err := s.keys.NewAccount(s.cfg.Password)
	if err != nil {
		return nil, err
	}

	return &keystore.Account{
		Addr: acc.Address,
		URL:  acc.URL,
	}, nil
}

func (s *KeyStore) SignTx(_ context.Context, addr gethcommon.Address, tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error) {
	return s.keys.SignTxWithPassphrase(
		gethaccounts.Account{Address: addr},
		s.cfg.Password,
		tx,
		chainID,
	)
}
