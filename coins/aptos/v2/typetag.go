package v2

import (
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"strings"
)

type TypeTagType uint64

const (
	TypeTagBool    TypeTagType = 0
	TypeTagU8      TypeTagType = 1
	TypeTagU64     TypeTagType = 2
	TypeTagU128    TypeTagType = 3
	TypeTagAddress TypeTagType = 4
	TypeTagSigner  TypeTagType = 5
	TypeTagVector  TypeTagType = 6
	TypeTagStruct  TypeTagType = 7
	TypeTagU16     TypeTagType = 8
	TypeTagU32     TypeTagType = 9
	TypeTagU256    TypeTagType = 10
)

// TypeTagImpl is an interface describing all the different types of TypeTag.  Unfortunately because of how serialization
// works, a wrapper TypeTag struct is needed to handle the differentiation between types
type TypeTagImpl interface {
	bcs.Struct
	GetType() TypeTagType
	String() string
}

// TypeTag is a wrapper around a TypeTagImpl e.g. BoolTag or U8Tag for the purpose of serialization and deserialization
type TypeTag struct {
	Value TypeTagImpl
}

func (tt *TypeTag) MarshalBCS(bcs *bcs.Serializer) {
	bcs.Uleb128(uint32(tt.Value.GetType()))
	bcs.Struct(tt.Value)
}

func (tt *TypeTag) UnmarshalBCS(des *bcs.Deserializer) {
	variant := des.Uleb128()
	switch TypeTagType(variant) {
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
		des.Struct(tt.Value)
	case TypeTagStruct:
		tt.Value = &StructTag{}
		des.Struct(tt.Value)
	default:
		des.SetError(fmt.Errorf("unknown TypeTag enum %d", variant))
	}
}

func (tt *TypeTag) String() string {
	return tt.Value.String()
}

type SignerTag struct{}

func (xt *SignerTag) String() string {
	return "signer"
}

func (xt *SignerTag) GetType() TypeTagType {
	return TypeTagSigner
}

func (xt *SignerTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *SignerTag) UnmarshalBCS(_ *bcs.Deserializer) {}

type AddressTag struct{}

func (xt *AddressTag) String() string {
	return "address"
}

func (xt *AddressTag) GetType() TypeTagType {
	return TypeTagAddress
}

func (xt *AddressTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *AddressTag) UnmarshalBCS(_ *bcs.Deserializer) {}

type BoolTag struct{}

func (xt *BoolTag) String() string {
	return "bool"
}

func (xt *BoolTag) GetType() TypeTagType {
	return TypeTagBool
}

func (xt *BoolTag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *BoolTag) UnmarshalBCS(_ *bcs.Deserializer) {}

type U8Tag struct{}

func (xt *U8Tag) String() string {
	return "u8"
}

func (xt *U8Tag) GetType() TypeTagType {
	return TypeTagU8
}

func (xt *U8Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U8Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

type U16Tag struct{}

func (xt *U16Tag) String() string {
	return "u16"
}

func (xt *U16Tag) GetType() TypeTagType {
	return TypeTagU16
}

func (xt *U16Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U16Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

type U32Tag struct{}

func (xt *U32Tag) String() string {
	return "u32"
}

func (xt *U32Tag) GetType() TypeTagType {
	return TypeTagU32
}

func (xt *U32Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U32Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

type U64Tag struct{}

func (xt *U64Tag) String() string {
	return "u64"
}

func (xt *U64Tag) GetType() TypeTagType {
	return TypeTagU64
}

func (xt *U64Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U64Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

type U128Tag struct{}

func (xt *U128Tag) String() string {
	return "u128"
}

func (xt *U128Tag) GetType() TypeTagType {
	return TypeTagU128
}

func (xt *U128Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U128Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

type U256Tag struct{}

func (xt *U256Tag) String() string {
	return "u256"
}

func (xt *U256Tag) GetType() TypeTagType {
	return TypeTagU256
}
func (xt *U256Tag) MarshalBCS(_ *bcs.Serializer)     {}
func (xt *U256Tag) UnmarshalBCS(_ *bcs.Deserializer) {}

type VectorTag struct {
	TypeParam TypeTag
}

func (xt *VectorTag) GetType() TypeTagType {
	return TypeTagVector
}

func (xt *VectorTag) String() string {
	out := strings.Builder{}
	out.WriteString("vector<")
	out.WriteString(xt.TypeParam.Value.String())
	out.WriteString(">")
	return out.String()
}

func (xt *VectorTag) MarshalBCS(serializer *bcs.Serializer) {
	serializer.Struct(&xt.TypeParam)
}

func (xt *VectorTag) UnmarshalBCS(deserializer *bcs.Deserializer) {
	var tag TypeTag
	tag.UnmarshalBCS(deserializer)
	xt.TypeParam = tag
}

type StructTag struct {
	Address    AccountAddress
	Module     string
	Name       string
	TypeParams []TypeTag
}

func (xt *StructTag) GetType() TypeTagType {
	return TypeTagStruct
}

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
func (xt *StructTag) MarshalBCS(serializer *bcs.Serializer) {
	xt.Address.MarshalBCS(serializer)
	serializer.WriteString(xt.Module)
	serializer.WriteString(xt.Name)
	bcs.SerializeSequence(xt.TypeParams, serializer)
}
func (xt *StructTag) UnmarshalBCS(deserializer *bcs.Deserializer) {
	xt.Address.UnmarshalBCS(deserializer)
	xt.Module = deserializer.ReadString()
	xt.Name = deserializer.ReadString()
	xt.TypeParams = bcs.DeserializeSequence[TypeTag](deserializer)
}

func NewTypeTag(inner TypeTagImpl) TypeTag {
	return TypeTag{
		Value: inner,
	}
}

func NewVectorTag(inner TypeTagImpl) *VectorTag {
	return &VectorTag{
		TypeParam: NewTypeTag(inner),
	}
}

func NewStringTag() *StructTag {
	return &StructTag{
		Address:    AccountOne,
		Module:     "string",
		Name:       "String",
		TypeParams: []TypeTag{},
	}
}

func NewOptionTag(inner TypeTagImpl) *StructTag {
	return &StructTag{
		Address:    AccountOne,
		Module:     "option",
		Name:       "Option",
		TypeParams: []TypeTag{NewTypeTag(inner)},
	}
}

func NewObjectTag(inner TypeTagImpl) *StructTag {
	return &StructTag{
		Address:    AccountOne,
		Module:     "object",
		Name:       "Object",
		TypeParams: []TypeTag{NewTypeTag(inner)},
	}
}
