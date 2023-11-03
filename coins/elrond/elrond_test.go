package elrond

import (
	"crypto/ed25519"
	"encoding/hex"
	"encoding/json"
	"github.com/btcsuite/btcd/btcutil/bech32"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestNewAddress(t *testing.T) {
	pk, _ := hex.DecodeString("27d57eb22fc218b83e9ea2da55746d9318ba6b89cfa31b797e7296bf8a66e4f1")
	address, err := AddressFromSeed(hex.EncodeToString(pk))
	require.Nil(t, err)
	expected := "erd1g7rq9se53mmr5zq8hmz2jpttsgwx7wfx8dgmp9qj2mc7dnj99cys33f33x"
	require.Equal(t, expected, address)
	ret := ValidateAddress(address)
	require.True(t, ret)
}

func TestBuild(t *testing.T) {
	pk, _ := hex.DecodeString("27d57eb22fc218b83e9ea2da55746d9318ba6b89cfa31b797e7296bf8a66e4f1")
	pk2, _ := hex.DecodeString("b6b77f7440d6c88b634bd95bf0da0a8660e781e88915f30319646170a57c4de1")
	privateKey := ed25519.NewKeyFromSeed(pk)
	toAddress, _ := AddressFromSeed(hex.EncodeToString(pk2))
	args := ArgCreateTransaction{
		Nonce:    3,
		Value:    "10000000000000000",
		RcvAddr:  toAddress,
		GasPrice: 1000000000,
		GasLimit: 50000,
		ChainID:  "T",
		Version:  2,
		Options:  1,
	}
	builder := NewTxBuilder(&privateKey)
	tran, _ := builder.build(args)
	signedTxBs, _ := json.Marshal(tran)
	expected := "{\"nonce\":3,\"value\":\"10000000000000000\",\"receiver\":\"erd1lu89n007dkpy3s3za7wwguasttw6xlqqlvx20h0wwm778h5pfe8s0faesp\",\"sender\":\"erd1g7rq9se53mmr5zq8hmz2jpttsgwx7wfx8dgmp9qj2mc7dnj99cys33f33x\",\"gasPrice\":1000000000,\"gasLimit\":50000,\"chainID\":\"T\",\"version\":2,\"options\":1,\"signature\":\"00f01ebc1ac6317df2492bd4ae692125b6b5c07448b5eae43ab3af6a56d529f6afc05f9acadb66a942f7014b06e15388c9e39eac19a5f405637a1947a2e45704\"}"
	require.Equal(t, expected, string(signedTxBs))

	dataBuilder := NewTxDataBuilder().
		Function("function").
		ArgBigInt(big.NewInt(15)).
		ArgInt64(14).
		ArgAddress("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8").
		ArgHexString("eeff00").
		ArgBytes([]byte("aa")).
		ArgBigInt(big.NewInt(0))
	args.Value = zeroString
	args.Data, _ = dataBuilder.ToDataBytes()
	builder2 := NewTxBuilder(&privateKey)
	tranWithData, _ := builder2.build(args)
	signedTxWithDataBs, _ := json.Marshal(tranWithData)
	expected = "{\"nonce\":3,\"value\":\"0\",\"receiver\":\"erd1lu89n007dkpy3s3za7wwguasttw6xlqqlvx20h0wwm778h5pfe8s0faesp\",\"sender\":\"erd1g7rq9se53mmr5zq8hmz2jpttsgwx7wfx8dgmp9qj2mc7dnj99cys33f33x\",\"gasPrice\":1000000000,\"gasLimit\":50000,\"data\":\"ZnVuY3Rpb25AMGZAMGVAYjJhMTE1NTVjZTUyMWU0OTQ0ZTA5YWIxNzU0OWQ4NWI0ODdkY2QyNmM4NGI1MDE3YTM5ZTMxYTM2NzA4ODliYUBlZWZmMDBANjE2MUAwMA==\",\"chainID\":\"T\",\"version\":2,\"options\":1,\"signature\":\"fdf46a866b934c262d91da08d847f99bb0fcaa1b3c4443983a95bb43f9f3cdb25985a262ae54b397baad1e33c3671a55d83d4061bcea0979eb9bdba414279600\"}"
	require.Equal(t, expected, string(signedTxWithDataBs))
}

func TestTxDataBuilder_AllGoodArguments(t *testing.T) {
	_, bytes, _ := bech32.DecodeToBase256("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	builder := NewTxDataBuilder().
		Function("function").
		ArgBigInt(big.NewInt(15)).
		ArgInt64(14).
		ArgAddress("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8").
		ArgHexString("eeff00").
		ArgBytes([]byte("aa")).
		ArgBigInt(big.NewInt(0))

	// function@0f@0e@b2a11555ce521e4944e09ab17549d85b487dcd26c84b5017a39e31a3670889ba@eeff00@6161@00
	expectedTxData := "function@" + hex.EncodeToString([]byte{15}) +
		"@" + hex.EncodeToString([]byte{14}) + "@" +
		hex.EncodeToString(bytes) + "@eeff00@" +
		hex.EncodeToString([]byte("aa")) + "@00"
	txData, err := builder.ToDataString()
	require.Nil(t, err)
	t.Log("txData : ", txData)
	txDataBytes, err := builder.ToDataBytes()
	require.Nil(t, err)
	expected := hex.EncodeToString([]byte(expectedTxData))
	require.Equal(t, expected, hex.EncodeToString(txDataBytes))
}

func TestTransfer(t *testing.T) {
	pk := "27d57eb22fc218b83e9ea2da55746d9318ba6b89cfa31b797e7296bf8a66e4f1"
	pk2, _ := hex.DecodeString("b6b77f7440d6c88b634bd95bf0da0a8660e781e88915f30319646170a57c4de1")
	toAddress, err := AddressFromSeed(hex.EncodeToString(pk2))
	require.Nil(t, err)
	args := ArgCreateTransaction{
		Nonce:    3,
		Value:    "10000000000000000", // decimal 18
		RcvAddr:  toAddress,
		GasPrice: 1000000000,
		GasLimit: 50000,
		ChainID:  "T",
		Version:  2,
		Options:  1,
	}
	signedTx, err := Transfer(args, pk)
	require.Nil(t, err)
	expected := "{\"nonce\":3,\"value\":\"10000000000000000\",\"receiver\":\"erd1lu89n007dkpy3s3za7wwguasttw6xlqqlvx20h0wwm778h5pfe8s0faesp\",\"sender\":\"erd1g7rq9se53mmr5zq8hmz2jpttsgwx7wfx8dgmp9qj2mc7dnj99cys33f33x\",\"gasPrice\":1000000000,\"gasLimit\":50000,\"chainID\":\"T\",\"version\":2,\"options\":1,\"signature\":\"00f01ebc1ac6317df2492bd4ae692125b6b5c07448b5eae43ab3af6a56d529f6afc05f9acadb66a942f7014b06e15388c9e39eac19a5f405637a1947a2e45704\"}"
	require.Equal(t, expected, signedTx)
}
