// Copyright (c) 2013-2015 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package base58

import (
	"fmt"
	"math/big"
)

//go:generate go run genalphabet.go

var bigRadix = [...]*big.Int{
	big.NewInt(0),
	big.NewInt(58),
	big.NewInt(58 * 58),
	big.NewInt(58 * 58 * 58),
	big.NewInt(58 * 58 * 58 * 58),
	big.NewInt(58 * 58 * 58 * 58 * 58),
	big.NewInt(58 * 58 * 58 * 58 * 58 * 58),
	big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58),
	big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58 * 58),
	big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58),
	bigRadix10,
}

var bigRadix10 = big.NewInt(58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58 * 58) // 58^10

// Decode decodes a modified base58 string to a byte slice.
func Decode(b string) []byte {
	answer := big.NewInt(0)
	scratch := new(big.Int)

	// Calculating with big.Int is slow for each iteration.
	//    x += b58[b[i]] * j
	//    j *= 58
	//
	// Instead we can try to do as much calculations on int64.
	// We can represent a 10 digit base58 number using an int64.
	//
	// Hence we'll try to convert 10, base58 digits at a time.
	// The rough idea is to calculate `t`, such that:
	//
	//   t := b58[b[i+9]] * 58^9 ... + b58[b[i+1]] * 58^1 + b58[b[i]] * 58^0
	//   x *= 58^10
	//   x += t
	//
	// Of course, in addition, we'll need to handle boundary condition when `b` is not multiple of 58^10.
	// In that case we'll use the bigRadix[n] lookup for the appropriate power.
	for t := b; len(t) > 0; {
		n := len(t)
		if n > 10 {
			n = 10
		}

		total := uint64(0)
		for _, v := range t[:n] {
			tmp := b58[v]
			if tmp == 255 {
				return []byte("")
			}
			total = total*58 + uint64(tmp)
		}

		answer.Mul(answer, bigRadix[n])
		scratch.SetUint64(total)
		answer.Add(answer, scratch)

		t = t[n:]
	}

	tmpval := answer.Bytes()

	var numZeros int
	for numZeros = 0; numZeros < len(b); numZeros++ {
		if b[numZeros] != alphabetIdx0 {
			break
		}
	}
	flen := numZeros + len(tmpval)
	val := make([]byte, flen)
	copy(val[numZeros:], tmpval)

	return val
}

// Encode encodes a byte slice to a modified base58 string.
func Encode(b []byte) string {
	x := new(big.Int)
	x.SetBytes(b)

	// maximum length of output is log58(2^(8*len(b))) == len(b) * 8 / log(58)
	maxlen := int(float64(len(b))*1.365658237309761) + 1
	answer := make([]byte, 0, maxlen)
	mod := new(big.Int)
	for x.Sign() > 0 {
		// Calculating with big.Int is slow for each iteration.
		//    x, mod = x / 58, x % 58
		//
		// Instead we can try to do as much calculations on int64.
		//    x, mod = x / 58^10, x % 58^10
		//
		// Which will give us mod, which is 10 digit base58 number.
		// We'll loop that 10 times to convert to the answer.

		x.DivMod(x, bigRadix10, mod)
		if x.Sign() == 0 {
			// When x = 0, we need to ensure we don't add any extra zeros.
			m := mod.Int64()
			for m > 0 {
				answer = append(answer, alphabet[m%58])
				m /= 58
			}
		} else {
			m := mod.Int64()
			for i := 0; i < 10; i++ {
				answer = append(answer, alphabet[m%58])
				m /= 58
			}
		}
	}

	// leading zero bytes
	for _, i := range b {
		if i != 0 {
			break
		}
		answer = append(answer, alphabetIdx0)
	}

	// reverse
	alen := len(answer)
	for i := 0; i < alen/2; i++ {
		answer[i], answer[alen-1-i] = answer[alen-1-i], answer[i]
	}

	return string(answer)
}

