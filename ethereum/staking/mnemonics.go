package staking

import (
	"github.com/tyler-smith/go-bip39"
)

func GenerateRandomMnemonics() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}
	return bip39.NewMnemonic(entropy)
}

func Seed(mnemonic, mnemonicPassphrase string) ([]byte, error) {
	if ok := bip39.IsMnemonicValid(mnemonic); !ok {
		return nil, bip39.ErrInvalidMnemonic
	}

	return bip39.NewSeed(mnemonic, mnemonicPassphrase), nil
}
