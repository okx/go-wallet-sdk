package near

import (
	"encoding/base64"
	"github.com/okx/go-wallet-sdk/coins/near/serialize"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestCreateTransaction(t *testing.T) {
	privateKey := "b9ec4d26ab5bec8df4314a9e3b8fc3f9c96d410b4cd13caa675018dcfc7916cceefbba85caaa14cb87b83314d5b86895f2d4b7633e29012e65bfb037c885c804"
	val := 0.222
	to := "ggasii.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)
	addr, err := PrivateKeyToAddr(privateKey)
	require.NoError(t, err)
	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
	require.NoError(t, err)
	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
	require.NoError(t, err)
	amount := decimal.NewFromFloat(val).Shift(24)
	ta, err := serialize.CreateTransfer(amount.String())
	require.NoError(t, err)
	tx.SetAction(ta)
	txData, err := tx.Serialize()
	require.NoError(t, err)
	txBase58 := base58.Encode(txData)
	sig, err := SignTransaction(txBase58, privateKey)
	require.NoError(t, err)
	stx, err := CreateSignedTransaction(tx, sig)
	require.NoError(t, err)
	stxData, err := stx.Serialize()
	require.NoError(t, err)
	b64Data := base64.StdEncoding.EncodeToString(stxData)
	expected := "QAAAAGVlZmJiYTg1Y2FhYTE0Y2I4N2I4MzMxNGQ1Yjg2ODk1ZjJkNGI3NjMzZTI5MDEyZTY1YmZiMDM3Yzg4NWM4MDQA7vu6hcqqFMuHuDMU1bholfLUt2M+KQEuZb+wN8iFyAQBAAAAAAAAAA4AAABnZ2FzaWkudGVzdG5ldLPinQIWXUUnnN9Qmtou83BpsylI4Fb+ZStWsef3s/kNAQAAAAMAAMAOl7HkpAIvAAAAAAAAACE/E/jQF9vlZSvRNf3Dnrr9Tm+gPB4s4wvE46LM18lgPtighyOfczJQMwhTJjFBe5xzBWbq3CJVhUYK21a9nQ0="
	assert.Equal(t, expected, b64Data)
}

func TestContactTransaction(t *testing.T) {
	privateKey := "b9ec4d26ab5bec8df4314a9e3b8fc3f9c96d410b4cd13caa675018dcfc7916cceefbba85caaa14cb87b83314d5b86895f2d4b7633e29012e65bfb037c885c804"
	val := 0.222
	to := "ft.examples.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)
	argsStr := `{"account_id": "serhii.testnet"}`
	gas := big.NewInt(1)
	addr := "ggasii.testnet"
	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
	require.NoError(t, err)
	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
	require.NoError(t, err)
	amount := decimal.NewFromFloat(val).Shift(24).BigInt()
	ta, err := serialize.CreateFunctionCall("storage_balance_of", []byte(argsStr), gas, amount)
	require.NoError(t, err)
	tx.SetAction(ta)
	txData, err := tx.Serialize()
	require.NoError(t, err)
	txHash := base58.Encode(txData)
	sig, err := SignTransaction(txHash, privateKey)
	require.NoError(t, err)
	stx, err := CreateSignedTransaction(tx, sig)
	require.NoError(t, err)
	stxData, err := stx.Serialize()
	require.NoError(t, err)
	b64Data := base64.StdEncoding.EncodeToString(stxData)
	expected := "DgAAAGdnYXNpaS50ZXN0bmV0AO77uoXKqhTLh7gzFNW4aJXy1LdjPikBLmW/sDfIhcgEAQAAAAAAAAATAAAAZnQuZXhhbXBsZXMudGVzdG5ldLPinQIWXUUnnN9Qmtou83BpsylI4Fb+ZStWsef3s/kNAQAAAAISAAAAc3RvcmFnZV9iYWxhbmNlX29mIAAAAHsiYWNjb3VudF9pZCI6ICJzZXJoaWkudGVzdG5ldCJ9AQAAAAAAAAAAAMAOl7HkpAIvAAAAAAAAACvEiv+vj1JDfHnrGZZ9vQlVvKCb2Bqsqe2KBB3ZhyM1YcWRR6WvjWVWpBmiXHt48xUf8ePtVcKdc0BNau8bJQM="
	assert.Equal(t, expected, b64Data)
}
func TestCalTxHash(t *testing.T) {
	signedTx := `QAAAAGQ3Mzg4OGEyNjE5Yzc3NjE3MzVmMjNjNzk4NTM2MTQ1ZGZhODdmOTMwNmI1ZjIxMjc1ZWI0YjFhN2JhOTcxYjkA1ziIomGcd2FzXyPHmFNhRd+of5MGtfISdetLGnupcbnjWwQAAAAAAEAAAAA4OWY5Nzc1ODU5ZWQzNDY3OGVhNDhlOWExYWViMjAyY2Q0YzI5ZGNlMTViZTA2NTJiOWY1MGUyMmEwYjY3ZWY3r4iB+lQhXiP818JF0LPDjkAFNvOeVJ/lAoe14WgEF6cBAAAAAwAA4ntBSX/LsDkAAAAAAAAAJv0PcmRmmTopCCBHfD2GNR3IKgmLzEL0K70jwXkjwXqbESEFCVaymK9VP/o9bFoPYeU+AFW92TyPy1fssMHaDQ==`
	hash, err := CalTxHash(signedTx, true)
	assert.NoError(t, err)
	expected := `GG4fSECycX149KgPp6YSNMUx7x26sTPUbDgFYxCNMMh7`
	assert.Equal(t, expected, hash)
}
