package staking

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDepositDataMarshalUnmarshal(t *testing.T) {
	var datas []*DepositData
	f, err := os.Open("testdata/deposit_data.json")
	require.NoError(t, err, "Open")

	err = json.NewDecoder(f).Decode(&datas)
	require.NoError(t, err, "Decode")

	_, err = json.Marshal(datas)
	require.NoError(t, err, "Marshal")
}

// newDepositData creates a DepositData object
func newDepositData(
	pubkey []byte,
	withdrawalCredentials beaconcommon.Root,
	amount beaconcommon.Gwei,
	version beaconcommon.Version,
) (*DepositData, error) {
	pubKey := new(beaconcommon.BLSPubkey)
	err := pubKey.UnmarshalText(pubkey)
	if err != nil {
		return nil, err
	}

	return &DepositData{
		DepositData: beaconcommon.DepositData{
			Pubkey:                *pubKey,
			WithdrawalCredentials: withdrawalCredentials,
			Amount:                amount,
		},
		Version: version,
	}, nil
}

func TestDepositDataSignAndVerifySignature(t *testing.T) {
	vkeys, err := GenerateValidatorKeys(
		"zebra sight furnace type elder speak spy beach parent snack million puppy mobile royal ski walnut awful dry culture orphan tourist throw expire shock",
		"",
		4,
		false,
		nil,
	)
	require.NoError(t, err)
	require.Len(t, vkeys, 4)

	withdrawalCreds := &beaconcommon.Root{}
	err = withdrawalCreds.UnmarshalText([]byte("0100000000000000000000007e654d251da770a068413677967f6d3ea2fea9e4"))
	require.NoError(t, err)

	var tests = []struct {
		desc              string
		vkey              *ValidatorKey
		version           beaconcommon.Version
		amount            beaconcommon.Gwei
		expectedSignature string
	}{
		{
			"32 ETH on mainnet",
			vkeys[0],
			beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			beaconcommon.Gwei(32000000000),
			"0x996d2810d937e70bf546ae3249b05122cb91f784449372a73875225d2023981a927d0d060bc81435d8bb75ff2e2ffd5b043c60fd31c9c658385b25568b2bb3c9b72809d525d11ed7184a099f5251130329f01f24656bcb659f78c29c04d0b63e",
		},
		{
			"32 ETH on prater",
			vkeys[1],
			beaconcommon.Version{0x00, 0x00, 0x10, 0x20},
			beaconcommon.Gwei(32000000000),
			"0xa3d8bebeb089824afdef39e47373119c86dfcc6fcb64460b5567aa0cc98b763b8dd73d508d8c569658ac70dbccf01d1206381c0d309cfa268e39392f1a12e3390b1f5fdaa4dbcf43afd400dc8ca9685fe20963af6dc63b4d8a94aeb6b61451b6",
		},
		{
			"1 ETH on mainnet",
			vkeys[2],
			beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			beaconcommon.Gwei(1000000000),
			"0x8c08c6aaf362891d920e9882e02ca8565ed67c6cad9bcf6055f62c8c32842ff4c50afe9f871dfa1c6ab48a309402898f00b5a98bf62eef9327000bbe833835bc3b500638bc815e7c8ce793acd36e525ea3235d3efde0dcbdb5894fd4528b5e3a",
		},
		{
			"16 ETH on kiln",
			vkeys[3],
			beaconcommon.Version{0x70, 0x00, 0x00, 0x69},
			beaconcommon.Gwei(16000000000),
			"0x8f00c2f2aff0c0ecde2a797ba709accc98833a3a051e538a9f8ff805386363a05dcb5db1401f4f165b53b6ec547919b20ae14de319a07de4d826c3dcf689b2fea59d5ea099f7d1e573d4e1e72f9f9e4b21293b30f978132cf4ffeb2ad95e8580",
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("#%v", i), func(t *testing.T) {
			depositData, err := newDepositData(
				[]byte(tt.vkey.Pubkey),
				*withdrawalCreds,
				tt.amount,
				tt.version,
			)
			require.NoError(t, err)

			depositData, err = depositData.Sign(tt.vkey)
			require.NoError(t, err)

			assert.Equal(t, fmt.Sprintf("0x%v", tt.vkey.Pubkey), depositData.Pubkey.String())
			assert.Equal(t, tt.amount, depositData.Amount)
			assert.Equal(t, *withdrawalCreds, depositData.WithdrawalCredentials)
			assert.Equal(
				t,
				tt.expectedSignature,
				depositData.Signature.String(),
			)

			isValid, err := depositData.VerifySignature()
			require.NoError(t, err)
			assert.True(t, isValid)
		})
	}

}
