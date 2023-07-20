package brc20

import (
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/okx/go-wallet-sdk/util"
	"testing"
)

// tb1pp6v2zc4dfxrx0c6xmh340u9w958w2frf7t8m6w...
// cQosyLdyUyieNEmSmWRxV7PdCWMzJPm3iH4w4Xv8zk...
func TestNewTapRootAddress(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("//todo please replace your hex key"))
	address, err := NewTapRootAddress(privKey, &chaincfg.TestNet3Params)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(address)

	wif, _ := btcutil.NewWIF(privKey, &chaincfg.TestNet3Params, true)
	fmt.Println(wif)
}

// please replace your hex key
// 20462fa3f0eefce7d6fa0363a2f3b3a84dbde4039deab02eb254c28e49df4a711fac0063036f7264010118746578742f706c61696e3b636861727365743d7574662d3800367b2270223a226272632d3230222c226f70223a227472616e73666572222c227469636b223a226f726469222c22616d74223a2231227d68
// c0462fa3f0eefce7d6fa0363a2f3b3a84dbde4039deab02eb254c28e...
// tb1pmwus5lpxnnet6wcyqtevls07y7u8h5wun7q7p9jglk707y2czfns...
func TestNewTapRootAddressWithScript(t *testing.T) {
	privKey, _ := btcec.PrivKeyFromBytes(util.RemoveZeroHex("//todo please replace your hex key"))
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
