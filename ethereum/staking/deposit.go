package staking

import (
	"encoding/hex"
	"encoding/json"

	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/tree"
)

type DepositData struct {
	*beaconcommon.DepositData

	Version beaconcommon.Version
}

func (data *DepositData) MarshalJSON() ([]byte, error) {
	type depositData struct {
		Pubkey                beaconcommon.BLSPubkey    `json:"pubkey"`
		WithdrawalCredentials beaconcommon.Root         `json:"withdrawal_credentials"`
		Amount                beaconcommon.Gwei         `json:"amount"`
		Signature             beaconcommon.BLSSignature `json:"signature"`
		Version               beaconcommon.Version      `json:"fork_version"`
		DepositMessageRoot    beaconcommon.Root         `json:"deposit_message_root"`
		DepositDataRoot       beaconcommon.Root         `json:"deposit_data_root"`
	}

	d := &depositData{
		Pubkey:                data.Pubkey,
		WithdrawalCredentials: data.WithdrawalCredentials,
		Amount:                data.Amount,
		Signature:             data.Signature,
		Version:               data.Version,
		DepositMessageRoot:    data.DepositData.MessageRoot(),
		DepositDataRoot:       data.DepositData.HashTreeRoot(tree.GetHashFn()),
	}

	return json.Marshal(d)
}

func ComputeDepositData(
	vkey *ValidatorKey,
	withdrawalCredentials beaconcommon.Root,
	amount beaconcommon.Gwei,
	version beaconcommon.Version,
) (*DepositData, error) {
	// Initialize DepositData
	pubKey := new(beaconcommon.BLSPubkey)
	err := pubKey.UnmarshalText([]byte(vkey.Pubkey))
	if err != nil {
		return nil, err
	}

	depositData := &beaconcommon.DepositData{
		Pubkey:                *pubKey,
		WithdrawalCredentials: withdrawalCredentials,
		Amount:                amount,
	}

	// Compute domain for deposit
	dom := beaconcommon.ComputeDomain(
		beaconcommon.DOMAIN_DEPOSIT,
		version,
		beaconcommon.Root{},
	)

	// Compute DepositMessage root to be signed
	depositMsgRoot := beaconcommon.ComputeSigningRoot(
		depositData.MessageRoot(),
		dom,
	)

	// Sign DepositMessage root
	sig, err := sign(vkey, depositMsgRoot)
	if err != nil {
		return nil, err
	}

	depositData.Signature = *sig

	return &DepositData{
		DepositData: depositData,
		Version:     version,
	}, nil
}

func sign(vkey *ValidatorKey, root beaconcommon.Root) (*beaconcommon.BLSSignature, error) {
	// Sign DepositMessage root
	sigHexBytes := vkey.PrivKey.Sign(
		root[:],
	).Marshal()

	sigBytes := make([]byte, 2*len(sigHexBytes))
	_ = hex.Encode(sigBytes, sigHexBytes)

	sig := new(beaconcommon.BLSSignature)
	err := sig.UnmarshalText(sigBytes)
	if err != nil {
		return nil, err
	}

	return sig, err
}
