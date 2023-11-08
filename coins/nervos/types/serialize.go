package types

import (
	"bytes"
	"encoding/binary"
	"errors"
	"reflect"
)

const u32Size uint = 4

type Serializer interface {
	Serialize() ([]byte, error)
}

func (w *WitnessArgs) Serialize() ([]byte, error) {
	l, err := SerializeOptionBytes(w.Lock)
	if err != nil {
		return nil, err
	}

	i, err := SerializeOptionBytes(w.InputType)
	if err != nil {
		return nil, err
	}

	o, err := SerializeOptionBytes(w.OutputType)
	if err != nil {
		return nil, err
	}

	return SerializeTable([][]byte{l, i, o}), nil
}

func SerializeUint(n uint) []byte {
	b := make([]byte, u32Size)
	binary.LittleEndian.PutUint32(b, uint32(n))

	return b
}

func SerializeUint64(n uint64) []byte {
	b := make([]byte, 8)
	binary.LittleEndian.PutUint64(b, n)

	return b
}

// SerializeArray serialize array
func SerializeArray(items []Serializer) ([][]byte, error) {
	ret := make([][]byte, len(items))
	for i := 0; i < len(items); i++ {
		data, err := items[i].Serialize()
		if err != nil {
			return nil, err
		}

		ret[i] = data
	}

	return ret, nil
}

// SerializeStruct serialize struct
func SerializeStruct(fields [][]byte) []byte {
	b := new(bytes.Buffer)

	for i := 0; i < len(fields); i++ {
		b.Write(fields[i])
	}

	return b.Bytes()
}

// SerializeOptionBytes serialize option
func SerializeOptionBytes(o []byte) ([]byte, error) {
	if o == nil || reflect.ValueOf(o).IsNil() {
		return []byte{}, nil
	}

	return SerializeBytes(o), nil
}

// SerializeBytes serialize bytes
// There are two steps of serializing a bytes:
//   Serialize the length as a 32 bit unsigned integer in little-endian.
//   Serialize all items in it.
func SerializeBytes(items []byte) []byte {
	// Empty fix vector bytes
	if len(items) == 0 {
		return []byte{00, 00, 00, 00}
	}

	l := SerializeUint(uint(len(items)))

	b := new(bytes.Buffer)

	b.Write(l)
	b.Write(items)

	return b.Bytes()
}

// SerializeTable serialize table
// The serializing steps are same as table:
//    Serialize the full size in bytes as a 32 bit unsigned integer in little-endian.
//    Serialize all offset of fields as 32 bit unsigned integer in little-endian.
//    Serialize all fields in it in the order they are declared.
func SerializeTable(fields [][]byte) []byte {
	size := u32Size
	offsets := make([]uint, len(fields))

	// Calculate first offset then loop for rest items offsets
	offsets[0] = u32Size + u32Size*uint(len(fields))
	for i := 0; i < len(fields); i++ {
		size += u32Size + uint(len(fields[i]))

		if i != 0 {
			offsets[i] = offsets[i-1] + uint(len(fields[i-1]))
		}
	}

	b := new(bytes.Buffer)

	b.Write(SerializeUint(size))

	for i := 0; i < len(fields); i++ {
		b.Write(SerializeUint(offsets[i]))
	}

	for i := 0; i < len(fields); i++ {
		b.Write(fields[i])
	}

	return b.Bytes()
}

// Serialize dep type
func (t DepType) Serialize() ([]byte, error) {
	if t == DepTypeCode {
		return []byte{00}, nil
	} else if t == DepTypeDepGroup {
		return []byte{01}, nil
	}
	return nil, errors.New("invalid dep group")
}

// Serialize cell dep
func (d *CellDep) Serialize() ([]byte, error) {
	o, err := d.OutPoint.Serialize()
	if err != nil {
		return nil, err
	}

	dd, err := d.DepType.Serialize()
	if err != nil {
		return nil, err
	}

	return SerializeStruct([][]byte{o, dd}), nil
}

// Serialize cell input
func (i *CellInput) Serialize() ([]byte, error) {
	s := SerializeUint64(i.Since)

	o, err := i.PreviousOutput.Serialize()
	if err != nil {
		return nil, err
	}

	return SerializeStruct([][]byte{s, o}), nil
}

func (h Hash) Serialize() ([]byte, error) {
	return h.Bytes(), nil
}

// Serialize outpoint
func (o *OutPoint) Serialize() ([]byte, error) {
	h, err := o.TxHash.Serialize()
	if err != nil {
		return nil, err
	}

	i := SerializeUint(o.Index)

	b := new(bytes.Buffer)

	b.Write(h)
	b.Write(i)

	return b.Bytes(), nil
}

func (t ScriptHashType) Serialize() ([]byte, error) {
	if t == HashTypeData {
		return []byte{00}, nil
	} else if t == HashTypeType {
		return []byte{01}, nil
	} else if t == HashTypeData1 {
		return []byte{02}, nil
	}
	return nil, errors.New("invalid script hash type")
}

// Serialize script
func (script *Script) Serialize() ([]byte, error) {
	h, err := script.CodeHash.Serialize()
	if err != nil {
		return nil, err
	}

	t, err := script.HashType.Serialize()
	if err != nil {
		return nil, err
	}

	a := SerializeBytes(script.Args)

	return SerializeTable([][]byte{h, t, a}), nil
}

// Serialize cell output
func (o *CellOutput) Serialize() ([]byte, error) {
	c := SerializeUint64(o.Capacity)

	l, err := o.Lock.Serialize()
	if err != nil {
		return nil, err
	}

	t, err := SerializeOption(o.Type)
	if err != nil {
		return nil, err
	}

	return SerializeTable([][]byte{c, l, t}), nil
}

// SerializeOption serialize option
func SerializeOption(o Serializer) ([]byte, error) {
	if o == nil || reflect.ValueOf(o).IsNil() {
		return []byte{}, nil
	}

	return o.Serialize()
}

// SerializeFixVec serialize fixvec vector
// There are two steps of serializing a fixvec:
//   Serialize the length as a 32 bit unsigned integer in little-endian.
//   Serialize all items in it.
func SerializeFixVec(items [][]byte) []byte {
	// Empty fix vector bytes
	if len(items) == 0 {
		return []byte{00, 00, 00, 00}
	}

	l := SerializeUint(uint(len(items)))

	b := new(bytes.Buffer)

	b.Write(l)

	for i := 0; i < len(items); i++ {
		b.Write(items[i])
	}

	return b.Bytes()
}

// SerializeDynVec serialize dynvec
// There are three steps of serializing a dynvec:
//    Serialize the full size in bytes as a 32 bit unsigned integer in little-endian.
//    Serialize all offset of items as 32 bit unsigned integer in little-endian.
//    Serialize all items in it.
func SerializeDynVec(items [][]byte) []byte {
	// Start with u32Size
	size := u32Size

	// Empty dyn vector, just return size's bytes
	if len(items) == 0 {
		return SerializeUint(size)
	}

	offsets := make([]uint, len(items))

	// Calculate first offset then loop for rest items offsets
	offsets[0] = size + u32Size*uint(len(items))
	for i := 0; i < len(items); i++ {
		size += u32Size + uint(len(items[i]))

		if i != 0 {
			offsets[i] = offsets[i-1] + uint(len(items[i-1]))
		}
	}

	b := new(bytes.Buffer)

	b.Write(SerializeUint(size))

	for i := 0; i < len(items); i++ {
		b.Write(SerializeUint(offsets[i]))
	}

	for i := 0; i < len(items); i++ {
		b.Write(items[i])
	}

	return b.Bytes()
}
