package brc20

import (
	"math/big"
	"strconv"
)

func ConvertToBigInt(v string) *big.Int {
	b := new(big.Int)
	b.SetString(v, 10)
	return b
}

func ConvertToUint32(v string) uint32 {
	i, _ := strconv.Atoi(v)
	return uint32(i)
}
