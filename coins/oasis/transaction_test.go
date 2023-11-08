package oasis

import (
	"encoding/base64"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/bech32"
	"gitlab.okg.com/wallet-sign-core/go-parent-sdk/crypto/cbor"
	"math/big"
	"testing"
)

func TestCreateTx(t *testing.T) {
	pk := "fb181e94e95cc6bedd2da03e6c4aca9951053f3e9865945dbc8975a6afd217c3ad55bbb7c192b8ecfeb6ad18bbd7681c0923f472d5b0c212fbde33008005ad61"
	chainId := "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535"
	toAddr := "oasis1qqx0wgxjwlw3jwatuwqj6582hdm9rjs4pcnvzz66"
	amount := big.NewInt(100000000)
	feeAmount := big.NewInt(0)
	gas := uint64(2000)
	nonce := uint64(7)

	_, toBytes, err := bech32.DecodeToBase256(toAddr)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	to := [21]byte{}
	copy(to[:], toBytes)

	transfer := Transfer{
		To:     to,
		Amount: amount.Bytes(),
	}

	tx := NewTx(nonce, gas, feeAmount, transfer)
	signedTx := SignTransaction(pk, chainId, tx)

	signedTxBytes, err := cbor.Marshal(signedTx)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	t.Logf("signed tx : %s", base64.StdEncoding.EncodeToString(signedTxBytes))
}

func TestCreateTransferTx(t *testing.T) {
	pk := "fb181e94e95cc6bedd2da03e6c4aca9951053f3e9865945dbc8975a6afd217c3ad55bbb7c192b8ecfeb6ad18bbd7681c0923f472d5b0c212fbde33008005ad61"
	chainId := "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535"
	toAddr := "oasis1qqx0wgxjwlw3jwatuwqj6582hdm9rjs4pcnvzz66"
	amount := big.NewInt(100000000)
	feeAmount := big.NewInt(0)
	gas := uint64(2000)
	nonce := uint64(8)

	tx := NewTransferTx(nonce, gas, feeAmount, toAddr, amount)
	signedTx := SignTransaction(pk, chainId, tx)

	signedTxBytes, err := cbor.Marshal(signedTx)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	t.Logf("signed tx : %s", base64.StdEncoding.EncodeToString(signedTxBytes))
}
