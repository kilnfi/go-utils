//go:build !integration
// +build !integration

package gethkeystore

import (
	"context"
	"math/big"
	"testing"

	gethcommon "github.com/ethereum/go-ethereum/common"
	gethtypes "github.com/ethereum/go-ethereum/core/types"
	keystore "github.com/kilnfi/go-utils/keystore"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInterface(t *testing.T) {
	assert.Implements(t, (*keystore.Store)(nil), new(KeyStore))
}

func TestSignTx(t *testing.T) {
	dir := t.TempDir()
	keys := New(&Config{
		Path:     dir,
		Password: "test-pwd",
	})

	acc, err := keys.CreateAccount(context.TODO())
	require.NoError(t, err)

	ok, err := keys.HasAccount(context.TODO(), acc.Addr)
	require.NoError(t, err)
	assert.True(t, ok)

	tx := gethtypes.NewTx(
		&gethtypes.DynamicFeeTx{},
	)
	v, r, s := tx.RawSignatureValues()
	require.Equal(t, v, big.NewInt(0))
	require.Equal(t, r, big.NewInt(0))
	require.Equal(t, s, big.NewInt(0))
	require.Equal(t, big.NewInt(0), tx.ChainId())

	tx, err = keys.SignTx(context.TODO(), acc.Addr, tx, big.NewInt(1))
	require.NoError(t, err)
	_, r, s = tx.RawSignatureValues()

	assert.NotEqual(t, r, big.NewInt(0))
	assert.NotEqual(t, s, big.NewInt(0))
	assert.Equal(t, big.NewInt(1), tx.ChainId())
}

func TestSignTxMissingAddress(t *testing.T) {
	dir := t.TempDir()
	keys := New(&Config{
		Path:     dir,
		Password: "test-pwd",
	})

	addr := gethcommon.HexToAddress("0x027f72bc0ca063e40577d30d336d52fd7bcc7375")
	ok, err := keys.HasAccount(context.TODO(), addr)
	require.NoError(t, err)
	assert.False(t, ok)

	tx := gethtypes.NewTx(
		&gethtypes.DynamicFeeTx{},
	)
	_, err = keys.SignTx(context.TODO(), addr, tx, big.NewInt(1))
	require.Error(t, err)
	assert.Equal(t, "no key for address \"0x027f72Bc0CA063E40577D30D336D52Fd7bCC7375\"", err.Error())
}
