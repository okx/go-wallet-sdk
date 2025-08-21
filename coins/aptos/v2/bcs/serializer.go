package bcs

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
	"math/big"
	"slices"
)

// Serializer is a holding type to serialize a set of items into one shared buffer
//
//	serializer := &Serializer{}
//	serializer.U64(uint64(10))
//	serializedBytes := serializer.ToBytes()
type Serializer struct {
	out bytes.Buffer // current serialized bytes
	err error        // any error that has occurred during serialization
}

// Serialize serializes a single item
//
//	type MyStruct struct {
//	  num uint64
//	}
//
//	func (str *MyStruct) MarshalBCS(ser *Serialize) {
//		ser.U64(num)
//	}
//
//	struct := &MyStruct{ num: 100 }
//	bytes, _ := Serialize(struct)
func Serialize(value Marshaler) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		value.MarshalBCS(ser)
	})
}

// Error the error if serialization has failed at any point
func (ser *Serializer) Error() error {
	return ser.err
}

// SetError If the data is well-formed but nonsense, [Marshaler.MarshalBCS] code can set the error.
func (ser *Serializer) SetError(err error) {
	ser.err = err
}

// Bool serialize a bool into a single byte, 0x01 for true and 0x00 for false
func (ser *Serializer) Bool(v bool) {
	if v {
		ser.U8(1)
	} else {
		ser.U8(0)
	}
}

func serializeUInt[T uint16 | uint32 | uint64](ser *Serializer, size uint, v T, serialize func(slice []byte, num T)) {
	ub := make([]byte, size)
	serialize(ub, v)
	ser.out.Write(ub)
}

func (ser *Serializer) serializeUBigInt(size uint, v *big.Int) {
	ub := make([]byte, size)
	v.FillBytes(ub)
	// Reverse, since big.Int outputs bytes in BigEndian
	slices.Reverse(ub)
	ser.out.Write(ub)
}

// U8 serialize a byte
func (ser *Serializer) U8(v uint8) {
	ser.out.WriteByte(v)
}

// U16 serialize an unsigned 16-bit integer in little-endian format
func (ser *Serializer) U16(v uint16) {
	serializeUInt(ser, 2, v, binary.LittleEndian.PutUint16)
}

// U32 serialize an unsigned 32-bit integer in little-endian format
func (ser *Serializer) U32(v uint32) {
	serializeUInt(ser, 4, v, binary.LittleEndian.PutUint32)
}

// U64 serialize an unsigned 64-bit integer in little-endian format
func (ser *Serializer) U64(v uint64) {
	serializeUInt(ser, 8, v, binary.LittleEndian.PutUint64)
}

// U128 serialize an unsigned 128-bit integer in little-endian format
func (ser *Serializer) U128(v big.Int) {
	ser.serializeUBigInt(16, &v)
}

// U256 serialize an unsigned 256-bit integer in little-endian format
func (ser *Serializer) U256(v big.Int) {
	ser.serializeUBigInt(32, &v)
}

// Uleb128 serialize an unsigned 32-bit integer as an Uleb128.  This is used specifically for sequence lengths, and enums.
func (ser *Serializer) Uleb128(val uint32) {
	for val>>7 != 0 {
		b, err := util.Uint32ToU8(val & 0xFF)
		if err != nil {
			ser.SetError(err)
			return
		}
		ser.out.WriteByte(b | 0x80)
		val >>= 7
	}
	b, err := util.Uint32ToU8(val & 0xFF)
	if err != nil {
		ser.SetError(err)
		return
	}
	ser.out.WriteByte(b)
}

// WriteBytes serialize an array of bytes with its length first as an Uleb128.
func (ser *Serializer) WriteBytes(v []byte) {
	length, err := util.IntToU32(len(v))
	if err != nil {
		ser.SetError(err)
		return
	}
	ser.Uleb128(length)
	ser.out.Write(v)
}

// WriteString similar to [Serializer.WriteBytes] using the UTF-8 byte representation of the string
func (ser *Serializer) WriteString(v string) {
	ser.WriteBytes([]byte(v))
}

// FixedBytes similar to [Serializer.WriteBytes], but it forgoes the length header.
// This is useful if you know the fixed length size of the data, such as AccountAddress
func (ser *Serializer) FixedBytes(v []byte) {
	ser.out.Write(v)
}

