package v2

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

var scriptU8Bytes = []byte{0, 1}
var scriptU16Bytes = []byte{6, 2, 0}
var scriptU32Bytes = []byte{7, 3, 0, 0, 0}
var scriptU64Bytes = []byte{1, 4, 0, 0, 0, 0, 0, 0, 0}
var scriptU128Bytes = []byte{2, 5, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var scriptU256Bytes = []byte{8, 6, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
var scriptBoolBytes = []byte{5, 0}
var scriptAddressBytes = []byte{3, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 4}
var scriptVectorU8Bytes = []byte{4, 5, 1, 2, 3, 4, 5}
var scriptVectorStringBytes = []byte{9, 5, 2, 1, 49, 1, 50}
var scriptSerializedBytes = []byte{9, 2, 1, 2}

func TestScript_MarshalBCS(t *testing.T) {
	validateBytes := func(input ScriptArgument, expected []byte) {
		var ser bcs.Serializer
		input.MarshalBCS(&ser)
		assert.Equal(t, expected, ser.ToBytes())
	}
	validateBytes(ScriptArgument{Variant: ScriptArgumentU8, Value: uint8(1)}, scriptU8Bytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentU16, Value: uint16(2)}, scriptU16Bytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentU32, Value: uint32(3)}, scriptU32Bytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentU64, Value: uint64(4)}, scriptU64Bytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentU128, Value: *big.NewInt(5)}, scriptU128Bytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentU256, Value: *big.NewInt(6)}, scriptU256Bytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentBool, Value: false}, scriptBoolBytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentAddress, Value: AccountFour}, scriptAddressBytes)

	validateBytes(ScriptArgument{Variant: ScriptArgumentU8Vector, Value: []byte{1, 2, 3, 4, 5}}, scriptVectorU8Bytes)
	validateBytes(ScriptArgument{Variant: ScriptArgumentSerialized, Value: "0x0102"}, scriptSerializedBytes)

}
