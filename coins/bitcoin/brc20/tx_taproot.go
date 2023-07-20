package brc20

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/ecdsa"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/util"
	"strings"
)

var (
	maxChange, _ = btcutil.NewAmount(0.01)
)

type TransactionBuilder struct {
	inputs  []Input
	outputs []Output
	params  *chaincfg.Params
	tx      *wire.MsgTx
}

type Input struct {
	txId          string
	vOut          uint32
	privateKeyHex string
	address       string
	value         string
	inscription   *Inscription
}

type Inscription struct {
	contentType string
	body        []byte
}

type Output struct {
	address string
	amount  string
}

func NewTxBuild(version int32, params *chaincfg.Params) *TransactionBuilder {
	if params == nil {
		params = &chaincfg.MainNetParams
	}
	builder := &TransactionBuilder{
		inputs:  nil,
		outputs: nil,
		params:  params,
		tx:      &wire.MsgTx{Version: version, LockTime: 0},
	}
	return builder
}

func NewTxBuildV1(params *chaincfg.Params) *TransactionBuilder {
	return NewTxBuild(1, params)
}

func NewInscription(contentType string, body []byte) *Inscription {
	return &Inscription{contentType: contentType, body: body}
}

func (build *TransactionBuilder) AddInput(txId string, vOut uint32, privateKeyHex string, address string, value string, inscription *Inscription) {
	input := Input{txId: txId, vOut: vOut, privateKeyHex: privateKeyHex, address: address, value: value, inscription: inscription}
	build.inputs = append(build.inputs, input)
}

func (build *TransactionBuilder) AddOutput(address string, amount string) {
	output := Output{address: address, amount: amount}
	build.outputs = append(build.outputs, output)
}

func (build *TransactionBuilder) checkChangeValue() bool {
	inSum := int64(0)
	for _, input := range build.inputs {
		inSum += util.ConvertToBigInt(input.value).Int64()
	}

	outSum := int64(0)
	for _, output := range build.outputs {
		outSum += util.ConvertToBigInt(output.amount).Int64()
	}

	change := inSum - outSum
	return change < int64(maxChange.ToUnit(btcutil.AmountSatoshi))
}

func getAddressOutputScript(address string, params *chaincfg.Params) ([]byte, error) {
	decAddress, err := btcutil.DecodeAddress(address, params)
	if err != nil {
		return nil, err
	}
	return txscript.PayToAddrScript(decAddress)
}

func calSegWitHashNew(amount int64, tx *wire.MsgTx, txHash *txscript.TxSigHashes, i int, scriptCode []byte) []byte {
	var sigHash bytes.Buffer
	var bVersion [4]byte
	binary.LittleEndian.PutUint32(bVersion[:], uint32(tx.Version))
	sigHash.Write(bVersion[:])
	sigHash.Write(txHash.HashPrevOutsV0[:])
	sigHash.Write(txHash.HashSequenceV0[:])
	txIn := tx.TxIn[i]
	sigHash.Write(txIn.PreviousOutPoint.Hash[:])
	var bIndex [4]byte
	binary.LittleEndian.PutUint32(bIndex[:], txIn.PreviousOutPoint.Index)
	sigHash.Write(bIndex[:])
	sigHash.Write(scriptCode)
	var bAmount [8]byte
	binary.LittleEndian.PutUint64(bAmount[:], uint64(amount))
	sigHash.Write(bAmount[:])
	var bSequence [4]byte
	binary.LittleEndian.PutUint32(bSequence[:], txIn.Sequence)
	sigHash.Write(bSequence[:])
	sigHash.Write(txHash.HashOutputsV0[:])
	var bLockTime [4]byte
	binary.LittleEndian.PutUint32(bLockTime[:], tx.LockTime)
	sigHash.Write(bLockTime[:])
	var bHashType [4]byte
	binary.LittleEndian.PutUint32(bHashType[:], uint32(txscript.SigHashAll))
	sigHash.Write(bHashType[:])
	hash := chainhash.DoubleHashB(sigHash.Bytes())
	return hash
}

func CreateControlBlock(privateKey *btcec.PrivateKey, inscriptionScript []byte) ([]byte, error) {
	proof := &txscript.TapscriptProof{
		TapLeaf:  txscript.NewBaseTapLeaf(schnorr.SerializePubKey(privateKey.PubKey())),
		RootNode: txscript.NewBaseTapLeaf(inscriptionScript),
	}
	controlBlock := proof.ToControlBlock(privateKey.PubKey())
	return controlBlock.ToBytes()
}

func IsTaprootAddress(address string, params *chaincfg.Params) (bool, error) {
	addr, err := btcutil.DecodeAddress(address, params)
	if err != nil {
		return false, err
	}
	switch addr.(type) {
	case *btcutil.AddressTaproot:
		return true, nil
	default:
		return false, nil
	}
}

