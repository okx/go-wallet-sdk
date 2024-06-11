package transaction_builder

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/aptos_types"
	"github.com/okx/go-wallet-sdk/coins/aptos/serde"
	"math/big"
	"reflect"
	"strconv"
)

func serializeArg(argVal any, argType aptos_types.TypeTag, serializer serde.Serializer) error {
	return serializeArgInner(argVal, argType, serializer, 0)
}

func serializeArgInner(argVal any, argType aptos_types.TypeTag, serializer serde.Serializer, depth int) error {
	switch argType.(type) {
	case *aptos_types.TypeTagBool:
		if v, ok := argVal.(bool); ok {
			return serializer.SerializeBool(v)
		}
		if v, ok := argVal.(string); ok {
			if v == "false" {
				return serializer.SerializeBool(false)
			}
			if v == "true" {
				return serializer.SerializeBool(true)
			}
		}
	case *aptos_types.TypeTagU8:
		if v, ok := argVal.(uint8); ok {
			return serializer.SerializeU8(v)
		}
		if v, ok := argVal.(int); ok && v == int(uint8(v)) {
			return serializer.SerializeU8(uint8(v))
		}
		if v, ok := argVal.(float64); ok && v == float64(uint8(v)) {
			return serializer.SerializeU8(uint8(v))
		}
		if v, ok := argVal.(string); ok {
			u, err := strconv.ParseUint(v, 10, 8)
			if err != nil {
				return err
			}
			return serializer.SerializeU8(uint8(u))
		}
	case *aptos_types.TypeTagU16:
		if v, ok := argVal.(uint16); ok {
			return serializer.SerializeU16(v)
		}
		if v, ok := argVal.(int); ok && v == int(uint16(v)) {
			return serializer.SerializeU16(uint16(v))
		}
		if v, ok := argVal.(float64); ok && v == float64(uint16(v)) {
			return serializer.SerializeU16(uint16(v))
		}
		if v, ok := argVal.(string); ok {
			u, err := strconv.ParseUint(v, 10, 16)
			if err != nil {
				return err
			}
			return serializer.SerializeU16(uint16(u))
		}
	case *aptos_types.TypeTagU32:
		if v, ok := argVal.(uint32); ok {
			return serializer.SerializeU32(v)
		}
		if v, ok := argVal.(int); ok && v == int(uint32(v)) {
			return serializer.SerializeU32(uint32(v))
		}
		if v, ok := argVal.(float64); ok && v == float64(uint32(v)) {
			return serializer.SerializeU32(uint32(v))
		}
		if v, ok := argVal.(string); ok {
			u, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				return err
			}
			return serializer.SerializeU32(uint32(u))
		}
	case *aptos_types.TypeTagU64:
		if v, ok := argVal.(uint64); ok {
			return serializer.SerializeU64(v)
		}
		if v, ok := argVal.(int); ok && v >= 0 {
			return serializer.SerializeU64(uint64(v))
		}
		if v, ok := argVal.(float64); ok && v >= 0 {
			return serializer.SerializeU64(uint64(v))
		}
		if v, ok := argVal.(string); ok {
			u, err := strconv.ParseUint(v, 10, 64)
			if err != nil {
				return err
			}
			return serializer.SerializeU64(u)
		}
	case *aptos_types.TypeTagU128:
		if v, ok := argVal.(serde.Uint128); ok {
			return serializer.SerializeU128(v)
		}
		if v, ok := argVal.(*big.Int); ok {
			vv, err := serde.FromBig(v)
			if err != nil {
				return err
			}
			return serializer.SerializeU128(*vv)
		}
		if v, ok := argVal.(int); ok && v >= 0 {
			vv, err := serde.FromBig(big.NewInt(int64(v)))
			if err != nil {
				return err
			}
			return serializer.SerializeU128(*vv)
		}
		if v, ok := argVal.(float64); ok && v >= 0 {
			vv, err := serde.FromBig(big.NewInt(int64(v)))
			if err != nil {
				return err
			}
			return serializer.SerializeU128(*vv)
		}
		if v, ok := argVal.(string); ok {
			if bigV, ok := big.NewInt(0).SetString(v, 10); ok {
				vv, err := serde.FromBig(bigV)
				if err != nil {
					return err
				}
				return serializer.SerializeU128(*vv)
			}
		}

	case *aptos_types.TypeTagU256:
		if v, ok := argVal.(serde.Uint256); ok {
			return serializer.SerializeU256(v)
		}
		if v, ok := argVal.(*big.Int); ok {
			vv, err := serde.BigInt2U256(v)
			if err != nil {
				return err
			}
			return serializer.SerializeU256(*vv)
		}
		if v, ok := argVal.(int); ok && v >= 0 {
			vv, err := serde.BigInt2U256(big.NewInt(int64(v)))
			if err != nil {
				return err
			}
			return serializer.SerializeU256(*vv)
		}
		if v, ok := argVal.(float64); ok && v >= 0 {
			vv, err := serde.BigInt2U256(big.NewInt(int64(v)))
			if err != nil {
				return err
			}
			return serializer.SerializeU256(*vv)
		}
		if v, ok := argVal.(string); ok {
			if bigV, ok := big.NewInt(0).SetString(v, 10); ok {
				vv, err := serde.BigInt2U256(bigV)
				if err != nil {
					return err
				}
				return serializer.SerializeU256(*vv)
			}
		}
	case *aptos_types.TypeTagAddress:
		if v, ok := argVal.(aptos_types.AccountAddress); ok {
			return serializer.SerializeFixedBytes(v[:])
		}
		if v, ok := argVal.(string); ok {
			addr, err := aptos_types.FromHex(v)
			if err != nil {
				return err
			}
			return serializer.SerializeFixedBytes(addr[:])
		}
	case *aptos_types.TypeTagVector:
		argTp, ok := argType.(*aptos_types.TypeTagVector)
		if !ok {
			return fmt.Errorf("invalid vector %v", argType)
		}
		return serializeVector(argVal, argTp, serializer, depth)
	case *aptos_types.TypeTagStruct:
		tag := argType.(*aptos_types.TypeTagStruct)
		structType := tag.ShortFunctionName()
		if structType == "0x1::string::String" {
			v, ok := argVal.(string)
			if !ok {
				return fmt.Errorf("invalid structType %v", argVal)
			}
			return serializer.SerializeStr(v)
		} else if structType == "0x1::object::Object" {
			if v, ok := argVal.(aptos_types.AccountAddress); ok {
				return serializer.SerializeFixedBytes(v[:])
			}
			if v, ok := argVal.(string); ok {
				addr, err := aptos_types.FromHex(v)
				if err != nil {
					return err
				}
				return serializer.SerializeFixedBytes(addr[:])
			}
		} else if structType == "0x1::option::Option" {
			if len(tag.Value.TypeArgs) != 1 {
				return fmt.Errorf("option has the wrong number of type arguments %d", len(tag.Value.TypeArgs))
			}
			return serializeOption(argVal, tag.Value.TypeArgs[0], serializer, depth)
		}
		return errors.New("unsupported struct type in function argument")
	default:
		return fmt.Errorf("unsupported arg type %v", argType)
	}
	return fmt.Errorf("invalid argument %v", argVal)
}

