package keystore

import (
	"github.com/herumi/bls-eth-go-binary/bls"
)

func init() {
	if err := bls.Init(bls.BLS12_381); err != nil {
		panic(err)
	}
}
