package base

import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
)

func EncodeCompactU16Length(bytes *[]byte, ln int) {
	rem_len := ln
	for {
		elem := rem_len & 0x7f
		rem_len >>= 7
		if rem_len == 0 {
			*bytes = append(*bytes, byte(elem))
			break
		} else {
			elem |= 0x80
			*bytes = append(*bytes, byte(elem))
		}
	}
}

// TypeID defines the internal representation of an instruction type ID
// (or account type, etc. in anchor programs)
// and it's used to associate instructions to decoders in the variant tracker.
type TypeID [8]byte

type BaseVariant struct {
	TypeID TypeID
	Impl   interface{}
}

var NoTypeIDDefaultID = TypeIDFromUint8(0)

// TypeIDFromBytes converts a []byte to a TypeID.
// The provided slice must be 8 bytes long or less.
func TypeIDFromBytes(slice []byte) (id TypeID) {
	// TODO: panic if len(slice) > 8 ???
	copy(id[:], slice)
	return id
}

// Uint32 parses the TypeID to a uint32.
func (vid TypeID) Uint32() uint32 {
	return Uint32FromTypeID(vid, binary.LittleEndian)
}

// Uint8 parses the TypeID to a Uint8.
func (vid TypeID) Uint8() uint8 {
	return Uint8FromTypeID(vid)
}

func (vid TypeID) Bytes() []byte {
	return vid[:]
}

// Uint32FromTypeID parses a TypeID bytes to a uint32.
func Uint32FromTypeID(vid TypeID, order binary.ByteOrder) (out uint32) {
	out = order.Uint32(vid[:])
	return out
}

// Uint32FromTypeID parses a TypeID bytes to a uint8.
func Uint8FromTypeID(vid TypeID) (out uint8) {
	return vid[0]
}

// TypeIDFromUint32 converts a uint8 to a TypeID.
func TypeIDFromUint8(v uint8) TypeID {
	return TypeIDFromBytes([]byte{v})
}

var TypeSize = struct {
	Bool int
	Byte int

	Int8  int
	Int16 int

	Uint8   int
	Uint16  int
	Uint32  int
	Uint64  int
	Uint128 int

	Float32 int
	Float64 int

	PublicKey int
	Signature int

	Tstamp         int
	BlockTimestamp int

	CurrencyName int
}{
	Byte: 1,
	Bool: 1,

	Int8:  1,
	Int16: 2,

	Uint8:   1,
	Uint16:  2,
	Uint32:  4,
	Uint64:  8,
	Uint128: 16,

	Float32: 4,
	Float64: 8,
}

// TypeIDFromUint32 converts a uint32 to a TypeID.
func TypeIDFromUint32(v uint32, bo binary.ByteOrder) TypeID {
	out := make([]byte, TypeSize.Uint32)
	bo.PutUint32(out, v)
	return TypeIDFromBytes(out)
}

type option struct {
	OptionalField bool
	SizeOfSlice   *int
	Order         binary.ByteOrder
}

var defaultByteOrder = binary.LittleEndian

func newDefaultOption() *option {
	return &option{
		OptionalField: false,
		Order:         defaultByteOrder,
	}
}

func (o *option) isOptional() bool {
	return o.OptionalField
}

func (o *option) hasSizeOfSlice() bool {
	return o.SizeOfSlice != nil
}

func (o *option) getSizeOfSlice() int {
	return *o.SizeOfSlice
}

func (o *option) setSizeOfSlice(size int) *option {
	o.SizeOfSlice = &size
	return o
}
func (o *option) setIsOptional(isOptional bool) *option {
	o.OptionalField = isOptional
	return o
}

type Encoder struct {
	output io.Writer
	count  int

	currentFieldOpt *option
}

func (e *Encoder) Encode(v interface{}) (err error) {
	return e.encodeBin(reflect.ValueOf(v), nil)
}

func NewBinEncoder(writer io.Writer) *Encoder {
	return &Encoder{
		output: writer,
		count:  0,
	}
}

func (e *Encoder) toWriter(bytes []byte) (err error) {
	e.count += len(bytes)
	_, err = e.output.Write(bytes)
	return
}

// Written returns the count of bytes written.
func (e *Encoder) Written() int {
	return e.count
}

func (e *Encoder) WriteBytes(b []byte, writeLength bool) error {
	if writeLength {
		if err := e.WriteLength(len(b)); err != nil {
			return err
		}
	}
	if len(b) == 0 {
		return nil
	}
	return e.toWriter(b)
}

func (e *Encoder) WriteLength(length int) error {
	return e.WriteUVarInt(length)
}

