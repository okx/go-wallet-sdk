/*
Author: https://github.com/zksync-sdk/zksync-go
*
*/
package core

import "errors"

type Bits struct {
	bits []bool
}

func NewBits(len uint) *Bits {
	return &Bits{bits: make([]bool, len)}
}

func (b *Bits) Len() uint {
	return uint(len(b.bits))
}

func (b *Bits) SetBit(i uint, v bool) {
	// len is not checked, can cause panic
	b.bits[i] = v
}

func (b *Bits) GetBit(i uint) bool {
	// len is not checked, can cause panic
	return b.bits[i]
}

func (b *Bits) Clone() *Bits {
	clone := NewBits(b.Len())
	copy(clone.bits, b.bits)
	return clone
}

func (b *Bits) Append(a *Bits) *Bits {
	b.bits = append(b.bits, a.bits...)
	return b
}

func (b *Bits) Reverse() *Bits {
	for i, j := 0, len(b.bits)-1; i < j; i, j = i+1, j-1 {
		b.bits[i], b.bits[j] = b.bits[j], b.bits[i]
	}
	return b
}

func (b *Bits) String() (s string) {
	for _, v := range b.bits {
		if v {
			s += `1`
		} else {
			s += `0`
		}
	}
	return
}

func (b *Bits) ToBytesBE() ([]byte, error) {
	bits := len(b.bits)
	if bits%8 != 0 {
		return nil, errors.New("wrong number of bits to pack")
	}
	bytes := len(b.bits) / 8
	res := make([]byte, bytes)
	for i, b := range b.bits {
		if b {
			byteIdx := i / 8
			bitIdx := 7 - i%8
			res[byteIdx] = res[byteIdx] | 1<<bitIdx
		}
	}
	return res, nil
}

func (b *Bits) FromBytesBE(bytes []byte) *Bits {
	for i, v := range bytes {
		b.SetBit(uint(i*8+0), v&0x80 > 0)
		b.SetBit(uint(i*8+1), v&0x40 > 0)
		b.SetBit(uint(i*8+2), v&0x20 > 0)
		b.SetBit(uint(i*8+3), v&0x10 > 0)
		b.SetBit(uint(i*8+4), v&0x08 > 0)
		b.SetBit(uint(i*8+5), v&0x04 > 0)
		b.SetBit(uint(i*8+6), v&0x02 > 0)
		b.SetBit(uint(i*8+7), v&0x01 > 0)
	}
	return b
}
