package flow

import (
	"github.com/okx/go-wallet-sdk/coins/flow/core"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGenerateKeyPair(t *testing.T) {
	privKey, pubKey := GenerateKeyPair()
	t.Log("privKey : ", privKey)
	t.Log("pubKey : ", pubKey)
}

func TestSignTx(t *testing.T) {
	_, pubKey := GenerateKeyPair()
	payerAddr := "0b65ef5c755c9117"
	payerSequenceNumber := uint64(12)
	referenceBlockIDHex := "d83f8a740f774665016cbc34221fa1b1a0f430fe938297e2265afeee84bd19f4"
	gasLimit := uint64(9999)
	tx := CreateNewAccountTx(pubKey, payerAddr, referenceBlockIDHex, payerSequenceNumber, gasLimit)
	signPrivKeyHex := "986b514eec3705d809868611722574bba6d7829cb557dcbfea18b47b203321ed"
	signAddr := "0x0b65ef5c755c9117"
	err := SignTx(signAddr, signPrivKeyHex, tx)
	require.Nil(t, err)
	txBytes, err := core.TransactionToHTTP(*tx)
	if err != nil {
		// todo
	}
	require.Nil(t, err)
	t.Log("tx : ", string(txBytes))
}

func TestSignTransferTx(t *testing.T) {
	amount := float64(1)
	toAddr := "0x0b65ef5c755c9117"
	payer := "0x7a1fa92ef1acbe3c"
	referenceBlockIDHex := "5e62a0eb9505be3499fc321df3afc705f5483fd4409b940df3cabb66988117ce"
	payerSequenceNumber := uint64(2)
	gasLimit := uint64(9999)
	tx := CreateTransferFlowTx(amount, toAddr, payer, referenceBlockIDHex, payerSequenceNumber, gasLimit)
	signPrivKeyHex := "3eabec25b247b2f2e83dee958d77732a1a6a848383ac0dd9d4b0e97c18ee7259"
	signAddr := "0x7a1fa92ef1acbe3c"
	err := SignTx(signAddr, signPrivKeyHex, tx)
	require.Nil(t, err)
	txBytes, err := core.TransactionToHTTP(*tx)
	if err != nil {
		// todo
	}
	require.Nil(t, err)
	t.Log("tx : ", string(txBytes))
}
