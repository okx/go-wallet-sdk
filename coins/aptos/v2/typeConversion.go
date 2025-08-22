package v2

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/common"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/api"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
	"math/big"
	"strconv"
	"strings"
)

func EntryFunctionFromAbi(abi any, moduleAddress AccountAddress, moduleName string, functionName string, typeArgs []any, args []any, options ...any) (*EntryFunction, error) {
	var function *api.MoveFunction
	switch abi := abi.(type) {
	case *api.MoveModule:
		moduleAbi := abi
		// Find function
		for _, fun := range moduleAbi.ExposedFunctions {
			if fun.Name == functionName {
				if !fun.IsEntry {
					return nil, fmt.Errorf("function %s is not a entry function in module %s", functionName, moduleAbi.Name)
				}

				function = fun
				break
			}
		}
	case *api.MoveFunction:
		function = abi
	default:
		return nil, fmt.Errorf("unknown abi type: %T", abi)
	}

	if function == nil {
		return nil, fmt.Errorf("entry function %s not found in module %s", functionName, moduleName)
	}

	// Check type args length matches
	if len(typeArgs) != len(function.GenericTypeParams) {
		return nil, fmt.Errorf("entry function %s does not have the correct number of type arguments for function %s", functionName, functionName)
	}

	// Convert TypeTag, *TypeTag, and string to TypeTag
	// TODO: Check properties of generic type?
	convertedTypeArgs := make([]TypeTag, 0)
	for _, typeArg := range typeArgs {
		if typeArg == "" {
			continue
		}
		tag, err := ConvertTypeTag(typeArg)
		if err != nil {
			return nil, err
		}
		convertedTypeArgs = append(convertedTypeArgs, *tag)
	}

	// Convert string types to actual types
	argTypes := make([]TypeTag, 0)
	for _, typeStr := range function.Params {
		typeArg, err := ParseTypeTag(typeStr)
		if err != nil {
			return nil, err
		}

		// If it's `signer` or `&signer` need to skip
		// TODO: only skip at the beginning
		switch innerArg := typeArg.Value.(type) {
		case *SignerTag:
			// Skip
			continue
		case *ReferenceTag:
			switch innerArg.TypeParam.Value.(type) {
			case *SignerTag:
				// Skip
				continue
			default:
				argTypes = append(argTypes, *typeArg)
			}
		default:
			argTypes = append(argTypes, *typeArg)
		}
	}

	// Check args length matches
	if len(args) != len(argTypes) {
		return nil, fmt.Errorf("entry function %s does not have the correct number of arguments for function %s", functionName, functionName)
	}

	convertedArgs := make([][]byte, len(args))
	for i, arg := range args {
		b, err := ConvertArg(argTypes[i], arg, argTypes, options...)
		if err != nil {
			return nil, err
		}
		convertedArgs[i] = b
	}

	return &EntryFunction{
		Module: ModuleId{
			Address: moduleAddress,
			Name:    moduleName,
		},
		Function: functionName,
		ArgTypes: convertedTypeArgs,
		Args:     convertedArgs,
	}, nil
}

func ConvertTypeTag(typeArg any) (*TypeTag, error) {
	switch typeArg := typeArg.(type) {
	case TypeTag:
		return &typeArg, nil
	case *TypeTag:
		if typeArg == nil {
			return nil, fmt.Errorf("invalid type tag %s, cannot be nil", typeArg)
		}
		return typeArg, nil
	case string:
		return ParseTypeTag(typeArg)
	default:
		return nil, errors.New("invalid type tag type")
	}
}

func ConvertToU8(arg any) (uint8, error) {
	switch arg := arg.(type) {
	case float64:
		return util.Float64ToU8(arg)
	case int:
		return util.IntToU8(arg)
	case uint:
		return util.UintToU8(arg)
	case uint8:
		return arg, nil
	case big.Int:
		return util.UintToU8(uint(arg.Uint64()))
	case *big.Int:
		if arg == nil {
			return 0, errors.New("cannot convert to uint8, input is nil")
		}
		return util.UintToU8(uint(arg.Uint64()))
	case string:
		// Convert the number
		u64, err := convertStr2Uint64(arg, 8)
		if err != nil {
			return 0, err
		}
		return uint8(u64), nil
	default:
		return 0, fmt.Errorf("cannot convert to uint8, input is %T", arg)
	}
}

func convertStr2Uint64(n string, bitSize int) (uint64, error) {
	if common.IsHexString(n) {
		r, err := strconv.ParseUint(strings.TrimPrefix(n, "0x"), 16, bitSize)
		if err != nil {
			return 0, err
		}
		return r, nil
	} else {
		u64, err := strconv.ParseUint(n, 10, bitSize)
		if err != nil {
			return 0, err
		}
		return u64, nil
	}
}

