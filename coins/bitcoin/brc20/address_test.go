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
	t.Log(address)
	wif, err := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	require.Nil(t, err)
	t.Log(wif)
	assert.Equal(t, wif.String(), "cNNWSuCmy77rPGnqm31JSdvFapf1ZVeDsUAmLXPNHmsdTEyw4eTj")
}

func TestNewTapRootAddressWithScript(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"))
	wif, _ := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	t.Log(wif)
	script, _ := CreateInscriptionScript(
		privKey,
		"text/plain;charset=utf-8",
		[]byte(fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "transfer", "ordi", "1")))
	address, err := NewTapRootAddressWithScript(privKey, script, &chaincfg.TestNet3Params)
	require.Nil(t, err)
	controlBlockBytes, _ := CreateControlBlock(privKey, script)
	t.Log("privKey : ", hex.EncodeToString(privKey.Serialize()))
	t.Log("script : ", hex.EncodeToString(script))
	t.Log("controlBlockBytes : ", hex.EncodeToString(controlBlockBytes))
	t.Log("address : ", address)
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
	controlBlockBytes, err := CreateControlBlock(privKey, script)
	require.Nil(t, err)
	t.Log("privKey : ", hex.EncodeToString(privKey.Serialize()))
	t.Log("script : ", hex.EncodeToString(script))
	t.Log("controlBlockBytes : ", hex.EncodeToString(controlBlockBytes))
	t.Log("address : ", address)
}
