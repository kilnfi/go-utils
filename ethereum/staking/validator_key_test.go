package staking

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateValidatorKeys(t *testing.T) {
	// Below test are based on result from the eth2-deposit-cli
	keys, err := GenerateValidatorKeys(
		"zebra sight furnace type elder speak spy beach parent snack million puppy mobile royal ski walnut awful dry culture orphan tourist throw expire shock",
		"",
		5,
		false,
		nil,
	)
	require.NoError(t, err)
	assert.Len(t, keys, 5)
	assert.Equal(
		t,
		"a1778bf6acdc4a2a13ddc736fc5a8cb4558547cf050e8b99619642bc3c1e9cc339a659b46cd13f89db90a74da189c5c7",
		hex.EncodeToString(keys[0].PrivKey.PublicKey().Marshal()),
	)
	assert.Equal(
		t,
		"8b81bd5bb8605fd565c63758a8d729c9f95e453492326d9d9be15bbfca54d0acbc404340de2d2659ac3679dd6a9672ce",
		hex.EncodeToString(keys[1].PrivKey.PublicKey().Marshal()),
	)
	assert.Equal(
		t,
		"8f49b87f3b16fc4f8784a831a48f1164f9a3380fa4dba0da992abd6bf2c78cfc0504a57272d8b28daa77fbd928b0b8ae",
		hex.EncodeToString(keys[2].PrivKey.PublicKey().Marshal()),
	)
	assert.Equal(
		t,
		"ac6a8140b913070ebab4f814cecf25291d5d09c3dabf08b983fa47aa7611d3a1974b0ae484aff218dbfe4d57b3b8232d",
		hex.EncodeToString(keys[3].PrivKey.PublicKey().Marshal()),
	)
	assert.Equal(
		t,
		"922713b9ad7edb0886997ae937e58323b4b4b440e3be77412e25d79827a217722e21aa1ff63e554c1df1dbe959f46e48",
		hex.EncodeToString(keys[4].PrivKey.PublicKey().Marshal()),
	)
}
