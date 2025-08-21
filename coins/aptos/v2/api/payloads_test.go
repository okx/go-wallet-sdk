package api

import (
	"encoding/json"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPayload_EntryFunction(t *testing.T) {
	testJson := `{
      "function": "0x1::object::transfer",
      "type_arguments": [
        "0x4::token::Token"
      ],
      "arguments": [
        {
          "inner": "0xe954c1985dd1da6eabd76554f33c6c8e0831eaaf06e97198a0973b7ef3a712ca"
        },
        "0x116fb1e503bfa08d1f5237206dd9645c944dfe31913e61836388e15824d68573"
      ],
      "type": "entry_function_payload"
    }`
	data := &TransactionPayload{}
	err := json.Unmarshal([]byte(testJson), &data)
	assert.NoError(t, err)
	assert.Equal(t, data.Type, TransactionPayloadVariantEntryFunction)
	payload := data.Inner.(*TransactionPayloadEntryFunction)

	assert.Equal(t, "0x1::object::transfer", payload.Function)
	assert.Len(t, payload.TypeArguments, 1)
	assert.Equal(t, "0x4::token::Token", payload.TypeArguments[0])
	assert.Len(t, payload.Arguments, 2)
}

func TestPayload_Script(t *testing.T) {
	testJson := `{
  "code": {
    "bytecode": "0xa11ceb0b0500000008010008020804030c150421020523100733500883012006a30114000000010002000301050800030403010002060105010001070002000008000200010403060c050301050001060c01080001030d6170746f735f6163636f756e740a6170746f735f636f696e04636f696e067369676e65720a616464726573735f6f66094170746f73436f696e0762616c616e6365046d696e74087472616e7366657200000000000000000000000000000000000000000000000000000000000000010308a0860100000000000308ffffffffffffffff000001170a0011000c030a03380007010a02170700172304120a000b030a0207001611020b000b010b02110302",
    "abi": {
      "name": "main",
      "visibility": "public",
      "is_entry": true,
      "is_view": false,
      "generic_type_params": [],
      "params": [
        "&signer",
        "address",
        "u64"
      ],
      "return": []
    }
  },
  "type_arguments": [],
  "arguments": [
    "0x978c213990c4833df71548df7ce49d54c759d6b6d932de22b24d56060b7af2aa",
    "100000000"
  ],
  "type": "script_payload"
}`
	data := &TransactionPayload{}
	err := json.Unmarshal([]byte(testJson), &data)
	assert.NoError(t, err)
	assert.Equal(t, data.Type, TransactionPayloadVariantScript)
	payload := data.Inner.(*TransactionPayloadScript)

	assert.Len(t, payload.Code.Bytecode, 263)
	assert.Equal(t, "main", payload.Code.Abi.Name)
	assert.Len(t, payload.TypeArguments, 0)
	assert.Len(t, payload.Arguments, 2)
}

func TestPayload_Multisig(t *testing.T) {
	testJson := `{
  "payload": {
      "function": "0x1::object::transfer",
      "type_arguments": [
        "0x4::token::Token"
      ],
      "arguments": [
        {
          "inner": "0xe954c1985dd1da6eabd76554f33c6c8e0831eaaf06e97198a0973b7ef3a712ca"
        },
        "0x116fb1e503bfa08d1f5237206dd9645c944dfe31913e61836388e15824d68573"
      ],
      "type": "entry_function_payload"
  },
  "multisig_address": "0x1",
  "type": "multisig_payload"
}`
	data := &TransactionPayload{}
	err := json.Unmarshal([]byte(testJson), &data)
	assert.NoError(t, err)
	assert.Equal(t, data.Type, TransactionPayloadVariantMultisig)
	payload := data.Inner.(*TransactionPayloadMultisig)
	assert.Equal(t, types.AccountOne, *payload.MultisigAddress)
}

func TestPayload_ModuleBundle(t *testing.T) {
	//todo
	//	testJson := `{
	//  "type": "module_bundle_payload"
	//}`
	//	data := &TransactionPayload{}
	//	err := json.Unmarshal([]byte(testJson), data)
	//	assert.NoError(t, err)
	//	assert.Equal(t, data.Type, TransactionPayloadVariantModuleBundle)
}

func TestPayload_Unknown(t *testing.T) {
	testJson := `{
  "something": true,
  "type": "new_payload"
}`
	data := &TransactionPayload{}
	err := json.Unmarshal([]byte(testJson), &data)
	assert.NoError(t, err)
	assert.Equal(t, data.Type, TransactionPayloadVariantUnknown)
	payload := data.Inner.(*TransactionPayloadUnknown)

	assert.Equal(t, payload.Type, "new_payload")
	assert.Equal(t, payload.Payload["something"].(bool), true)
}
