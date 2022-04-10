package staking

import (
	"encoding/hex"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/skillz-blockchain/go-utils/crypto/bls" // nolint
	e2types "github.com/wealdtech/go-eth2-types/v2"
	util "github.com/wealdtech/go-eth2-util"
)

type ValidatorKey struct {
	UUID string

	PrivKey *e2types.BLSPrivateKey

	Pubkey             string
	MnemonicPassphrase string
	MnemonicPassword   string
	Path               string

	Desc string
}

func GenerateValidatorKey(seed []byte, path, desc string) (*ValidatorKey, error) {
	privKey, err := util.PrivateKeyFromSeedAndPath(seed, path)
	if err != nil {
		return nil, err
	}

	return &ValidatorKey{
		UUID:    uuid.New().String(),
		PrivKey: privKey,
		Pubkey:  hex.EncodeToString(privKey.PublicKey().Marshal()),
		Path:    path,
		Desc:    desc,
	}, nil
}

func GenerateValidatorKeys(mnemonicPassphrase, mnemonicPassword string, count int, storeMnemo bool, cb func(string) error) (keys []*ValidatorKey, err error) {
	seed, err := Seed(mnemonicPassphrase, mnemonicPassword)
	if err != nil {
		return nil, err
	}

	keys = make([]*ValidatorKey, count)
	for i := 0; i < count; i++ {
		keys[i], err = GenerateValidatorKey(
			seed,
			fmt.Sprintf("m/12381/3600/%d/0/0", i), // Set path as EIP-2334 format (c.f https://eips.ethereum.org/EIPS/eip-2334)
			"",
		)

		if storeMnemo {
			keys[i].MnemonicPassphrase = mnemonicPassphrase
			keys[i].MnemonicPassword = mnemonicPassword
		}

		if err != nil {
			return nil, err
		}

		if cb != nil {
			err = cb(keys[i].Pubkey)
			if err != nil {
				return nil, err
			}
		}
	}

	return keys, nil
}

func ValidatorKeyFromBytes(privkey []byte) (*ValidatorKey, error) {
	pkey, err := e2types.BLSPrivateKeyFromBytes(privkey)
	if err != nil {
		return nil, err
	}

	return &ValidatorKey{
		PrivKey: pkey,
		Pubkey:  hex.EncodeToString(pkey.PublicKey().Marshal()),
	}, nil
}

func ValidatorKeyFromString(privkey string) (*ValidatorKey, error) {
	b, err := hex.DecodeString(privkey)
	if err != nil {
		return nil, err
	}

	return ValidatorKeyFromBytes(b)
}
