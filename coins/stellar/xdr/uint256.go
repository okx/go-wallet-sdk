package xdr

import (
	"math/big"
)

// String returns a display friendly form of the uint256
func (u Uint256) String() string {
	return new(big.Int).SetBytes(u[:]).String()
}

func (s Uint256) Equals(o Uint256) bool {
	if len(s) != len(o) {
		return false
	}
	for i := 0; i < len(s); i++ {
		if s[i] != o[i] {
			return false
		}
	}
	return true
}
