package serde

import (
	"bytes"
	"errors"
)

// BinarySerializer is a partial implementation of the `Serializer` interface.
// It is used as an embedded struct by the Bincode and BCS serializers.
type BinarySerializer struct {
	Buffer               bytes.Buffer
	containerDepthBudget uint64
}

func NewBinarySerializer(max_container_depth uint64) *BinarySerializer {
	s := new(BinarySerializer)
	s.containerDepthBudget = max_container_depth
	return s
}

func (s *BinarySerializer) IncreaseContainerDepth() error {
	if s.containerDepthBudget == 0 {
		return errors.New("exceeded maximum container depth")
	}
	s.containerDepthBudget -= 1
	return nil
}

func (s *BinarySerializer) DecreaseContainerDepth() {
	s.containerDepthBudget += 1
}

// SerializeBytes `serializeLen` to be provided by the extending struct.
func (s *BinarySerializer) SerializeBytes(value []byte, serializeLen func(uint64) error) error {
	err := serializeLen(uint64(len(value)))
	if err != nil {
		return err
	}
	s.Buffer.Write(value)
	return nil
}

func (s *BinarySerializer) SerializeFixedBytes(value []byte) error {
	s.Buffer.Write(value)
	return nil
}

// SerializeStr `serializeLen` to be provided by the extending struct.
func (s *BinarySerializer) SerializeStr(value string, serializeLen func(uint64) error) error {
	return s.SerializeBytes([]byte(value), serializeLen)
}

func (s *BinarySerializer) SerializeBool(value bool) error {
	if value {
		return s.Buffer.WriteByte(1)
	}
	return s.Buffer.WriteByte(0)
}

func (s *BinarySerializer) SerializeUnit(value struct{}) error {
	return nil
}

// SerializeChar is unimplemented.
func (s *BinarySerializer) SerializeChar(value rune) error {
	return errors.New("unimplemented")
}

func (s *BinarySerializer) SerializeU8(value uint8) error {
	s.Buffer.WriteByte(value)
	return nil
}

func (s *BinarySerializer) SerializeU16(value uint16) error {
	s.Buffer.WriteByte(byte(value))
	s.Buffer.WriteByte(byte(value >> 8))
	return nil
}

func (s *BinarySerializer) SerializeU32(value uint32) error {
	s.Buffer.WriteByte(byte(value))
	s.Buffer.WriteByte(byte(value >> 8))
	s.Buffer.WriteByte(byte(value >> 16))
	s.Buffer.WriteByte(byte(value >> 24))
	return nil
}

func (s *BinarySerializer) SerializeU64(value uint64) error {
	s.Buffer.WriteByte(byte(value))
	s.Buffer.WriteByte(byte(value >> 8))
	s.Buffer.WriteByte(byte(value >> 16))
	s.Buffer.WriteByte(byte(value >> 24))
	s.Buffer.WriteByte(byte(value >> 32))
	s.Buffer.WriteByte(byte(value >> 40))
	s.Buffer.WriteByte(byte(value >> 48))
	s.Buffer.WriteByte(byte(value >> 56))
	return nil
}

func (s *BinarySerializer) SerializeU128(value Uint128) error {
	err := s.SerializeU64(value.Low)
	if err != nil {
		return err
	}
	err = s.SerializeU64(value.High)
	if err != nil {
		return err
	}
	return nil
}

func (s *BinarySerializer) SerializeU256(value Uint256) error {
	err := s.SerializeU128(value.Low)
	if err != nil {
		return err
	}
	err = s.SerializeU128(value.High)
	if err != nil {
		return err
	}
	return nil
}

func (s *BinarySerializer) SerializeI8(value int8) error {
	err := s.SerializeU8(uint8(value))
	if err != nil {
		return err
	}
	return nil
}

func (s *BinarySerializer) SerializeI16(value int16) error {
	err := s.SerializeU16(uint16(value))
	if err != nil {
		return err
	}
	return nil
}

func (s *BinarySerializer) SerializeI32(value int32) error {
	err := s.SerializeU32(uint32(value))
	if err != nil {
		return err
	}
	return nil
}

func (s *BinarySerializer) SerializeI64(value int64) error {
	err := s.SerializeU64(uint64(value))
	if err != nil {
		return err
	}
	return nil
}

func (s *BinarySerializer) SerializeI128(value Int128) error {
	err := s.SerializeU64(value.Low)
	if err != nil {
		return err
	}
	err = s.SerializeI64(value.High)
	if err != nil {
		return err
	}
	return nil
}

func (s *BinarySerializer) SerializeOptionTag(value bool) error {
	return s.SerializeBool(value)
}

func (s *BinarySerializer) GetBufferOffset() uint64 {
	return uint64(s.Buffer.Len())
}

func (s *BinarySerializer) GetBytes() []byte {
	return s.Buffer.Bytes()
}
