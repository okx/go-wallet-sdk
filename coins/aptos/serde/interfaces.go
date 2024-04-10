package serde

type Serializer interface {
	SerializeStr(value string) error

	SerializeBytes(value []byte) error

	SerializeFixedBytes(value []byte) error

	SerializeBool(value bool) error

	SerializeUnit(value struct{}) error

	SerializeChar(value rune) error

	SerializeF32(value float32) error

	SerializeF64(value float64) error

	SerializeU8(value uint8) error

	SerializeU16(value uint16) error

	SerializeU32(value uint32) error

	SerializeU64(value uint64) error

	SerializeU128(value Uint128) error

	SerializeU256(value Uint256) error

	SerializeI8(value int8) error

	SerializeI16(value int16) error

	SerializeI32(value int32) error

	SerializeI64(value int64) error

	SerializeI128(value Int128) error

	SerializeLen(value uint64) error

	SerializeVariantIndex(value uint32) error

	SerializeOptionTag(value bool) error

	GetBufferOffset() uint64

	SortMapEntries(offsets []uint64)

	GetBytes() []byte

	IncreaseContainerDepth() error

	DecreaseContainerDepth()
}

type Deserializer interface {
	DeserializeStr() (string, error)

	DeserializeBytes() ([]byte, error)

	DeserializeBool() (bool, error)

	DeserializeUnit() (struct{}, error)

	DeserializeChar() (rune, error)

	DeserializeF32() (float32, error)

	DeserializeF64() (float64, error)

	DeserializeU8() (uint8, error)

	DeserializeU16() (uint16, error)

	DeserializeU32() (uint32, error)

	DeserializeU64() (uint64, error)

	DeserializeU128() (Uint128, error)

	DeserializeI8() (int8, error)

	DeserializeI16() (int16, error)

	DeserializeI32() (int32, error)

	DeserializeI64() (int64, error)

	DeserializeI128() (Int128, error)

	DeserializeLen() (uint64, error)

	DeserializeVariantIndex() (uint32, error)

	DeserializeOptionTag() (bool, error)

	GetBufferOffset() uint64

	CheckThatKeySlicesAreIncreasing(key1, key2 Slice) error

	IncreaseContainerDepth() error

	DecreaseContainerDepth()
}

type Slice struct {
	Start uint64
	End   uint64
}
