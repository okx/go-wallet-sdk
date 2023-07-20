package base58

import (
	"crypto/sha256"
	"github.com/btcsuite/btcutil/base58"
)

func Encode(input []byte) string {
	return base58.Encode(input)
}

func Decode(b string) []byte {
	return base58.Decode(b)
}

func CheckEncode(input []byte, version byte) string {
	return base58.CheckEncode(input, version)
}

func checksum(input []byte) (cksum [4]byte) {
	h := sha256.Sum256(input)
	h2 := sha256.Sum256(h[:])
	copy(cksum[:], h2[:4])
	return
}

func CheckEncodeRaw(input []byte) string {
	b := make([]byte, 0, len(input)+4)
	b = append(b, input...)
	cksum := checksum(b)
	b = append(b, cksum[:]...)
	return base58.Encode(b)
}

func CheckDecode(input string) (result []byte, version byte, err error) {
	return base58.CheckDecode(input)
}
