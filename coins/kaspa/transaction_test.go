package kaspa

import (
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTransfer(t *testing.T) {
	var txInputs []*TxInput
	txInputs = append(txInputs, &TxInput{
		TxId:       "120c5410cc4512f29da50a8befc67c1cfbf7bb4f594ef91c14741150d8dadd24",
		Index:      0,
		Address:    "kaspa:qrcnkrtrjptghtrntvyqkqafj06f9tamn0pnqvelmt2vmz68yp4gqj5lnal2h",
		Amount:     "900000",
		PrivateKey: "b827bb46d921bde498a535999d7554071045f02e4fdfdebb10b08583f1c6afbe",
	})
	txData := &TxData{
		TxInputs:      txInputs,
		ToAddress:     "kaspa:qqvxjssnw024e93vykhzd8d7t6dua2sx8ak4mj7xm8s9370yevxcv0jgl2xfj", // 443642da97444e52af9eb35e3d32d6270f47d255854b63299b29f21c1ded4c7c
		Amount:        "100000",
		Fee:           "10000",
		ChangeAddress: "kaspa:qrcnkrtrjptghtrntvyqkqafj06f9tamn0pnqvelmt2vmz68yp4gqj5lnal2h",
		MinOutput:     "546",
	}

	signedTx, err := Transfer(txData)
	if err != nil {
		// todo
	}
	require.NoError(t, err)
	res := &struct {
		TxId string `json:"txId"`
	}{}
	err = json.Unmarshal([]byte(signedTx), res)
	require.NoError(t, err)
	expected := "0dcbc57ae8b4be9c0769bdfffd54db09f7b36f048ba19d43f09c14b323a2b0d8"
	require.Equal(t, expected, res.TxId)
}

func TestSignMessage(t *testing.T) {
	privateKey, err := hex.DecodeString("b7e151628aed2a6abf7158809cf4f3c762e7160f38b4da56a784d9045190cfef")
	assert.Nil(t, err)
	signature, err := SignMessage("Hello Kaspa!", privateKey)
	assert.Nil(t, err)
	assert.Equal(t, "a106673fbb90b19f9ff55a0a40ec7ad56933ae0cf0170503886c59564044f93b1fe29233933790c70d4718e448cbe45ae212908b5f62d061feda048c16184964", signature)
}
