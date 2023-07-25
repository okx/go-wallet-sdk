package brc20

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/okx/go-wallet-sdk/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewTapRootAddress(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"))
	address, err := NewTapRootAddress(privKey, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(address)

	wif, _ := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	fmt.Println(wif)
	assert.Equal(t, wif.String(), "cNNWSuCmy77rPGnqm31JSdvFapf1ZVeDsUAmLXPNHmsdTEyw4eTj")
}

func TestNewTapRootAddressWithScript(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("1790962db820729606cd7b255ace1ac5ebb129ac8e9b2d8534d022194ab25b37"))
	wif, _ := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	fmt.Println(wif)
	script, _ := CreateInscriptionScript(
		privKey,
		"text/plain;charset=utf-8",
		[]byte(fmt.Sprintf(`{"p":"brc-20","op":"%s","tick":"%s","amt":"%s"}`, "transfer", "ordi", "1")))
	address, err := NewTapRootAddressWithScript(privKey, script, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	controlBlockBytes, _ := CreateControlBlock(privKey, script)

	fmt.Println(hex.EncodeToString(privKey.Serialize()))
	fmt.Println(hex.EncodeToString(script))
	fmt.Println(hex.EncodeToString(controlBlockBytes))
	fmt.Println(address)
}
