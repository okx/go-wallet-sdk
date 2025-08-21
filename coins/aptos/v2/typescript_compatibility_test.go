package v2

import (
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/bcs"
	"github.com/okx/go-wallet-sdk/coins/aptos/v2/internal/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	MULTI_AGENT_NO_FEE_PAYER = "0x000000000000000000000000000000000000000000000000000000000000000a01000000000000000200000000000000000000000000000000000000000000000000000000000000010d6170746f735f6163636f756e74087472616e73666572000040420f0000000000640000000000000000f15365000000000101000000000000000000000000000000000000000000000000000000000000000400"
	MULTI_AGENT_FEE_PAYER    = "0x000000000000000000000000000000000000000000000000000000000000000a01000000000000000200000000000000000000000000000000000000000000000000000000000000010d6170746f735f6163636f756e74087472616e73666572000040420f0000000000640000000000000000f153650000000001010000000000000000000000000000000000000000000000000000000000000004010000000000000000000000000000000000000000000000000000000000000002"
	SINGLE_SIGNER_FEE_PAYER  = "0x000000000000000000000000000000000000000000000000000000000000000a01000000000000000200000000000000000000000000000000000000000000000000000000000000010d6170746f735f6163636f756e74087472616e73666572000040420f0000000000640000000000000000f15365000000000100010000000000000000000000000000000000000000000000000000000000000002"
)

func TestMultiAgentFromTypescript_NoFeePayer(t *testing.T) {
	t.Parallel()
	rawTxnWithData := &RawTransactionWithData{}
	bytes, err := util.ParseHex(MULTI_AGENT_NO_FEE_PAYER)
	require.NoError(t, err)
	des := bcs.NewDeserializer(bytes)
	rawTxnWithData.UnmarshalTypeScriptBCS(des)
	require.NoError(t, des.Error())

	assert.Equal(t, MultiAgentRawTransactionWithDataVariant, rawTxnWithData.Variant)
	inner, ok := rawTxnWithData.Inner.(*MultiAgentRawTransactionWithData)
	require.True(t, ok)
	assert.Equal(t, AccountTen, inner.RawTxn.Sender)
	assert.Equal(t, AccountFour, inner.SecondarySigners[0])

	// Test serializing back
	data, err := bcs.SerializeSingle(func(ser *bcs.Serializer) {
		rawTxnWithData.MarshalTypeScriptBCS(ser)
	})
	require.NoError(t, err)
	assert.Equal(t, bytes, data)
}

func TestMultiAgentFromTypescript_FeePayer(t *testing.T) {
	t.Parallel()
	rawTxnWithData := &RawTransactionWithData{}
	bytes, err := util.ParseHex(MULTI_AGENT_FEE_PAYER)
	require.NoError(t, err)
	des := bcs.NewDeserializer(bytes)
	rawTxnWithData.UnmarshalTypeScriptBCS(des)
	require.NoError(t, des.Error())

	assert.Equal(t, MultiAgentWithFeePayerRawTransactionWithDataVariant, rawTxnWithData.Variant)
	inner, ok := rawTxnWithData.Inner.(*MultiAgentWithFeePayerRawTransactionWithData)
	require.True(t, ok)
	assert.Equal(t, AccountTen, inner.RawTxn.Sender)
	assert.Len(t, inner.SecondarySigners, 1)
	assert.Equal(t, AccountFour, inner.SecondarySigners[0])
	assert.Equal(t, &AccountTwo, inner.FeePayer)

	// Test serializing back
	data, err := bcs.SerializeSingle(func(ser *bcs.Serializer) {
		rawTxnWithData.MarshalTypeScriptBCS(ser)
	})
	require.NoError(t, err)
	assert.Equal(t, bytes, data)
}

func TestMultiAgentFromTypescript_SingleSignerFeePayer(t *testing.T) {
	t.Parallel()
	rawTxnWithData := &RawTransactionWithData{}
	bytes, err := util.ParseHex(SINGLE_SIGNER_FEE_PAYER)
	require.NoError(t, err)
	des := bcs.NewDeserializer(bytes)
	rawTxnWithData.UnmarshalTypeScriptBCS(des)
	require.NoError(t, des.Error())

	assert.Equal(t, MultiAgentWithFeePayerRawTransactionWithDataVariant, rawTxnWithData.Variant)
	inner, ok := rawTxnWithData.Inner.(*MultiAgentWithFeePayerRawTransactionWithData)
	require.True(t, ok)
	assert.Equal(t, AccountTen, inner.RawTxn.Sender)
	assert.Empty(t, inner.SecondarySigners)
	assert.Equal(t, &AccountTwo, inner.FeePayer)

	// Test serializing back
	data, err := bcs.SerializeSingle(func(ser *bcs.Serializer) {
		rawTxnWithData.MarshalTypeScriptBCS(ser)
	})
	require.NoError(t, err)
	assert.Equal(t, bytes, data)
}
