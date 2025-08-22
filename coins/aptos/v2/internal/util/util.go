package util

import (
	"encoding/hex"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/common"
	"golang.org/x/crypto/sha3"
	"math"
	"math/big"
	"strconv"
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

// BytesToHex converts a byte slice to a hex string with a leading 0x
func BytesToHex(bytes []byte) string {
	return "0x" + hex.EncodeToString(bytes)
}

// StrToUint64 converts a string to a uint64
func StrToUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, 10, 64)
}

// StrToBigInt converts a string to a big.Int for u128 and u256 values
func StrToBigInt(val string) (*big.Int, error) {
	num := &big.Int{}
	var ok bool
	if common.IsHexString(val) {
		num, ok = num.SetString(strings.TrimPrefix(val, "0x"), 16)
	} else {
		num, ok = num.SetString(val, 10)
	}
	if !ok {
		return nil, fmt.Errorf("num %s is not an integer", val)
	}
	return num, nil
}

func UintToU8(u uint) (uint8, error) {
	if u > math.MaxUint8 {
		return 0, fmt.Errorf("u %d is greater than %d", u, math.MaxUint8)
	}

	val := uint8(u)
	return val, nil
}

func UintToU16(u uint) (uint16, error) {
	if u > math.MaxUint16 {
		return 0, fmt.Errorf("u %d is greater than %d", u, uint16(math.MaxUint16))
	}

	val := uint16(u)
	return val, nil
}

func UintToU32(u uint) (uint32, error) {
	if u > math.MaxUint32 {
		return 0, fmt.Errorf("u %d is greater than %d", u, uint32(math.MaxUint32))
	}
	val := uint32(u)
	return val, nil
}

func UintToUBigInt(u uint) (*big.Int, error) {
	str := strconv.FormatUint(uint64(u), 10)
	return StrToBigInt(str)
}

func Float64ToU8(u float64) (uint8, error) {
	switch {
	case u > math.MaxUint8:
		return 0, fmt.Errorf("u %f is greater than %d", u, math.MaxUint8)
	case u < 0:
		return 0, fmt.Errorf("u %f is less than 0", u)
	case math.Floor(u) != u:
		return 0, fmt.Errorf("u %f should be an integer", u)
	default:
		val := uint8(u)
		return val, nil
	}
}

func Float64ToU16(u float64) (uint16, error) {
	switch {
	case u > math.MaxUint16:
		return 0, fmt.Errorf("u %f is greater than %d", u, uint16(math.MaxUint16))
	case u < 0:
		return 0, fmt.Errorf("u %f is less than 0", u)
	case math.Floor(u) != u:
		return 0, fmt.Errorf("u %f should be an integer", u)
	default:
		val := uint16(u)
		return val, nil
	}
}

func Float64ToU32(u float64) (uint32, error) {
	switch {
	case u > math.MaxUint32:
		return 0, fmt.Errorf("u %f is greater than %d", u, uint32(math.MaxUint32))
	case u < 0:
		return 0, fmt.Errorf("u %f is less than 0", u)
	case math.Floor(u) != u:
		return 0, fmt.Errorf("u %f should be an integer", u)
	default:
		val := uint32(u)
		return val, nil
	}
}

func Float64ToU64(u float64) (uint64, error) {
	switch {
	case u > math.MaxUint64:
		return 0, fmt.Errorf("u %f is greater than %d", u, uint64(math.MaxUint64))
	case u < 0:
		return 0, fmt.Errorf("u %f is less than 0", u)
	case math.Floor(u) != u:
		return 0, fmt.Errorf("u %f should be an integer", u)
	default:
		val := uint64(u)
		return val, nil
	}
}

func IntToU8(u int) (uint8, error) {
	if u > math.MaxUint8 {
		return 0, fmt.Errorf("u %d is greater than %d", u, math.MaxUint8)
	} else if u < 0 {
		return 0, fmt.Errorf("u %d is less than 0", u)
	}

	val := uint8(u)
	return val, nil
}

func IntToU16(u int) (uint16, error) {
	if u > math.MaxUint16 {
		return 0, fmt.Errorf("u %d is greater than %d", u, math.MaxUint16)
	} else if u < 0 {
		return 0, fmt.Errorf("u %d is less than 0", u)
	}

	val := uint16(u)
	return val, nil
}

func IntToU32(u int) (uint32, error) {
	//if u > math.MaxUint32 {
	//	return 0, fmt.Errorf("u %d is greater than %d", uint32(u), uint32(math.MaxUint32))
	//} else
	if u < 0 {
		return 0, fmt.Errorf("u %d is less than 0", u)
	}

	val := uint32(u)
	return val, nil
}

func IntToU64(u int) (uint64, error) {
	if u < 0 {
		return 0, fmt.Errorf("u %d is less than 0", u)
	}

	val := uint64(u)
	return val, nil
}

func IntToUBigInt(u int) (*big.Int, error) {
	if u < 0 {
		return nil, fmt.Errorf("u %d is less than 0", u)
	}

	return big.NewInt(int64(u)), nil
}

func Float64ToBigInt(u float64) (*big.Int, error) {
	if u < 0 {
		return nil, fmt.Errorf("%f is less than 0", u)
	}

	return big.NewInt(0).SetUint64(uint64(u)), nil
}

func Uint32ToU8(u uint32) (uint8, error) {
	if u > math.MaxUint8 {
		return 0, fmt.Errorf("u %d is greater than %d", u, math.MaxUint8)
	}
	return uint8(u), nil
}
