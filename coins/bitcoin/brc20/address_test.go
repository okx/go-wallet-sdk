package brc20

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestNewTapRootAddress(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"))
	address, err := NewTapRootAddress(privKey, &chaincfg.TestNet3Params)
	require.Nil(t, err)

	assert.Equal(t, address, "tb1pflq6z6mdduna235j3k3wn8tu6r39d4lc5celw9c7tfu6agp2yxvqfpyzqh")
	wif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	require.Nil(t, err)
	assert.Equal(t, wif.String(), "cNNWSuCmy77rPGnqm31JSdvFapf1ZVeDsUAmLXPNHmsdTEyw4eTj")
}

func TestNewTapRootAddressWithScript(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"))
	wif, _ := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	assert.Equal(t, "cNNWSuCmy77rPGnqm31JSdvFapf1ZVeDsUAmLXPNHmsdTEyw4eTj", wif.String())
	script, _ := CreateInscriptionScript(
		privKey,
		"text/plain;charset=utf-8",
		[]byte(fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "transfer", "ordi", "1")))
	address, err := NewTapRootAddressWithScript(privKey, script, &chaincfg.TestNet3Params)
	require.Nil(t, err)
	controlBlockBytes, _ := CreateControlBlock(privKey, script)
	assert.Equal(t, "1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37", hex.EncodeToString(privKey.Serialize()))
	assert.Equal(t, "201053e9ef0295d334b6bb22e20cc717eb1a16a546f692572c8830b4bc14c13676ac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800367b2270223a226272632d3230222c226f70223a227472616e73666572222c227469636b223a226f726469222c22616d74223a2231227d68", hex.EncodeToString(script))
	assert.Equal(t, "c11053e9ef0295d334b6bb22e20cc717eb1a16a546f692572c8830b4bc14c13676", hex.EncodeToString(controlBlockBytes))
	assert.Equal(t, "tb1p60xzvksp8jmsngwjgfaxu7kz5l3yn3yzswsj26np7w276ryz8vjqlj35nn", address)
}

func TestNewTapRootAddressWithScriptWithPubKey(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("604a9c5b807b8ef912e7a02321a66be93df4e13c4c0ef4e3ad6d8fc590e4ccd7"))
	pubKeyBytes := schnorr.SerializePubKey(privKey.PubKey())
	script, err := CreateInscriptionScriptWithPubKey(
		pubKeyBytes,
		"text/plain;charset=utf-8",
		[]byte(fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "transfer", "ordi", "1")))
	require.Nil(t, err)
	address := NewTapRootAddressWithScriptWithPubKey(pubKeyBytes, script, &chaincfg.MainNetParams)
	controlBlockBytes, _ := CreateControlBlock(privKey, script)

	assert.Equal(t, "604a9c5b807b8ef912e7a02321a66be93df4e13c4c0ef4e3ad6d8fc590e4ccd7", hex.EncodeToString(privKey.Serialize()))
	assert.Equal(t, "20462fa3f0eefce7d6fa0363a2f3b3a84dbde4039deab02eb254c28e49df4a711fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800367b2270223a226272632d3230222c226f70223a227472616e73666572222c227469636b223a226f726469222c22616d74223a2231227d68", hex.EncodeToString(script))
	assert.Equal(t, "c0462fa3f0eefce7d6fa0363a2f3b3a84dbde4039deab02eb254c28e49df4a711f", hex.EncodeToString(controlBlockBytes))
	assert.Equal(t, "bc1pmwus5lpxnnet6wcyqtevls07y7u8h5wun7q7p9jglk707y2czfns6hk0gp", address)
}
