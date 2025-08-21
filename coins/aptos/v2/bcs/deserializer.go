package bcs

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"slices"
)

// Deserializer is a type to deserialize a known set of bytes.
// The reader must know the types, as the format is not self-describing.
//
// Use [NewDeserializer] to initialize the Deserializer
//
//	bytes := []byte{0x01}
//	deserializer := NewDeserializer(bytes)
//	num := deserializer.U8()
//	if deserializer.Error() != nil {
//		return deserializer.Error()
//	}
type Deserializer struct {
	source []byte // Underlying data to parse
	pos    int    // Current position in the buffer
	err    error  // Any error that has happened so far
}

// NewDeserializer creates a new Deserializer from a byte array.
func NewDeserializer(bytes []byte) *Deserializer {
	return &Deserializer{
		source: bytes,
		pos:    0,
		err:    nil,
	}
}

// Deserialize deserializes a single item from bytes.
//
// This function will error if there are remaining bytes.
func Deserialize(dest Unmarshaler, bytes []byte) error {
	des := Deserializer{
		source: bytes,
		pos:    0,
		err:    nil,
	}
	des.Struct(dest)
	if des.err != nil {
		return des.err
	}
	if des.Remaining() > 0 {
		return fmt.Errorf("deserialize failed: remaining %d byte(s)", des.Remaining())
	}
	return nil
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
//
//	bytes := []byte{0x01, 0x02}
//	deserializer := NewDeserializer(bytes)
//	num := deserializer.U8()
//	deserializer.Remaining == 1
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

func deserializeUint[T uint8 | uint16 | uint32 | uint64](des *Deserializer, typeName string, size int, decode func(slice []byte) T) T {
	end := des.pos + size
	if end > len(des.source) {
		des.setError("not enough bytes remaining to deserialize %s", typeName)
		return T(0)
	}
	out := decode(des.source[des.pos:end])
	des.pos = end
	return out
}

func (des *Deserializer) deserializeUBigint(typeName string, size int) big.Int {
	end := des.pos + size
	if end > len(des.source) {
		des.setError("not enough bytes remaining to deserialize %s", typeName)
		return *big.NewInt(-1)
	}
	bytesBigEndian := make([]byte, size)
	copy(bytesBigEndian, des.source[des.pos:end])
	des.pos = end
	slices.Reverse(bytesBigEndian)
	var out big.Int
	out.SetBytes(bytesBigEndian)
	return out
}

// U8 deserializes a single unsigned 8-bit integer
func (des *Deserializer) U8() uint8 {
	return deserializeUint(des, "u8", 1, func(slice []byte) uint8 {
		return slice[0]
	})
}

// U16 deserializes a single unsigned 16-bit integer
func (des *Deserializer) U16() uint16 {
	return deserializeUint(des, "u16", 2, binary.LittleEndian.Uint16)
}

// U32 deserializes a single unsigned 32-bit integer
func (des *Deserializer) U32() uint32 {
	return deserializeUint(des, "u32", 4, binary.LittleEndian.Uint32)
}

// U64 deserializes a single unsigned 64-bit integer

func (des *Deserializer) U64() uint64 {
	return deserializeUint(des, "u64", 8, binary.LittleEndian.Uint64)
}

// U128 deserializes a single unsigned 128-bit integer
func (des *Deserializer) U128() big.Int {
	return des.deserializeUBigint("u128", 16)
}

// U256 deserializes a single unsigned 256-bit integer
func (des *Deserializer) U256() big.Int {
	return des.deserializeUBigint("u256", 32)
}

// Uleb128 deserializes a 32-bit integer from a variable length [Unsigned LEB128]
//
// [Unsigned LEB128]: https://en.wikipedia.org/wiki/LEB128#Unsigned_LEB128
func (des *Deserializer) Uleb128() uint32 {
	const maxU32 = uint64(0xFFFFFFFF)
	var out uint64
	shift := 0

	for out < maxU32 {
		// Ensure we still have bytes to process
		if des.pos >= len(des.source) {
			des.setError("not enough bytes remaining to deserialize uleb128")
			return 0
		}

		// Append the next byte
		val := des.source[des.pos]
		out |= uint64(val&0x7f) << shift
		des.pos++

		// If at any point the highest bit is not set, there are no more bytes to read
		if (val & 0x80) == 0 {
			break
		}

		shift += 7
	}

	// If all bytes have 0x80, then we have an invalid uleb128
	if out > maxU32 {
		des.setError("uleb128 is invalid as it goes higher than the max u32 value")
		return 0
	}

	return uint32(out)
}

// ReadBytes reads bytes prefixed with a length
func (des *Deserializer) ReadBytes() []byte {
	length := des.Uleb128()
	if des.err != nil {
		return nil
	}

	dest := make([]byte, length)
	des.readBytes("bytes", int(length), dest)
	return dest
}

// ReadString reads UTF-8 bytes prefixed with a length
func (des *Deserializer) ReadString() string {
	return string(des.ReadBytes())
}

// ReadFixedBytes reads bytes not-prefixed with a length
func (des *Deserializer) ReadFixedBytes(length int) []byte {
	out := make([]byte, length)
	des.ReadFixedBytesInto(out)
	return out
}

// ReadFixedBytesInto reads bytes not-prefixed with a length into a byte array
func (des *Deserializer) ReadFixedBytesInto(dest []byte) {
	length := len(dest)
	des.readBytes("fixedBytes", length, dest)
}

func (des *Deserializer) readBytes(typeName string, length int, dest []byte) {
	end := des.pos + length
	if end > len(des.source) {
		des.setError("not enough bytes remaining to deserialize %s", typeName)
		return
	}
	copy(dest, des.source[des.pos:end])
	des.pos = end
}

// Struct reads an Unmarshaler implementation from bcs bytes
//
// This is used for handling types outside the provided primitives
func (des *Deserializer) Struct(v Unmarshaler) {
	if v == nil {
		des.setError("cannot deserialize into nil")
		return
	}
	v.UnmarshalBCS(des)
}

// DeserializeSequence deserializes an Unmarshaler implementation array
//
// This lets you deserialize a whole sequence of [Unmarshaler], and will fail if any member fails.
// All sequences are prefixed with an Uleb128 length.
func DeserializeSequence[T any](des *Deserializer) []T {
	return DeserializeSequenceWithFunction(des, func(des *Deserializer, out *T) {
		mv, ok := any(out).(Unmarshaler)
		if ok {
			mv.UnmarshalBCS(des)
		} else {
			// If it isn't of type Unmarshaler, we pass up an error
			des.setError("type is not Unmarshaler")
		}
	})
}

// DeserializeSequenceWithFunction deserializes any array with the given function
//
// This lets you deserialize a whole sequence of any type, and will fail if any member fails.
// All sequences are prefixed with an Uleb128 length.
func DeserializeSequenceWithFunction[T any](des *Deserializer, deserialize func(des *Deserializer, out *T)) []T {
	length := des.Uleb128()
	if des.Error() != nil {
		return nil
	}
	out := make([]T, length)
	for i := range length {
		deserialize(des, &out[i])

		if des.Error() != nil {
			des.setError("could not deserialize sequence[%d] member of %w", i, des.Error())
			return nil
		}
	}
	return out
}

// DeserializeOption deserializes an optional value
//
// # Under the hood, this is represented as a 0 or 1 length array
//
// Here's an example for handling an optional value:
//
//	// For a Some(10) value
//	bytes == []byte{0x01, 0x0A}
//	des := NewDeserializer(bytes)
//	output := DeserializeOption(des, nil, func(des *Deserializer, out *uint8) {
//		out = des.U8()
//	})
//	// output == &10
//
//	// For a None value
//	bytes2 == []byte{0x00}
//	des2 := NewDeserializer(bytes2)
//	output := DeserializeOption(des2, nil, func(des *Deserializer, out *uint8) {
//		out = des.U8()
//	})
//	// output == nil
func DeserializeOption[T any](des *Deserializer, deserialize func(des *Deserializer, out *T)) *T {
	array := DeserializeSequenceWithFunction(des, deserialize)
	switch len(array) {
	case 0:
		// None
		return nil
	case 1:
		// Some
		return &array[0]
	default:
		des.setError("expected 0 or 1 element as an option, got %d", len(array))
	}
	return nil
}

// setError overrides the previous error, this can only be called from within the bcs package
func (des *Deserializer) setError(msg string, args ...any) {
	if des.err != nil {
		return
	}
	des.err = fmt.Errorf(msg, args...)
}

// DeserializeBytes deserializes a byte array prefixed with a length
func (des *Deserializer) Serialized() *Serialized {
	return &Serialized{des.ReadBytes()}
}
