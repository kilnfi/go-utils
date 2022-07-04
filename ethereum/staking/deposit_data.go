package staking

import (
	"encoding/hex"
	"encoding/json"
	"fmt"

	ethcl "github.com/kilnfi/go-utils/ethereum/consensus"
	beaconcommon "github.com/protolambda/zrnt/eth2/beacon/common"
	"github.com/protolambda/ztyp/tree"
	e2types "github.com/wealdtech/go-eth2-types/v2"
)

type DepositData struct {
	beaconcommon.DepositData

	Version beaconcommon.Version
}

func (data *DepositData) Network() string {
	n, _ := ethcl.Network(data.Version)
	return n
}

func (data *DepositData) Sign(
	vkey *ValidatorKey,
) (*DepositData, error) {
	if len(data.Pubkey) == 0 {
		pubKey := new(beaconcommon.BLSPubkey)
		err := pubKey.UnmarshalText(vkey.PrivKey.PublicKey().Marshal())
		if err != nil {
			return nil, err
		}
		data.Pubkey = *pubKey
	} else if "0x"+hex.EncodeToString(vkey.PrivKey.PublicKey().Marshal()) != data.Pubkey.String() {
		return nil, fmt.Errorf("signing keys does not match data public key")
	}

	// Sign DepositMessage root
	sig, err := sign(vkey, data.SigningRoot())
	if err != nil {
		return nil, err
	}

	data.Signature = *sig

	return data, nil
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

func (data *DepositData) VerifySignature() (bool, error) {
	sig, err := e2types.BLSSignatureFromBytes(data.Signature[:])
	if err != nil {
		return false, err
	}

	pubkey, err := e2types.BLSPublicKeyFromBytes(data.Pubkey[:])
	if err != nil {
		return false, err
	}

	root := data.SigningRoot()
	return sig.Verify(root[:], pubkey), nil
}

// Compute DepositMessage root to be signed
func (data *DepositData) SigningRoot() beaconcommon.Root {
	return beaconcommon.ComputeSigningRoot(
		data.MessageRoot(),
		DepositDomain(data.Version),
	)
}

type jsonDepositData struct {
	Pubkey                beaconcommon.BLSPubkey    `json:"pubkey"`
	WithdrawalCredentials beaconcommon.Root         `json:"withdrawal_credentials"`
	Amount                beaconcommon.Gwei         `json:"amount"`
	Signature             beaconcommon.BLSSignature `json:"signature"`
	Version               beaconcommon.Version      `json:"fork_version"`
	Network               string                    `json:"network_name,omitempty"`
	DepositMessageRoot    beaconcommon.Root         `json:"deposit_message_root"`
	DepositDataRoot       beaconcommon.Root         `json:"deposit_data_root"`
}

func (data *DepositData) MarshalJSON() ([]byte, error) {
	d := &jsonDepositData{
		Pubkey:                data.Pubkey,
		WithdrawalCredentials: data.WithdrawalCredentials,
		Amount:                data.Amount,
		Signature:             data.Signature,
		Version:               data.Version,
		Network:               data.Network(),
		DepositMessageRoot:    data.DepositData.MessageRoot(),
		DepositDataRoot:       data.DepositData.HashTreeRoot(tree.GetHashFn()),
	}

	return json.Marshal(d)
}

func (data *DepositData) UnmarshalJSON(b []byte) error {
	d := new(jsonDepositData)
	err := json.Unmarshal(b, d)
	if err != nil {
		return err
	}

	data.Pubkey = d.Pubkey
	data.WithdrawalCredentials = d.WithdrawalCredentials
	data.Amount = d.Amount
	data.Signature = d.Signature
	data.Version = d.Version

	// Validates `deposit_message_root` and `deposit_data_root`
	if (d.DepositMessageRoot != beaconcommon.Root{}) && (d.DepositMessageRoot != data.MessageRoot()) {
		return fmt.Errorf("invalid `deposit_message_root` for `deposit_data`")
	}

	if (d.DepositDataRoot != beaconcommon.Root{}) && (d.DepositDataRoot != data.HashTreeRoot(tree.GetHashFn())) {
		return fmt.Errorf("invalid `deposit_data_root` for `deposit_data`")
	}

	return nil
}

// Return the bls domain for deposit
func DepositDomain(version beaconcommon.Version) beaconcommon.BLSDomain {
	return beaconcommon.ComputeDomain(
		beaconcommon.DOMAIN_DEPOSIT,
		version,
		beaconcommon.Root{},
	)
}
