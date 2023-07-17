package util

import (
	"encoding/hex"
	"strings"
)

// delete the 0x from the front
func RemoveZeroHex(s string) []byte {
	if strings.HasPrefix(s, "0x") || strings.HasPrefix(s, "0X") {
		s = s[2:]
	}
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return Hex2Bytes(s)
}

func Hex2Bytes(str string) []byte {
	h, _ := hex.DecodeString(str)
	return h
}

func EncodeHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func EncodeHexWith0x(bytes []byte) string {
	return "0x" + EncodeHex(bytes)
}

func DecodeHexString(hexString string) ([]byte, error) {
	if strings.HasPrefix(hexString, "0x") || strings.HasPrefix(hexString, "0X") {
		hexString = hexString[2:]
	}
	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}
	return hex.DecodeString(hexString)
}
