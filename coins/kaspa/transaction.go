package kaspa

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/model/externalapi"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/consensushashing"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/constants"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/subnetworks"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/transactionid"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/txscript"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/consensus/utils/utxo"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/dagconfig"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/domain/miningmanager/mempool"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/util"
	"github.com/okx/go-wallet-sdk/coins/kaspa/kaspad/util/txmass"
	"golang.org/x/crypto/blake2b"
	"strconv"
)

type TxInput struct {
	TxId       string `json:"txId"`
	Index      uint32 `json:"index"`
	Address    string `json:"address"`
	Amount     string `json:"amount"`
	PrivateKey string `json:"privateKey"`
}

type TxData struct {
	TxInputs      []*TxInput `json:"txInputs"`
	ToAddress     string     `json:"toAddress"`
	Amount        string     `json:"amount"`
	Fee           string     `json:"fee"`
	ChangeAddress string     `json:"changeAddress"`
	MinOutput     string     `json:"minOutput"`
}

type TransactionResult struct {
	Tx   string `json:"tx"`
	TxId string `json:"txId"`
}

type TransactionMessage struct {
	Transaction *Transaction `json:"transaction"`
	AllowOrphan bool         `json:"allowOrphan"`
}

type Transaction struct {
	Version      uint32               `json:"version"`
	Inputs       []*TransactionInput  `json:"inputs"`
	Outputs      []*TransactionOutput `json:"outputs"`
	LockTime     uint64               `json:"lockTime"`
	SubnetworkId string               `json:"subnetworkId"`
}

type TransactionInput struct {
	PreviousOutpoint *Outpoint `json:"previousOutpoint"`
	SignatureScript  string    `json:"signatureScript"`
	Sequence         uint64    `json:"sequence"`
	SigOpCount       uint32    `json:"sigOpCount"`
}

type Outpoint struct {
	TransactionId string `json:"transactionId"`
	Index         uint32 `json:"index"`
}

type TransactionOutput struct {
	Amount          uint64           `json:"amount"`
	ScriptPublicKey *ScriptPublicKey `json:"scriptPublicKey"`
}

type ScriptPublicKey struct {
	Version         uint32 `json:"version"`
	ScriptPublicKey string `json:"scriptPublicKey"`
}

func Transfer(txData *TxData) (string, error) {
	return TransferWithNetParams(txData, dagconfig.MainnetParams)
}

