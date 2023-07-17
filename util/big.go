package util

import "math/big"

func ConvertToBigInt(v string) *big.Int {
	b := new(big.Int)
	b.SetString(v, 10)
	return b
}
