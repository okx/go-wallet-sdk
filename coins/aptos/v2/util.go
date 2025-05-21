package v2

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
)

// -- Note these are copied from internal/util/util.go to prevent package loops, but still allow devs to use it

// ParseHex Convenience function to deal with 0x at the beginning of hex strings
func ParseHex(hexStr string) ([]byte, error) {
	// This had to be redefined separately to get around a package loop
	return util.ParseHex(hexStr)
}

// Sha3256Hash takes a hash of the given sets of bytes
func Sha3256Hash(bytes [][]byte) (output []byte) {
	return util.Sha3256Hash(bytes)
}
