package types

import (
	"encoding/hex"
	"fmt"
)

// BCSable interface for `ToBCS`
type BCSable interface {
	BcsSerialize() ([]byte, error)
}

// ToBCS serialize given `BCSable` into BCS bytes.
// It panics if bcs serialization failed.
func ToBCS(t BCSable) []byte {
	ret, err := t.BcsSerialize()
	if err != nil {
		panic(fmt.Sprintf("bcs serialize failed: %v", err.Error()))
	}
	return ret
}

// ToHex serialize given `BCSable` into BCS bytes and then return hex-encoded string
// It panics if bcs serialization failed.
func ToHex(t BCSable) string {
	return hex.EncodeToString(ToBCS(t))
}