func ConvertToU16(arg any) (uint16, error) {
	switch arg := arg.(type) {
	case float64:
		return util.Float64ToU16(arg)
	case int:
		return util.IntToU16(arg)
	case uint:
		return util.UintToU16(arg)
	case uint16:
		return arg, nil
	case big.Int:
		return util.UintToU16(uint(arg.Uint64()))
	case *big.Int:
		if arg == nil {
			return 0, errors.New("cannot convert to uint16, input is nil")
		}
		return util.UintToU16(uint(arg.Uint64()))
	case string:
		// Convert the number
		u64, err := convertStr2Uint64(arg, 16)
		if err != nil {
			return 0, err
		}
		return uint16(u64), nil
	default:
		return 0, fmt.Errorf("invalid input type for uint16: %T", arg)
	}
}

func ConvertToU32(arg any) (uint32, error) {
	switch arg := arg.(type) {
	case float64:
		return util.Float64ToU32(arg)
	case int:
		return util.IntToU32(arg)
	case uint:
		return util.UintToU32(arg)
	case uint32:
		return arg, nil
	case big.Int:
		return util.UintToU32(uint(arg.Uint64()))
	case *big.Int:
		if arg == nil {
			return 0, errors.New("cannot convert to uint32, input is nil")
		}
		return util.UintToU32(uint(arg.Uint64()))
	case string:
		// Convert the number
		u64, err := convertStr2Uint64(arg, 32)
		if err != nil {
			return 0, err
		}
		return uint32(u64), nil
	default:
		return 0, fmt.Errorf("invalid input type for uint32: %T", arg)
	}
}

func ConvertToU64(arg any) (uint64, error) {
	switch arg := arg.(type) {
	case float64:
		return util.Float64ToU64(arg)
	case int:
		return util.IntToU64(arg)
	case uint:
		return uint64(arg), nil
	case uint64:
		return arg, nil
	case big.Int:
		return arg.Uint64(), nil
	case *big.Int:
		if arg == nil {
			return 0, errors.New("cannot convert to uint64, input is nil")
		}
		return arg.Uint64(), nil
	case string:
		return convertStr2Uint64(arg, 64)
	default:
		return 0, fmt.Errorf("invalid input type for uint64: %T", arg)
	}
}

// TODO: Check bounds of bigints
func ConvertToU128(arg any) (*big.Int, error) {
	switch arg := arg.(type) {
	case float64:
		return util.Float64ToBigInt(arg)
	case int:
		return util.IntToUBigInt(arg)
	case uint:
		return util.UintToUBigInt(arg)
	case big.Int:
		return &arg, nil
	case *big.Int:
		if arg == nil {
			return nil, errors.New("cannot convert to uint128, input is nil")
		}
		return arg, nil
	case string:
		// Convert the number
		return util.StrToBigInt(arg)
	default:
		return nil, fmt.Errorf("invalid input type for uint128: %T", arg)
	}
}

// TODO: Check bounds of bigints
func ConvertToU256(arg any) (*big.Int, error) {
	switch arg := arg.(type) {
	case float64:
		return util.Float64ToBigInt(arg)
	case int:
		return util.IntToUBigInt(arg)
	case uint:
		return util.UintToUBigInt(arg)
	case big.Int:
		return &arg, nil
	case *big.Int:
		if arg == nil {
			return nil, errors.New("cannot convert to uint256, input is nil")
		}
		return arg, nil
	case string:
		// Convert the number
		return util.StrToBigInt(arg)
	default:
		return nil, fmt.Errorf("invalid input type for uint256: %T", arg)
	}
}

func ConvertToBool(arg any) (bool, error) {
	switch arg := arg.(type) {
	case bool:
		return arg, nil
	case string:
		switch arg {
		case "true":
			return true, nil
		case "false":
			return false, nil
		default:
			return false, errors.New("invalid boolean input for bool")
		}
	default:
		return false, fmt.Errorf("invalid input type for bool: %T", arg)
	}
}

func ConvertToAddress(arg any) (*AccountAddress, error) {
	switch arg := arg.(type) {
	case AccountAddress:
		addr := arg
		return &addr, nil
	case *AccountAddress:
		if arg == nil {
			return nil, errors.New("invalid account address, nil")
		}
		return arg, nil
	case string:
		addr := &AccountAddress{}
		err := addr.ParseStringRelaxed(arg)
		if err != nil {
			return nil, err
		}
		return addr, nil
	default:
		return nil, fmt.Errorf("invalid input type for address: %T", arg)
	}
}

