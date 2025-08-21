package v2

import (
	"errors"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"regexp"
	"strconv"
	"strings"
)

//region TypeTag

// TypeTagVariant is an enum representing the different types of TypeTag
type TypeTagVariant uint32

const (
	TypeTagBool      TypeTagVariant = 0   // Represents the bool type in Move BoolTag
	TypeTagU8        TypeTagVariant = 1   // Represents the u8 type in Move U8Tag
	TypeTagU64       TypeTagVariant = 2   // Represents the u64 type in Move U64Tag
	TypeTagU128      TypeTagVariant = 3   // Represents the u128 type in Move U128Tag
	TypeTagAddress   TypeTagVariant = 4   // Represents the address type in Move AddressTag
	TypeTagSigner    TypeTagVariant = 5   // Represents the signer type in Move SignerTag
	TypeTagVector    TypeTagVariant = 6   // Represents the vector type in Move VectorTag
	TypeTagStruct    TypeTagVariant = 7   // Represents the struct type in Move StructTag
	TypeTagU16       TypeTagVariant = 8   // Represents the u16 type in Move U16Tag
	TypeTagU32       TypeTagVariant = 9   // Represents the u32 type in Move U32Tag
	TypeTagU256      TypeTagVariant = 10  // Represents the u256 type in Move U256Tag
	TypeTagGeneric   TypeTagVariant = 254 // Represents a generic type in Move GenericTag
	TypeTagReference TypeTagVariant = 255 // Represents the reference type in Move ReferenceTag
)

// TypeTagImpl is an interface describing all the different types of [TypeTag].  Unfortunately because of how serialization
// works, a wrapper TypeTag struct is needed to handle the differentiation between types
type TypeTagImpl interface {
	bcs.Struct
	// GetType returns the TypeTagVariant for this [TypeTag]
	GetType() TypeTagVariant
	// String returns the canonical Move string representation of this [TypeTag]
	String() string
}

// TypeTag is a wrapper around a [TypeTagImpl] e.g. [BoolTag] or [U8Tag] for the purpose of serialization and deserialization
// Implements:
//   - [bcs.Struct]
type TypeTag struct {
	Value TypeTagImpl
}

type parseInfo struct {
	expectedTypes int
	types         []TypeTag
	str           string
}

