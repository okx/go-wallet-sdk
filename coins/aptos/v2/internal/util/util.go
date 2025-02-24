package util

import (
	"encoding/hex"
	"golang.org/x/crypto/sha3"
	"strings"
)

func Sha3256Hash(bytes [][]byte) (output []byte) {
	hasher := sha3.New256()
	for _, b := range bytes {
		hasher.Write(b)
	}
	return hasher.Sum([]byte{})
}

// ParseHex Convenience function to deal with 0x at the beginning of hex strings
func ParseHex(hexStr string) ([]byte, error) {
	if strings.HasPrefix(hexStr, "0x") {
		hexStr = hexStr[2:]
	}
	return hex.DecodeString(hexStr)
}
