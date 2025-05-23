package flow

import (
	"encoding/hex"
	"fmt"
	"github.com/okx/go-wallet-sdk/coins/flow/core"
)

const DefaultAccountIndex = 0

func CreateNewAccountTx(publicKeyHex, payer, referenceBlockIDHex string, payerSequenceNumber, gasLimit uint64) *core.Transaction {
	payerAddress := core.HexToAddress(payer)
	pubKeyBytes, _ := hex.DecodeString(publicKeyHex)
	pubKeyParam := new(core.FlowPublicKey).FromBytes(pubKeyBytes).ToParam()
	tx := &core.Transaction{
		Script:           []byte(core.CreateAccountTpl),
		Arguments:        [][]byte{[]byte(pubKeyParam), []byte(core.ContractsParam)},
		ReferenceBlockID: core.HexToID(referenceBlockIDHex),
		GasLimit:         gasLimit,
		ProposalKey: core.ProposalKey{
			Address:        payerAddress,
			KeyIndex:       DefaultAccountIndex,
			SequenceNumber: payerSequenceNumber,
		},
		Payer:              payerAddress,
		Authorizers:        []core.Address{payerAddress},
		PayloadSignatures:  nil,
		EnvelopeSignatures: nil,
	}
	return tx
}

func CreateTransferFlowTx(amount float64, toAddr, payer, referenceBlockIDHex string, payerSequenceNumber, gasLimit uint64) *core.Transaction {
	payerAddress := core.HexToAddress(payer)
	amountParam := fmt.Sprintf(`{"type":"UFix64","value":"%f"}`, amount)
	addrParam := fmt.Sprintf(`{"type":"Address","value":"%s"}`, toAddr)
	tx := &core.Transaction{
		Script:           []byte(core.TransferFlowTpl),
		Arguments:        [][]byte{[]byte(amountParam), []byte(addrParam)},
		ReferenceBlockID: core.HexToID(referenceBlockIDHex),
		GasLimit:         gasLimit,
		ProposalKey: core.ProposalKey{
			Address:        payerAddress,
			KeyIndex:       DefaultAccountIndex,
			SequenceNumber: payerSequenceNumber,
		},
		Payer:              payerAddress,
		Authorizers:        []core.Address{payerAddress},
		PayloadSignatures:  nil,
		EnvelopeSignatures: nil,
	}
	return tx
}

func CreateTx(script, payer, referenceBlockIDHex string, args []string, payerSequenceNumber, gasLimit uint64) *core.Transaction {
	payerAddress := core.HexToAddress(payer)
	arguments := make([][]byte, len(args))
	for i, arg := range args {
		arguments[i] = []byte(arg)
	}
	tx := &core.Transaction{
		Script:           []byte(script),
		Arguments:        arguments,
		ReferenceBlockID: core.HexToID(referenceBlockIDHex),
		GasLimit:         gasLimit,
		ProposalKey: core.ProposalKey{
			Address:        payerAddress,
			KeyIndex:       DefaultAccountIndex,
			SequenceNumber: payerSequenceNumber,
		},
		Payer:              payerAddress,
		Authorizers:        []core.Address{payerAddress},
		PayloadSignatures:  nil,
		EnvelopeSignatures: nil,
	}
	return tx
}
