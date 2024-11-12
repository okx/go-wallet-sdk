package starknet

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"regexp"
)

const (
	FIELD_GEN   int    = 3
	FIELD_PRIME string = "3618502788666131213697322783095070105623107215331596699973092056135872020481"
)

var (
	MaxFelt     = StrToFelt(FIELD_PRIME)
	asciiRegexp = regexp.MustCompile(`^([[:graph:]]|[[:space:]]){1,31}$`)
)

const FeltLength = 32

// Felt represents Field Element or Felt from cairo.
type Felt [FeltLength]byte

// Big converts a Felt to its big.Int representation.
func (f Felt) Big() *big.Int { return new(big.Int).SetBytes(f[:]) }

// Bytes gets the byte representation of the Felt.
func (f Felt) Bytes() []byte { return f[:] }

// StrToFelt converts a string containing a decimal, hexadecimal or UTF8 charset into a Felt.
func StrToFelt(str string) Felt {
	var f Felt
	f.strToFelt(str)
	return f
}

func (f *Felt) strToFelt(str string) bool {
	if b, ok := new(big.Int).SetString(str, 0); ok {
		b.FillBytes(f[:])
		return ok
	}

	if asciiRegexp.MatchString(str) {
		hexStr := hex.EncodeToString([]byte(str))
		if b, ok := new(big.Int).SetString(hexStr, 16); ok {
			b.FillBytes(f[:])
			return ok
		}
	}

	return false
}

// BigToFelt converts a big.Int to its Felt representation.
func BigToFelt(b *big.Int) Felt {
	var f Felt
	b.FillBytes(f[:])
	return f
}

// BytesToFelt converts a []byte to its Felt representation.
func BytesToFelt(b []byte) Felt {
	var f Felt
	copy(f[:], b)
	return f
}

// String converts a Felt into its 'short string' representation.
func (f Felt) ShortString() string {
	str := string(f.Big().Bytes())
	if asciiRegexp.MatchString(str) {
		return str
	}
	return ""
}

// String converts a Felt into its hexadecimal string representation and implement fmt.Stringer.
func (f Felt) String() string {
	return fmt.Sprintf("0x%x", f.Big())
}
