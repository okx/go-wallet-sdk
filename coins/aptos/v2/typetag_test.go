package v2

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTypeTag(t *testing.T) {
	// This is unfortunate with references
	nested := NewTypeTag(NewOptionTag(NewVectorTag(NewObjectTag(NewStringTag()))))

	assert.Equal(t, "0x1::option::Option<vector<0x1::object::Object<0x1::string::String>>>", nested.String())

	ser := &bcs.Serializer{}
	ser.Struct(&nested)
	assert.NoError(t, ser.Error())

	bytes := ser.ToBytes()

	des := bcs.NewDeserializer(bytes)
	tag := &TypeTag{}
	des.Struct(tag)
	assert.NoError(t, des.Error())

	// Check the deserialized is correct
	assert.Equal(t, &nested, tag)
}

func TestTypeTagIdentities(t *testing.T) {
	checkVariant(t, &AddressTag{}, TypeTagAddress, "address")
	checkVariant(t, &SignerTag{}, TypeTagSigner, "signer")
	checkVariant(t, &BoolTag{}, TypeTagBool, "bool")
	checkVariant(t, &U8Tag{}, TypeTagU8, "u8")
	checkVariant(t, &U16Tag{}, TypeTagU16, "u16")
	checkVariant(t, &U32Tag{}, TypeTagU32, "u32")
	checkVariant(t, &U64Tag{}, TypeTagU64, "u64")
	checkVariant(t, &U128Tag{}, TypeTagU128, "u128")
	checkVariant(t, &U256Tag{}, TypeTagU256, "u256")

	checkVariant(t, NewVectorTag(&U8Tag{}), TypeTagVector, "vector<u8>")
	checkVariant(t, NewStringTag(), TypeTagStruct, "0x1::string::String")
}

func checkVariant[T TypeTagImpl](t *testing.T, tag T, expectedType TypeTagVariant, expectedString string) {
	assert.Equal(t, expectedType, tag.GetType())
	assert.Equal(t, expectedString, tag.String())

	// Serialize and deserialize test
	tt := NewTypeTag(tag)
	bytes, err := bcs.Serialize(&tt)
	assert.NoError(t, err)
	var newTag TypeTag
	err = bcs.Deserialize(&newTag, bytes)
	assert.NoError(t, err)
	assert.Equal(t, tt, newTag)
}

func TestStructTag(t *testing.T) {
	st := StructTag{
		Address: AccountOne,
		Module:  "coin",
		Name:    "CoinStore",
		TypeParams: []TypeTag{
			{Value: &StructTag{
				Address:    AccountOne,
				Module:     "aptos_coin",
				Name:       "AptosCoin",
				TypeParams: nil,
			}},
		},
	}
	assert.Equal(t, "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin>", st.String())
	var aa3 AccountAddress
	err := aa3.ParseStringRelaxed("0x3")
	assert.NoError(t, err)
	st.TypeParams = append(st.TypeParams, TypeTag{Value: &StructTag{
		Address:    aa3,
		Module:     "other",
		Name:       "thing",
		TypeParams: nil,
	}})
	assert.Equal(t, "0x1::coin::CoinStore<0x1::aptos_coin::AptosCoin,0x3::other::thing>", st.String())
}
