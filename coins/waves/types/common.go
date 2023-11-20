/**
MIT License

Copyright (c) 2018 WavesPlatform

*/

package types

import "encoding/binary"

func PutBool(buf []byte, b bool) {
	if b {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
}

// PutBytesWithUInt16Len prepends given buf with 2 bytes of it's length.
func PutBytesWithUInt16Len(buf []byte, data []byte) {
	sl := uint16(len(data))
	binary.BigEndian.PutUint16(buf, sl)
	copy(buf[2:], data)
}

// PutStringWithUInt16Len writes to the buffer `buf` two bytes of the string `s` length followed with the bytes of the string `s`.
func PutStringWithUInt16Len(buf []byte, s string) {
	sl := uint16(len(s))
	binary.BigEndian.PutUint16(buf, sl)
	copy(buf[2:], s)
}
