package base

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"unicode/utf8"
)

// Decoder implements the EOS unpacking, similar to FC_BUFFER
type Decoder struct {
	data []byte
	pos  int

	currentFieldOpt *option

	encoding Encoding
}

// Reset resets the decoder to decode a new message.
func (dec *Decoder) Reset(data []byte) {
	dec.data = data
	dec.pos = 0
	dec.currentFieldOpt = nil
}

func (dec *Decoder) IsBorsh() bool {
	return dec.encoding.IsBorsh()
}

func (dec *Decoder) IsBin() bool {
	return dec.encoding.IsBin()
}

func (dec *Decoder) IsCompactU16() bool {
	return dec.encoding.IsCompactU16()
}

func NewDecoderWithEncoding(data []byte, enc Encoding) *Decoder {
	if !isValidEncoding(enc) {
		panic(fmt.Sprintf("provided encoding is not valid: %s", enc))
	}
	return &Decoder{
		data:     data,
		encoding: enc,
	}
}

// SetEncoding sets the encoding scheme to use for decoding.
func (dec *Decoder) SetEncoding(enc Encoding) {
	dec.encoding = enc
}

func NewBinDecoder(data []byte) *Decoder {
	return NewDecoderWithEncoding(data, EncodingBin)
}

func (dec *Decoder) Decode(v interface{}) (err error) {
	switch dec.encoding {
	case EncodingBin:
		return dec.decodeWithOptionBin(v, nil)
	default:
		panic(fmt.Errorf("encoding not implemented: %s", dec.encoding))
	}
}

var ErrVarIntBufferSize = errors.New("varint: invalid buffer size")

