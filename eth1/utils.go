package eth1

import (
	"fmt"
	"math/big"

	gethhexutil "github.com/ethereum/go-ethereum/common/hexutil"
)

// FromBlockNumArg decodes a string into a big.Int block number
func FromBlockNumArg(s string) (*big.Int, error) {
	switch {
	case s == "pending":
		return big.NewInt(-1), nil
	case s == "latest":
		return nil, nil
	default:
		b, err := DecodeBig(s)
		if err != nil {
			return nil, fmt.Errorf("invalid block number: %v", err)
		}
		return b, nil
	}
}

// ToBlockNumArg transforms a big.Int into a block string representation
func ToBlockNumArg(number *big.Int) string {
	switch {
	case number == nil:
		return "latest"
	case number.Cmp(big.NewInt(-1)) == 0:
		return "pending"
	default:
		return EncodeBig(number)
	}
}

// DecodeBig decodes either
// - a hex with 0x prefix
// - a decimal
// - "" (decoded to <nil>)
func DecodeBig(s string) (*big.Int, error) {
	switch {
	case s == "":
		return nil, nil
	case Has0xPrefix(s):
		return gethhexutil.DecodeBig(s)
	default:
		b, ok := new(big.Int).SetString(s, 10)
		if !ok {
			return nil, fmt.Errorf("invalid number %q", s)
		}
		return b, nil
	}
}

// EncodeBig encodes either
// - >0 to a hex with 0x prefix
// - <0 to a hex with -0x prefix
// - <nil> to ""
func EncodeBig(b *big.Int) string {
	switch {
	case b == nil:
		return ""
	default:
		return gethhexutil.EncodeBig(b)
	}
}

// Has0xPrefix returns either input starts with a 0x prefix
func Has0xPrefix(input string) bool {
	return len(input) >= 2 && input[0] == '0' && (input[1] == 'x' || input[1] == 'X')
}