func ParseTypeTag(inputStr string) (*TypeTag, error) {
	inputRunes := []rune(inputStr)
	// Represents the stack of types currently being processed
	saved := make([]parseInfo, 0)
	// Represents the inner types for a type tag e.g. '0x1::coin::Coin<InnerType>'
	innerTypes := make([]TypeTag, 0)
	// Represents the current parsed types in a comma list e.g. 'u8, u8'
	curTypes := make([]TypeTag, 0)
	// The current character index of the whole string
	cur := 0
	// The current working string as type name
	currentStr := ""
	// The expected types based on the number of commas
	expectedTypes := 1

	// Iterate through characters, handling border conditions, we don't use a range because we sometimes skip ahead
	for cur < len(inputRunes) {
		r := inputRunes[cur]
		//println(fmt.Printf("%c | %s | %s\n", r, currentStr, util.PrettyJson(saved)))

		switch r {
		case '<':
			// Start of a type argument, save the current state
			saved = append(saved, parseInfo{
				expectedTypes: expectedTypes,
				types:         curTypes,
				str:           currentStr,
			})

			// Clear current state
			currentStr = ""
			curTypes = make([]TypeTag, 0)
			expectedTypes = 1
		case '>':
			// End of type arguments, process last type, if there is no type string then don't parse it
			if currentStr != "" {
				newType, err := ParseTypeTagInner(currentStr, innerTypes)
				if err != nil {
					return nil, err
				}

				curTypes = append(curTypes, *newType)
			}

			// If there's nothing left there were too many '>'
			savedLength := len(saved)
			if savedLength == 0 {
				return nil, errors.New("no inner types found")
			}

			// Ensure commas match types
			if expectedTypes != len(curTypes) {
				return nil, errors.New("inner type count mismatch, too many commas")
			}

			// Pop off stack
			savedPop := saved[savedLength-1]
			saved = saved[:savedLength-1]

			innerTypes = curTypes
			curTypes = savedPop.types
			currentStr = savedPop.str
			expectedTypes = savedPop.expectedTypes
		case ',':
			if len(saved) == 0 {
				return nil, fmt.Errorf("unexpected comma at top level type")
			}
			if len(currentStr) == 0 {
				return nil, fmt.Errorf("unexpected comma, TypeTag is missing")
			}

			newType, err := ParseTypeTagInner(currentStr, innerTypes)
			if err != nil {
				return nil, err
			}
			innerTypes = make([]TypeTag, 0)
			curTypes = append(curTypes, *newType)
			currentStr = ""
			expectedTypes += 1
		case ' ':
			// TODO whitespace, do we include tabs, etc.
			parsedTypeTag := false

			if len(currentStr) != 0 {
				// parse type tag, and push it on the current types
				newType, err := ParseTypeTagInner(currentStr, innerTypes)
				if err != nil {
					return nil, err
				}

				innerTypes = make([]TypeTag, 0)
				curTypes = append(curTypes, *newType)
				currentStr = ""
				parsedTypeTag = true
			}

			// Skip any additional whitespace
			for cur < len(inputRunes) {
				if inputRunes[cur] != ' ' {
					break
				}
				cur += 1
			}

			// Next char must be a comma or a closing > if something was parsed before it
			nextChar := inputRunes[cur]
			if cur < len(inputRunes) && parsedTypeTag && nextChar != ',' && nextChar != '>' {
				return nil, fmt.Errorf("unexpected character at top level type")
			}

			// Skip over incrementing, we already did it above
			continue
		default:
			currentStr += string(r)
		}

		cur += 1
	}

	if len(saved) > 0 {
		return nil, fmt.Errorf("missing type argument close '>'")
	}

	switch len(curTypes) {
	case 0:
		return ParseTypeTagInner(currentStr, innerTypes)
	case 1:
		if currentStr == "" {
			return &curTypes[0], nil
		}
		return nil, fmt.Errorf("unexpected comma ','")
	default:
		return nil, fmt.Errorf("unexpected whitespace")
	}
}

func ParseTypeTagInner(input string, types []TypeTag) (*TypeTag, error) {
	str := strings.TrimSpace(input)

	//println(fmt.Printf("-- %s | %s\n", input, util.PrettyJson(types)))
	// TODO: for now we aren't going to lowercase this

	// Handle primitive types
	switch str {
	case "bool":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &BoolTag{}}, nil
	case "u8":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &U8Tag{}}, nil
	case "u16":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &U16Tag{}}, nil
	case "u32":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &U32Tag{}}, nil
	case "u64":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &U64Tag{}}, nil
	case "u128":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &U128Tag{}}, nil
	case "u256":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &U256Tag{}}, nil
	case "address":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &AddressTag{}}, nil
	case "signer":
		if len(types) > 0 {
			return nil, fmt.Errorf("invalid type tag, primitive with generics")
		}
		return &TypeTag{Value: &SignerTag{}}, nil
	case "vector":
		if len(types) != 1 {
			return nil, fmt.Errorf("unexpected number of types for vector, expected 1, got %d", len(types))
		}
		return &TypeTag{Value: &VectorTag{TypeParam: types[0]}}, nil
	default:
		// If it's a reference
		if strings.HasPrefix(str, "&") {
			actualType, _ := strings.CutPrefix(str, "&")
			inner, err := ParseTypeTagInner(actualType, types)
			if err != nil {
				return nil, err
			}
			return &TypeTag{Value: &ReferenceTag{TypeParam: *inner}}, nil
		}

		// If it's generic
		if strings.HasPrefix(str, "T") {
			numStr := strings.TrimPrefix(str, "T")
			num, err := strconv.ParseUint(numStr, 10, 32)
			if err != nil {
				return nil, err
			}
			return &TypeTag{Value: &GenericTag{Num: num}}, nil
		}

		parts := strings.Split(str, "::")
		if len(parts) != 3 {
			// TODO: More informative message
			return nil, errors.New("invalid type tag")
		}

		// Validate struct address
		address := &AccountAddress{}
		//println("PARTS:", util.PrettyJson(parts))
		err := address.ParseStringWithPrefixRelaxed(parts[0])
		if err != nil {
			return nil, errors.New("invalid type tag struct address")
		}

		// Validate module
		module := parts[1]
		moduleValid, err := regexp.MatchString("^[a-zA-Z_0-9]+$", module)
		if !moduleValid || err != nil {
			return nil, errors.New("invalid type tag struct module")
		}

		// Validate name
		name := parts[2]
		nameValid, err := regexp.MatchString("^[a-zA-Z_0-9]+$", name)
		if !nameValid || err != nil {
			return nil, errors.New("invalid type tag struct name")
		}

		return &TypeTag{Value: &StructTag{
			Address:    *address,
			Module:     module,
			Name:       name,
			TypeParams: types,
		}}, nil
	}
}

