package bcs

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
)

// Serializer is a holding type to serialize a set of items into one shared buffer
type Serializer struct {
	out bytes.Buffer
	err error
}

// Serialize serializes a single item
func Serialize(value Marshaler) (bytes []byte, err error) {
	var ser Serializer
	value.MarshalBCS(&ser)
	err = ser.Error()
	if err != nil {
		return
	}
	bytes = ser.ToBytes()
	return
}

// Error the error if serialization has failed at any point
func (ser *Serializer) Error() error {
	return ser.err
}

// SetError If the data is well-formed but nonsense, MarshalBCS() code can set error
func (ser *Serializer) SetError(err error) {
	ser.err = err
}

// Bool serialize a bool into a single byte
func (ser *Serializer) Bool(v bool) {
	if v {
		ser.U8(1)
	} else {
		ser.U8(0)
	}
}

// U8 serialize a byte
func (ser *Serializer) U8(v uint8) {
	ser.out.WriteByte(v)
}

// U16 serialize an unsigned 16 bit integer
func (ser *Serializer) U16(v uint16) {
	var ub [2]byte
	binary.LittleEndian.PutUint16(ub[:], v)
	ser.out.Write(ub[:])
}

// U32 serialize an unsigned 32 bit integer
func (ser *Serializer) U32(v uint32) {
	var ub [4]byte
	binary.LittleEndian.PutUint32(ub[:], v)
	ser.out.Write(ub[:])
}

// U64 serialize an unsigned 64 bit integer
func (ser *Serializer) U64(v uint64) {
	var ub [8]byte
	binary.LittleEndian.PutUint64(ub[:], v)
	ser.out.Write(ub[:])
}

// U128 serialize an unsigned 128 bit integer
func (ser *Serializer) U128(v big.Int) {
	var ub [16]byte
	v.FillBytes(ub[:])
	reverse(ub[:])
	ser.out.Write(ub[:])
}

// U256 serialize an unsigned 256 bit integer
func (ser *Serializer) U256(v big.Int) {
	var ub [32]byte
	v.FillBytes(ub[:])
	reverse(ub[:])
	ser.out.Write(ub[:])
}

// Uleb128 serialize an unsigned 32-bit integer as an Uleb128.  This is used specifically for sequence lengths
func (ser *Serializer) Uleb128(v uint32) {
	for v > 0x80 {
		nb := uint8(v & 0x7f)
		ser.out.WriteByte(0x80 | nb)
		v = v >> 7
	}
	ser.out.WriteByte(uint8(v & 0x7f))
}

// WriteBytes serialize an array of bytes with its length first as an uleb128
func (ser *Serializer) WriteBytes(v []byte) {
	ser.Uleb128(uint32(len(v)))
	ser.out.Write(v)
}

// WriteString similar to WriteBytes using the UTF-8 byte representation of the string
func (ser *Serializer) WriteString(v string) {
	ser.WriteBytes([]byte(v))
}

// FixedBytes similar to WriteBytes, but it forgoes the length header.  This is useful if you know the fixed length
// size of the data, such as AccountAddress
func (ser *Serializer) FixedBytes(v []byte) {
	ser.out.Write(v)
}

// Struct uses custom serialization for a Marshaler implementation
func (ser *Serializer) Struct(v Marshaler) {
	v.MarshalBCS(ser)
}

// ToBytes outputs the encoded bytes
func (ser *Serializer) ToBytes() []byte {
	return ser.out.Bytes()
}

// Reset clears the serializer to be reused
func (ser *Serializer) Reset() {
	ser.out.Reset()
	ser.err = nil
}

// SerializeSequence serializes a sequence of Marshaler implemented types.  Prefixed with the length of the sequence
func SerializeSequence[AT []T, T any](array AT, ser *Serializer) {
	ser.Uleb128(uint32(len(array)))
	for i, v := range array {
		// Check if by value is Marshaler
		mv, ok := any(v).(Marshaler)
		if ok {
			mv.MarshalBCS(ser)
			continue
		}
		// Check if by reference is Marshaler
		mv, ok = any(&v).(Marshaler)
		if ok {
			mv.MarshalBCS(ser)
			continue
		}
		ser.SetError(fmt.Errorf("could not serialize sequence[%d] member of %T", i, v))
		return
	}
}