// ConvertToVectorU8 returns the BCS encoded version of the bytes
func ConvertToVectorU8(arg any, options ...any) ([]byte, error) {
	compatibilityMode := false
	for _, option := range options {
		if compatMode, ok := option.(CompatibilityMode); ok {
			compatibilityMode = bool(compatMode)
		}
	}

	// Convert input to normalized byte array
	switch arg := arg.(type) {
	// Special case, handle hex string
	case string:
		// Parse string as UTF-8 bytes, in compatibility mode
		if compatibilityMode {
			return bcs.SerializeBytes([]byte(arg))
		}

		bytes, err := util.ParseHex(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeBytes(bytes)
	default:
		return convertToVectorInner(VectorTag{TypeParam: TypeTag{Value: &U8Tag{}}}, arg, []TypeTag{}, options...)
	}
}

// convertToVectorInner specifically is a wrapper to handle the many possible vector input types
func convertToVectorInner(vectorTag VectorTag, arg any, generics []TypeTag, options ...any) ([]byte, error) {
	switch arg := arg.(type) {
	case []any:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []string:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []uint:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []uint8:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []uint16:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []uint32:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []uint64:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []int:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []big.Int:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []*big.Int:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []bool:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []AccountAddress:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	case []*AccountAddress:
		return convertToVectorInnerTyped(vectorTag, arg, generics, options...)
	default:
		return nil, fmt.Errorf("invalid input type for struct vector %T", arg)
	}
}

// convertToVectorInnerTyped handles typed array access to convert a vector
func convertToVectorInnerTyped[T any](vectorTag VectorTag, arg []T, generics []TypeTag, options ...any) ([]byte, error) {
	if arg == nil {
		return nil, errors.New("cannot convert nil to vector")
	}
	length, err := util.IntToU32(len(arg))
	if err != nil {
		return nil, err
	}
	buffer, err := bcs.SerializeUleb128(length)
	if err != nil {
		return nil, err
	}
	for _, item := range arg {
		val, err := ConvertArg(vectorTag.TypeParam, item, generics, options...)
		if err != nil {
			return nil, err
		}
		buffer = append(buffer, val...)
	}
	return buffer, nil
}

func ConvertToVector(vectorTag VectorTag, arg any, generics []TypeTag, options ...any) ([]byte, error) {
	// We have to switch based on type, thanks Golang
	switch innerType := vectorTag.TypeParam.Value.(type) {
	case *U8Tag:
		// Special case to handle hex
		return ConvertToVectorU8(arg, options...)
	case *GenericTag:
		if innerType.Num >= uint64(len(generics)) {
			return nil, errors.New("generic number out of bounds")
		}
		genericType := generics[innerType.Num]
		return convertToVectorInner(VectorTag{TypeParam: genericType}, arg, generics, options...)
	case *ReferenceTag:
		return convertToVectorInner(VectorTag{TypeParam: innerType.TypeParam}, arg, generics, options...)
	default:
		return convertToVectorInner(vectorTag, arg, generics, options...)
	}
}

// CompatibilityMode enables compatibility with the TS SDK in behavior
// This includes "0x00" as an None option
// And string interpreted as bytes instead of hex
type CompatibilityMode bool

func ConvertArg(typeArg TypeTag, arg any, generics []TypeTag, options ...any) ([]byte, error) {
	switch innerType := typeArg.Value.(type) {
	case *U8Tag:
		num, err := ConvertToU8(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeU8(num)
	case *U16Tag:
		num, err := ConvertToU16(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeU16(num)
	case *U32Tag:
		num, err := ConvertToU32(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeU32(num)
	case *U64Tag:
		num, err := ConvertToU64(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeU64(num)
	case *U128Tag:
		num, err := ConvertToU128(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeU128(*num)
	case *U256Tag:
		num, err := ConvertToU256(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeU256(*num)
	case *BoolTag:
		bo, err := ConvertToBool(arg)
		if err != nil {
			return nil, err
		}
		return bcs.SerializeBool(bo)
	case *AddressTag, *SignerTag:
		a, err := ConvertToAddress(arg)
		if err != nil {
			return nil, err
		}
		return bcs.Serialize(a)
	case *GenericTag:
		// Convert based on number
		genericNum := innerType.Num
		if genericNum >= uint64(len(generics)) {
			return nil, errors.New("generic number out of bounds")
		}

		tag := generics[genericNum]
		return ConvertArg(tag, arg, generics, options...)
	case *ReferenceTag:
		// Convert based on inner type
		return ConvertArg(innerType.TypeParam, arg, generics, options...)
	case *VectorTag:
		return ConvertToVector(*innerType, arg, generics, options...)
	case *StructTag:
		structTag := innerType
		// TODO: We should be able to support custom structs, but for now only support known
		if AccountOne == structTag.Address {
			switch structTag.Module {
			case "object":
				if structTag.Name == "Object" {
					// TODO: Move to function
					// Handle as address, inner type doesn't matter
					// TODO: Improve error message
					return ConvertArg(TypeTag{&AddressTag{}}, arg, generics, options...)
				}
			case "string":
				if structTag.Name == "String" {
					// Handle as string, we won't let bytes as an input for now here
					switch arg := arg.(type) {
					case string:
						return bcs.SerializeBytes([]byte(arg))
					default:
						return nil, errors.New("invalid input type for 0x1::string::String")
					}
				}
			case "option":
				if structTag.Name == "Option" {
					// Check it has the proper inner type
					if len(structTag.TypeParams) != 1 {
						return nil, errors.New("invalid input type for option, must have exactly one type arg")
					}
					// Get inner type
					typeParam := structTag.TypeParams[0]

					return ConvertToOption(typeParam, arg, generics, options...)
				}
			}
		}
	default:
		return nil, fmt.Errorf("unknown type argument %T", innerType)
	}

	return nil, errors.New("failed to convert type argument")
}

func convertCompatibilitySerializedType(typeParam TypeTag, arg bcs.Deserializer, generics []TypeTag) ([]byte, error) {
	switch innerType := typeParam.Value.(type) {
	case *U8Tag:
		return bcs.SerializeU8(arg.U8())
	case *U16Tag:
		return bcs.SerializeU16(arg.U16())
	case *U32Tag:
		return bcs.SerializeU32(arg.U32())
	case *U64Tag:
		return bcs.SerializeU64(arg.U64())
	case *U128Tag:
		return bcs.SerializeU128(arg.U128())
	case *U256Tag:
		return bcs.SerializeU256(arg.U256())
	case *BoolTag:
		return bcs.SerializeBool(arg.Bool())
	case *AddressTag:
		return bcs.SerializeBytes(arg.ReadFixedBytes(32))
	case *SignerTag:
		return nil, errors.New("signer is not supported")
	case *GenericTag:
		genericNum := innerType.Num
		if genericNum >= uint64(len(generics)) {
			return nil, errors.New("generic number out of bounds")
		}
		genericType := generics[genericNum]
		return convertCompatibilitySerializedType(genericType, arg, generics)
	case *ReferenceTag:
		return convertCompatibilitySerializedType(innerType.TypeParam, arg, generics)
	case *VectorTag:
		length := arg.Uleb128()
		buffer, err := bcs.SerializeUleb128(length)
		if err != nil {
			return nil, err
		}
		for range int(length) {
			b, err := convertCompatibilitySerializedType(innerType.TypeParam, arg, generics)
			if err != nil {
				return nil, err
			}
			buffer = append(buffer, b...)
		}
		return buffer, nil
	case *StructTag:
		return convertCompatibilitySerializedType(typeParam, arg, generics)
	default:
		return nil, errors.New("unknown type")
	}
}

func ConvertToOption(typeParam TypeTag, arg any, generics []TypeTag, options ...any) ([]byte, error) {
	compatibilityMode := false
	for _, option := range options {
		if compatMode, ok := option.(CompatibilityMode); ok {
			compatibilityMode = bool(compatMode)
		}
	}

	if arg == nil {
		return bcs.SerializeU8(0)
	}

	if compatibilityMode {
		if typedArg, ok := arg.(string); ok {
			if len(typedArg) >= 2 && typedArg[:2] == "0x" {
				typedArg = typedArg[2:]
			}
			bytes, err := hex.DecodeString(typedArg)
			if err != nil {
				return nil, err
			}
			des := bcs.NewDeserializer(bytes)
			length := des.Uleb128()
			if length == 0 {
				return bcs.SerializeU8(0)
			} else {
				b := []byte{1}
				buffer, err := convertCompatibilitySerializedType(typeParam, *des, generics)
				if err != nil {
					return nil, err
				}
				return append(b, buffer...), nil
			}
		}
	}

	b := []byte{1}
	buffer, err := ConvertArg(typeParam, arg, generics, options...)
	if err != nil {
		return nil, err
	}
	return append(b, buffer...), nil
}