func (build *TransactionBuilder) Build() (string, error) {
	if len(build.inputs) == 0 || len(build.outputs) == 0 {
		return "", fmt.Errorf("input or output is empty")
	}

	if !build.checkChangeValue() {
		return "", fmt.Errorf("change amount exceed max")
	}

	tx := build.tx
	var prevPkScripts [][]byte
	prevOuts := txscript.NewMultiPrevOutFetcher(nil)
	for i := 0; i < len(build.inputs); i++ {
		input := build.inputs[i]
		hash, err := chainhash.NewHashFromStr(input.txId)
		if err != nil {
			return "", err
		}

		outPoint := wire.NewOutPoint(hash, input.vOut)
		txIn := wire.NewTxIn(outPoint, nil, nil)
		tx.TxIn = append(tx.TxIn, txIn)

		pkScript, err := getAddressOutputScript(input.address, build.params)
		if err != nil {
			return "", err
		}

		prevPkScripts = append(prevPkScripts, pkScript)

		prevOuts.AddPrevOut(*outPoint, &wire.TxOut{
			Value:    util.ConvertToBigInt(input.value).Int64(),
			PkScript: pkScript,
		})
	}

	for i := 0; i < len(build.outputs); i++ {
		output := build.outputs[i]
		script, err := getAddressOutputScript(output.address, build.params)
		if err != nil {
			return "", err
		}
		txOut := wire.NewTxOut(util.ConvertToBigInt(output.amount).Int64(), script)
		tx.TxOut = append(tx.TxOut, txOut)
	}

	for i := 0; i < len(build.inputs); i++ {
		input := build.inputs[i]
		privateBytes := util.RemoveZeroHex(input.privateKeyHex)
		prvKey, pubKey := btcec.PrivKeyFromBytes(privateBytes)

		oneIndex := strings.LastIndexByte(input.address, '1')
		isTaproot, err := IsTaprootAddress(input.address, build.params)
		if err != nil {
			return "", err
		}
		isSegWit := oneIndex > 1 && strings.ToLower(input.address[:oneIndex]) == build.params.Bech32HRPSegwit
		amount := util.ConvertToBigInt(input.value).Int64()
		// the address is taproot address
		if isTaproot {
			//create taproot script
			inscriptionScript, err := CreateInscriptionScript(prvKey, input.inscription.contentType, input.inscription.body)
			if err != nil {
				return "", err
			}
			controlBlockWitness, err := CreateControlBlock(prvKey, inscriptionScript)
			if err != nil {
				return "", err
			}
			hash, err := txscript.CalcTapscriptSignaturehash(txscript.NewTxSigHashes(tx, prevOuts), txscript.SigHashDefault, tx, i, prevOuts, txscript.NewBaseTapLeaf(inscriptionScript))
			if err != nil {
				return "", err
			}
			signature, err := schnorr.Sign(prvKey, hash)
			if err != nil {
				return "", err
			}
			tx.TxIn[i].Witness = wire.TxWitness{signature.Serialize(), inscriptionScript, controlBlockWitness}
		} else if isSegWit {
			// for  SegWit address
			sigHashes := txscript.NewTxSigHashes(tx, prevOuts)
			//p2pkh code
			scriptStr := fmt.Sprintf("1976a914%s88ac", hex.EncodeToString(btcutil.Hash160(pubKey.SerializeCompressed())))
			scriptCode, err := hex.DecodeString(scriptStr)
			if err != nil {
				return "", err
			}
			hash := calSegWitHashNew(amount, tx, sigHashes, i, scriptCode)
			signature := ecdsa.Sign(prvKey, hash)
			sign := append(signature.Serialize(), byte(txscript.SigHashAll))
			tx.TxIn[i].Witness = wire.TxWitness{sign, pubKey.SerializeCompressed()}
		} else {
			if strings.HasPrefix(input.address, "1") {
				// legacy address
				pkScript := prevPkScripts[i]
				script, err := txscript.SignatureScript(tx, i, pkScript, txscript.SigHashAll, prvKey, true)
				if err != nil {
					return "", err
				}
				tx.TxIn[i].SignatureScript = script
			} else {
				// P2SH address - Multi-signature address (not supported) && Segregated Witness compatible address
				sigHashes := txscript.NewTxSigHashes(tx, prevOuts)
				//P2PKH
				scriptStr := fmt.Sprintf("1976a914%s88ac", hex.EncodeToString(btcutil.Hash160(pubKey.SerializeCompressed())))
				scriptCode, err := hex.DecodeString(scriptStr)
				if err != nil {
					return "", err
				}
				hash := calSegWitHashNew(amount, tx, sigHashes, i, scriptCode)
				// sign
				signature := ecdsa.Sign(prvKey, hash)
				sign := append(signature.Serialize(), byte(txscript.SigHashAll))
				tx.TxIn[i].Witness = wire.TxWitness{sign, pubKey.SerializeCompressed()}

				ha := btcutil.Hash160(pubKey.SerializeCompressed())
				var redeemScript []byte
				redeemScript = append(redeemScript, 0x16)
				redeemScript = append(redeemScript, 0)
				redeemScript = append(redeemScript, 20)
				redeemScript = append(redeemScript, ha...)
				tx.TxIn[i].SignatureScript = redeemScript
			}
		}
	}
	buf := bytes.NewBuffer(make([]byte, 0, tx.SerializeSize()))
	if err := tx.Serialize(buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}
