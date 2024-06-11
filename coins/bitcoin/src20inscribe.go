package bitcoin

import (
	"crypto/rc4"
	"encoding/hex"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/util"
	"strings"
)

const PART_LEN = 31
const Five = 5 //It has to be multiplied by five

type Src20InscriptionRequest struct {
	CommitTxPrevOutputList PrevOutputs      `json:"commitTxPrevOutputList"`
	CommitFeeRate          int64            `json:"commitFeeRate"`
	InscriptionData        *InscriptionData `json:"inscriptionDataList"`
	RevealOutValue         int64            `json:"revealOutValue"`
	Address                string           `json:"address"`
	DustSize               int64            `json:"dustSize"`
}

type Src20InscriptionTool struct {
	Network                   *chaincfg.Params
	CommitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrivateKeyList    []*btcec.PrivateKey
	RevealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrevOutputList    []*PrevOutput
	CommitTx                  *wire.MsgTx
	MustCommitTxFee           int64
	CommitAddrs               []string
}

func NewSrc20InscriptionTool(network *chaincfg.Params, request *Src20InscriptionRequest) (*Src20InscriptionTool, error) {
	var commitTxPrivateKeyList []*btcec.PrivateKey
	for _, prevOutput := range request.CommitTxPrevOutputList {
		privateKeyWif, err := btcutil.DecodeWIF(prevOutput.PrivateKey)
		if err != nil {
			return nil, err
		}
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, privateKeyWif.PrivKey)
	}
	tool := &Src20InscriptionTool{
		Network:                   network,
		CommitTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrivateKeyList:    commitTxPrivateKeyList,
		RevealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrevOutputList:    request.CommitTxPrevOutputList,
		CommitAddrs:               make([]string, 0),
	}
	return tool, tool._initTool(network, request)
}

func (tool *Src20InscriptionTool) _initTool(network *chaincfg.Params, request *Src20InscriptionRequest) error {
	revealOutValue := DefaultRevealOutValue
	if request.RevealOutValue > 0 {
		revealOutValue = request.RevealOutValue
	}
	dustSize := DefaultMinChangeValue
	if request.DustSize > 0 {
		dustSize = request.DustSize
	}
	err := tool.buildCommitTx(request.CommitTxPrevOutputList, request.InscriptionData, request.Address, revealOutValue, request.CommitFeeRate, dustSize)
	if err != nil {
		return err
	}
	err = tool.signCommitTx()
	if err != nil {
		return errors.New("sign commit tx error")
	}
	return err
}