// String gives the canonical TypeTag string value used in Move
func (tt *TypeTag) String() string {
	return tt.Value.String()
}

//region TypeTag bcs.Struct

// MarshalBCS serializes the TypeTag to bytes
//
// Implements:
//   - [bcs.Marshaler]
func (tt *TypeTag) MarshalBCS(ser *bcs.Serializer) {
	ser.Uleb128(uint32(tt.Value.GetType()))
	ser.Struct(tt.Value)
}

// UnmarshalBCS deserializes the TypeTag from bytes
//
// Implements:
//   - [bcs.Unmarshaler]
func (tt *TypeTag) UnmarshalBCS(des *bcs.Deserializer) {
	variant := TypeTagVariant(des.Uleb128())
	switch variant {
	case TypeTagAddress:
		tt.Value = &AddressTag{}
	case TypeTagSigner:
		tt.Value = &SignerTag{}
	case TypeTagBool:
		tt.Value = &BoolTag{}
	case TypeTagU8:
		tt.Value = &U8Tag{}
	case TypeTagU16:
		tt.Value = &U16Tag{}
	case TypeTagU32:
		tt.Value = &U32Tag{}
	case TypeTagU64:
		tt.Value = &U64Tag{}
	case TypeTagU128:
		tt.Value = &U128Tag{}
	case TypeTagU256:
		tt.Value = &U256Tag{}
	case TypeTagVector:
		tt.Value = &VectorTag{}
	case TypeTagStruct:
		tt.Value = &StructTag{}
	default:
		des.SetError(fmt.Errorf("unknown TypeTag enum %d", variant))
		return
	}
	des.Struct(tt.Value)
}

//endregion
//endregion

//region SignerTag

// SignerTag represents the signer type in Move
type SignerTag struct{}

//region SignerTag TypeTagImpl

func (xt *SignerTag) String() string {
	return "signer"
}

func (xt *SignerTag) GetType() TypeTagVariant {
	return TypeTagSigner
}

//endregion

//region SignerTag bcs.Struct

