package v2

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/util"
	"math/big"
)

// Script A Move script as compiled code as a transaction
type Script struct {
	Code     []byte           // The compiled script bytes
	ArgTypes []TypeTag        // The types of the arguments
	Args     []ScriptArgument // The arguments
}

//region Script TransactionPayloadImpl

func (s *Script) PayloadType() TransactionPayloadVariant {
	return TransactionPayloadVariantScript
}

func (s *Script) ExecutableType() TransactionExecutableVariant {
	return TransactionExecutableVariantScript
}

//endregion

//region Script bcs.Struct

func (s *Script) MarshalBCS(ser *bcs.Serializer) {
	ser.WriteBytes(s.Code)
	bcs.SerializeSequence(s.ArgTypes, ser)
	bcs.SerializeSequence(s.Args, ser)
}

func (s *Script) UnmarshalBCS(des *bcs.Deserializer) {
	s.Code = des.ReadBytes()
	s.ArgTypes = bcs.DeserializeSequence[TypeTag](des)
	s.Args = bcs.DeserializeSequence[ScriptArgument](des)
}

//endregion
//endregion

//region ScriptArgument

// ScriptArgumentVariant the type of the script argument.  If there isn't a value here, it is not supported.
//
// Note that the only vector supported is vector<u8>
type ScriptArgumentVariant uint32

const (
	ScriptArgumentU8         ScriptArgumentVariant = 0 // u8 type argument
	ScriptArgumentU64        ScriptArgumentVariant = 1 // u64 type argument
	ScriptArgumentU128       ScriptArgumentVariant = 2 // u128 type argument
	ScriptArgumentAddress    ScriptArgumentVariant = 3 // address type argument
	ScriptArgumentU8Vector   ScriptArgumentVariant = 4 // vector<u8> type argument
	ScriptArgumentBool       ScriptArgumentVariant = 5 // bool type argument
	ScriptArgumentU16        ScriptArgumentVariant = 6 // u16 type argument
	ScriptArgumentU32        ScriptArgumentVariant = 7 //	u32 type argument
	ScriptArgumentU256       ScriptArgumentVariant = 8 //	u256 type argument
	ScriptArgumentSerialized ScriptArgumentVariant = 9 //	u256 type argument
)

// ScriptArgument a Move script argument, which encodes its type with it
type ScriptArgument struct {
	Variant ScriptArgumentVariant // The type of the argument
	Value   any                   // The value of the argument
}

//region ScriptArgument bcs.Struct
// TODO: consider making a separate function to parse the value at input time rather than build time

func (sa *ScriptArgument) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(sa.Variant))
	switch sa.Variant {
	case ScriptArgumentU8:
		value, ok := (sa.Value).(uint8)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentU8, must be uint8", sa.Value))
		}
		ser.U8(value)
	case ScriptArgumentU16:
		value, ok := (sa.Value).(uint16)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentU16, must be uint16", sa.Value))
		}
		ser.U16(value)
	case ScriptArgumentU32:
		value, ok := (sa.Value).(uint32)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentU32, must be uint32", sa.Value))
		}
		ser.U32(value)
	case ScriptArgumentU64:
		value, ok := (sa.Value).(uint64)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentU64, must be uint64", sa.Value))
		}
		ser.U64(value)
	case ScriptArgumentU128:
		value, ok := (sa.Value).(big.Int)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgument128, must be big.Int", sa.Value))
		}
		ser.U128(value)
	case ScriptArgumentU256:
		value, ok := (sa.Value).(big.Int)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgument256, must be big.Int", sa.Value))
		}
		ser.U256(value)
	case ScriptArgumentAddress:
		addr, ok := (sa.Value).(AccountAddress)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentAddress, must be AccountAddress", sa.Value))
		}
		ser.Struct(&addr)
	case ScriptArgumentU8Vector:
		bytes, ok := (sa.Value).([]byte)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentU8Vector, must be []byte", sa.Value))
		}
		ser.WriteBytes(bytes)
	case ScriptArgumentBool:
		value, ok := (sa.Value).(bool)
		if !ok {
			ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentBool, must be bool", sa.Value))
		}
		ser.Bool(value)
	case ScriptArgumentSerialized:
		value, ok := (sa.Value).([]byte)
		if ok {
			ser.WriteBytes(value)
		} else {
			valueStr, ok := (sa.Value).(string)
			if !ok {
				ser.SetError(fmt.Errorf("invalid input type (%T) for ScriptArgumentSerialized, must be string or bytearray", sa.Value))
			} else {
				ser.WriteBytes(util.RemoveZeroHex(valueStr))
			}
		}
	}
}

func (sa *ScriptArgument) UnmarshalBCS(des *bcs.Deserializer) {
	sa.Variant = ScriptArgumentVariant(des.Uleb128())
	switch sa.Variant {
	case ScriptArgumentU8:
		sa.Value = des.U8()
	case ScriptArgumentU16:
		sa.Value = des.U16()
	case ScriptArgumentU32:
		sa.Value = des.U32()
	case ScriptArgumentU64:
		sa.Value = des.U64()
	case ScriptArgumentU128:
		sa.Value = des.U128()
	case ScriptArgumentU256:
		sa.Value = des.U256()
	case ScriptArgumentAddress:
		aa := AccountAddress{}
		aa.UnmarshalBCS(des)
		sa.Value = aa
	case ScriptArgumentU8Vector:
		sa.Value = des.ReadBytes()
	case ScriptArgumentBool:
		sa.Value = des.Bool()
	case ScriptArgumentSerialized:
		sa.Value = des.ReadBytes()
	}
}
