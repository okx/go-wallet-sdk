package oasis

import (
	"encoding/base64"
	"github.com/okx/go-wallet-sdk/crypto/bech32"
	"github.com/okx/go-wallet-sdk/crypto/cbor"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestCreateTx(t *testing.T) {
	pk := "a30a45ef8c019d22b7e8d18f11197677bff80ff4d2f23ab9ac14bdbac32c86e7baf40754ed3843e0464f814c3c605d8c36500cfb6892e2bd441839102f4200ed"
	chainId := "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535"
	toAddr := "oasis1qqx0wgxjwlw3jwatuwqj6582hdm9rjs4pcnvzz66"
	amount := big.NewInt(100000000)
	feeAmount := big.NewInt(0)
	gas := uint64(2000)
	nonce := uint64(7)

	_, toBytes, err := bech32.DecodeToBase256(toAddr)
	require.NoError(t, err)

	to := [21]byte{}
	copy(to[:], toBytes)

	transfer := Transfer{
		To:     to,
		Amount: amount.Bytes(),
	}
	tx := NewTx(nonce, gas, feeAmount, transfer)
	signedTx := SignTransaction(pk, chainId, tx)
	signedTxBytes, err := cbor.Marshal(signedTx)
	require.NoError(t, err)
	signedTxData := base64.StdEncoding.EncodeToString(signedTxBytes)
	exp := "onN1bnRydXN0ZWRfcmF3X3ZhbHVlWF+kY2ZlZaJjZ2FzGQfQZmFtb3VudEBkYm9keaJidG9VAAz3INJ33Rk7q+OBLVDqu3ZRyhUOZmFtb3VudEQF9eEAZW5vbmNlB2ZtZXRob2Rwc3Rha2luZy5UcmFuc2ZlcmlzaWduYXR1cmWianB1YmxpY19rZXlYILr0B1TtOEPgRk+BTDxgXYw2UAz7aJLivUQYORAvQgDtaXNpZ25hdHVyZVhAlZ3qt97MuDmu03vNRru/ZZLY+FguRxmcrBoaQAqQJkwHMrr7SO9hleEtvHQsx6kFHsvRhf04SGGStHG/ABJ7Cg=="
	require.Equal(t, exp, signedTxData)

}

func TestCreateTransferTx(t *testing.T) {
	pk := "a30a45ef8c019d22b7e8d18f11197677bff80ff4d2f23ab9ac14bdbac32c86e7baf40754ed3843e0464f814c3c605d8c36500cfb6892e2bd441839102f4200ed"
	chainId := "b11b369e0da5bb230b220127f5e7b242d385ef8c6f54906243f30af63c815535"
	toAddr := "oasis1qqx0wgxjwlw3jwatuwqj6582hdm9rjs4pcnvzz66"
	amount := big.NewInt(100000000)
	feeAmount := big.NewInt(0)
	gas := uint64(2000)
	nonce := uint64(8)
	tx := NewTransferTx(nonce, gas, feeAmount, toAddr, amount)
	signedTx := SignTransaction(pk, chainId, tx)
	signedTxBytes, err := cbor.Marshal(signedTx)
	require.NoError(t, err)
	signedTxData := base64.StdEncoding.EncodeToString(signedTxBytes)
	exp := "onN1bnRydXN0ZWRfcmF3X3ZhbHVlWF+kY2ZlZaJjZ2FzGQfQZmFtb3VudEBkYm9keaJidG9VAAz3INJ33Rk7q+OBLVDqu3ZRyhUOZmFtb3VudEQF9eEAZW5vbmNlCGZtZXRob2Rwc3Rha2luZy5UcmFuc2ZlcmlzaWduYXR1cmWianB1YmxpY19rZXlYILr0B1TtOEPgRk+BTDxgXYw2UAz7aJLivUQYORAvQgDtaXNpZ25hdHVyZVhAO2e0X7MlIl8e6x9U+iVsC0SELQckGsuLNwZyEnyio9z5yO2p39JEmeKTI6RGtquxTh9eOdAvPw20Bc5h6AXDDA=="
	require.Equal(t, exp, signedTxData)

}