// Struct uses custom serialization for a [Marshaler] implementation.
func (ser *Serializer) Struct(v Marshaler) {
	if v == nil {
		ser.SetError(errors.New("cannot marshal nil"))
		return
	}
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

// SerializeSequence serializes a sequence of [Marshaler] implemented types.  Prefixed with the length of the sequence.
//
// It works with both array values by reference and by value:
//
//	type MyStruct struct {
//		num uint64
//	}
//
//	func (str *MyStruct) MarshalBCS(ser *Serialize) {
//		ser.U64(num)
//	}
//
//	myArray := []MyStruct{
//		MyStruct{num: 0},
//		MyStruct{num: 1},
//		MyStruct{num: 2},
//	}
//
//	serializer := &Serializer{}
//	SerializeSequence(myArray, ser)
//	bytes := serializer.ToBytes()
func SerializeSequence[AT []T, T any](array AT, ser *Serializer) {
	SerializeSequenceWithFunction(array, ser, func(ser *Serializer, item T) {
		// Check if by value is Marshaler
		mv, ok := any(item).(Marshaler)
		if ok {
			mv.MarshalBCS(ser)
			return
		}
		// Check if by reference is Marshaler
		mv, ok = any(&item).(Marshaler)
		if ok {
			mv.MarshalBCS(ser)
			return
		}
		// If neither works, let's pass an error up
		ser.SetError(errors.New("type or reference of type is not Marshaler"))
	})
}

// SerializeSequenceWithFunction allows custom serialization of a sequence, which can be useful for non-bcs.Struct types
//
//	array := []string{"hello", "blockchain"}
//
//	SerializeSequenceWithFunction(array, func(ser *Serializer, item string) {
//		ser.WriteString(item)
//	}
func SerializeSequenceWithFunction[AT []T, T any](array AT, ser *Serializer, serialize func(ser *Serializer, item T)) {
	length, err := util.IntToU32(len(array))
	if err != nil {
		ser.SetError(err)
		return
	}
	ser.Uleb128(length)
	for i, v := range array {
		serialize(ser, v)
		// Exit early if there's an error
		if ser.Error() != nil {
			ser.SetError(fmt.Errorf("could not serialize sequence[%d] member of %T %w", i, v, ser.Error()))
			return
		}
	}
}

// SerializeSequenceOnly serializes a sequence into a single value using [SerializeSequence]
//
//	type MyStruct struct {
//		num uint64
//	}
//
//	func (str *MyStruct) MarshalBCS(ser *Serialize) {
//		ser.U64(num)
//	}
//
//	myArray := []MyStruct{
//		MyStruct{num: 0},
//		MyStruct{num: 1},
//		MyStruct{num: 2},
//	}
//
//	bytes, err := SerializeSequenceOnly(myArray)
func SerializeSequenceOnly[AT []T, T any](input AT) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		SerializeSequence(input, ser)
	})
}

// SerializeBool Serializes a single boolean
//
//	bytes, _ := SerializeBool(true)
func SerializeBool(input bool) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.Bool(input)
	})
}

// SerializeU8 Serializes a single uint8
//
//	bytes, _ := SerializeU8(uint8(200))
func SerializeU8(input uint8) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.U8(input)
	})
}

// SerializeU16 Serializes a single uint16
//
//	bytes, _ := SerializeU16(uint16(50000))
func SerializeU16(input uint16) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.U16(input)
	})
}

// SerializeU32 Serializes a single uint32
//
//	bytes, _ := SerializeU32(uint32(50000))
func SerializeU32(input uint32) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.U32(input)
	})
}

// SerializeU64 Serializes a single uint64
//
//	bytes, _ := SerializeU64(uint64(20))
func SerializeU64(input uint64) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.U64(input)
	})
}

// SerializeU128 Serializes a single uint128
//
//	u128 := big.NewInt(1)
//	bytes, _ := SerializeU128(u128)
func SerializeU128(input big.Int) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.U128(input)
	})
}

// SerializeU256 Serializes a single uint256
//
//	u256 := big.NewInt(1)
//	bytes, _ := SerializeU256(u256)
func SerializeU256(input big.Int) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.U256(input)
	})
}

// SerializeUleb128 Serializes a single uleb128
func SerializeUleb128(input uint32) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.Uleb128(input)
	})
}

// SerializeBytes Serializes a single byte array
//
//	input := []byte{0x1, 0x2}
//	bytes, _ := SerializeBytes(input)
func SerializeBytes(input []byte) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.WriteBytes(input)
	})
}

// SerializeSingle is a convenience function, to not have to create a serializer to serialize one value
//
// Here's an example for handling a nested byte array
//
//	input := [][]byte{[]byte{0x1}, []byte{0x2}}
//	bytes, _ := SerializeSingle(func(ser *Serializer) {
//		ser.Uleb128(len(input))
//		for _, list := range input {
//			ser.WriteBytes(list)
//		}
//	})
func SerializeSingle(marshal func(ser *Serializer)) ([]byte, error) {
	ser := &Serializer{}
	marshal(ser)
	err := ser.Error()
	if err != nil {
		return nil, err
	}
	return ser.ToBytes(), nil
}

// SerializeOption serializes an optional value
//
// # Under the hood, this is represented as a 0 or 1 length array
//
// Here's an example for handling an optional value:
//
//	// For a Some(10) value
//	input := uint8(10)
//	ser := &Serializer{}
//	bytes, _ := SerializeOption(ser, &input, func(ser *Serializer, item uint8) {
//		ser.U8(item)
//	})
//	// bytes == []byte{0x01,0x0A}
//
//	// For a None value
//	ser2 := &Serializer{}
//	bytes2, _ := SerializeOption(ser2, nil, func(ser *Serializer, item uint8) {
//		ser.U8(item)
//	})
//	// bytes2 == []byte{0x00}
func SerializeOption[T any](ser *Serializer, input *T, serialize func(ser *Serializer, item T)) {
	if input == nil {
		SerializeSequenceWithFunction([]T{}, ser, serialize)
	} else {
		SerializeSequenceWithFunction([]T{*input}, ser, serialize)
	}
}

// Serialized wraps already serialized bytes with BCS serialization
// It prepends the byte array with its length (encoded as Uleb128) and writes the bytes
// Primarily used for handling nested serialization structures where bytes are already in BCS format
func (ser *Serializer) Serialized(s Serialized) {
	ser.WriteBytes(s.Value)
}

func SerializeSerialized(input Serialized) ([]byte, error) {
	return SerializeSingle(func(ser *Serializer) {
		ser.Serialized(input)
	})
}
