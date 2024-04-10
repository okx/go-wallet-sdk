package serde

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

type Uint128 struct {
	High uint64
	Low  uint64
}

type Uint256 struct {
	High Uint128
	Low  Uint128
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

func ZeroUint256() Uint256 {
	return Uint256{Low: Zero()}
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

func MaxUint256() Uint256 {
	return Uint256{
		Low:  Max(),
		High: Max(),
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
func BigInt2U256(n *big.Int) (*Uint256, error) {
	m, ok := BigIntX2Uint256(n)
	if !ok {
		return nil, errors.New("BigInt2U256 not ok")
	}
	return &m, nil
}
func BigIntX2Uint256(i *big.Int) (Uint256, bool) {
	switch {
	case i == nil:
		return ZeroUint256(), true // assuming nil === 0
	case i.Sign() < 0:
		return ZeroUint256(), false // value cannot be negative!
	case i.BitLen() > 256:
		return MaxUint256(), false // value overflows 256bit!
	}

	t := new(big.Int)

	if i.BitLen() <= 128 {
		lo, ok := FromBigX(i)
		return Uint256{Low: lo}, ok
	}
	low := new(big.Int).SetBytes(i.Bytes()[16:])
	lo, ok := FromBigX(low)
	hi, ok := FromBigX(t.Rsh(i, 128))
	return Uint256{
		Low:  lo,
		High: hi,
	}, ok
}
