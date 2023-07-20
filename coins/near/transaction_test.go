package near

import (
	"encoding/base64"
	"encoding/hex"
	"github.com/okx/go-wallet-sdk/coins/near/serialize"
	"github.com/okx/go-wallet-sdk/crypto/base58"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"math/big"
	"testing"
)

func TestCreateTransaction(t *testing.T) {

	privateKey := "//todo please replace your hex key"
	val := 0.222
	to := "ggasii.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)

	addr, err := PrivateKeyToAddr(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
	if err != nil {
		t.Fatal(err)
	}

	amount := decimal.NewFromFloat(val).Shift(24)
	ta, err := serialize.CreateTransfer(amount.String())
	if err != nil {
		t.Fatal(err)
	}
	tx.SetAction(ta)
	txData, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	txBase58 := base58.Encode(txData)
	sig, err := SignTransaction(txBase58, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	stx, err := CreateSignedTransaction(tx, sig)
	if err != nil {
		t.Fatal(err)
	}
	stxData, err := stx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	b64Data := base64.StdEncoding.EncodeToString(stxData)
	t.Logf("transaction data : %s", b64Data)
	assert.Equal(t, "QAAAAGQyNWNmZGFlMGY5ODMyZTk4YmJkYzg3ZjNhMTU2YmI3NjVjZDk5NjRlMDA4NzhiZjY2ZGE3NDU5MTUzN2UwYTkA0lz9rg+YMumLvch/OhVrt2XNmWTgCHi/Ztp0WRU34KkBAAAAAAAAAA4AAABnZ2FzaWkudGVzdG5ldLPinQIWXUUnnN9Qmtou83BpsylI4Fb+ZStWsef3s/kNAQAAAAMAAMAOl7HkpAIvAAAAAAAAAOWsbqH7odK8g6Sw84lwbt6/xNNRCziw0mUpyvof/rcC9yCZ2ujjTVAeWIcKgSJ+CbzUmazccvBZ7YHjgdBstQc=", b64Data)
}

func TestContactTransaction(t *testing.T) {

	privateKey := "//todo please replace your hex key"
	val := 0.222
	to := "ft.examples.testnet"
	blockHash := "D7CPxgTXyRKYTSYuwAiRwDJH5RdHz7DwPt4EViptAW4L"
	nonce := int64(1)

	argsStr := `{"account_id": "serhii.testnet"}`
	gas := big.NewInt(1)

	addr := "ggasii.testnet"

	publicKeyHex, err := PrivateKeyToPublicKeyHex(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	tx, err := CreateTransaction(addr, to, publicKeyHex, blockHash, nonce)
	if err != nil {
		t.Fatal(err)
	}

	amount := decimal.NewFromFloat(val).Shift(24).BigInt()
	println(amount.String())
	ta, err := serialize.CreateFunctionCall("storage_balance_of", []byte(argsStr), gas, amount)
	if err != nil {
		t.Fatal(err)
	}

	tx.SetAction(ta)
	txData, err := tx.Serialize()
	if err != nil {
		t.Fatal(err)
	}

	txHash := base58.Encode(txData)
	sig, err := SignTransaction(txHash, privateKey)
	if err != nil {
		t.Fatal(err)
	}
	stx, err := CreateSignedTransaction(tx, sig)
	if err != nil {
		t.Fatal(err)
	}
	stxData, err := stx.Serialize()
	if err != nil {
		t.Fatal(err)
	}
	println(hex.EncodeToString(stxData))

	b64Data := base64.StdEncoding.EncodeToString(stxData)
	t.Logf("transaction data : %s", b64Data)

}
