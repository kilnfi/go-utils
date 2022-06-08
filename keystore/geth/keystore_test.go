package gethkeystore

import (
	"context"
	"math/big"
	"testing"

	gethtypes "github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSignTx(t *testing.T) {
	dir := t.TempDir()
	keys := New(&Config{
		Path:     dir,
		Password: "test-pwd",
	})

	acc, err := keys.CreateAccount(context.TODO())
	require.NoError(t, err)

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
