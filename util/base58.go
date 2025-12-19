package util

import "github.com/okx/go-wallet-sdk/crypto/base58"

// Decode base58 check encoded string with version bytes
func CheckDecodeDoubleV(s string) (result []byte, version [2]byte, err error) {
	decoded, vByte, err := base58.CheckDecode(s)
	switch err {
	case base58.ErrChecksum:
		err = base58.ErrChecksum
		return
	case base58.ErrInvalidFormat:
		err = base58.ErrInvalidFormat
		return
	default:
		return
	case nil:
	}
	if len(decoded) < 1 {
		err = base58.ErrInvalidFormat
		return
	}
	return decoded[1:], [2]byte{vByte, decoded[0]}, nil
}

// Decode base58 encoded string with error
func DecodeBase58(s string) ([]byte, error) {
	if s == "" {
		return []byte{}, nil
	}
	return base58.DecodeAlphabet(s, base58.BTCAlphabet)
}
