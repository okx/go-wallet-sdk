package bitcoin

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/util"
)

type TransactionBuilder struct {
	inputs    []Input
	outputs   []Output
	netParams *chaincfg.Params
	tx        *wire.MsgTx
}

type Input struct {
	txId          string
	vOut          uint32
	privateKeyHex string
	redeemScript  string
}

type Output struct {
	address string
	amount  int64
}

func NewTxBuild(version int32, netParams *chaincfg.Params) *TransactionBuilder {
	if netParams == nil {
		netParams = &chaincfg.MainNetParams
	}
	builder := &TransactionBuilder{
		inputs:    nil,
		outputs:   nil,
		netParams: netParams,
		tx:        &wire.MsgTx{Version: version, LockTime: 0},
	}
	return builder
}

func (build *TransactionBuilder) AddInput(txId string, vOut uint32, privateKeyHex string, redeemScript string) {
	input := Input{txId: txId, vOut: vOut, privateKeyHex: privateKeyHex, redeemScript: redeemScript}
	build.inputs = append(build.inputs, input)
}

func (build *TransactionBuilder) AddOutput(address string, amount int64) {
	output := Output{address: address, amount: amount}
	build.outputs = append(build.outputs, output)
}

func (build *TransactionBuilder) SingleBuild() (string, error) {
	if len(build.inputs) == 0 || len(build.outputs) == 0 {
		return "", errors.New("invalid inputs or outputs")
	}

	tx := build.tx
	var scriptArray [][]byte
	var ecKeyArray []btcec.PrivateKey
	for i := 0; i < len(build.inputs); i++ {
		input := build.inputs[i]
		privateBytes, err := hex.DecodeString(input.privateKeyHex)
		if err != nil {
			return "", err
		}
		prvKey, publicKey := btcec.PrivKeyFromBytes(privateBytes)
		var signatureScript []byte
		if input.redeemScript == "" {
			addPub, err := btcutil.NewAddressPubKey(publicKey.SerializeCompressed(), &chaincfg.MainNetParams)
			if err != nil {
				return "", err
			}
			decodeAddress, err := btcutil.DecodeAddress(addPub.EncodeAddress(), &chaincfg.MainNetParams)
			if err != nil {
				return "", err
			}
			signatureScript, err = txscript.PayToAddrScript(decodeAddress)
			if err != nil {
				return "", err
			}
		} else {
			signatureScript, err = hex.DecodeString(input.redeemScript)
			if err != nil {
				return "", err
			}
		}
		scriptArray = append(scriptArray, signatureScript)
		ecKeyArray = append(ecKeyArray, *prvKey)

		hash, err := chainhash.NewHashFromStr(input.txId)
		if err != nil {
			return "", err
		}
		outPoint := wire.NewOutPoint(hash, input.vOut)
		txIn := wire.NewTxIn(outPoint, signatureScript, nil)
		tx.TxIn = append(tx.TxIn, txIn)
	}

	for i := 0; i < len(build.outputs); i++ {
		output := build.outputs[i]
		address, err := btcutil.DecodeAddress(output.address, build.netParams)
		if err != nil {
			return "", err
		}
		script, err := txscript.PayToAddrScript(address)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(output.amount, script)
		tx.TxOut = append(tx.TxOut, txOut)
	}

	for i := 0; i < len(build.inputs); i++ {
		ecKey := ecKeyArray[i]
		redeemScript := scriptArray[i]
		sigHash, err := txscript.CalcSignatureHash(redeemScript, txscript.SigHashAll, tx, i)
		if err != nil {
			return "", err
		}
		sign := ecdsa.Sign(&ecKey, sigHash)
		builder := txscript.NewScriptBuilder()
		if build.inputs[i].redeemScript != "" { // for multiple-sign
			builder.AddOp(txscript.OP_FALSE)
		} else {
			redeemScript = ecKey.PubKey().SerializeCompressed()
		}
		sig1 := append(sign.Serialize(), byte(txscript.SigHashAll))
		scriptBuilder, err := builder.AddData(sig1).AddData(redeemScript).Script()
		if err != nil {
			return "", err
		}
		tx.TxIn[i].SignatureScript = scriptBuilder
	}
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

// Second signature
func NewTxFromHex(txHex string) (*wire.MsgTx, error) {
	txBytes, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	reader := bytes.NewReader(txBytes)
	tx := &wire.MsgTx{}
	err = tx.Deserialize(reader)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func MultiSignBuild(tx *wire.MsgTx, priKeyList []string) (string, error) {
	txIns := tx.TxIn
	if len(txIns) != len(priKeyList) {
		return "", errors.New("invalid prikey list")
	}
	for i := 0; i < len(txIns); i++ {
		txIn := txIns[i]
		scriptList, err := txscript.PushedData(txIn.SignatureScript) // [][]  sign+script
		if err != nil {
			return "", err
		}
		redeemScript := scriptList[len(scriptList)-1]
		privateBytes, err := hex.DecodeString(priKeyList[i])
		if err != nil {
			return "", err
		}
		ecKey, _ := btcec.PrivKeyFromBytes(privateBytes)
		sigHash, err := txscript.CalcSignatureHash(redeemScript, txscript.SigHashAll, tx, i)
		if err != nil {
			return "", err
		}
		sign := ecdsa.Sign(ecKey, sigHash)
		sig2 := append(sign.Serialize(), byte(txscript.SigHashAll))

		builder := txscript.NewScriptBuilder()
		for i := 0; i < len(scriptList)-1; i++ {
			builder.AddData(scriptList[i])
		}
		scriptBuilder, err := builder.AddData(sig2).AddData(redeemScript).Script()
		if err != nil {
			return "", err
		}
		tx.TxIn[i].SignatureScript = scriptBuilder
	}
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (build *TransactionBuilder) UnSignedTx(pubKeyMap map[int]string) (string, map[int]string, error) {
	if len(build.inputs) == 0 || len(build.outputs) == 0 {
		return "", nil, fmt.Errorf("input or output miss")
	}

	tx := build.tx
	var scriptArray [][]byte
	for i := 0; i < len(build.inputs); i++ {
		input := build.inputs[i]
		var signatureScript []byte
		addPub, err := btcutil.NewAddressPubKey(util.RemoveZeroHex(pubKeyMap[i]), &chaincfg.MainNetParams)
		if err != nil {
			return "", nil, err
		}
		decodeAddress, err := btcutil.DecodeAddress(addPub.EncodeAddress(), &chaincfg.MainNetParams)
		if err != nil {
			return "", nil, err
		}
		signatureScript, err = txscript.PayToAddrScript(decodeAddress)
		if err != nil {
			return "", nil, err
		}
		scriptArray = append(scriptArray, signatureScript)

		hash, err := chainhash.NewHashFromStr(input.txId)
		if err != nil {
			return "", nil, err
		}
		outPoint := wire.NewOutPoint(hash, input.vOut)
		txIn := wire.NewTxIn(outPoint, signatureScript, nil)
		tx.TxIn = append(tx.TxIn, txIn)
	}

	for i := 0; i < len(build.outputs); i++ {
		output := build.outputs[i]
		address, err := btcutil.DecodeAddress(output.address, build.netParams)
		if err != nil {
			return "", nil, err
		}
		script, err := txscript.PayToAddrScript(address)
		if err != nil {
			return "", nil, err
		}
		txOut := wire.NewTxOut(output.amount, script)
		tx.TxOut = append(tx.TxOut, txOut)
	}

	hashes := make(map[int]string)
	for i := 0; i < len(build.inputs); i++ {
		redeemScript := scriptArray[i]
		sigHash, err := txscript.CalcSignatureHash(redeemScript, txscript.SigHashAll, tx, i)
		if err != nil {
			return "", nil, err
		}
		hashes[i] = hex.EncodeToString(sigHash)

		builder := txscript.NewScriptBuilder()
		sig1 := append(make([]byte, 70), byte(txscript.SigHashAll))
		scriptBuilder, err := builder.AddData(sig1).AddData(redeemScript).Script()
		if err != nil {
			return "", nil, err
		}
		tx.TxIn[i].SignatureScript = scriptBuilder
	}
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err := tx.Serialize(buf)
	if err != nil {
		return "", nil, err
	}
	return hex.EncodeToString(buf.Bytes()), hashes, nil
}

func SignTx(raw string, pubKeyMap map[int]string, signatureMap map[int]string) (string, error) {
	txBytes, err := hex.DecodeString(raw)
	if err != nil {
		return "", err
	}
	reader := bytes.NewReader(txBytes)
	tx := &wire.MsgTx{}
	err = tx.Deserialize(reader)
	if err != nil {
		return "", err
	}

	if len(tx.TxIn) != len(signatureMap) {
		return "", fmt.Errorf("signature miss")
	}

	for i := 0; i < len(tx.TxIn); i++ {
		builder := txscript.NewScriptBuilder()
		publicKey, err := btcec.ParsePubKey(util.RemoveZeroHex(pubKeyMap[i]))
		if err != nil {
			return "", err
		}
		redeemScript := publicKey.SerializeCompressed()
		sig1 := append(util.RemoveZeroHex(signatureMap[i]), byte(txscript.SigHashAll))
		scriptBuilder, err := builder.AddData(sig1).AddData(redeemScript).Script()
		if err != nil {
			return "", err
		}
		tx.TxIn[i].SignatureScript = scriptBuilder
	}
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	err = tx.Serialize(buf)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}