// EncodeAlphabet encodes the passed bytes into a base58 encoded string with the
// passed alphabet.
func EncodeAlphabet(bin []byte, alphabet *Alphabet) string {
	return FastBase58EncodingAlphabet(bin, alphabet)
}

// DecodeAlphabet decodes the base58 encoded bytes using the given b58 alphabet.
func DecodeAlphabet(str string, alphabet *Alphabet) ([]byte, error) {
	return FastBase58DecodingAlphabet(str, alphabet)
}

// FastBase58EncodingAlphabet encodes the passed bytes into a base58 encoded
// string with the passed alphabet.
func FastBase58EncodingAlphabet(bin []byte, alphabet *Alphabet) string {
	zero := alphabet.encode[0]

	binsz := len(bin)
	var i, j, zcount, high int
	var carry uint32

	for zcount < binsz && bin[zcount] == 0 {
		zcount++
	}

	size := (binsz-zcount)*138/100 + 1
	var buf = make([]uint32, size)

	high = size - 1
	for i = zcount; i < binsz; i++ {
		j = size - 1
		for carry = uint32(bin[i]); j > high || carry != 0; j-- {
			carry += buf[j] << 8
			buf[j] = carry % 58
			carry /= 58
		}
		high = j
	}

	for j = 0; j < size && buf[j] == 0; j++ {
	}

	var b58 = make([]byte, size-j+zcount)

	if zcount != 0 {
		for i = 0; i < zcount; i++ {
			b58[i] = zero
		}
	}

	for i = zcount; j < size; i++ {
		b58[i] = alphabet.encode[buf[j]]
		j++
	}

	return string(b58)
}

// FastBase58DecodingAlphabet decodes the base58 encoded bytes using the given
// b58 alphabet.
func FastBase58DecodingAlphabet(str string, alphabet *Alphabet) ([]byte, error) {
	if len(str) == 0 {
		return nil, fmt.Errorf("zero length string")
	}

	var (
		t, c   uint64
		zmask  uint32
		zcount int

		b58u  = []rune(str)
		b58sz = len(b58u)

		outisz    = (b58sz + 3) >> 2
		binu      = make([]byte, (b58sz+3)*3)
		bytesleft = b58sz & 3

		zero = rune(alphabet.encode[0])
	)

	if bytesleft > 0 {
		zmask = 0xffffffff << uint32(bytesleft*8)
	} else {
		bytesleft = 4
	}

	var outi = make([]uint32, outisz)

	for i := 0; i < b58sz && b58u[i] == zero; i++ {
		zcount++
	}

	for _, r := range b58u {
		if r > 127 {
			return nil, fmt.Errorf("high-bit set on invalid digit")
		}
		if alphabet.decode[r] == -1 {
			return nil, fmt.Errorf("invalid base58 digit (%q)", r)
		}

		c = uint64(alphabet.decode[r])

		for j := outisz - 1; j >= 0; j-- {
			t = uint64(outi[j])*58 + c
			c = (t >> 32) & 0x3f
			outi[j] = uint32(t & 0xffffffff)
		}

		if c > 0 {
			return nil, fmt.Errorf("output number too big (carry to the next int32)")
		}

		if outi[0]&zmask != 0 {
			return nil, fmt.Errorf("output number too big (last int32 filled too far)")
		}
	}

	var j, cnt int
	for j, cnt = 0, 0; j < outisz; j++ {
		for mask := byte(bytesleft-1) * 8; mask <= 0x18; mask, cnt = mask-8, cnt+1 {
			binu[cnt] = byte(outi[j] >> mask)
		}
		if j == 0 {
			bytesleft = 4 // because it could be less than 4 the first time through
		}
	}

	for n, v := range binu {
		if v > 0 {
			start := n - zcount
			if start < 0 {
				start = 0
			}
			return binu[start:cnt], nil
		}
	}
	return binu[:cnt], nil
}
