package bitcoin

import (
	"bytes"
	"encoding/hex"
	"errors"

	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

func GenInput(txId string, vOut uint32, privateKeyHex string, redeemScript string, address string, amount int64) Input {
	return Input{
		txId:          txId,
		vOut:          vOut,
		privateKeyHex: privateKeyHex,
		redeemScript:  redeemScript,
		address:       address,
		amount:        amount,
	}
}

func (build *TransactionBuilder) signatureScriptFromPublicKey(publicKey *btcec.PublicKey) (
	signatureScript []byte, err error) {
	addPub, err := btcutil.NewAddressPubKey(publicKey.SerializeCompressed(), build.netParams)
	if err != nil {
		return
	}

	decodeAddress, err := btcutil.DecodeAddress(addPub.EncodeAddress(), build.netParams)
	if err != nil {
		return
	}

	signatureScript, err = txscript.PayToAddrScript(decodeAddress)

	return
}

func (build *TransactionBuilder) SingleBuild2() (string, error) {
	if len(build.inputs) == 0 || len(build.outputs) == 0 {
		return "", errors.New("invalid inputs or outputs")
	}

	tx := build.tx
	var scriptArray [][]byte

	for i := 0; i < len(build.inputs); i++ {
		input := build.inputs[i]

		var signatureScript []byte
		var err error

		if input.redeemScript != "" {
			signatureScript, err = hex.DecodeString(input.redeemScript)
		}

		if err != nil {
			return "", err
		}

		scriptArray = append(scriptArray, signatureScript)

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
		redeemScript := scriptArray[i]
		if len(redeemScript) != 0 {
			builder := txscript.NewScriptBuilder()

			signatureScript, err := builder.AddOp(txscript.OP_FALSE).AddData(redeemScript).Script()
			if err != nil {
				return "", err
			}

			tx.TxIn[i].SignatureScript = signatureScript
		}
	}

	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}

	return hex.EncodeToString(buf.Bytes()), nil
}

func multiSignAddOne(tx *wire.MsgTx, inputIndex int, privateKeyHex string) (err error) {
	scriptList, err := txscript.PushedData(tx.TxIn[inputIndex].SignatureScript) // [][]  sign+script
	if err != nil {
		return
	}

	redeemScript := scriptList[len(scriptList)-1]

	sigHash, err := txscript.CalcSignatureHash(redeemScript, txscript.SigHashAll, tx, inputIndex)
	if err != nil {
		return
	}

	wif, err := btcutil.DecodeWIF(privateKeyHex)
	if err != nil {
		return
	}

	ecKey := wif.PrivKey

	sign := ecdsa.Sign(ecKey, sigHash)
	sig2 := append(sign.Serialize(), byte(txscript.SigHashAll))

	builder := txscript.NewScriptBuilder()

	for i := 0; i < len(scriptList)-1; i++ {
		builder.AddData(scriptList[i])
	}

	scriptBuilder, err := builder.AddData(sig2).AddData(redeemScript).Script()
	if err != nil {
		return
	}

	tx.TxIn[inputIndex].SignatureScript = scriptBuilder

	return
}

func MultiSignBuildTx(tx *wire.MsgTx, inputs []Input, multiSignPriKeyList map[int][]string,
	privateKeyList map[int]string, network *chaincfg.Params) (err error) {
	if len(inputs) != len(tx.TxIn) {
		err = errors.New("invalid args")

		return
	}

	txIns := tx.TxIn

	prevOutFetcher := txscript.NewMultiPrevOutFetcher(nil)

	for i := 0; i < len(txIns); i++ {
		keys := multiSignPriKeyList[i]
		input := inputs[i]

		var signatureScript []byte

		if input.redeemScript != "" {
			signatureScript, err = hex.DecodeString(input.redeemScript)
		} else {
			signatureScript, err = AddrToPkScript(input.address, network)
		}

		var txHash *chainhash.Hash

		txHash, err = chainhash.NewHashFromStr(input.txId)
		if err != nil {
			return
		}

		outPoint := wire.NewOutPoint(txHash, input.vOut)

		txOut := wire.NewTxOut(input.amount, signatureScript)
		prevOutFetcher.AddPrevOut(*outPoint, txOut)

		if len(keys) == 0 {
			continue
		}

		for _, key := range keys {
			err = multiSignAddOne(tx, i, key)
			if err != nil {
				return
			}
		}
	}

	txSigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)

	for i, in := range tx.TxIn {
		prevOut := prevOutFetcher.FetchPrevOutput(in.PreviousOutPoint)

		privateKeyS, ok := privateKeyList[i]
		if !ok {
			continue
		}

		var privateKeyWif *btcutil.WIF

		privateKeyWif, err = btcutil.DecodeWIF(privateKeyS)
		if err != nil {
			return
		}

		err = SignTxInput1(privateKeyWif.PrivKey, tx, i, txSigHashes, prevOut.PkScript, prevOut.Value)
		if err != nil {
			return err
		}
	}

	return
}
