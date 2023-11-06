package eos

import (
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/eos/types"
)

func HexToHexBytes(data string) types.HexBytes {
	bytes, _ := hex.DecodeString(data)
	return types.HexBytes(bytes)
}

func HexToChecksum256(data string) types.Checksum256 {
	return types.Checksum256(HexToHexBytes(data))
}
