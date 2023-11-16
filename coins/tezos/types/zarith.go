/**
The MIT License (MIT) Copyright (c) 2021-2023 Blockwatch Data Inc.
*/

package types

import (
	"bytes"
	"strconv"
)

type N int64

func (n N) Int64() int64 {
	return int64(n)
}

func (n N) String() string {
	return strconv.FormatInt(int64(n), 10)
}

func (n *N) SetInt64(i int64) *N {
	*n = N(i)
	return n
}

func (n N) EncodeBuffer(buf *bytes.Buffer) error {
	x := int64(n)
	for x >= 0x80 {
		buf.WriteByte(byte(x) | 0x80)
		x >>= 7
	}
	buf.WriteByte(byte(x))
	return nil
}