func (xt *SignerTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *SignerTag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region AddressTag

// AddressTag represents the address type in Move
type AddressTag struct{}

//region AddressTag TypeTagImpl

func (xt *AddressTag) String() string {
	return "address"
}

func (xt *AddressTag) GetType() TypeTagVariant {
	return TypeTagAddress
}

//endregion

//region AddressTag bcs.Struct

func (xt *AddressTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *AddressTag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region BoolTag

// BoolTag represents the bool type in Move
type BoolTag struct{}

//region BoolTag TypeTagImpl

func (xt *BoolTag) String() string {
	return "bool"
}

func (xt *BoolTag) GetType() TypeTagVariant {
	return TypeTagBool
}

//endregion

//region BoolTag bcs.struct

func (xt *BoolTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *BoolTag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region U8Tag

// U8Tag represents the u8 type in Move
type U8Tag struct{}

//region U8Tag TypeTagImpl

func (xt *U8Tag) String() string {
	return "u8"
}

func (xt *U8Tag) GetType() TypeTagVariant {
	return TypeTagU8
}

//endregion

//region U8Tag bcs.Struct

func (xt *U8Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U8Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region U16Tag

// U16Tag represents the u16 type in Move
type U16Tag struct{}

//region U16Tag TypeTagImpl

func (xt *U16Tag) String() string {
	return "u16"
}

func (xt *U16Tag) GetType() TypeTagVariant {
	return TypeTagU16
}

//endregion

//region U16Tag bcs.Struct

func (xt *U16Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U16Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region U32Tag

// U32Tag represents the u32 type in Move
type U32Tag struct{}

//region U32Tag TypeTagImpl

func (xt *U32Tag) String() string {
	return "u32"
}

func (xt *U32Tag) GetType() TypeTagVariant {
	return TypeTagU32
}

//endregion

//region U32Tag bcs.Struct

func (xt *U32Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U32Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region U64Tag

// U64Tag represents the u64 type in Move
type U64Tag struct{}

//region U64Tag TypeTagImpl

func (xt *U64Tag) String() string {
	return "u64"
}

func (xt *U64Tag) GetType() TypeTagVariant {
	return TypeTagU64
}

//endregion

//region U64Tag bcs.Struct

func (xt *U64Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U64Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region U128Tag

// U128Tag represents the u128 type in Move
type U128Tag struct{}

//region U128Tag TypeTagImpl

func (xt *U128Tag) String() string {
	return "u128"
}

func (xt *U128Tag) GetType() TypeTagVariant {
	return TypeTagU128
}

//endregion

//region U128Tag bcs.Struct

func (xt *U128Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U128Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region U256Tag

// U256Tag represents the u256 type in Move
type U256Tag struct{}

//region U256Tag TypeTagImpl

func (xt *U256Tag) String() string {
	return "u256"
}

func (xt *U256Tag) GetType() TypeTagVariant {
	return TypeTagU256
}

//endregion

//region U256Tag bcs.Struct

func (xt *U256Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U256Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion
//endregion

//region VectorTag

// VectorTag represents the vector<T> type in Move, where T is another [TypeTag]
type VectorTag struct {
	TypeParam TypeTag // TypeParam is the type of the elements in the vector
}

//region VectorTag TypeTagImpl

func (xt *VectorTag) GetType() TypeTagVariant {
	return TypeTagVector
}

func (xt *VectorTag) String() string {
	out := strings.Builder{}
	out.WriteString("vector<")
	out.WriteString(xt.TypeParam.Value.String())
	out.WriteString(">")
	return out.String()
}

//endregion

//region TypeTagVector bcs.Struct

func (xt *VectorTag) MarshalBCS(ser *bcs.Serializer) {
	ser.Struct(&xt.TypeParam)
}

func (xt *VectorTag) UnmarshalBCS(des *bcs.Deserializer) {
	var tag TypeTag
	tag.UnmarshalBCS(des)
	xt.TypeParam = tag
}

//endregion
//endregion

//region StructTag

// StructTag represents an on-chain struct of the form address::module::name<T1,T2,...> and each T is a [TypeTag]
type StructTag struct {
	Address    AccountAddress // Address is the address of the module
	Module     string         // Module is the name of the module
	Name       string         // Name is the name of the struct
	TypeParams []TypeTag      // TypeParams are the TypeTags of the type parameters
}

//region StructTag TypeTagImpl

func (xt *StructTag) GetType() TypeTagVariant {
	return TypeTagStruct
}

// String outputs to the form address::module::name<type1, type2> e.g.
// 0x1::string::String or 0x42::my_mod::MultiType<u8,0x1::string::String>
func (xt *StructTag) String() string {
	out := strings.Builder{}
	out.WriteString(xt.Address.String())
	out.WriteString("::")
	out.WriteString(xt.Module)
	out.WriteString("::")
	out.WriteString(xt.Name)
	if len(xt.TypeParams) != 0 {
		out.WriteRune('<')
		for i, tp := range xt.TypeParams {
			if i != 0 {
				out.WriteRune(',')
			}
			out.WriteString(tp.String())
		}
		out.WriteRune('>')
	}
	return out.String()
}

//endregion

//region StructTag bcs.Struct

func (xt *StructTag) MarshalBCS(ser *bcs.Serializer) {
	xt.Address.MarshalBCS(ser)
	ser.WriteString(xt.Module)
	ser.WriteString(xt.Name)
	bcs.SerializeSequence(xt.TypeParams, ser)
}
func (xt *StructTag) UnmarshalBCS(des *bcs.Deserializer) {
	xt.Address.UnmarshalBCS(des)
	xt.Module = des.ReadString()
	xt.Name = des.ReadString()
	xt.TypeParams = bcs.DeserializeSequence[TypeTag](des)
}

//endregion
//endregion

//region ReferenceTag

// ReferenceTag represents a reference of a type in Move
type ReferenceTag struct {
	TypeParam TypeTag
}

//region ReferenceTag TypeTagImpl

func (xt *ReferenceTag) String() string {
	out := strings.Builder{}
	out.WriteString("&")
	out.WriteString(xt.TypeParam.Value.String())
	return out.String()
}

func (xt *ReferenceTag) GetType() TypeTagVariant {
	return TypeTagReference
}

//endregion

//region Reference bcs.Struct

// TODO: Do we need a proper serialization here
func (xt *ReferenceTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *ReferenceTag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion

//region GenericTag

// GenericTag represents a generic of a type in Move
type GenericTag struct {
	Num uint64
}

//region GenericTag TypeTagImpl

func (xt *GenericTag) String() string {
	out := strings.Builder{}
	out.WriteString("T")
	out.WriteString(strconv.FormatUint(xt.Num, 10))
	return out.String()
}

func (xt *GenericTag) GetType() TypeTagVariant {
	return TypeTagGeneric
}

//endregion

//region Generic bcs.Struct

// TODO: Do we need a proper serialization here
func (xt *GenericTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *GenericTag) UnmarshalBCS(_ *bcs.Deserializer) {}

//endregion

//region TypeTag helpers

// NewTypeTag wraps a TypeTagImpl in a TypeTag
func NewTypeTag(inner TypeTagImpl) TypeTag {
	return TypeTag{
		Value: inner,
	}
}

// NewVectorTag creates a TypeTag for vector<inner>
func NewVectorTag(inner TypeTagImpl) *VectorTag {
	return &VectorTag{
		TypeParam: NewTypeTag(inner),
	}
}

// NewStringTag creates a TypeTag for 0x1::string::String
func NewStringTag() *StructTag {
	return &StructTag{
		Address:    AccountOne,
		Module:     "string",
		Name:       "String",
		TypeParams: []TypeTag{},
	}
}

// NewOptionTag creates a 0x1::option::Option TypeTag based on an inner type
func NewOptionTag(inner TypeTagImpl) *StructTag {
	return &StructTag{
		Address:    AccountOne,
		Module:     "option",
		Name:       "Option",
		TypeParams: []TypeTag{NewTypeTag(inner)},
	}
}

// NewObjectTag creates a 0x1::object::Object TypeTag based on an inner type
func NewObjectTag(inner TypeTagImpl) *StructTag {
	return &StructTag{
		Address:    AccountOne,
		Module:     "object",
		Name:       "Object",
		TypeParams: []TypeTag{NewTypeTag(inner)},
	}
}

// AptosCoinTypeTag is the TypeTag for 0x1::aptos_coin::AptosCoin
var AptosCoinTypeTag = TypeTag{&StructTag{
	Address: AccountOne,
	Module:  "aptos_coin",
	Name:    "AptosCoin",
}}

//endregion