func TransferWithNetParams(txData *TxData, params dagconfig.Params) (string, error) {
	var totalInput uint64
	var inputs []*externalapi.DomainTransactionInput
	for _, input := range txData.TxInputs {
		txIdBytes, err := hex.DecodeString(input.TxId)
		if err != nil {
			return "", err
		}
		transactionID, err := transactionid.FromBytes(txIdBytes)
		if err != nil {
			return "", err
		}
		outpoint := externalapi.DomainOutpoint{
			TransactionID: *transactionID,
			Index:         input.Index,
		}
		fromAddress, err := util.DecodeAddress(input.Address, params.Prefix)
		if err != nil {
			return "", err
		}
		scriptPublicKey, err := txscript.PayToAddrScript(fromAddress)
		if err != nil {
			return "", err
		}
		inputs = append(inputs, &externalapi.DomainTransactionInput{
			PreviousOutpoint: outpoint,
			SigOpCount:       1,
			UTXOEntry: utxo.NewUTXOEntry(
				StrToUint64(input.Amount),
				scriptPublicKey,
				false,
				0,
			),
		})
		totalInput += StrToUint64(input.Amount)
	}

	var outputs []*externalapi.DomainTransactionOutput
	toAddress, err := util.DecodeAddress(txData.ToAddress, params.Prefix)
	if err != nil {
		return "", err
	}
	scriptPublicKey, err := txscript.PayToAddrScript(toAddress)
	if err != nil {
		return "", err
	}
	outputs = append(outputs, &externalapi.DomainTransactionOutput{
		Value:           StrToUint64(txData.Amount),
		ScriptPublicKey: scriptPublicKey,
	})

	// change
	toAmount := StrToUint64(txData.Amount)
	fee := StrToUint64(txData.Fee)
	minOutput := StrToUint64(txData.MinOutput)
	if minOutput == 0 {
		minOutput = uint64(546)
	}
	if totalInput >= toAmount+fee+minOutput {
		change := totalInput - (toAmount + fee)
		changeAddress, err := util.DecodeAddress(txData.ChangeAddress, params.Prefix)
		if err != nil {
			return "", err
		}
		changeScriptPublicKey, err := txscript.PayToAddrScript(changeAddress)
		if err != nil {
			return "", err
		}
		changeOutput := &externalapi.DomainTransactionOutput{
			Value:           change,
			ScriptPublicKey: changeScriptPublicKey,
		}
		outputs = append(outputs, changeOutput)
	}

	domainTransaction := &externalapi.DomainTransaction{
		Version:      constants.MaxTransactionVersion,
		Inputs:       inputs,
		Outputs:      outputs,
		LockTime:     0,
		SubnetworkID: subnetworks.SubnetworkIDNative,
		Gas:          0,
		Payload:      nil,
	}

	// sign
	for i, input := range domainTransaction.Inputs {
		prvKeyBytes, err := hex.DecodeString(txData.TxInputs[i].PrivateKey)
		if err != nil {
			return "", err
		}
		prvKey, _ := btcec.PrivKeyFromBytes(prvKeyBytes)

		signatureScript, err := txscript.SignatureScript(domainTransaction, i, consensushashing.SigHashAll, prvKey,
			&consensushashing.SighashReusedValues{})
		if err != nil {
			return "", err
		}
		input.SignatureScript = signatureScript
	}

	txMassCalculator := txmass.NewCalculator(params.MassPerTxByte, params.MassPerScriptPubKeyByte, params.MassPerSigOp)
	transactionMass := txMassCalculator.CalculateTransactionMass(domainTransaction)
	if transactionMass > mempool.MaximumStandardTransactionMass {
		return "", errors.New("exceeding the maximum mass allowed for transaction")
	}

	tx, err := serialize(domainTransaction)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(&TransactionResult{
		Tx:   tx,
		TxId: consensushashing.TransactionID(domainTransaction).String(),
	})
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func serialize(tx *externalapi.DomainTransaction) (string, error) {
	inputs := make([]*TransactionInput, len(tx.Inputs))
	for i, input := range tx.Inputs {
		inputs[i] = &TransactionInput{
			PreviousOutpoint: &Outpoint{
				TransactionId: input.PreviousOutpoint.TransactionID.String(),
				Index:         input.PreviousOutpoint.Index,
			},
			SignatureScript: hex.EncodeToString(input.SignatureScript),
			Sequence:        input.Sequence,
			SigOpCount:      uint32(input.SigOpCount),
		}
	}

	outputs := make([]*TransactionOutput, len(tx.Outputs))
	for i, output := range tx.Outputs {
		outputs[i] = &TransactionOutput{
			Amount: output.Value,
			ScriptPublicKey: &ScriptPublicKey{
				Version:         uint32(output.ScriptPublicKey.Version),
				ScriptPublicKey: hex.EncodeToString(output.ScriptPublicKey.Script),
			},
		}
	}

	transactionMessage := &TransactionMessage{
		Transaction: &Transaction{
			Version:      uint32(tx.Version),
			Inputs:       inputs,
			Outputs:      outputs,
			LockTime:     tx.LockTime,
			SubnetworkId: tx.SubnetworkID.String(),
		},
		AllowOrphan: false,
	}

	jsonBytes, err := json.Marshal(transactionMessage)
	if err != nil {
		return "", err
	}
	return string(jsonBytes), nil
}

func deserialize(txJSON string) (*externalapi.DomainTransaction, error) {
	var transactionMessage TransactionMessage
	err := json.Unmarshal([]byte(txJSON), &transactionMessage)
	if err != nil {
		return nil, err
	}

	tx := transactionMessage.Transaction
	if tx == nil {
		return nil, errors.New("transaction is nil")
	}

	// Parse inputs
	inputs := make([]*externalapi.DomainTransactionInput, len(tx.Inputs))
	for i, input := range tx.Inputs {
		if input.PreviousOutpoint == nil {
			return nil, errors.New("previous outpoint is nil")
		}

		// Parse transaction ID
		txIDBytes, err := hex.DecodeString(input.PreviousOutpoint.TransactionId)
		if err != nil {
			return nil, err
		}
		transactionID, err := transactionid.FromBytes(txIDBytes)
		if err != nil {
			return nil, err
		}

		// Parse signature script
		signatureScript, err := hex.DecodeString(input.SignatureScript)
		if err != nil {
			return nil, err
		}

		inputs[i] = &externalapi.DomainTransactionInput{
			PreviousOutpoint: externalapi.DomainOutpoint{
				TransactionID: *transactionID,
				Index:         input.PreviousOutpoint.Index,
			},
			SignatureScript: signatureScript,
			Sequence:        input.Sequence,
			SigOpCount:      uint8(input.SigOpCount),
		}
	}

	// Parse outputs
	outputs := make([]*externalapi.DomainTransactionOutput, len(tx.Outputs))
	for i, output := range tx.Outputs {
		if output.ScriptPublicKey == nil {
			return nil, errors.New("script public key is nil")
		}

		// Parse script public key
		scriptBytes, err := hex.DecodeString(output.ScriptPublicKey.ScriptPublicKey)
		if err != nil {
			return nil, err
		}

		outputs[i] = &externalapi.DomainTransactionOutput{
			Value: output.Amount,
			ScriptPublicKey: &externalapi.ScriptPublicKey{
				Script:  scriptBytes,
				Version: uint16(output.ScriptPublicKey.Version),
			},
		}
	}

	// Parse SubnetworkID
	subnetworkID, err := subnetworks.FromString(tx.SubnetworkId)
	if err != nil {
		return nil, err
	}

	return &externalapi.DomainTransaction{
		Version:      uint16(tx.Version),
		Inputs:       inputs,
		Outputs:      outputs,
		LockTime:     tx.LockTime,
		SubnetworkID: *subnetworkID,
		Gas:          0,
		Payload:      nil,
	}, nil
}

func CalTxHash(rawTx string) (string, error) {
	tx, err := deserialize(rawTx)
	if err != nil {
		return "", err
	}
	return consensushashing.TransactionID(tx).String(), nil
}

func StrToUint64(s string) uint64 {
	v, _ := strconv.ParseUint(s, 10, 64)
	return v
}

func SignMessage(message string, privateKey []byte) (string, error) {
	blake2b256, err := blake2b.New256([]byte("PersonalMessageSigningHash"))
	if err != nil {
		return "", err
	}

	blake2b256.Write([]byte(message))
	hash := blake2b256.Sum(nil)

	prvKey, _ := btcec.PrivKeyFromBytes(privateKey)

	signature, err := schnorr.Sign(prvKey, hash)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(signature.Serialize()), nil
}
