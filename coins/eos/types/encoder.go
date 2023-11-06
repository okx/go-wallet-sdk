package types

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/eoscanada/eos-go/ecc"

	"io"
	"math"
	"reflect"
)

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

	Checksum160 int
	Checksum256 int
	Checksum512 int

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

	Checksum160: 20,
	Checksum256: 32,
	Checksum512: 64,

	PublicKey: 34,
	Signature: 66,

	Tstamp:         8,
	BlockTimestamp: 4,

	CurrencyName: 7,
}

type MarshalerBinary interface {
	MarshalBinary(encoder *Encoder) error
}

func MarshalBinary(v interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)
	encoder := NewEncoder(buf)
	err := encoder.Encode(v)
	return buf.Bytes(), err
}

type Encoder struct {
	output io.Writer
	Order  binary.ByteOrder
	count  int
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		output: w,
		Order:  binary.LittleEndian,
		count:  0,
	}
}

func (e *Encoder) toWriter(bytes []byte) (err error) {
	e.count += len(bytes)
	_, err = e.output.Write(bytes)
	return
}

func (e *Encoder) writeUVarInt(v int) (err error) {
	buf := make([]byte, 8)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) writeByteArray(b []byte) error {
	if err := e.writeUVarInt(len(b)); err != nil {
		return err
	}
	return e.toWriter(b)
}

func (e *Encoder) writeString(s string) (err error) {
	return e.writeByteArray([]byte(s))
}

func (e *Encoder) writeByte(b byte) (err error) {
	return e.toWriter([]byte{b})
}

func (e *Encoder) writeBool(b bool) (err error) {
	var out byte
	if b {
		out = 1
	}
	return e.writeByte(out)
}

func (e *Encoder) writeName(name Name) error {
	val, err := StringToName(string(name))
	if err != nil {
		return fmt.Errorf("writeName: %w", err)
	}
	return e.writeUint64(val)
}

func (e *Encoder) writeSignature(s ecc.Signature) (err error) {
	err = s.Validate()
	if err != nil {
		return fmt.Errorf("invalid signature: %w", err)
	}

	if err = e.writeByte(byte(s.Curve)); err != nil {
		return
	}

	return e.toWriter(s.Content)
}