func serializeOption(argVal any, argType aptos_types.TypeTag, serializer serde.Serializer, depth int) error {
	rv := reflect.ValueOf(argVal)
	if argVal == nil || rv.IsNil() {
		err := serializer.SerializeVariantIndex(0)
		if err != nil {
			return err
		}
	} else {
		err := serializer.SerializeVariantIndex(1)
		if err != nil {
			return err
		}
		return serializeArgInner(argVal, argType, serializer, depth+1)
	}
	return nil
}

func serializeVector(argVal any, argType *aptos_types.TypeTagVector, serializer serde.Serializer, depth int) error {
	itemType := argType.Value
	switch itemType.(type) {
	case *aptos_types.TypeTagU8:
		if v, ok := argVal.([]byte); ok {
			return serializer.SerializeFixedBytes(v)
		}

		if v, ok := argVal.(string); ok {
			if aptos_types.IsHexString(v) {
				vb, err := hex.DecodeString(v[2:])
				if err != nil {
					return err
				}
				return serializer.SerializeBytes(vb)
			}
			return serializer.SerializeStr(v)
		}
	}
	rv := reflect.ValueOf(argVal)
	if rv.Kind() != reflect.Array && rv.Kind() != reflect.Slice {
		return errors.New("invalid vector args")
	}
	length := rv.Len()
	if err := serializer.SerializeVariantIndex(uint32(length)); err != nil {
		return err
	}
	for i := 0; i < length; i++ {
		if err := serializeArgInner(rv.Index(i).Interface(), itemType, serializer, depth+1); err != nil {
			return err
		}
	}
	return nil
}