func (tool *Src20InscriptionTool) buildCommitTx(commitTxPrevOutputList PrevOutputs, inscriptionData *InscriptionData, changeAddress string, revealOutValue, commitFeeRate int64, minChangeValue int64) error {
	bf := make([]byte, 0, len(inscriptionData.ContentType)+len(inscriptionData.Body))
	bf = append(bf, inscriptionData.ContentType...)
	bf = append(bf, inscriptionData.Body...)
	for bf[len(bf)-1] == 0 {
		bf = bf[0 : len(bf)-1]
	}
	l := 2 + len(bf)
	total := l
	if l%62 != 0 {
		total = (l + 62 - l%62)
	}
	data := make([]byte, total)
	data[0], data[1] = byte(len(bf)/256), byte(len(bf)%256)
	copy(data[2:], bf)
	key, err := util.DecodeHexString(commitTxPrevOutputList[0].TxId)
	if err != nil {
		return err
	}
	c, err := rc4.NewCipher(key)
	if err != nil {
		return err
	}
	buf := make([]byte, len(data))
	c.XORKeyStream(buf, data)
	totalSenderAmount := btcutil.Amount(0)

	tx := wire.NewMsgTx(DefaultTxVersion)
	revealAddrPkScript, err := AddrToPkScript(inscriptionData.RevealAddr, tool.Network)
	if err != nil {
		return err
	}
	tx.AddTxOut(wire.NewTxOut(revealOutValue, revealAddrPkScript))
	totalRevealPrevOutputValue := revealOutValue
	minF := func(b []byte) int {
		if len(b) > PART_LEN {
			return PART_LEN
		}
		return len(b)
	}
	for len(buf) > 0 {
		buf1 := buf[0:minF(buf)]
		first := hex.EncodeToString(buf1)
		if len(first) < 62 {
			first = first + strings.Repeat("0", 62-len(first))
		}
		buf = buf[len(buf1):]
		buf2 := buf[0:minF(buf)]
		second := hex.EncodeToString(buf2)
		if len(second) < 62 {
			second = second + strings.Repeat("0", 62-len(second))
		}
		buf = buf[len(buf2):]
		pubkeys := []string{"03" + first + "00", "02" + second + "00", "020202020202020202020202020202020202020202020202020202020202020202"}
		builder := txscript.NewScriptBuilder().AddOp(txscript.OP_1)
		for _, v := range pubkeys {
			key, err := util.DecodeHexString(v)
			if err != nil {
				return err
			}
			builder.AddData(key)
		}
		builder.AddOp(txscript.OP_3).AddOp(txscript.OP_CHECKMULTISIG)
		pkScript, err := builder.Script()
		if err != nil {
			return err
		}
		tx.AddTxOut(wire.NewTxOut(revealOutValue, pkScript))
		totalRevealPrevOutputValue += revealOutValue
	}
	for _, prevOutput := range commitTxPrevOutputList {
		txHash, err := chainhash.NewHashFromStr(prevOutput.TxId)
		if err != nil {
			return err
		}
		outPoint := wire.NewOutPoint(txHash, prevOutput.VOut)
		pkScript, err := AddrToPkScript(prevOutput.Address, tool.Network)
		if err != nil {
			return err
		}
		txOut := wire.NewTxOut(prevOutput.Amount, pkScript)
		tool.CommitTxPrevOutputFetcher.AddPrevOut(*outPoint, txOut)

		in := wire.NewTxIn(outPoint, nil, nil)
		in.Sequence = DefaultSequenceNum
		tx.AddTxIn(in)

		totalSenderAmount += btcutil.Amount(prevOutput.Amount)
	}
	changePkScript, err := AddrToPkScript(changeAddress, tool.Network)
	if err != nil {
		return err
	}
	tx.AddTxOut(wire.NewTxOut(0, changePkScript))
	txForEstimate := wire.NewMsgTx(DefaultTxVersion)
	txForEstimate.TxIn = tx.TxIn
	txForEstimate.TxOut = tx.TxOut
	if err := Sign(txForEstimate, tool.CommitTxPrivateKeyList, tool.CommitTxPrevOutputFetcher); err != nil {
		return err
	}

	view, _ := commitTxPrevOutputList.UtxoViewpoint(tool.Network)
	vsize := GetTxVirtualSizeByView(btcutil.NewTx(txForEstimate), view)
	fee := btcutil.Amount(vsize) * btcutil.Amount(commitFeeRate)
	changeAmount := totalSenderAmount - btcutil.Amount(totalRevealPrevOutputValue) - fee
	if int64(changeAmount) >= minChangeValue {
		tx.TxOut[len(tx.TxOut)-1].Value = int64(changeAmount)
	} else {
		tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
		txForEstimate.TxOut = txForEstimate.TxOut[:len(txForEstimate.TxOut)-1]
		feeWithoutChange := btcutil.Amount(GetTxVirtualSizeByView(btcutil.NewTx(txForEstimate), view)) * btcutil.Amount(commitFeeRate)
		if totalSenderAmount-btcutil.Amount(totalRevealPrevOutputValue)-feeWithoutChange < 0 {
			tool.MustCommitTxFee = int64(btcutil.Amount(totalRevealPrevOutputValue) + fee)
			return errors.New("insufficient balance")
		}
	}
	tool.CommitTx = tx
	return nil
}

func GetSigOps(tx *btcutil.Tx, view UtxoViewpoint) (f int64) {
	defer func() {
		if r := recover(); r != nil {
			f = 0
		}
	}()
	sigops, err := GetSigOpCost(tx, false, view, true, true)
	if err != nil {
		return 0
	}
	return int64(sigops) * Five
}

func (tool *Src20InscriptionTool) signCommitTx() error {
	return Sign(tool.CommitTx, tool.CommitTxPrivateKeyList, tool.CommitTxPrevOutputFetcher)
}

func (tool *Src20InscriptionTool) GetCommitTxHex() (string, error) {
	return GetTxHex(tool.CommitTx)
}

func (tool *Src20InscriptionTool) CalculateFee() (int64, []int64) {
	commitTxFee := int64(0)
	for _, in := range tool.CommitTx.TxIn {
		commitTxFee += tool.CommitTxPrevOutputFetcher.FetchPrevOutput(in.PreviousOutPoint).Value
	}
	for _, out := range tool.CommitTx.TxOut {
		commitTxFee -= out.Value
	}
	return commitTxFee, make([]int64, 0)
}

func Src20Inscribe(network *chaincfg.Params, request *Src20InscriptionRequest) (*InscribeTxs, error) {
	tool, err := NewSrc20InscriptionTool(network, request)
	if err != nil && err.Error() == "insufficient balance" {
		return &InscribeTxs{
			CommitTx:    "",
			RevealTxs:   []string{},
			CommitTxFee: tool.MustCommitTxFee,
			CommitAddrs: tool.CommitAddrs,
		}, nil
	}

	if err != nil {
		return nil, err
	}

	commitTx, err := tool.GetCommitTxHex()
	if err != nil {
		return nil, err
	}

	commitTxFee, revealTxFees := tool.CalculateFee()

	return &InscribeTxs{
		CommitTx:     commitTx,
		CommitTxFee:  commitTxFee,
		RevealTxs:    make([]string, 0),
		RevealTxFees: revealTxFees,
		CommitAddrs:  tool.CommitAddrs,
	}, nil
}
