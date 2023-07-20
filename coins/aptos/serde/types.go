package serde

import (
	"fmt"
	"math"
	"math/big"
)

type Uint128 struct {
	High uint64
	Low  uint64
}

type Int128 struct {
	High int64
	Low  uint64
}

// FromBig converts *big.Int to 128-bit Uint128 value ignoring overflows.
// If input integer is nil or negative then return Zero.
// If input interger overflows 128-bit then return Max.
func FromBig(i *big.Int) (*Uint128, error) {
	u, r := FromBigX(i)
	if !r {
		return nil, fmt.Errorf("u128 from big int error")
	}
	return &u, nil
}

// Zero is the lowest possible Uint128 value.
func Zero() Uint128 {
	return From64(0)
}

// From64 converts 64-bit value v to a Uint128 value.
// Upper 64-bit half will be zero.
func From64(v uint64) Uint128 {
	return Uint128{Low: v}
}

// Max is the largest possible Uint128 value.
func Max() Uint128 {
	return Uint128{
		Low:  math.MaxUint64,
		High: math.MaxUint64,
	}
}

// FromBigX converts *big.Int to 128-bit Uint128 value (eXtended version).
// Provides ok successful flag as a second return value.
// If input integer is negative or overflows 128-bit then ok=false.
// If input is nil then zero 128-bit returned.
func FromBigX(i *big.Int) (Uint128, bool) {
	switch {
	case i == nil:
		return Zero(), true // assuming nil === 0
	case i.Sign() < 0:
		return Zero(), false // value cannot be negative!
	case i.BitLen() > 128:
		return Max(), false // value overflows 128-bit!
	}

	// Note, actually result of big.Int.Uint64 is undefined
	// if stored value is greater than 2^64
	// but we assume that it just gets lower 64 bits.
	t := new(big.Int)
	lo := i.Uint64()
	hi := t.Rsh(i, 64).Uint64()
	return Uint128{
		Low:  lo,
		High: hi,
	}, true
}
