package gethkeystore

import (
	"context"
	"fmt"
	"math/big"

	gethaccounts "github.com/ethereum/go-ethereum/accounts"
	gethkeystore "github.com/ethereum/go-ethereum/accounts/keystore"
	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	gethcrypto "github.com/ethereum/go-ethereum/crypto"
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

// Import a secp256k1 private key in hexadecimal format
func (s *KeyStore) Import(_ context.Context, hexkey string) (*keystore.Account, error) {
	priv, err := gethcrypto.HexToECDSA(hexkey)
	if err != nil {
		return nil, err
	}

	acc, err := s.keys.ImportECDSA(priv, s.cfg.Password)
	if err != nil {
		return nil, err
	}

	return &keystore.Account{
		Addr: acc.Address,
		URL:  acc.URL,
	}, nil
}

func (s *KeyStore) SignTx(_ context.Context, addr gethcommon.Address, tx *gethtypes.Transaction, chainID *big.Int) (*gethtypes.Transaction, error) {
	if !s.keys.HasAddress(addr) {
		return nil, fmt.Errorf("no key for address %q", addr.String())
	}
	return s.keys.SignTxWithPassphrase(
		gethaccounts.Account{Address: addr},
		s.cfg.Password,
		tx,
		chainID,
	)
}

func (s *KeyStore) HasAccount(_ context.Context, addr gethcommon.Address) (bool, error) {
	return s.keys.HasAddress(addr), nil
}
