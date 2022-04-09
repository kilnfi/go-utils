package staking

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestKeystoreManager(t *testing.T) {
	var (
		mnemonic = "forest engage two brief ketchup gaze corn approve about lady uncle ball rhythm eternal alley box very evil tribe guard shoulder open venture curve"
		password = "test2022"
	)
	mngr := NewKeystoreManager()

	vkeys, err := mngr.GenerateValidatorKeys(mnemonic, 1, false, nil)
	require.NoError(t, err)
	vkey := vkeys[0]
	require.Len(t, vkeys, 1)
	assert.Equal(
		t,
		"8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021",
		hex.EncodeToString(vkey.PrivKey.PublicKey().Marshal()),
	)

	t.Run("scrypt", func(t *testing.T) {
		ks, err := mngr.EncryptToScryptKeystore(vkey, password)
		require.NoError(t, err)

		decryptedVkey, err := mngr.DecryptFromKeystore(ks, password)
		require.NoError(t, err)

		assert.Equal(t, vkey, decryptedVkey)
	})

	t.Run("pbkdf2", func(t *testing.T) {
		ks, err := mngr.EncryptToPbkdf2Keystore(vkey, password)
		require.NoError(t, err)

		decryptedVkey, err := mngr.DecryptFromKeystore(ks, password)
		require.NoError(t, err)

		assert.Equal(t, vkey, decryptedVkey)
	})
}
