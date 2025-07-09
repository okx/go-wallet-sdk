package util

import (
	"encoding/hex"
	"regexp"
)

func RemoveZeroHex(s string) []byte {
	return DecodeHexStringPad(s)
}

func EncodeHex(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func EncodeHexWithPrefix(bytes []byte) string {
	return "0x" + EncodeHex(bytes)
}

func DecodeHexString(hexString string) []byte {
	bytes, _ := DecodeHexStringErr(hexString)
	return bytes
}

func DecodeHexStringErr(hexString string) ([]byte, error) {
	return hex.DecodeString(RemoveHexPrefix(hexString))
}

func DecodeHexStringPad(s string) []byte {
	bytes, _ := DecodeHexStringPadErr(s)
	return bytes
}
func DecodeHexStringPadErr(hexString string) ([]byte, error) {
	return hex.DecodeString(RemoveHexPrefixPad(hexString))
}

func RemoveHexPrefixPad(s string) string {
	s = RemoveHexPrefix(s)
	if len(s)%2 == 1 {
		s = "0" + s
	}
	return s
}

func RemoveHexPrefix(s string) string {
	if len(s) > 1 {
		if s[0:2] == "0x" || s[0:2] == "0X" {
			s = s[2:]
		}
	}
	return s
}

func IsHexStringRelaxed(s string) bool {
	res, err := regexp.MatchString(`(?i)^0x[0-9a-f]+$|^[0-9a-fA-F]+$`, s)
	if err != nil {
		return false
	}
	return res
}
