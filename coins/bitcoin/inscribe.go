package bitcoin

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcec/v2/schnorr"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

type InscriptionData struct {
	ContentType string `json:"contentType"`
	Body        []byte `json:"body"`
	RevealAddr  string `json:"revealAddr"`
}

type PrevOutput struct {
	TxId       string `json:"txId"`
	VOut       uint32 `json:"vOut"`
	Amount     int64  `json:"amount"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type InscriptionRequest struct {
	CommitTxPrevOutputList []*PrevOutput     `json:"commitTxPrevOutputList"`
	CommitFeeRate          int64             `json:"commitFeeRate"`
	RevealFeeRate          int64             `json:"revealFeeRate"`
	InscriptionDataList    []InscriptionData `json:"inscriptionDataList"`
	RevealOutValue         int64             `json:"revealOutValue"`
	ChangeAddress          string            `json:"changeAddress"`
	MinChangeValue         int64             `json:"minChangeValue"`
}

type inscriptionTxCtxData struct {
	PrivateKey              *btcec.PrivateKey
	InscriptionScript       []byte
	CommitTxAddress         string
	CommitTxAddressPkScript []byte
	ControlBlockWitness     []byte
	RevealTxPrevOutput      *wire.TxOut
}

type InscriptionBuilder struct {
	Network                   *chaincfg.Params
	CommitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrivateKeyList    []*btcec.PrivateKey
	InscriptionTxCtxDataList  []*inscriptionTxCtxData
	RevealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrevOutputList    []*PrevOutput
	RevealTx                  []*wire.MsgTx
	CommitTx                  *wire.MsgTx
	MustCommitTxFee           int64
	MustRevealTxFees          []int64
	CommitAddrs               []string
}

type InscribeTxs struct {
	CommitTx     string   `json:"commitTx"`
	RevealTxs    []string `json:"revealTxs"`
	CommitTxFee  int64    `json:"commitTxFee"`
	RevealTxFees []int64  `json:"revealTxFees"`
	CommitAddrs  []string `json:"commitAddrs"`
}

const (
	DefaultTxVersion      = 2
	DefaultSequenceNum    = 0xfffffffd
	DefaultRevealOutValue = int64(546)
	DefaultMinChangeValue = int64(546)

	MaxStandardTxWeight = 4000000 / 10
	WitnessScaleFactor  = 4
)

func NewInscriptionTool(network *chaincfg.Params, request *InscriptionRequest) (*InscriptionBuilder, error) {
	var commitTxPrivateKeyList []*btcec.PrivateKey
	for _, prevOutput := range request.CommitTxPrevOutputList {
		privateKeyWif, err := btcutil.DecodeWIF(prevOutput.PrivateKey)
		if err != nil {
			return nil, err
		}
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, privateKeyWif.PrivKey)
	}
	tool := &InscriptionBuilder{
		Network:                   network,
		CommitTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrivateKeyList:    commitTxPrivateKeyList,
		InscriptionTxCtxDataList:  make([]*inscriptionTxCtxData, len(request.InscriptionDataList)),
		RevealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrevOutputList:    request.CommitTxPrevOutputList,
	}
	return tool, tool.initTool(network, request)
}

func (builder *InscriptionBuilder) initTool(network *chaincfg.Params, request *InscriptionRequest) error {
	destinations := make([]string, len(request.InscriptionDataList))
	revealOutValue := DefaultRevealOutValue
	if request.RevealOutValue > 0 {
		revealOutValue = request.RevealOutValue
	}
	minChangeValue := DefaultMinChangeValue
	if request.MinChangeValue > 0 {
		minChangeValue = request.MinChangeValue
	}
	for i := 0; i < len(request.InscriptionDataList); i++ {
		inscriptionTxCtxData, err := newInscriptionTxCtxData(network, request, i)
		if err != nil {
			return err
		}
		builder.InscriptionTxCtxDataList[i] = inscriptionTxCtxData
		destinations[i] = request.InscriptionDataList[i].RevealAddr
	}
	totalRevealPrevOutputValue, err := builder.buildEmptyRevealTx(destinations, revealOutValue, request.RevealFeeRate)
	if err != nil {
		return err
	}
	err = builder.buildCommitTx(request.CommitTxPrevOutputList, request.ChangeAddress, totalRevealPrevOutputValue, request.CommitFeeRate, minChangeValue)
	if err != nil {
		return err
	}
	err = builder.signCommitTx()
	if err != nil {
		return errors.New("sign commit tx error")
	}
	err = builder.completeRevealTx()
	if err != nil {
		return err
	}
	return nil
}

func newInscriptionTxCtxData(network *chaincfg.Params, inscriptionRequest *InscriptionRequest, indexOfInscriptionDataList int) (*inscriptionTxCtxData, error) {
	privateKeyWif, err := btcutil.DecodeWIF(inscriptionRequest.CommitTxPrevOutputList[0].PrivateKey)
	if err != nil {
		return nil, err
	}
	privateKey := privateKeyWif.PrivKey

	inscriptionBuilder := txscript.NewScriptBuilder().
		AddData(schnorr.SerializePubKey(privateKey.PubKey())).
		AddOp(txscript.OP_CHECKSIG).
		AddOp(txscript.OP_FALSE).
		AddOp(txscript.OP_IF).
		AddData([]byte("ord")).
		AddOp(txscript.OP_DATA_1).
		AddOp(txscript.OP_DATA_1).
		AddData([]byte(inscriptionRequest.InscriptionDataList[indexOfInscriptionDataList].ContentType)).
		AddOp(txscript.OP_0)
	maxChunkSize := 520
	// use taproot to skip txscript.MaxScriptSize 10000
	bodySize := len(inscriptionRequest.InscriptionDataList[indexOfInscriptionDataList].Body)
	for i := 0; i < bodySize; i += maxChunkSize {
		end := i + maxChunkSize
		if end > bodySize {
			end = bodySize
		}

		inscriptionBuilder.AddFullData(inscriptionRequest.InscriptionDataList[indexOfInscriptionDataList].Body[i:end])
	}
	inscriptionScript, err := inscriptionBuilder.Script()
	if err != nil {
		return nil, err
	}
	inscriptionScript = append(inscriptionScript, txscript.OP_ENDIF)

	proof := &txscript.TapscriptProof{
		TapLeaf:  txscript.NewBaseTapLeaf(schnorr.SerializePubKey(privateKey.PubKey())),
		RootNode: txscript.NewBaseTapLeaf(inscriptionScript),
	}

	controlBlock := proof.ToControlBlock(privateKey.PubKey())
	controlBlockWitness, err := controlBlock.ToBytes()
	if err != nil {
		return nil, err
	}

	tapHash := proof.RootNode.TapHash()
	commitTxAddress, err := btcutil.NewAddressTaproot(schnorr.SerializePubKey(txscript.ComputeTaprootOutputKey(privateKey.PubKey(), tapHash[:])), network)
	if err != nil {
		return nil, err
	}
	commitTxAddressPkScript, err := txscript.PayToAddrScript(commitTxAddress)
	if err != nil {
		return nil, err
	}

	return &inscriptionTxCtxData{
		PrivateKey:              privateKey,
		InscriptionScript:       inscriptionScript,
		CommitTxAddress:         commitTxAddress.EncodeAddress(),
		CommitTxAddressPkScript: commitTxAddressPkScript,
		ControlBlockWitness:     controlBlockWitness,
	}, nil
}

func (builder *InscriptionBuilder) buildEmptyRevealTx(destination []string, revealOutValue, revealFeeRate int64) (int64, error) {
	addTxInTxOutIntoRevealTx := func(tx *wire.MsgTx, index int) error {
		in := wire.NewTxIn(&wire.OutPoint{Index: uint32(index)}, nil, nil)
		in.Sequence = DefaultSequenceNum
		tx.AddTxIn(in)
		scriptPubKey, err := AddrToPkScript(destination[index], builder.Network)
		if err != nil {
			return err
		}
		out := wire.NewTxOut(revealOutValue, scriptPubKey)
		tx.AddTxOut(out)
		return nil
	}

	totalPrevOutputValue := int64(0)
	total := len(builder.InscriptionTxCtxDataList)
	revealTx := make([]*wire.MsgTx, total)
	mustRevealTxFees := make([]int64, total)
	commitAddrs := make([]string, total)
	for i := 0; i < total; i++ {
		tx := wire.NewMsgTx(DefaultTxVersion)
		err := addTxInTxOutIntoRevealTx(tx, i)
		if err != nil {
			return 0, err
		}
		prevOutputValue := revealOutValue + int64(tx.SerializeSize())*revealFeeRate
		emptySignature := make([]byte, 64)
		emptyControlBlockWitness := make([]byte, 33)
		fee := (int64(wire.TxWitness{emptySignature, builder.InscriptionTxCtxDataList[i].InscriptionScript, emptyControlBlockWitness}.SerializeSize()+2+3) / 4) * revealFeeRate
		prevOutputValue += fee
		builder.InscriptionTxCtxDataList[i].RevealTxPrevOutput = &wire.TxOut{
			PkScript: builder.InscriptionTxCtxDataList[i].CommitTxAddressPkScript,
			Value:    prevOutputValue,
		}
		totalPrevOutputValue += prevOutputValue
		revealTx[i] = tx
		mustRevealTxFees[i] = int64(tx.SerializeSize())*revealFeeRate + fee
		commitAddrs[i] = builder.InscriptionTxCtxDataList[i].CommitTxAddress
	}
	builder.RevealTx = revealTx
	builder.MustRevealTxFees = mustRevealTxFees
	builder.CommitAddrs = commitAddrs

	return totalPrevOutputValue, nil
}

func (builder *InscriptionBuilder) buildCommitTx(commitTxPrevOutputList []*PrevOutput, changeAddress string, totalRevealPrevOutputValue, commitFeeRate int64, minChangeValue int64) error {
	totalSenderAmount := btcutil.Amount(0)
	tx := wire.NewMsgTx(DefaultTxVersion)
	changePkScript, err := AddrToPkScript(changeAddress, builder.Network)
	if err != nil {
		return err
	}
	for _, prevOutput := range commitTxPrevOutputList {
		txHash, err := chainhash.NewHashFromStr(prevOutput.TxId)
		if err != nil {
			return err
		}
		outPoint := wire.NewOutPoint(txHash, prevOutput.VOut)
		pkScript, err := AddrToPkScript(prevOutput.Address, builder.Network)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(prevOutput.Amount, pkScript)
		builder.CommitTxPrevOutputFetcher.AddPrevOut(*outPoint, txOut)

		in := wire.NewTxIn(outPoint, nil, nil)
		in.Sequence = DefaultSequenceNum
		tx.AddTxIn(in)

		totalSenderAmount += btcutil.Amount(prevOutput.Amount)
	}
	for i := range builder.InscriptionTxCtxDataList {
		tx.AddTxOut(builder.InscriptionTxCtxDataList[i].RevealTxPrevOutput)
	}

	tx.AddTxOut(wire.NewTxOut(0, changePkScript))

	txForEstimate := wire.NewMsgTx(DefaultTxVersion)
	txForEstimate.TxIn = tx.TxIn
	txForEstimate.TxOut = tx.TxOut
	if err = Sign(txForEstimate, builder.CommitTxPrivateKeyList, builder.CommitTxPrevOutputFetcher); err != nil {
		return err
	}

	fee := btcutil.Amount(GetTxVirtualSize(btcutil.NewTx(txForEstimate))) * btcutil.Amount(commitFeeRate)
	changeAmount := totalSenderAmount - btcutil.Amount(totalRevealPrevOutputValue) - fee
	if int64(changeAmount) >= minChangeValue {
		tx.TxOut[len(tx.TxOut)-1].Value = int64(changeAmount)
	} else {
		tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
		if changeAmount < 0 {
			txForEstimate.TxOut = txForEstimate.TxOut[:len(txForEstimate.TxOut)-1]
			feeWithoutChange := btcutil.Amount(GetTxVirtualSize(btcutil.NewTx(txForEstimate))) * btcutil.Amount(commitFeeRate)
			if totalSenderAmount-btcutil.Amount(totalRevealPrevOutputValue)-feeWithoutChange < 0 {
				builder.MustCommitTxFee = int64(fee)
				return errors.New("insufficient balance")
			}
		}
	}
	builder.CommitTx = tx
	return nil
}

func (builder *InscriptionBuilder) completeRevealTx() error {
	for i := range builder.InscriptionTxCtxDataList {
		builder.RevealTxPrevOutputFetcher.AddPrevOut(wire.OutPoint{
			Hash:  builder.CommitTx.TxHash(),
			Index: uint32(i),
		}, builder.InscriptionTxCtxDataList[i].RevealTxPrevOutput)
		builder.RevealTx[i].TxIn[0].PreviousOutPoint.Hash = builder.CommitTx.TxHash()
	}
	for i := range builder.InscriptionTxCtxDataList {
		revealTx := builder.RevealTx[i]
		witnessArray, err := txscript.CalcTapscriptSignaturehash(txscript.NewTxSigHashes(revealTx, builder.RevealTxPrevOutputFetcher),
			txscript.SigHashDefault, revealTx, 0, builder.RevealTxPrevOutputFetcher, txscript.NewBaseTapLeaf(builder.InscriptionTxCtxDataList[i].InscriptionScript))
		if err != nil {
			return err
		}
		signature, err := schnorr.Sign(builder.InscriptionTxCtxDataList[i].PrivateKey, witnessArray)
		if err != nil {
			return err
		}
		witness := wire.TxWitness{signature.Serialize(), builder.InscriptionTxCtxDataList[i].InscriptionScript, builder.InscriptionTxCtxDataList[i].ControlBlockWitness}
		builder.RevealTx[i].TxIn[0].Witness = witness
	}
	// check tx max tx wight
	for i, tx := range builder.RevealTx {
		revealWeight := GetTransactionWeight(btcutil.NewTx(tx))
		if revealWeight > MaxStandardTxWeight {
			return errors.New(fmt.Sprintf("reveal(index %d) transaction weight greater than %d (MAX_STANDARD_TX_WEIGHT): %d", i, MaxStandardTxWeight, revealWeight))
		}
	}
	return nil
}

func (builder *InscriptionBuilder) signCommitTx() error {
	return Sign(builder.CommitTx, builder.CommitTxPrivateKeyList, builder.CommitTxPrevOutputFetcher)
}

func Sign(tx *wire.MsgTx, privateKeys []*btcec.PrivateKey, prevOutFetcher *txscript.MultiPrevOutFetcher) error {
	for i, in := range tx.TxIn {
		prevOut := prevOutFetcher.FetchPrevOutput(in.PreviousOutPoint)
		txSigHashes := txscript.NewTxSigHashes(tx, prevOutFetcher)
		privKey := privateKeys[i]
		if txscript.IsPayToTaproot(prevOut.PkScript) {
			witness, err := txscript.TaprootWitnessSignature(tx, txSigHashes, i, prevOut.Value, prevOut.PkScript, txscript.SigHashDefault, privKey)
			if err != nil {
				return err
			}
			in.Witness = witness
		} else if txscript.IsPayToPubKeyHash(prevOut.PkScript) {
			sigScript, err := txscript.SignatureScript(tx, i, prevOut.PkScript, txscript.SigHashAll, privKey, true)
			if err != nil {
				return err
			}
			in.SignatureScript = sigScript
		} else {
			pubKeyBytes := privKey.PubKey().SerializeCompressed()
			script, err := PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
			if err != nil {
				return err
			}
			amount := prevOut.Value
			witness, err := txscript.WitnessSignature(tx, txSigHashes, i, amount, script, txscript.SigHashAll, privKey, true)
			if err != nil {
				return err
			}
			in.Witness = witness

			if txscript.IsPayToScriptHash(prevOut.PkScript) {
				redeemScript, err := PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
				if err != nil {
					return err
				}
				in.SignatureScript = append([]byte{byte(len(redeemScript))}, redeemScript...)
			}
		}
	}

	return nil
}

func GetTxHex(tx *wire.MsgTx) (string, error) {
	var buf bytes.Buffer
	if err := tx.Serialize(&buf); err != nil {
		return "", err
	}
	return hex.EncodeToString(buf.Bytes()), nil
}

func (builder *InscriptionBuilder) GetCommitTxHex() (string, error) {
	return GetTxHex(builder.CommitTx)
}

func (builder *InscriptionBuilder) GetRevealTxHexList() ([]string, error) {
	txHexList := make([]string, len(builder.RevealTx))
	for i := range builder.RevealTx {
		txHex, err := GetTxHex(builder.RevealTx[i])
		if err != nil {
			return nil, err
		}
		txHexList[i] = txHex
	}
	return txHexList, nil
}

func (builder *InscriptionBuilder) CalculateFee() (int64, []int64) {
	commitTxFee := int64(0)
	for _, in := range builder.CommitTx.TxIn {
		commitTxFee += builder.CommitTxPrevOutputFetcher.FetchPrevOutput(in.PreviousOutPoint).Value
	}
	for _, out := range builder.CommitTx.TxOut {
		commitTxFee -= out.Value
	}
	revealTxFees := make([]int64, 0)
	for _, tx := range builder.RevealTx {
		revealTxFee := int64(0)
		for i, in := range tx.TxIn {
			revealTxFee += builder.RevealTxPrevOutputFetcher.FetchPrevOutput(in.PreviousOutPoint).Value
			revealTxFee -= tx.TxOut[i].Value
			revealTxFees = append(revealTxFees, revealTxFee)
		}
	}
	return commitTxFee, revealTxFees
}

func Inscribe(network *chaincfg.Params, request *InscriptionRequest) (*InscribeTxs, error) {
	tool, err := NewInscriptionTool(network, request)
	if err != nil && err.Error() == "insufficient balance" {
		return &InscribeTxs{
			CommitTx:     "",
			RevealTxs:    []string{},
			CommitTxFee:  tool.MustCommitTxFee,
			RevealTxFees: tool.MustRevealTxFees,
			CommitAddrs:  tool.CommitAddrs,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	commitTx, err := tool.GetCommitTxHex()
	if err != nil {
		return nil, err
	}
	revealTxs, err := tool.GetRevealTxHexList()
	if err != nil {
		return nil, err
	}

	commitTxFee, revealTxFees := tool.CalculateFee()

	return &InscribeTxs{
		CommitTx:     commitTx,
		RevealTxs:    revealTxs,
		CommitTxFee:  commitTxFee,
		RevealTxFees: revealTxFees,
		CommitAddrs:  tool.CommitAddrs,
	}, nil
}

// GetTransactionWeight computes the value of the weight metric for a given
// transaction. Currently the weight metric is simply the sum of the
// transactions's serialized size without any witness data scaled
// proportionally by the WitnessScaleFactor, and the transaction's serialized
// size including any witness data.
func GetTransactionWeight(tx *btcutil.Tx) int64 {
	msgTx := tx.MsgTx()

	baseSize := msgTx.SerializeSizeStripped()
	totalSize := msgTx.SerializeSize()

	// (baseSize * 3) + totalSize
	return int64((baseSize * (WitnessScaleFactor - 1)) + totalSize)
}

// GetTxVirtualSize computes the virtual size of a given transaction. A
// transaction's virtual size is based off its weight, creating a discount for
// any witness data it contains, proportional to the current
// blockchain.WitnessScaleFactor value.
func GetTxVirtualSize(tx *btcutil.Tx) int64 {
	// vSize := (weight(tx) + 3) / 4
	//       := (((baseSize * 3) + totalSize) + 3) / 4
	// We add 3 here as a way to compute the ceiling of the prior arithmetic
	// to 4. The division by 4 creates a discount for wit witness data.
	return (GetTransactionWeight(tx) + (WitnessScaleFactor - 1)) / WitnessScaleFactor
}
