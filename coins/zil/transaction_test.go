package zil

import (
	json2 "encoding/json"
	"github.com/stretchr/testify/require"
	"math/big"
	"testing"
)

func TestSignTransaction(t *testing.T) {
	privateKey := "c0dc46b9f9d6ef1c88dff3f1e82adc61cb11d77ab76a8d66338f14c2722cb4d8"
	to := "zil1fwh4ltdguhde9s7nysnp33d5wye6uqpugufkz7"
	gasPrice := "2000000000"
	amount := big.NewInt(10000000000)
	gasLimit := big.NewInt(50)
	nonce := 2
	chainId := 333
	tx := CreateTransferTransaction(to, gasPrice, amount, gasLimit, nonce, chainId)
	err := SignTransaction(privateKey, tx)
	require.NoError(t, err)
	payload := tx.ToTransactionPayload()
	json, err := json2.Marshal(payload)
	require.NoError(t, err)
	t.Log(string(json))
}
