package bcs

import (
	"encoding/binary"
	"fmt"
	"math/big"
)

// Deserializer is a type to deserialize a known set of bytes.
// The reader must know the types, as the format is not self-describing
type Deserializer struct {
	source []byte
	pos    int
	err    error
}

// NewDeserializer creates a new Deserializer from a byte array
func NewDeserializer(bytes []byte) *Deserializer {
	return &Deserializer{
		source: bytes,
		pos:    0,
		err:    nil,
	}
}

// Deserialize deserializes a single item
func Deserialize(dest Unmarshaler, bytes []byte) error {
	des := Deserializer{
		source: bytes,
		pos:    0,
		err:    nil,
	}
	dest.UnmarshalBCS(&des)
	return des.err
}

// Error If there has been any error, return it
func (des *Deserializer) Error() error {
	return des.err
}

// SetError If the data is well-formed but nonsense, UnmarshalBCS() code can set error
func (des *Deserializer) SetError(err error) {
	des.err = err
}

// Remaining tells the remaining bytes, which can be useful if there were more bytes than expected
func (des *Deserializer) Remaining() int {
	return len(des.source) - des.pos
}

// Bool deserializes a single byte as a bool
func (des *Deserializer) Bool() bool {
	if des.pos >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize bool")
		return false
	}

	out := false
	switch des.U8() {
	case 0:
		out = false
	case 1:
		out = true
	default:
		des.setError("bad bool at [%des]: %x", des.pos-1, des.source[des.pos-1])
	}
	return out
}

// U8 deserializes a single unsigned 8-bit integer
func (des *Deserializer) U8() uint8 {
	if des.pos >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize u8")
		return 0
	}
	out := des.source[des.pos]
	des.pos++
	return out
}

// U16 deserializes a single unsigned 16-bit integer
func (des *Deserializer) U16() uint16 {
	if des.pos+1 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize u16")
		return 0
	}
	out := binary.LittleEndian.Uint16(des.source[des.pos : des.pos+2])
	des.pos += 2
	return out
}

// U32 deserializes a single unsigned 32-bit integer
func (des *Deserializer) U32() uint32 {
	if des.pos+3 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize u32")
		return 0
	}
	out := binary.LittleEndian.Uint32(des.source[des.pos : des.pos+4])
	des.pos += 4
	return out
}

// U64 deserializes a single unsigned 64-bit integer
func (des *Deserializer) U64() uint64 {
	if des.pos+7 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize u64")
		return 0
	}
	out := binary.LittleEndian.Uint64(des.source[des.pos : des.pos+8])
	des.pos += 8
	return out
}

// U128 deserializes a single unsigned 128-bit integer
func (des *Deserializer) U128() big.Int {
	if des.pos+15 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize u128")
		return *big.NewInt(-1)
	}
	var bytesBigEndian [16]byte
	copy(bytesBigEndian[:], des.source[des.pos:des.pos+16])
	des.pos += 16
	reverse(bytesBigEndian[:])
	var out big.Int
	out.SetBytes(bytesBigEndian[:])
	return out
}

// U256 deserializes a single unsigned 256-bit integer
func (des *Deserializer) U256() big.Int {
	if des.pos+31 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize u256")
		return *big.NewInt(-1)
	}
	var bytesBigEndian [32]byte
	copy(bytesBigEndian[:], des.source[des.pos:des.pos+32])
	des.pos += 32
	reverse(bytesBigEndian[:])
	var out big.Int
	out.SetBytes(bytesBigEndian[:])
	return out
}

// Uleb128 deserializes a 32-bit integer from a variable length uleb128
func (des *Deserializer) Uleb128() uint32 {
	var out uint32 = 0
	shift := 0

	for {
		if des.pos >= len(des.source) {
			des.setError("not enough bytes remaining to deserialize uleb128")
			return 0
		}

		val := des.source[des.pos]
		out = out | (uint32(val&0x7f) << shift)
		des.pos++
		if (val & 0x80) == 0 {
			break
		}
		shift += 7
		// TODO: if shift is too much, error
	}

	return out
}

// ReadBytes reads bytes prefixed with a length
func (des *Deserializer) ReadBytes() []byte {
	length := des.Uleb128()
	if des.err != nil {
		return nil
	}
	if des.pos+int(length)-1 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize bytes")
		return nil
	}
	out := make([]byte, length)
	copy(out, des.source[des.pos:des.pos+int(length)])
	des.pos += int(length)
	return out
}

// ReadString reads UTF-8 bytes prefixed with a length
func (des *Deserializer) ReadString() string {
	return string(des.ReadBytes())
}

// ReadFixedBytes reads bytes not-prefixed with a length
func (des *Deserializer) ReadFixedBytes(length int) []byte {
	if des.pos+length-1 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize fixedBytes")
		return nil
	}
	out := make([]byte, length)
	des.ReadFixedBytesInto(out)
	return out
}

// ReadFixedBytesInto reads bytes not-prefixed with a length into a byte array
func (des *Deserializer) ReadFixedBytesInto(dest []byte) {
	length := len(dest)
	if des.pos+length-1 >= len(des.source) {
		des.setError("not enough bytes remaining to deserialize fixedBytes")
		return
	}
	copy(dest, des.source[des.pos:des.pos+length])
	des.pos += length
}

// Struct reads an Unmarshaler implementation from bcs bytes
func (des *Deserializer) Struct(v Unmarshaler) {
	v.UnmarshalBCS(des)
}

// DeserializeSequence deserializes an Unmarshaler implementation array
func DeserializeSequence[T any](des *Deserializer) []T {
	length := des.Uleb128()
	if des.Error() != nil {
		return nil
	}
	out := make([]T, length)
	for i := 0; i < int(length); i++ {
		v := &(out[i])
		mv, ok := any(v).(Unmarshaler)
		if ok {
			mv.UnmarshalBCS(des)
		} else {
			des.setError("could not deserialize sequence[%d] member of %T", i, v)
			return nil
		}
	}
	return out
}

// DeserializeMapToSlices returns two slices []K and []V of equal length that are equivalent to map[K]V but may represent types that are not valid Go map keys.
func DeserializeMapToSlices[K, V any](des *Deserializer) (keys []K, values []V) {
	count := des.Uleb128()
	keys = make([]K, 0, count)
	values = make([]V, 0, count)
	// todo go 1.22
	//for range count {
	for i := uint32(0); i < count; i++ {
		var nextK K
		var nextV V
		switch sv := any(&nextK).(type) {
		case Unmarshaler:
			sv.UnmarshalBCS(des)
		case *string:
			*sv = des.ReadString()
		}
		switch sv := any(&nextV).(type) {
		case Unmarshaler:
			sv.UnmarshalBCS(des)
		case *string:
			*sv = des.ReadString()
		case *[]byte:
			*sv = des.ReadBytes()
		}
		keys = append(keys, nextK)
		values = append(values, nextV)
	}
	return
}

// setError sets the deserialization error
func (des *Deserializer) setError(msg string, args ...any) {
	if des.err != nil {
		return
	}
	des.err = fmt.Errorf(msg, args...)
}
