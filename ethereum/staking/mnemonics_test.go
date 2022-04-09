package staking

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testMnemo = []struct {
	desc, entropy, mnemonic, seed string
}{
	{
		"#1",
		"ff61bfd257776d7d0ab35ffbab3b9c3893cb12e98cf7d6a2",
		"youth assume virus puzzle item salon client hip wing flush train illness device maximum plate page stove because",
		"d1478de1c0bb69e0b779ca6c190f7ab6465d3fdcf53278eda9362f41d1c9cad91b7c791c8b2526c4e3164347a87a930916e85832664bc099846be05aa5d2bccd",
	},
	{
		"#2",
		"ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		"zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo zoo vote",
		"dd48c104698c30cfe2b6142103248622fb7bb0ff692eebb00089b32d22484e1613912f0a5b694407be899ffd31ed3992c456cdf60f5d4564b8ba3f05a69890ad",
	},
	{
		"#3",
		"9a7141f29ada0be6add0c3b1f023fe8733f3afb01872e7d537302b3c6cf3cbe7",
		"omit meat lake cup pass viable resemble blur rapid library zero attack dish style scatter atom treat predict slot filter short ketchup convince vault",
		"af6f5a5aa3b6e561e96c010cc7b02e1106a2c59ffade462bab762b05558e9c416d9825e2d1355c39ddacb72d8c90274df3f65bf5026952827a72ffe334647e57",
	},
}

var testPassword = "TREZOR"

func TestSeed(t *testing.T) {
	for _, tt := range testMnemo {
		t.Run(tt.desc, func(t *testing.T) {
			seed, err := Seed(tt.mnemonic, testPassword)
			require.NoError(t, err)
			assert.Equal(t, tt.seed, hex.EncodeToString(seed))
		})
	}
}