func (dec *Decoder) ReadUvarint64() (uint64, error) {
	l, read := binary.Uvarint(dec.data[dec.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	dec.pos += read
	return l, nil
}

func (d *Decoder) ReadVarint64() (out int64, err error) {
	l, read := binary.Varint(d.data[d.pos:])
	if read <= 0 {
		return l, ErrVarIntBufferSize
	}
	d.pos += read
	return l, nil
}

func (dec *Decoder) ReadVarint32() (out int32, err error) {
	n, err := dec.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int32(n)
	return
}

func (dec *Decoder) ReadUvarint32() (out uint32, err error) {
	n, err := dec.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint32(n)
	return
}

func (dec *Decoder) ReadVarint16() (out int16, err error) {
	n, err := dec.ReadVarint64()
	if err != nil {
		return out, err
	}
	out = int16(n)
	return
}

func (dec *Decoder) ReadUvarint16() (out uint16, err error) {
	n, err := dec.ReadUvarint64()
	if err != nil {
		return out, err
	}
	out = uint16(n)
	return
}

func (dec *Decoder) ReadByteSlice() (out []byte, err error) {
	length, err := dec.ReadLength()
	if err != nil {
		return nil, err
	}

	if len(dec.data) < dec.pos+length {
		return nil, fmt.Errorf("byte array: varlen=%d, missing %d bytes", length, dec.pos+length-len(dec.data))
	}

	out = dec.data[dec.pos : dec.pos+length]
	dec.pos += length
	return
}

func (dec *Decoder) ReadLength() (length int, err error) {
	switch dec.encoding {
	case EncodingBin:
		val, err := dec.ReadUvarint64()
		if err != nil {
			return 0, err
		}
		if val > 0x7FFF_FFFF {
			return 0, io.ErrUnexpectedEOF
		}
		length = int(val)
	default:
		panic(fmt.Errorf("encoding not implemented: %s", dec.encoding))
	}
	return
}

type peekAbleByteReader interface {
	io.ByteReader
	Peek(n int) ([]byte, error)
}

func readNBytes(n int, reader *Decoder) ([]byte, error) {
	if n == 0 {
		return make([]byte, 0), nil
	}
	if n < 0 || n > 0x7FFF_FFFF {
		return nil, fmt.Errorf("invalid length n: %v", n)
	}
	if reader.pos+n > len(reader.data) {
		return nil, fmt.Errorf("not enough data: %d bytes missing", reader.pos+n-len(reader.data))
	}
	out := reader.data[reader.pos : reader.pos+n]
	reader.pos += n
	return out, nil
}

func discardNBytes(n int, reader *Decoder) error {
	if n == 0 {
		return nil
	}
	if n < 0 || n > 0x7FFF_FFFF {
		return fmt.Errorf("invalid length n: %v", n)
	}
	return reader.SkipBytes(uint(n))
}

func (d *Decoder) Read(buf []byte) (int, error) {
	if d.pos+len(buf) > len(d.data) {
		return 0, io.ErrShortBuffer
	}
	numCopied := copy(buf, d.data[d.pos:])
	d.pos += numCopied
	// must read exactly len(buf) bytes
	if numCopied != len(buf) {
		return 0, io.ErrUnexpectedEOF
	}
	return len(buf), nil
}

func (dec *Decoder) ReadNBytes(n int) (out []byte, err error) {
	return readNBytes(n, dec)
}

// ReadBytes reads a byte slice of length n.
func (dec *Decoder) ReadBytes(n int) (out []byte, err error) {
	return readNBytes(n, dec)
}

func (dec *Decoder) Discard(n int) (err error) {
	return discardNBytes(n, dec)
}

func (dec *Decoder) ReadTypeID() (out TypeID, err error) {
	discriminator, err := dec.ReadNBytes(8)
	if err != nil {
		return TypeID{}, err
	}
	return TypeIDFromBytes(discriminator), nil
}

func (dec *Decoder) Peek(n int) (out []byte, err error) {
	if n < 0 {
		err = fmt.Errorf("n not valid: %d", n)
		return
	}

	requiredSize := TypeSize.Byte * n
	if dec.Remaining() < requiredSize {
		err = fmt.Errorf("required [%d] bytes, remaining [%d]", requiredSize, dec.Remaining())
		return
	}

	out = dec.data[dec.pos : dec.pos+n]
	return
}

// ReadCompactU16 reads a compact u16 from the decoder.
func (dec *Decoder) ReadCompactU16() (out int, err error) {
	out, err = DecodeCompactU16LengthFromByteReader(dec)
	return
}

func (dec *Decoder) ReadByte() (out byte, err error) {
	if dec.Remaining() < TypeSize.Byte {
		err = fmt.Errorf("required [1] byte, remaining [%d]", dec.Remaining())
		return
	}

	out = dec.data[dec.pos]
	dec.pos++
	return
}

func (dec *Decoder) ReadBool() (out bool, err error) {
	if dec.Remaining() < TypeSize.Bool {
		err = fmt.Errorf("bool required [%d] byte, remaining [%d]", TypeSize.Bool, dec.Remaining())
		return
	}

	b, err := dec.ReadByte()

	if err != nil {
		err = fmt.Errorf("readBool, %s", err)
	}
	out = b != 0
	return
}

func (dec *Decoder) ReadUint8() (out uint8, err error) {
	out, err = dec.ReadByte()
	return
}

func (dec *Decoder) ReadInt8() (out int8, err error) {
	b, err := dec.ReadByte()
	out = int8(b)
	return
}

func (dec *Decoder) ReadUint16(order binary.ByteOrder) (out uint16, err error) {
	if dec.Remaining() < TypeSize.Uint16 {
		err = fmt.Errorf("uint16 required [%d] bytes, remaining [%d]", TypeSize.Uint16, dec.Remaining())
		return
	}

	out = order.Uint16(dec.data[dec.pos:])
	dec.pos += TypeSize.Uint16
	return
}

func (dec *Decoder) ReadInt16(order binary.ByteOrder) (out int16, err error) {
	n, err := dec.ReadUint16(order)
	out = int16(n)
	return
}

func (dec *Decoder) ReadInt64(order binary.ByteOrder) (out int64, err error) {
	n, err := dec.ReadUint64(order)
	out = int64(n)
	return
}

func (dec *Decoder) ReadUint32(order binary.ByteOrder) (out uint32, err error) {
	if dec.Remaining() < TypeSize.Uint32 {
		err = fmt.Errorf("uint32 required [%d] bytes, remaining [%d]", TypeSize.Uint32, dec.Remaining())
		return
	}

	out = order.Uint32(dec.data[dec.pos:])
	dec.pos += TypeSize.Uint32
	return
}

func (dec *Decoder) ReadInt32(order binary.ByteOrder) (out int32, err error) {
	n, err := dec.ReadUint32(order)
	out = int32(n)
	return
}

func (dec *Decoder) ReadUint64(order binary.ByteOrder) (out uint64, err error) {
	if dec.Remaining() < TypeSize.Uint64 {
		err = fmt.Errorf("decode: uint64 required [%d] bytes, remaining [%d]", TypeSize.Uint64, dec.Remaining())
		return
	}

	data, err := dec.ReadNBytes(TypeSize.Uint64)
	if err != nil {
		return 0, err
	}
	out = order.Uint64(data)
	return
}

func (dec *Decoder) ReadInt128(order binary.ByteOrder) (out Int128, err error) {
	v, err := dec.ReadUint128(order)
	if err != nil {
		return
	}
	return Int128(v), nil
}

func (dec *Decoder) ReadUint128(order binary.ByteOrder) (out Uint128, err error) {
	if dec.Remaining() < TypeSize.Uint128 {
		err = fmt.Errorf("uint128 required [%d] bytes, remaining [%d]", TypeSize.Uint128, dec.Remaining())
		return
	}

	data := dec.data[dec.pos : dec.pos+TypeSize.Uint128]

	if order == binary.LittleEndian {
		out.Lo = order.Uint64(data[:8])
		out.Hi = order.Uint64(data[8:])
	} else {
		// TODO: is this correct?
		out.Hi = order.Uint64(data[:8])
		out.Lo = order.Uint64(data[8:])
	}

	dec.pos += TypeSize.Uint128
	return
}

func (dec *Decoder) ReadFloat32(order binary.ByteOrder) (out float32, err error) {
	if dec.Remaining() < TypeSize.Float32 {
		err = fmt.Errorf("float32 required [%d] bytes, remaining [%d]", TypeSize.Float32, dec.Remaining())
		return
	}

	n := order.Uint32(dec.data[dec.pos:])
	out = math.Float32frombits(n)
	dec.pos += TypeSize.Float32

	if dec.IsBorsh() {
		if math.IsNaN(float64(out)) {
			return 0, errors.New("NaN for float not allowed")
		}
	}
	return
}

func (dec *Decoder) ReadFloat64(order binary.ByteOrder) (out float64, err error) {
	if dec.Remaining() < TypeSize.Float64 {
		err = fmt.Errorf("float64 required [%d] bytes, remaining [%d]", TypeSize.Float64, dec.Remaining())
		return
	}

	n := order.Uint64(dec.data[dec.pos:])
	out = math.Float64frombits(n)
	dec.pos += TypeSize.Float64

	if dec.IsBorsh() {
		if math.IsNaN(out) {
			return 0, errors.New("NaN for float not allowed")
		}
	}
	return
}

func (dec *Decoder) ReadFloat128(order binary.ByteOrder) (out Float128, err error) {
	value, err := dec.ReadUint128(order)
	if err != nil {
		return out, fmt.Errorf("float128: %s", err)
	}
	return Float128(value), nil
}

func (dec *Decoder) SafeReadUTF8String() (out string, err error) {
	data, err := dec.ReadByteSlice()
	out = strings.Map(fixUtf, string(data))
	return
}

func fixUtf(r rune) rune {
	if r == utf8.RuneError {
		return '�'
	}
	return r
}

func (dec *Decoder) ReadString() (out string, err error) {
	data, err := dec.ReadByteSlice()
	out = string(data)
	return
}

func (dec *Decoder) ReadRustString() (out string, err error) {
	length, err := dec.ReadUint64(binary.LittleEndian)
	if err != nil {
		return "", err
	}
	if length > 0x7FFF_FFFF {
		return "", io.ErrUnexpectedEOF
	}
	bytes, err := dec.ReadNBytes(int(length))
	if err != nil {
		return "", err
	}
	out = string(bytes)
	return
}

func (dec *Decoder) ReadCompactU16Length() (int, error) {
	val, err := DecodeCompactU16LengthFromByteReader(dec)
	return val, err
}

func (dec *Decoder) SkipBytes(count uint) error {
	if uint(dec.Remaining()) < count {
		return fmt.Errorf("request to skip %d but only %d bytes remain", count, dec.Remaining())
	}
	dec.pos += int(count)
	return nil
}

func (dec *Decoder) SetPosition(idx uint) error {
	if int(idx) < len(dec.data) {
		dec.pos = int(idx)
		return nil
	}
	return fmt.Errorf("request to set position to %d outsize of buffer (buffer size %d)", idx, len(dec.data))
}

func (dec *Decoder) Position() uint {
	return uint(dec.pos)
}

func (dec *Decoder) Remaining() int {
	return len(dec.data) - dec.pos
}

func (dec *Decoder) HasRemaining() bool {
	return dec.Remaining() > 0
}

// indirect walks down v allocating pointers as needed,
// until it gets to a non-pointer.
// if it encounters an Unmarshaler, indirect stops and returns that.
// if decodingNull is true, indirect stops at the last pointer so it can be set to nil.
//
// *Note* This is a copy of `encoding/json/decoder.go#indirect` of Golang 1.14.
//
// See here: https://github.com/golang/go/blob/go1.14.2/src/encoding/json/decode.go#L439
func indirect(v reflect.Value, decodingNull bool) (BinaryUnmarshaler, reflect.Value) {
	// Issue #24153 indicates that it is generally not a guaranteed property
	// that you may round-trip a reflect.Value by calling Value.Addr().Elem()
	// and expect the value to still be settable for values derived from
	// unexported embedded struct fields.
	//
	// The logic below effectively does this when it first addresses the value
	// (to satisfy possible pointer methods) and continues to dereference
	// subsequent pointers as necessary.
	//
	// After the first round-trip, we set v back to the original value to
	// preserve the original RW flags contained in reflect.Value.
	v0 := v
	haveAddr := false

	// If v is a named type and is addressable,
	// start with its address, so that if the type has pointer methods,
	// we find them.
	if v.Kind() != reflect.Ptr && v.Type().Name() != "" && v.CanAddr() {
		haveAddr = true
		v = v.Addr()
	}
	for {
		// Load value from interface, but only if the result will be
		// usefully addressable.
		if v.Kind() == reflect.Interface && !v.IsNil() {
			e := v.Elem()
			if e.Kind() == reflect.Ptr && !e.IsNil() && (!decodingNull || e.Elem().Kind() == reflect.Ptr) {
				haveAddr = false
				v = e
				continue
			}
		}

		if v.Kind() != reflect.Ptr {
			break
		}

		if v.Elem().Kind() != reflect.Ptr && decodingNull && v.CanSet() {
			break
		}

		// Prevent infinite loop if v is an interface pointing to its own address:
		//     var v interface{}
		//     v = &v
		if v.Elem().Kind() == reflect.Interface && v.Elem().Elem() == v {
			v = v.Elem()
			break
		}
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		if v.Type().NumMethod() > 0 && v.CanInterface() {
			if u, ok := v.Interface().(BinaryUnmarshaler); ok {
				return u, reflect.Value{}
			}
		}

		if haveAddr {
			v = v0 // restore original value after round-trip Value.Addr().Elem()
			haveAddr = false
		} else {
			v = v.Elem()
		}
	}
	return nil, v
}

func reflect_readArrayOfBytes(d *Decoder, l int, rv reflect.Value) error {
	buf, err := d.ReadNBytes(l)
	if err != nil {
		return err
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint8, but a custom type like [n]CustomUint8:
		if rv.Type().Elem() != typeOfUint8 {
			// if the type of the array is not [n]uint8, but a custom type like [n]CustomUint8:
			// then we need to convert each uint8 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint8, but a custom type like []CustomUint8:
		if rv.Type().Elem() != typeOfUint8 {
			// convert the []uint8 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

func reflect_readArrayOfUint16(d *Decoder, l int, rv reflect.Value, order binary.ByteOrder) error {
	buf := make([]uint16, l)
	for i := 0; i < l; i++ {
		n, err := d.ReadUint16(order)
		if err != nil {
			return err
		}
		buf[i] = n
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint16, but a custom type like [n]CustomUint16:
		if rv.Type().Elem() != typeOfUint16 {
			// if the type of the array is not [n]uint16, but a custom type like [n]CustomUint16:
			// then we need to convert each uint16 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint16, but a custom type like []CustomUint16:
		if rv.Type().Elem() != typeOfUint16 {
			// convert the []uint16 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

func reflect_readArrayOfUint32(d *Decoder, l int, rv reflect.Value, order binary.ByteOrder) error {
	buf := make([]uint32, l)
	for i := 0; i < l; i++ {
		n, err := d.ReadUint32(order)
		if err != nil {
			return err
		}
		buf[i] = n
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint32, but a custom type like [n]CustomUint32:
		if rv.Type().Elem() != typeOfUint32 {
			// if the type of the array is not [n]uint32, but a custom type like [n]CustomUint32:
			// then we need to convert each uint32 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint32, but a custom type like []CustomUint32:
		if rv.Type().Elem() != typeOfUint32 {
			// convert the []uint32 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

func init() {
	if typeOfByte != typeOfUint8 {
		panic("typeOfByte != typeOfUint8")
	}
}

var (
	typeOfByte   = reflect.TypeOf(byte(0))
	typeOfUint8  = reflect.TypeOf(uint8(0))
	typeOfUint16 = reflect.TypeOf(uint16(0))
	typeOfUint32 = reflect.TypeOf(uint32(0))
	typeOfUint64 = reflect.TypeOf(uint64(0))
)

func reflect_readArrayOfUint64(d *Decoder, l int, rv reflect.Value, order binary.ByteOrder) error {
	buf := make([]uint64, l)
	for i := 0; i < l; i++ {
		n, err := d.ReadUint64(order)
		if err != nil {
			return err
		}
		buf[i] = n
	}
	switch rv.Kind() {
	case reflect.Array:
		// if the type of the array is not [n]uint64, but a custom type like [n]CustomUint64:
		if rv.Type().Elem() != typeOfUint64 {
			// if the type of the array is not [n]uint64, but a custom type like [n]CustomUint64:
			// then we need to convert each uint64 to the custom type
			for i := 0; i < l; i++ {
				rv.Index(i).Set(reflect.ValueOf(buf[i]).Convert(rv.Index(i).Type()))
			}
		} else {
			reflect.Copy(rv, reflect.ValueOf(buf))
		}
	case reflect.Slice:
		// if the type of the slice is not []uint64, but a custom type like []CustomUint64:
		if rv.Type().Elem() != typeOfUint64 {
			// convert the []uint64 to the custom type
			customSlice := reflect.MakeSlice(rv.Type(), len(buf), len(buf))
			for i := 0; i < len(buf); i++ {
				customSlice.Index(i).SetUint(uint64(buf[i]))
			}
			rv.Set(customSlice)
		} else {
			rv.Set(reflect.ValueOf(buf))
		}
	default:
		return fmt.Errorf("unsupported kind: %s", rv.Kind())
	}
	return nil
}

// reflect_readArrayOfUint_ is used for reading arrays/slices of uints of any size.
func reflect_readArrayOfUint_(d *Decoder, l int, k reflect.Kind, rv reflect.Value, order binary.ByteOrder) error {
	switch k {
	// case reflect.Uint:
	// 	// switch on system architecture (32 or 64 bit)
	// 	if unsafe.Sizeof(uintptr(0)) == 4 {
	// 		return reflect_readArrayOfUint32(  d, l, rv, order)
	// 	}
	// 	return reflect_readArrayOfUint64(  d, l, rv, order)
	case reflect.Uint8:
		if l > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfBytes(d, l, rv)
	case reflect.Uint16:
		if l*2 > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfUint16(d, l, rv, order)
	case reflect.Uint32:
		if l*4 > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfUint32(d, l, rv, order)
	case reflect.Uint64:
		if l*8 > d.Remaining() {
			return io.ErrUnexpectedEOF
		}
		return reflect_readArrayOfUint64(d, l, rv, order)
	default:
		return fmt.Errorf("unsupported kind: %v", k)
	}
}

func (dec *Decoder) decodeWithOptionBin(v interface{}, option *option) (err error) {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return &InvalidDecoderError{reflect.TypeOf(v)}
	}

	// We decode rv not rv.Elem because the Unmarshaler interface
	// test must be applied at the top level of the value.
	err = dec.decodeBin(rv, option)
	if err != nil {
		return err
	}
	return nil
}

func (dec *Decoder) decodeBin(rv reflect.Value, opt *option) (err error) {
	if opt == nil {
		opt = newDefaultOption()
	}
	dec.currentFieldOpt = opt

	unmarshaler, rv := indirect(rv, opt.isOptional())

	if opt.isOptional() {
		isPresent, e := dec.ReadUint32(binary.LittleEndian)
		if e != nil {
			err = fmt.Errorf("decode: %s isPresent, %s", rv.Type().String(), e)
			return
		}

		if isPresent == 0 {
			rv.Set(reflect.Zero(rv.Type()))
			return
		}

		// we have ptr here we should not go get the element
		unmarshaler, rv = indirect(rv, false)
	}

	if unmarshaler != nil {
		return unmarshaler.UnmarshalWithDecoder(dec)
	}
	rt := rv.Type()

	switch rv.Kind() {
	case reflect.String:
		s, e := dec.ReadRustString()
		if e != nil {
			err = e
			return
		}
		rv.SetString(s)
		return
	case reflect.Uint8:
		var n byte
		n, err = dec.ReadByte()
		rv.SetUint(uint64(n))
		return
	case reflect.Int8:
		var n int8
		n, err = dec.ReadInt8()
		rv.SetInt(int64(n))
		return
	case reflect.Int16:
		var n int16
		n, err = dec.ReadInt16(opt.Order)
		rv.SetInt(int64(n))
		return
	case reflect.Int32:
		var n int32
		n, err = dec.ReadInt32(opt.Order)
		rv.SetInt(int64(n))
		return
	case reflect.Int64:
		var n int64
		n, err = dec.ReadInt64(opt.Order)
		rv.SetInt(int64(n))
		return
	case reflect.Uint16:
		var n uint16
		n, err = dec.ReadUint16(opt.Order)
		rv.SetUint(uint64(n))
		return
	case reflect.Uint32:
		var n uint32
		n, err = dec.ReadUint32(opt.Order)
		rv.SetUint(uint64(n))
		return
	case reflect.Uint64:
		var n uint64
		n, err = dec.ReadUint64(opt.Order)
		rv.SetUint(n)
		return
	case reflect.Float32:
		var n float32
		n, err = dec.ReadFloat32(opt.Order)
		rv.SetFloat(float64(n))
		return
	case reflect.Float64:
		var n float64
		n, err = dec.ReadFloat64(opt.Order)
		rv.SetFloat(n)
		return
	case reflect.Bool:
		var r bool
		r, err = dec.ReadBool()
		rv.SetBool(r)
		return
	case reflect.Interface:
		// skip
		return nil
	}
	switch rt.Kind() {
	case reflect.Array:
		l := rt.Len()

		switch k := rv.Type().Elem().Kind(); k {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if err := reflect_readArrayOfUint_(dec, l, k, rv, LE); err != nil {
				return err
			}
		default:
			for i := 0; i < l; i++ {
				if err = dec.decodeBin(rv.Index(i), nil); err != nil {
					return
				}
			}
		}
		return
	case reflect.Slice:
		var l int
		if opt.hasSizeOfSlice() {
			l = opt.getSizeOfSlice()
		} else {
			length, err := dec.ReadLength()
			if err != nil {
				return err
			}
			l = length
		}

		if l > dec.Remaining() {
			return io.ErrUnexpectedEOF
		}

		switch k := rv.Type().Elem().Kind(); k {
		case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if err := reflect_readArrayOfUint_(dec, l, k, rv, LE); err != nil {
				return err
			}
		default:
			rv.Set(reflect.MakeSlice(rt, 0, 0))
			for i := 0; i < l; i++ {
				// create new element of type rt:
				element := reflect.New(rt.Elem())
				// decode into element:
				if err = dec.decodeBin(element, nil); err != nil {
					return
				}
				// append to slice:
				rv.Set(reflect.Append(rv, element.Elem()))
			}
		}

	case reflect.Struct:
		if err = dec.decodeStructBin(rt, rv); err != nil {
			return
		}

	case reflect.Map:
		l, err := dec.ReadLength()
		if err != nil {
			return err
		}
		if l == 0 {
			// If the map has no content, keep it nil.
			return nil
		}
		rv.Set(reflect.MakeMap(rt))
		for i := 0; i < int(l); i++ {
			key := reflect.New(rt.Key())
			err := dec.decodeBin(key.Elem(), nil)
			if err != nil {
				return err
			}
			val := reflect.New(rt.Elem())
			err = dec.decodeBin(val.Elem(), nil)
			if err != nil {
				return err
			}
			rv.SetMapIndex(key.Elem(), val.Elem())
		}
		return nil

	default:
		return fmt.Errorf("decode: unsupported type %q", rt)
	}

	return
}

func (dec *Decoder) decodeStructBin(rt reflect.Type, rv reflect.Value) (err error) {
	l := rv.NumField()

	sizeOfMap := map[string]int{}
	seenBinaryExtensionField := false
	for i := 0; i < l; i++ {
		structField := rt.Field(i)
		fieldTag := parseFieldTag(structField.Tag)

		if fieldTag.Skip {
			continue
		}

		if !fieldTag.BinaryExtension && seenBinaryExtensionField {
			panic(fmt.Sprintf("the `bin:\"binary_extension\"` tags must be packed together at the end of struct fields, problematic field %q", structField.Name))
		}

		if fieldTag.BinaryExtension {
			seenBinaryExtensionField = true
			// FIXME: This works only if what is in `d.data` is the actual full data buffer that
			//        needs to be decoded. If there is for example two structs in the buffer, this
			//        will not work as we would continue into the next struct.
			//
			//        But at the same time, does it make sense otherwise? What would be the inference
			//        rule in the case of extra bytes available? Continue decoding and revert if it's
			//        not working? But how to detect valid errors?
			if len(dec.data[dec.pos:]) <= 0 {
				continue
			}
		}
		v := rv.Field(i)
		if !v.CanSet() {
			// This means that the field cannot be set, to fix this
			// we need to create a pointer to said field
			if !v.CanAddr() {
				// we cannot create a point to field skipping
				return fmt.Errorf("unable to decode a none setup struc field %q with type %q", structField.Name, v.Kind())
			}
			v = v.Addr()
		}

		if !v.CanSet() {
			continue
		}

		option := &option{
			OptionalField: fieldTag.Optional,
			Order:         fieldTag.Order,
		}

		if s, ok := sizeOfMap[structField.Name]; ok {
			option.setSizeOfSlice(s)
		}

		if err = dec.decodeBin(v, option); err != nil {
			return fmt.Errorf("error while decoding %q field: %w", structField.Name, err)
		}

		if fieldTag.SizeOf != "" {
			size := sizeof(structField.Type, v)
			sizeOfMap[fieldTag.SizeOf] = size
		}
	}
	return
}

// DecodeCompactU16Length decodes a "Compact-u16" length from the provided byte slice.
func DecodeCompactU16Length(bytes []byte) int {
	ln := 0
	size := 0
	for {
		elem := int(bytes[0])
		bytes = bytes[1:]
		ln |= (elem & 0x7f) << (size * 7)
		size += 1
		if (elem & 0x80) == 0 {
			break
		}
	}
	return ln
}

// DecodeCompactU16LengthFromByteReader decodes a "Compact-u16" length from the provided io.ByteReader.
func DecodeCompactU16LengthFromByteReader(reader io.ByteReader) (int, error) {
	ln := 0
	size := 0
	for {
		elemByte, err := reader.ReadByte()
		if err != nil {
			return 0, err
		}
		elem := int(elemByte)
		ln |= (elem & 0x7f) << (size * 7)
		size += 1
		if (elem & 0x80) == 0 {
			break
		}
	}
	return ln, nil
}

type BinaryUnmarshaler interface {
	UnmarshalWithDecoder(decoder *Decoder) error
}

// An InvalidDecoderError describes an invalid argument passed to Decoder.
// (The argument to Decoder must be a non-nil pointer.)
type InvalidDecoderError struct {
	Type reflect.Type
}

func (e *InvalidDecoderError) Error() string {
	if e.Type == nil {
		return "decoder: Decode(nil)"
	}

	if e.Type.Kind() != reflect.Ptr {
		return "decoder: Decode(non-pointer " + e.Type.String() + ")"
	}
	return "decoder: Decode(nil " + e.Type.String() + ")"
}
