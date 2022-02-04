package flag

import (
	"math/big"

	gethcommon "github.com/ethereum/go-ethereum/common"
	"github.com/spf13/pflag"

	"github.com/skillz-blockchain/go-utils/eth1"
)

// TransactOptsVar registers a set of custom flags for eth1.TransactOpts
func TransactOptsVar(f *pflag.FlagSet, txOpts *eth1.TransactOpts) {
	AddressVar(
		f,
		&txOpts.From,
		"from",
		gethcommon.Address{},
		`Optional account used to sign the transaction.
Expects an 0x prefixed Ethereum address`,
	)
	BigIntVarP(
		f,
		&txOpts.Nonce,
		"nonce",
		"n",
		nil,
		`Optional nonce to use for the transaction, if not set then use pending nonce.
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.Value,
		"value",
		"v",
		big.NewInt(0),
		`Optional funds to transfer along the transaction in Wei, if not set then no funds are transfered
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.GasPrice,
		"gas-price",
		"p",
		nil,
		`Optional gas price to use for the transaction execution in Wei. If not set then uses gas price oracle
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.GasFeeCap,
		"gas-fee-cap",
		"f",
		nil,
		`Optional gas fee cap to use for the EIP-1559 transaction execution in Wei. If not set then uses gas price oracle
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	BigIntVarP(
		f,
		&txOpts.GasTipCap,
		"gas-tip-cap",
		"t",
		nil,
		`Optional gas priority tip fee cap to use for the EIP-1559 transaction execution in Wei. If not set then uses gas price oracle
Expects either a decimal or an hex encoded value with 0x prefix`,
	)
	f.Uint64VarP(
		&txOpts.GasLimit,
		"gas-limit",
		"l",
		0,
		`Optional gas limit to set for the transaction execution. If not set then estimate gas
Expects either a decimal or an hex encoded value with 0x prefix`,
	)

	// We want NoSend to be true by default
	f.BoolVar(
		&txOpts.Send,
		"send",
		false,
		"If not set then performs all transaction steps but does not send the transaction",
	)
}

// // txOptsNonceValue is a type implenting pflag.Value interface
// type txOptsNonceValue struct {
// 	txOpts *eth1.TransactOpts
// }

// func (v *txOptsNonceValue) Set(s string) error {
// 	if s != "" {
// 		n, err := eth1.DecodeBig(s)
// 		if err != nil {
// 			return err
// 		}
// 		v.txOpts.Nonce = n
// 	}

// 	return nil
// }

// func (v *txOptsNonceValue) Type() string { return "string" }
// func (v *txOptsNonceValue) String() string {
// 	if v.txOpts.Nonce == nil {
// 		return ""
// 	}
// 	return gethhexutil.EncodeBig(v.txOpts.Nonce)
// }

// // txOptsValueValue is a type implenting pflag.Value interface
// type txOptsValueValue struct {
// 	txOpts *eth1.TransactOpts
// }

// func (v *txOptsValueValue) Set(s string) error {
// 	if s != "" {
// 		val, err := eth1.DecodeBig(s)
// 		if err != nil {
// 			return err
// 		}
// 		v.txOpts.Value = val
// 	}

// 	return nil
// }

// func (v *txOptsValueValue) Type() string { return "string" }
// func (v *txOptsValueValue) String() string {
// 	if v.txOpts.Value == nil {
// 		return ""
// 	}
// 	return gethhexutil.EncodeBig(v.txOpts.Value)
// }

// // txOptsGasPriceValue is a type implenting pflag.Value interface
// type txOptsGasPriceValue struct {
// 	txOpts *eth1.TransactOpts
// }

// func (v *txOptsGasPriceValue) Set(s string) error {
// 	if s != "" {
// 		p, err := eth1.DecodeBig(s)
// 		if err != nil {
// 			return err
// 		}
// 		v.txOpts.GasPrice = p
// 	}

// 	return nil
// }

// func (v *txOptsGasPriceValue) Type() string { return "string" }
// func (v *txOptsGasPriceValue) String() string {
// 	if v.txOpts.GasPrice == nil {
// 		return ""
// 	}
// 	return gethhexutil.EncodeBig(v.txOpts.GasPrice)
// }

// type txOptsGasFeeCapValue struct {
// 	txOpts *eth1.TransactOpts
// }

// func (v *txOptsGasFeeCapValue) Set(s string) error {
// 	if s != "" {
// 		p, err := eth1.DecodeBig(s)
// 		if err != nil {
// 			return err
// 		}
// 		v.txOpts.GasFeeCap = p
// 	}

// 	return nil
// }

// func (v *txOptsGasFeeCapValue) Type() string { return "string" }
// func (v *txOptsGasFeeCapValue) String() string {
// 	if v.txOpts.GasFeeCap == nil {
// 		return ""
// 	}
// 	return gethhexutil.EncodeBig(v.txOpts.GasFeeCap)
// }

// type txOptsGasTipCapValue struct {
// 	txOpts *eth1.TransactOpts
// }

// func (v *txOptsGasTipCapValue) Set(s string) error {
// 	if s != "" {
// 		p, err := eth1.DecodeBig(s)
// 		if err != nil {
// 			return err
// 		}
// 		v.txOpts.GasTipCap = p
// 	}

// 	return nil
// }

// func (v *txOptsGasTipCapValue) Type() string { return "string" }
// func (v *txOptsGasTipCapValue) String() string {
// 	if v.txOpts.GasTipCap == nil {
// 		return ""
// 	}
// 	return gethhexutil.EncodeBig(v.txOpts.GasTipCap)
// }

// type txOptsGasLimitValue struct {
// 	txOpts *eth1.TransactOpts
// }

// func (v *txOptsGasLimitValue) Set(s string) error {
// 	if s != "" {
// 		p, err := eth1.DecodeBig(s)
// 		if err != nil {
// 			return err
// 		}
// 		v.txOpts.GasLimit = p.Uint64()
// 	}

// 	return nil
// }

// func (v *txOptsGasLimitValue) Type() string { return "string" }
// func (v *txOptsGasLimitValue) String() string {
// 	return gethhexutil.EncodeBig(big.NewInt(int64(v.txOpts.GasLimit)))
// }
