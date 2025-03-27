package xdr

import "encoding/hex"

func (h Hash) HexString() string {
	return hex.EncodeToString(h[:])
}

func (s Hash) Equals(o Hash) bool {
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
