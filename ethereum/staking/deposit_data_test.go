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

func TestInvalidDepositDataMarshalUnmarshal(t *testing.T) {
	tests := []struct {
		desc string
		raw  []byte
	}{
		{
			desc: "invalid deposit_message_root",
			raw: []byte(`{
	"pubkey":"9161cc71f1f70a2a251fe7e820ec288fc47e23ed4d364ddd6728f1a4a742556082b32024942d9d5abb5d1b335e51dd44",
	"withdrawal_credentials":"0008bd79b392ab5a5ebab2be87ffd37d9ee4f4c14a04001ce268d135f4435f4a",
	"amount":32000000000,
	"signature":"93c08b211bd2419847b08be6462a80730c3c00d1dc19483010247357bbaffdb8a5189d4a3acd7c3f2b72e1aa48b40eb315950f5f34c02206b06b146db5aeb93fafad904bee3b3a0d73a8e0346cbb8b9fe2fef17738527beaeb7d6f7f7d0d8bf6",
	"deposit_message_root":"050d2a5175866a78ad51c601ca26fbb1a9bed708dd71896ac7f9641935035ff9",
	"deposit_data_root":"bc6d383f5255e7c0b291fcc2fba2fb617bf55fde2c6dcf9aa1f8de4648b1a514",
	"fork_version":"00001020",
}`),
		},
		{
			desc: "invalid deposit_data_root",
			raw: []byte(`{
	"pubkey":"9161cc71f1f70a2a251fe7e820ec288fc47e23ed4d364ddd6728f1a4a742556082b32024942d9d5abb5d1b335e51dd44",
	"withdrawal_credentials":"0008bd79b392ab5a5ebab2be87ffd37d9ee4f4c14a04001ce268d135f4435f4a",
	"amount":32000000000,
	"signature":"93c08b211bd2419847b08be6462a80730c3c00d1dc19483010247357bbaffdb8a5189d4a3acd7c3f2b72e1aa48b40eb315950f5f34c02206b06b146db5aeb93fafad904bee3b3a0d73a8e0346cbb8b9fe2fef17738527beaeb7d6f7f7d0d8bf6",
	"deposit_message_root":"aae1f11b1cdc047d494959441cabc7db49b1da0180f4f5219e9460ed269d9669",
	"deposit_data_root":"bdd8a5186ba7c9b01c23cb40291bf645d5ce2120d5bfe4b34ac88cc60caa06e7",
	"fork_version":"00001020",
}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			err := json.Unmarshal(tt.raw, new(DepositData))
			assert.Error(t, err)
		})
	}
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

func TestValidateDepositData(t *testing.T) {
	withdrawalCreds := new(beaconcommon.Root)
	err := json.Unmarshal(
		[]byte(`"00ba9a5425d1bb1d6af8223664da5ff39ca2ff110f917a2b0dc73b8f206f28a4"`),
		withdrawalCreds,
	)
	require.NoError(t, err)

	var tests = []struct {
		desc string

		raw             []byte
		withdrawalCreds beaconcommon.Root
		version         beaconcommon.Version
		amount          beaconcommon.Gwei

		expectedErr error
	}{
		{
			desc:            "all fields - valid",
			version:         beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			withdrawalCreds: *withdrawalCreds,
			amount:          beaconcommon.Gwei(32000000000),
			raw: []byte(`{
				"pubkey": "8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021", 
				"withdrawal_credentials": "00ba9a5425d1bb1d6af8223664da5ff39ca2ff110f917a2b0dc73b8f206f28a4", 
				"amount": 32000000000, 
				"signature": "abe8918b786effebb667e592ef43790b83f54fe4fcc6163f03d1fdb08ff4a75ed67d804659bb4adb44605d00dafa7b910b294e815c6e19166457ee9c22dab6b5fcbc6fd092617e8f3b0a462831cfdda0560fc006569b28d76b5efbd226139674",  
				"fork_version": "00000000", 
				"network_name": "mainnet"
			}`),
		},
		{
			desc:            "only pubkey and signature - valid",
			version:         beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			withdrawalCreds: *withdrawalCreds,
			amount:          beaconcommon.Gwei(32000000000),
			raw: []byte(`{
				"pubkey": "8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021", 
				"signature": "abe8918b786effebb667e592ef43790b83f54fe4fcc6163f03d1fdb08ff4a75ed67d804659bb4adb44605d00dafa7b910b294e815c6e19166457ee9c22dab6b5fcbc6fd092617e8f3b0a462831cfdda0560fc006569b28d76b5efbd226139674" 
			}`),
		},
		{
			desc:            "all fields - invalid withdrawal credentials",
			version:         beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			withdrawalCreds: *withdrawalCreds,
			amount:          beaconcommon.Gwei(32000000000),
			raw: []byte(`{
				"pubkey": "8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021", 
				"withdrawal_credentials": "00c363f07c4ba0ec792081858d6756497d69682990a9f84a27622b2887bd2c3c", 
				"amount": 32000000000, 
				"signature": "abe8918b786effebb667e592ef43790b83f54fe4fcc6163f03d1fdb08ff4a75ed67d804659bb4adb44605d00dafa7b910b294e815c6e19166457ee9c22dab6b5fcbc6fd092617e8f3b0a462831cfdda0560fc006569b28d76b5efbd226139674",  
				"fork_version": "00000000", 
				"network_name": "mainnet"
			}`),
			expectedErr: fmt.Errorf("invalid `withdrawal_credentials` 0x00c363f07c4ba0ec792081858d6756497d69682990a9f84a27622b2887bd2c3c at pos 0 (expected 0x00ba9a5425d1bb1d6af8223664da5ff39ca2ff110f917a2b0dc73b8f206f28a4)"),
		},
		{
			desc:            "all fields - invalid version",
			version:         beaconcommon.Version{0x10, 0x00, 0x00, 0x20},
			withdrawalCreds: *withdrawalCreds,
			amount:          beaconcommon.Gwei(32000000000),
			raw: []byte(`{
				"pubkey": "8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021", 
				"withdrawal_credentials": "00ba9a5425d1bb1d6af8223664da5ff39ca2ff110f917a2b0dc73b8f206f28a4", 
				"amount": 32000000000, 
				"signature": "abe8918b786effebb667e592ef43790b83f54fe4fcc6163f03d1fdb08ff4a75ed67d804659bb4adb44605d00dafa7b910b294e815c6e19166457ee9c22dab6b5fcbc6fd092617e8f3b0a462831cfdda0560fc006569b28d76b5efbd226139674",  
				"fork_version": "10000000", 
				"network_name": "mainnet"
			}`),
			expectedErr: fmt.Errorf("invalid `fork_version` 0x10000000 at pos 0 (expected 0x10000020)"),
		},
		{
			desc:            "all fields - invalid amount",
			version:         beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			withdrawalCreds: *withdrawalCreds,
			amount:          beaconcommon.Gwei(16000000000),
			raw: []byte(`{
				"pubkey": "8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021", 
				"withdrawal_credentials": "00ba9a5425d1bb1d6af8223664da5ff39ca2ff110f917a2b0dc73b8f206f28a4", 
				"amount": 32000000000, 
				"signature": "abe8918b786effebb667e592ef43790b83f54fe4fcc6163f03d1fdb08ff4a75ed67d804659bb4adb44605d00dafa7b910b294e815c6e19166457ee9c22dab6b5fcbc6fd092617e8f3b0a462831cfdda0560fc006569b28d76b5efbd226139674",  
				"fork_version": "00000000", 
				"network_name": "mainnet"
			}`),
			expectedErr: fmt.Errorf("invalid `amount` 32000000000 at pos 0 (expected 16000000000)"),
		},
		{
			desc:            "all fields - invalid signature",
			version:         beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			withdrawalCreds: *withdrawalCreds,
			amount:          beaconcommon.Gwei(32000000000),
			raw: []byte(`{
				"pubkey": "8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021", 
				"withdrawal_credentials": "00ba9a5425d1bb1d6af8223664da5ff39ca2ff110f917a2b0dc73b8f206f28a4", 
				"amount": 32000000000, 
				"signature": "ae6e6a8b2ca8988a870e0ce9f4feff916e0c3bcae32ace785e24ea0b4e2a896c152e12a15405df7caeff7694fa993106025e787ebc9faec11b2ab4a1da1ef4268819606f08dfb78151cd900c82e0c119db8c07d7ee04cf2efe3aaccbfc2ec4a8",  
				"fork_version": "00000000", 
				"network_name": "mainnet"
			}`),
			expectedErr: fmt.Errorf("invalid `signature` for `pubkey` at pos 0"),
		},
		{
			desc:            "only pubkey and signature - invalid signature",
			version:         beaconcommon.Version{0x00, 0x00, 0x00, 0x00},
			withdrawalCreds: *withdrawalCreds,
			amount:          beaconcommon.Gwei(32000000000),
			raw: []byte(`{
				"pubkey": "8ab1ee15397f8686d946a84479f58b16ea791a7817e923bc8567e6782f6797ca486b27bfbf1ae4f23d02d5aa54f1b021", 
				"signature": "ae6e6a8b2ca8988a870e0ce9f4feff916e0c3bcae32ace785e24ea0b4e2a896c152e12a15405df7caeff7694fa993106025e787ebc9faec11b2ab4a1da1ef4268819606f08dfb78151cd900c82e0c119db8c07d7ee04cf2efe3aaccbfc2ec4a8" 
			}`),
			expectedErr: fmt.Errorf("invalid `signature` for `pubkey` at pos 0"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			data := new(DepositData)
			err := json.Unmarshal(tt.raw, data)
			require.NoError(t, err)

			err = ValidateDepositData(
				tt.withdrawalCreds,
				tt.version,
				tt.amount,
				data,
			)
			if tt.expectedErr == nil {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			}
		})
	}
}