func (e *Encoder) writeUint16(i uint16) (err error) {
	buf := make([]byte, TypeSize.Uint16)
	binary.LittleEndian.PutUint16(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeInt16(i int16) (err error) {
	return e.writeUint16(uint16(i))
}

func (e *Encoder) writeInt32(i int32) (err error) {
	return e.writeUint32(uint32(i))
}

func (e *Encoder) writeUint32(i uint32) (err error) {
	buf := make([]byte, TypeSize.Uint32)
	binary.LittleEndian.PutUint32(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeUVarInt32(v uint32) (err error) {
	buf := make([]byte, binary.MaxVarintLen32)
	l := binary.PutUvarint(buf, uint64(v))
	return e.toWriter(buf[:l])
}

func (e *Encoder) writeInt64(i int64) (err error) {
	return e.writeUint64(uint64(i))
}

func (e *Encoder) writeUint64(i uint64) (err error) {
	buf := make([]byte, TypeSize.Uint64)
	binary.LittleEndian.PutUint64(buf, i)
	return e.toWriter(buf)
}

func (e *Encoder) writeJSONTime(tm JSONTime) (err error) {
	return e.writeUint32(uint32(tm.Unix()))
}

func (e *Encoder) writeFloat32(f float32) (err error) {
	i := math.Float32bits(f)
	buf := make([]byte, TypeSize.Uint32)
	binary.LittleEndian.PutUint32(buf, i)

	return e.toWriter(buf)
}
func (e *Encoder) writeFloat64(f float64) (err error) {
	i := math.Float64bits(f)
	buf := make([]byte, TypeSize.Uint64)
	binary.LittleEndian.PutUint64(buf, i)

	return e.toWriter(buf)
}

func (e *Encoder) writeAsset(asset Asset) (err error) {
	e.writeUint64(uint64(asset.Amount))
	e.writeByte(asset.Precision)

	symbol := make([]byte, 7, 7)

	copy(symbol[:], []byte(asset.Symbol.Symbol))
	return e.toWriter(symbol)
}

func (e *Encoder) writeChecksum256(checksum Checksum256) error {
	if len(checksum) == 0 {
		return e.toWriter(bytes.Repeat([]byte{0}, TypeSize.Checksum256))
	}
	return e.toWriter(checksum)
}

func (e *Encoder) writePublicKey(pk ecc.PublicKey) (err error) {
	err = pk.Validate()
	if err != nil {
		return fmt.Errorf("invalid public key: %w", err)
	}

	if err = e.writeByte(byte(pk.Curve)); err != nil {
		return err
	}

	return e.toWriter(pk.Content)
}

func (e *Encoder) writeActionData(actionData ActionData) (err error) {
	if actionData.Data != nil {
		//if reflect.TypeOf(actionData.Data) == reflect.TypeOf(&ActionData{}) {
		//	log.Fatal("pas cool")
		//}
		var d interface{}
		d = actionData.Data
		if reflect.TypeOf(d).Kind() == reflect.Ptr {
			d = reflect.ValueOf(actionData.Data).Elem().Interface()
		}

		if reflect.TypeOf(d).Kind() == reflect.String { //todo : this is a very bad ack ......
			data, err := hex.DecodeString(d.(string))
			if err != nil {
				return fmt.Errorf("ack, %s", err)
			}
			e.writeByteArray(data)
			return nil
		}

		raw, err := MarshalBinary(d)
		if err != nil {
			return err
		}
		return e.writeByteArray(raw)
	}

	return e.writeByteArray(actionData.HexData)
}

func (e *Encoder) Encode(v interface{}) (err error) {
	switch cv := v.(type) {
	case bool:
		return e.writeBool(cv)
	case Checksum256:
		return e.writeChecksum256(cv)
	case Varuint32:
		return e.writeUVarInt32(uint32(cv))
	case MarshalerBinary:
		return cv.MarshalBinary(e)
	case string:
		return e.writeString(cv)
	case CompressionType:
		return e.writeByte(byte(cv))
	case ecc.Signature:
		return e.writeSignature(cv)
	case ecc.PublicKey:
		return e.writePublicKey(cv)
	case Name:
		return e.writeName(cv)
	case AccountName:
		name := Name(cv)
		return e.writeName(name)
	case PermissionName:
		name := Name(cv)
		return e.writeName(name)
	case ActionName:
		name := Name(cv)
		return e.writeName(name)
	case Symbol:
		value, err := cv.ToUint64()
		if err != nil {
			return fmt.Errorf("encoding symbol: %w", err)
		}
		return e.writeUint64(value)
	case SymbolCode:
		return e.writeUint64(uint64(cv))
	case Asset:
		return e.writeAsset(cv)
	case byte:
		return e.writeByte(cv)
	case int8:
		return e.writeByte(byte(cv))
	case int16:
		return e.writeInt16(cv)
	case uint16:
		return e.writeUint16(cv)
	case int32:
		return e.writeInt32(cv)
	case uint32:
		return e.writeUint32(cv)
	case uint64:
		return e.writeUint64(cv)
	case int64:
		return e.writeInt64(cv)
	case float32:
		return e.writeFloat32(cv)
	case float64:
		return e.writeFloat64(cv)
	case JSONTime:
		return e.writeJSONTime(cv)
	case HexBytes:
		return e.writeByteArray(cv)
	case ActionData:
		return e.writeActionData(cv)
	case *ActionData:
		return e.writeActionData(*cv)

	case nil:
	default:
		rv := reflect.Indirect(reflect.ValueOf(v))
		t := rv.Type()

		switch t.Kind() {
		case reflect.Array:
			l := t.Len()
			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
		case reflect.Slice:
			l := rv.Len()
			if err = e.writeUVarInt(l); err != nil {
				return
			}
			for i := 0; i < l; i++ {
				if err = e.Encode(rv.Index(i).Interface()); err != nil {
					return
				}
			}
		case reflect.Struct:
			l := rv.NumField()

			for i := 0; i < l; i++ {
				field := t.Field(i)

				tag := field.Tag.Get("eos")
				if tag == "-" {
					continue
				}

				if v := rv.Field(i); t.Field(i).Name != "_" {
					if v.CanInterface() {
						isPresent := true
						if tag == "optional" {
							isPresent = !v.IsZero()
							e.writeBool(isPresent)
						}

						if isPresent {
							if err = e.Encode(v.Interface()); err != nil {
								return
							}
						}
					} else {
					}
				}
			}

		case reflect.Map:
			keyCount := len(rv.MapKeys())
			if err = e.writeUVarInt(keyCount); err != nil {
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
			return errors.New("Encode: unsupported type " + t.String())
		}
	}

	return
}