func (e *Encoder) WriteUVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) WriteVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutVarint(buf, int64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) WriteByte(b byte) (err error) {
	return e.toWriter([]byte{b})
}

func (e *Encoder) WriteBool(b bool) (err error) {
	var out byte
	if b {
		out = 1
	}
	return e.WriteByte(out)
}

func (e *Encoder) WriteUint8(i uint8) (err error) {
	return e.WriteByte(i)
}

func (e *Encoder) WriteUint16(i uint16, order binary.ByteOrder) (err error) {
	buf := make([]byte, TypeSize.Uint16)
	order.PutUint16(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) WriteInt16(i int16, order binary.ByteOrder) (err error) {
	return e.WriteUint16(uint16(i), order)
}

func (e *Encoder) WriteInt32(i int32, order binary.ByteOrder) (err error) {
	return e.WriteUint32(uint32(i), order)
}

func (e *Encoder) WriteUint32(i uint32, order binary.ByteOrder) (err error) {
	buf := make([]byte, TypeSize.Uint32)
	order.PutUint32(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) WriteInt64(i int64, order binary.ByteOrder) (err error) {
	return e.WriteUint64(uint64(i), order)
}

func (e *Encoder) WriteUint64(i uint64, order binary.ByteOrder) (err error) {
	buf := make([]byte, TypeSize.Uint64)
	order.PutUint64(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) WriteUint128(i Uint128, order binary.ByteOrder) (err error) {
	buf := make([]byte, TypeSize.Uint128)
	order.PutUint64(buf, i.Lo)
	order.PutUint64(buf[TypeSize.Uint64:], i.Hi)
	return e.toWriter(buf)
}

func (e *Encoder) WriteInt128(i Int128, order binary.ByteOrder) (err error) {
	buf := make([]byte, TypeSize.Uint128)
	order.PutUint64(buf, i.Lo)
	order.PutUint64(buf[TypeSize.Uint64:], i.Hi)
	return e.toWriter(buf)
}

func (e *Encoder) WriteFloat32(f float32, order binary.ByteOrder) (err error) {
	i := math.Float32bits(f)
	buf := make([]byte, TypeSize.Uint32)
	order.PutUint32(buf, i)

	return e.toWriter(buf)
}
func (e *Encoder) WriteFloat64(f float64, order binary.ByteOrder) (err error) {
	i := math.Float64bits(f)
	buf := make([]byte, TypeSize.Uint64)
	order.PutUint64(buf, i)

	return e.toWriter(buf)
}

func (e *Encoder) WriteString(s string) (err error) {
	return e.WriteBytes([]byte(s), true)
}

func (e *Encoder) WriteRustString(s string) (err error) {
	err = e.WriteUint64(uint64(len(s)), binary.LittleEndian)
	if err != nil {
		return err
	}
	return e.WriteBytes([]byte(s), false)
}

func (e *Encoder) WriteCompactU16Length(ln int) (err error) {
	buf := make([]byte, 0)
	EncodeCompactU16Length(&buf, ln)
	return e.toWriter(buf)
}

type fieldTag struct {
	SizeOf          string
	Skip            bool
	Order           binary.ByteOrder
	Optional        bool
	BinaryExtension bool

	IsBorshEnum bool
}

func parseFieldTag(tag reflect.StructTag) *fieldTag {
	t := &fieldTag{
		Order: defaultByteOrder,
	}
	tagStr := tag.Get("bin")
	for _, s := range strings.Split(tagStr, " ") {
		if strings.HasPrefix(s, "sizeof=") {
			tmp := strings.SplitN(s, "=", 2)
			t.SizeOf = tmp[1]
		} else if s == "big" {
			t.Order = binary.BigEndian
		} else if s == "little" {
			t.Order = binary.LittleEndian
		} else if s == "optional" {
			t.Optional = true
		} else if s == "binary_extension" {
			t.BinaryExtension = true
		} else if s == "-" {
			t.Skip = true
		}
	}

	// TODO: parse other borsh tags
	if strings.TrimSpace(tag.Get("borsh_skip")) == "true" {
		t.Skip = true
	}
	if strings.TrimSpace(tag.Get("borsh_enum")) == "true" {
		t.IsBorshEnum = true
	}
	return t
}

func isZero(rv reflect.Value) (b bool) {
	return rv.Kind() == 0
}

func sizeof(t reflect.Type, v reflect.Value) int {
	switch t.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return int(v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		n := int(v.Uint())
		// all the builtin array length types are native int
		// so this guards against weird truncation
		if n < 0 {
			return 0
		}
		return n
	default:
		panic("sizeof field ")
	}
}

type BinaryMarshaler interface {
	MarshalWithEncoder(encoder *Encoder) error
}

func (e *Encoder) encodeBin(rv reflect.Value, opt *option) (err error) {
	if opt == nil {
		opt = newDefaultOption()
	}
	e.currentFieldOpt = opt

	if opt.isOptional() {
		if rv.IsZero() {
			return e.WriteUint32(0, binary.LittleEndian)
		}
		err := e.WriteUint32(1, binary.LittleEndian)
		if err != nil {
			return err
		}
		// The optionality has been used; stop its propagation:
		opt.setIsOptional(false)
	}

	if isZero(rv) {
		return nil
	}

	if marshaler, ok := rv.Interface().(BinaryMarshaler); ok {
		return marshaler.MarshalWithEncoder(e)
	}

	switch rv.Kind() {
	case reflect.String:
		return e.WriteRustString(rv.String())
	case reflect.Uint8:
		return e.WriteByte(byte(rv.Uint()))
	case reflect.Int8:
		return e.WriteByte(byte(rv.Int()))
	case reflect.Int16:
		return e.WriteInt16(int16(rv.Int()), opt.Order)
	case reflect.Uint16:
		return e.WriteUint16(uint16(rv.Uint()), opt.Order)
	case reflect.Int32:
		return e.WriteInt32(int32(rv.Int()), opt.Order)
	case reflect.Uint32:
		return e.WriteUint32(uint32(rv.Uint()), opt.Order)
	case reflect.Uint64:
		return e.WriteUint64(rv.Uint(), opt.Order)
	case reflect.Int64:
		return e.WriteInt64(rv.Int(), opt.Order)
	case reflect.Float32:
		return e.WriteFloat32(float32(rv.Float()), opt.Order)
	case reflect.Float64:
		return e.WriteFloat64(rv.Float(), opt.Order)
	case reflect.Bool:
		return e.WriteBool(rv.Bool())
	case reflect.Ptr:
		return e.encodeBin(rv.Elem(), opt)
	case reflect.Interface:
		// skip
		return nil
	}

	rv = reflect.Indirect(rv)
	rt := rv.Type()
	switch rt.Kind() {
	case reflect.Array:
		l := rt.Len()
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			// if it's a [n]byte, accumulate and write in one command:
			arr := make([]byte, l)
			for i := 0; i < l; i++ {
				arr[i] = byte(rv.Index(i).Uint())
			}
			if err := e.WriteBytes(arr, false); err != nil {
				return err
			}
		} else {
			for i := 0; i < l; i++ {
				if err = e.encodeBin(rv.Index(i), nil); err != nil {
					return
				}
			}
		}
	case reflect.Slice:
		var l int
		if opt.hasSizeOfSlice() {
			l = opt.getSizeOfSlice()
		} else {
			l = rv.Len()
			if err = e.WriteUVarInt(l); err != nil {
				return
			}
		}
		// we would want to skip to the correct head_offset

		for i := 0; i < l; i++ {
			if err = e.encodeBin(rv.Index(i), nil); err != nil {
				return
			}
		}
	case reflect.Struct:
		if err = e.encodeStructBin(rt, rv); err != nil {
			return
		}

	case reflect.Map:
		keyCount := len(rv.MapKeys())

		if err = e.WriteUVarInt(keyCount); err != nil {
			return
		}

		for _, mapKey := range rv.MapKeys() {
			if err = e.Encode(mapKey.Interface()); err != nil {
				return
			}

			if err = e.Encode(rv.MapIndex(mapKey).Interface()); err != nil {
				return
			}
		}

	default:
		return fmt.Errorf("encode: unsupported type %q", rt)
	}
	return
}

func (e *Encoder) encodeStructBin(rt reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()

	sizeOfMap := map[string]int{}
	for i := 0; i < l; i++ {
		structField := rt.Field(i)
		fieldTag := parseFieldTag(structField.Tag)

		if fieldTag.Skip {
			continue
		}

		rv := rv.Field(i)

		if fieldTag.SizeOf != "" {
			sizeOfMap[fieldTag.SizeOf] = sizeof(structField.Type, rv)
		}

		if !rv.CanInterface() {
			continue
		}

		option := &option{
			OptionalField: fieldTag.Optional,
			Order:         fieldTag.Order,
		}

		if s, ok := sizeOfMap[structField.Name]; ok {
			option.setSizeOfSlice(s)
		}

		if err := e.encodeBin(rv, option); err != nil {
			return fmt.Errorf("error while encoding %q field: %w", structField.Name, err)
		}
	}
	return nil
}
