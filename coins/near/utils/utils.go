package utils

import (
	"fmt"
	"math/big"
)

func BigIntToUintBytes(i *big.Int, bytelen int) ([]byte, error) {
	if i.Sign() < 0 {
		return nil, fmt.Errorf("cannot encode a negative big.Int into an unsigned integer")
	}

	max := big.NewInt(0).Exp(big.NewInt(2), big.NewInt(int64(bytelen*8)), nil)
	if i.CmpAbs(max) > 0 {
		return nil, fmt.Errorf("cannot encode big.Int to []byte: given big.Int exceeds highest number "+
			"%v for an uint with %v bits", max, bytelen*8)
	}

	res := make([]byte, bytelen)

	bs := i.Bytes()
	copy(res[len(res)-len(bs):], bs)
	return res, nil
}

// Reverse reverses bytes in place (manipulates the underlying array)
func Reverse(b []byte) {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
}
