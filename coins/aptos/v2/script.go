package v2

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"math/big"
)

// Script A Move script as compiled code as a transaction
type Script struct {
	Code     []byte
	ArgTypes []TypeTag
	Args     []ScriptArgument
}

func (sc *Script) MarshalBCS(serializer *bcs.Serializer) {
	serializer.WriteBytes(sc.Code)
	bcs.SerializeSequence(sc.ArgTypes, serializer)
	bcs.SerializeSequence(sc.Args, serializer)
}

func (sc *Script) UnmarshalBCS(deserializer *bcs.Deserializer) {
	sc.Code = deserializer.ReadBytes()
	sc.ArgTypes = bcs.DeserializeSequence[TypeTag](deserializer)
	sc.Args = bcs.DeserializeSequence[ScriptArgument](deserializer)
}

// ScriptArgument a Move script argument, which encodes its type with it
// TODO: improve typing
type ScriptArgument struct {
	Variant ScriptArgumentVariant
	Value   any
}

type ScriptArgumentVariant uint8

const (
	ScriptArgumentU8       ScriptArgumentVariant = 0
	ScriptArgumentU64      ScriptArgumentVariant = 1
	ScriptArgumentU128     ScriptArgumentVariant = 2
	ScriptArgumentAddress  ScriptArgumentVariant = 3
	ScriptArgumentU8Vector ScriptArgumentVariant = 4
	ScriptArgumentBool     ScriptArgumentVariant = 5
	ScriptArgumentU16      ScriptArgumentVariant = 6
	ScriptArgumentU32      ScriptArgumentVariant = 7
	ScriptArgumentU256     ScriptArgumentVariant = 8
)

func (sa *ScriptArgument) MarshalBCS(bcs *bcs.Serializer) {
	bcs.U8(uint8(sa.Variant))
	switch sa.Variant {
	case ScriptArgumentU8:
		bcs.U8(sa.Value.(uint8))
	case ScriptArgumentU16:
		bcs.U16(sa.Value.(uint16))
	case ScriptArgumentU32:
		bcs.U32(sa.Value.(uint32))
	case ScriptArgumentU64:
		bcs.U64(sa.Value.(uint64))
	case ScriptArgumentU128:
		bcs.U128(sa.Value.(big.Int))
	case ScriptArgumentU256:
		bcs.U256(sa.Value.(big.Int))
	case ScriptArgumentAddress:
		addr := sa.Value.(AccountAddress)
		bcs.Struct(&addr)
	case ScriptArgumentU8Vector:
		bcs.WriteBytes(sa.Value.([]byte))
	case ScriptArgumentBool:
		bcs.Bool(sa.Value.(bool))
	}
}

func (sa *ScriptArgument) UnmarshalBCS(bcs *bcs.Deserializer) {
	variant := bcs.U8()
	switch ScriptArgumentVariant(variant) {
	case ScriptArgumentU8:
		sa.Value = bcs.U8()
	case ScriptArgumentU16:
		sa.Value = bcs.U16()
	case ScriptArgumentU32:
		sa.Value = bcs.U32()
	case ScriptArgumentU64:
		sa.Value = bcs.U64()
	case ScriptArgumentU128:
		sa.Value = bcs.U128()
	case ScriptArgumentU256:
		sa.Value = bcs.U256()
	case ScriptArgumentAddress:
		aa := AccountAddress{}
		aa.UnmarshalBCS(bcs)
		sa.Value = aa
	case ScriptArgumentU8Vector:
		sa.Value = bcs.ReadBytes()
	case ScriptArgumentBool:
		sa.Value = bcs.Bool()
	}
}
