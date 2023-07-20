package base

import "encoding/binary"

type Encoding int

const (
	EncodingBin Encoding = iota
	EncodingCompactU16
	EncodingBorsh
)

func (enc Encoding) String() string {
	switch enc {
	case EncodingBin:
		return "Bin"
	case EncodingCompactU16:
		return "CompactU16"
	case EncodingBorsh:
		return "Borsh"
	default:
		return ""
	}
}

func (en Encoding) IsBorsh() bool {
	return en == EncodingBorsh
}

func (en Encoding) IsBin() bool {
	return en == EncodingBin
}

func (en Encoding) IsCompactU16() bool {
	return en == EncodingCompactU16
}

func isValidEncoding(enc Encoding) bool {
	switch enc {
	case EncodingBin, EncodingCompactU16, EncodingBorsh:
		return true
	default:
		return false
	}
}

var LE binary.ByteOrder = binary.LittleEndian
var BE binary.ByteOrder = binary.BigEndian
