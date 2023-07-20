package bcs

import (
	"bytes"
	"errors"
	"sort"

	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
)

// `serializer` extends `serde.BinarySerializer` to implement `serde.Serializer`.
type serializer struct {
	serde.BinarySerializer
}

func NewSerializer() serde.Serializer {
	return &serializer{*serde.NewBinarySerializer(MaxContainerDepth)}
}

// SerializeF32 is unimplemented
func (s *serializer) SerializeF32(value float32) error {
	return errors.New("unimplemented")
}

// SerializeF64 is unimplemented
func (s *serializer) SerializeF64(value float64) error {
	return errors.New("unimplemented")
}

func (s *serializer) SerializeStr(value string) error {
	return s.BinarySerializer.SerializeStr(value, s.SerializeLen)
}

func (s *serializer) SerializeBytes(value []byte) error {
	return s.BinarySerializer.SerializeBytes(value, s.SerializeLen)
}

func (s *serializer) SerializeFixedBytes(value []byte) error {
	return s.BinarySerializer.SerializeFixedBytes(value)
}

func (s *serializer) SerializeLen(value uint64) error {
	if value > MaxSequenceLength {
		return errors.New("length is too large")
	}
	s.serializeU32AsUleb128(uint32(value))
	return nil
}

func (s *serializer) SerializeVariantIndex(value uint32) error {
	s.serializeU32AsUleb128(value)
	return nil
}

func (s *serializer) SortMapEntries(offsets []uint64) {
	if len(offsets) <= 1 {
		return
	}
	data := s.Buffer.Bytes()
	slices := make([]serde.Slice, len(offsets))
	for i, v := range offsets {
		var w uint64
		if i+1 < len(offsets) {
			w = offsets[i+1]
		} else {
			w = uint64(len(data))
		}
		slices[i] = serde.Slice{Start: v, End: w}
	}
	entries := map_entries{data, slices}
	sort.Sort(entries)
	buffer := make([]byte, len(data)-int(offsets[0]))
	current := buffer[0:0]
	for _, slice := range entries.slices {
		current = append(current, data[slice.Start:slice.End]...)
	}
	copy(data[offsets[0]:], current)
}

func (s *serializer) serializeU32AsUleb128(value uint32) {
	for value >= 0x80 {
		b := byte((value & 0x7f) | 0x80)
		_ = s.Buffer.WriteByte(b)
		value = value >> 7
	}
	_ = s.Buffer.WriteByte(byte(value))
}

type map_entries struct {
	data   []byte
	slices []serde.Slice
}

func (a map_entries) Len() int { return len(a.slices) }

func (a map_entries) Less(i, j int) bool {
	slice_i := a.data[a.slices[i].Start:a.slices[i].End]
	slice_j := a.data[a.slices[j].Start:a.slices[j].End]
	return bytes.Compare(slice_i, slice_j) < 0
}

func (a map_entries) Swap(i, j int) { a.slices[i], a.slices[j] = a.slices[j], a.slices[i] }
