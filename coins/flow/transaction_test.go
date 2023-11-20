package flow

import (
	"github.com/okx/go-wallet-sdk/coins/flow/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreateNewAccountTx(t *testing.T) {
	privKey, pubKey := GenerateKeyPair()
	t.Logf("private key hex : %s", privKey)

	payerAddr := "0b65ef5c755c9117"
	payerSequenceNumber := uint64(9)

	referenceBlockIDHex := "9c45198cc1deda9087ec2f57607c1b4d6ae59e32a7f4619f47b05d8edb6fe21a"

	gasLimit := uint64(999)
	tx := CreateNewAccountTx(pubKey, payerAddr, referenceBlockIDHex, payerSequenceNumber, gasLimit)

	bytes, err := core.TransactionToHTTP(*tx)
	if err != nil {
		t.Errorf(err.Error())
		return
	}

	t.Logf(string(bytes))
}

func TestCreateTransferFlowTx(t *testing.T) {
	amount := float64(1)
	toAddr := "0xbd7d04ba8666b4d2"
	payer := "0x0b65ef5c755c9117"
	referenceBlockIDHex := "d77ceec957d4036f44bc33aaa03049b4cb96c75f0bbedafb09c62ac7e9b604d2"
	payerSequenceNumber := uint64(8)
	gasLimit := uint64(9999)
	tx := CreateTransferFlowTx(amount, toAddr, payer, referenceBlockIDHex, payerSequenceNumber, gasLimit)
	txBytes, err := core.TransactionToHTTP(*tx)
	require.Nil(t, err)
	t.Log("tx : ", string(txBytes))
}
