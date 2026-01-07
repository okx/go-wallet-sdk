package util

import (
	"encoding/hex"
	"encoding/json"
	"math/big"
	"strconv"
)

func getValAndBase(v string) (string, int) {
	if len(v) > 1 && (v[:2] == "0x" || v[:2] == "0X") {
		return v[2:], 16
	}
	return v, 10
}

func ToUint32(v string) uint32 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseUint(x, base, 32)
	if err != nil {
		return 0
	}
	return uint32(i)
}

func ToInt32(v string) int32 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseInt(x, base, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

func ToUint64(v string) uint64 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseUint(x, base, 64)
	if err != nil {
		return 0
	}
	return i
}

func ToUint8(v string) uint8 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseUint(x, base, 8)
	if err != nil {
		return 0
	}
	return uint8(i)
}

func ToInt64(v string) int64 {
	x, base := getValAndBase(v)
	i, err := strconv.ParseInt(x, base, 64)
	if err != nil {
		return 0
	}
	return i
}

func ToBigInt(v string) *big.Int {
	x, base := getValAndBase(v)
	b := new(big.Int)
	b.SetString(x, base)
	return b
}

func ConvertToInt32(v string) int32 {
	i, err := strconv.ParseInt(v, 10, 32)
	if err != nil {
		return 0
	}
	return int32(i)
}

func ConvertToUint64(v string) uint64 {
	i, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func ConvertToUint8(v string) uint8 {
	i, err := strconv.ParseUint(v, 10, 8)
	if err != nil {
		return 0
	}
	return uint8(i)
}

func ConvertToInt64(v string) int64 {
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return 0
	}
	return i
}

func ConvertToBigInt(v string) *big.Int {
	b := new(big.Int)
	b.SetString(v, 10)
	return b
}

func ConvertHexToBigInt(v string) *big.Int {
	b := new(big.Int)
	b.SetString(v, 16)
	return b
}

func ConvertHexWithPrefixToBigInt(v string) *big.Int {
	v = hex.EncodeToString(DecodeHexStringPad(v))
	b := new(big.Int)
	b.SetString(v, 16)
	return b
}

func ConvertToJsonString(v any) (string, error) {
	bs, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}
