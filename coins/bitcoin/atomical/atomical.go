package atomical

import (
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"sort"
)

const (
	ErrCodeLessAtomicalAmt    = 2011400
	ErrCodeAtomicalChangeFail = 2011401
	ErrCodeVoutDust           = 2011402
	ErrCodeCommon             = 2011403
	ErrCodeUnknownAsset       = 2011404
	ErrInsufficientBalance    = 1000001
	ErrCodeMul                = 2011420
)

var (
	errUnknown            = &Err{ErrCode: ErrCodeCommon}
	ErInsufficientBalance = "Insufficient Balance"
)

var (
	NFT = "NFT"
	FT  = "FT"
)

type Err struct {
	ErrCode int         `json:"errCode"`
	Data    interface{} `json:"data"`
}

type ErrAtomicalIdAmt struct {
	AtomicalId string `json:"atomicalId"`
	Amount     int64  `json:"amount"`
}
type AtomicalData struct {
	AtomicalId string `json:"atomicalId"` // case：9527290d5f28479fa752f3eb9484ccbc5a951e2b2b5a49870318683e188e357ei0
	Type       string `json:"type"`       // FT | NFT
}

type AtomicalDatas []*AtomicalData

func (s AtomicalDatas) Sort() {
	if len(s) < 2 {
		return
	}
	sort.Sort(s)
}

func (s AtomicalDatas) Len() int { return len(s) }

func (x AtomicalDatas) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func (x AtomicalDatas) Less(i, j int) bool {
	return x[i].AtomicalId < x[j].AtomicalId
}

type PrevOutput struct {
	TxId       string        `json:"txId"`
	VOut       uint32        `json:"vOut"`
	Amount     int64         `json:"amount"`
	Address    string        `json:"address"`
	Data       AtomicalDatas `json:"data"`
	PrivateKey string        `json:"privateKey"`
}

type PrevOutputs []*PrevOutput

func (s PrevOutputs) UtxoViewpoint(net *chaincfg.Params) (bitcoin.UtxoViewpoint, error) {
	view := make(bitcoin.UtxoViewpoint, len(s))
	for _, v := range s {
		h, err := chainhash.NewHashFromStr(v.TxId)
		if err != nil {
			return nil, err
		}
		changePkScript, err := AddrToPkScript(v.Address, net)
		if err != nil {
			return nil, err
		}
		view[wire.OutPoint{Index: v.VOut, Hash: *h}] = changePkScript
	}
	return view, nil
}

type Output struct {
	Amount  int64         `json:"amount"`
	Address string        `json:"address"`
	Data    AtomicalDatas `json:"data"`
}

type AtomicalRequest struct {
	Inputs   PrevOutputs `json:"inputs"`
	Outputs  []*Output   `json:"outputs"`
	FeePerB  int64       `json:"feePerB"`
	Address  string      `json:"address"`
	DustSize int64       `json:"dustSize"`
}

func (a *AtomicalRequest) Check(network *chaincfg.Params, minChangeValue int64) *Err {
	if len(a.Inputs) == 0 || len(a.Outputs) == 0 {
		return &Err{ErrCode: ErrCodeCommon}
	}
	if _, err := btcutil.DecodeAddress(a.Address, network); err != nil {
		return &Err{ErrCode: ErrCodeCommon}
	}
	unspent, need := make(map[string]int64), make(map[string]int64)
	typs1, typs2 := make(map[string]struct{}), make(map[string]struct{})
	for _, v := range a.Inputs {
		if v == nil || v.Amount <= 0 {
			return &Err{ErrCode: ErrCodeCommon}
		}
		if _, err := btcutil.DecodeAddress(v.Address, network); err != nil {
			return &Err{ErrCode: ErrCodeCommon}
		}
		if len(v.Data) > 1 {
			return &Err{ErrCode: ErrCodeMul}
		}
		for _, v1 := range v.Data {
			if v1.Type != NFT && v1.Type != FT {
				return &Err{ErrCode: ErrCodeCommon}
			}
			if v1.Type == NFT {
				if _, ok := unspent[v1.AtomicalId]; ok {
					return &Err{ErrCode: ErrCodeCommon}
				}
				unspent[v1.AtomicalId] = 1
			} else {
				unspent[v1.AtomicalId] = unspent[v1.AtomicalId] + v.Amount
			}
			typs1[v1.Type] = struct{}{}
		}
	}
	if len(unspent) != 1 {
		return &Err{ErrCode: ErrCodeMul}
	}
	if len(typs1) != 1 {
		return &Err{ErrCode: ErrCodeMul}
	}
	for _, v := range a.Outputs {
		if v == nil || v.Amount <= 0 {
			return &Err{ErrCode: ErrCodeCommon}
		}
		if len(v.Data) == 0 {
			return &Err{ErrCode: ErrCodeCommon}
		}
		if len(v.Data) > 1 {
			return &Err{ErrCode: ErrCodeMul}
		}
		if _, err := btcutil.DecodeAddress(v.Address, network); err != nil {
			return &Err{ErrCode: ErrCodeCommon}
		}
		if v.Amount < minChangeValue {
			return &Err{ErrCode: ErrCodeVoutDust, Data: &ErrAtomicalIdAmt{
				AtomicalId: v.Data[0].AtomicalId, Amount: v.Amount,
			}}
		}
		for _, v1 := range v.Data {
			if v1.Type != NFT && v1.Type != FT {
				return &Err{ErrCode: ErrCodeCommon}
			}
			if v1.Type == NFT {
				if _, ok := need[v1.AtomicalId]; ok {
					return &Err{ErrCode: ErrCodeCommon}
				}
				need[v1.AtomicalId] = 1
			} else {
				need[v1.AtomicalId] = need[v1.AtomicalId] + v.Amount
			}
			typs2[v1.Type] = struct{}{}
		}
	}
	if len(need) != 1 {
		return &Err{ErrCode: ErrCodeMul}
	}
	if len(typs2) != 1 {
		return &Err{ErrCode: ErrCodeMul}
	}
	for k, _ := range typs2 {
		if _, ok := typs1[k]; !ok {
			return &Err{ErrCode: ErrCodeMul}
		}
	}
	for k, v := range need {
		if unspent[k] == 0 {
			return &Err{ErrCode: ErrCodeUnknownAsset}
		}
		if unspent[k] < v {
			return &Err{ErrCode: ErrCodeLessAtomicalAmt, Data: &ErrAtomicalIdAmt{
				AtomicalId: k, Amount: unspent[k] - v,
			}}
		}
		unspent[k] = unspent[k] - v
	}
	for k, v := range unspent {
		if v > 0 && v < minChangeValue {
			return &Err{ErrCode: ErrCodeAtomicalChangeFail, Data: &ErrAtomicalIdAmt{
				AtomicalId: k, Amount: v,
			}}
		}
	}
	return nil
}

type AtomicalTransferTool struct {
	CommitTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrivateKeyList    []*btcec.PrivateKey
	RevealTxPrevOutputFetcher *txscript.MultiPrevOutFetcher
	CommitTxPrevOutputList    []*PrevOutput
	CommitTx                  *wire.MsgTx
	MustCommitTxFee           int64
	CommitAddrs               []string
}

func NewAtomicalTransferTool(network *chaincfg.Params, request *AtomicalRequest) (*AtomicalTransferTool, *Err) {
	var commitTxPrivateKeyList []*btcec.PrivateKey

	for _, v := range request.Inputs {
		privateKeyWif, err := btcutil.DecodeWIF(v.PrivateKey)
		if err != nil {
			return nil, &Err{ErrCode: ErrCodeCommon}
		}
		commitTxPrivateKeyList = append(commitTxPrivateKeyList, privateKeyWif.PrivKey)
	}
	tool := &AtomicalTransferTool{
		CommitTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrivateKeyList:    commitTxPrivateKeyList,
		RevealTxPrevOutputFetcher: txscript.NewMultiPrevOutFetcher(nil),
		CommitTxPrevOutputList:    request.Inputs,
	}
	return tool, tool._initTool(network, request)
}

func (tool *AtomicalTransferTool) _initTool(network *chaincfg.Params, request *AtomicalRequest) *Err {
	minChangeValue := bitcoin.DefaultMinChangeValue
	if request.DustSize > 0 {
		minChangeValue = request.DustSize
	}
	if err := request.Check(network, minChangeValue); err != nil {
		return err
	}
	err := tool.buildCommitTx(network, request, minChangeValue)
	if err != nil {
		if err.Error() == ErInsufficientBalance {
			return &Err{ErrCode: ErrInsufficientBalance}
		}
		return errUnknown
	}
	err = tool.signCommitTx()
	if err != nil {
		return errUnknown
	}
	return nil
}

func (tool *AtomicalTransferTool) buildCommitTx(network *chaincfg.Params, request *AtomicalRequest, minChangeValue int64) error {
	totalSenderAmount := btcutil.Amount(0)
	tx := wire.NewMsgTx(bitcoin.DefaultTxVersion)
	unspent := make(map[string]int64)
	types := make(map[string]string)
	changePkScript, err := bitcoin.AddrToPkScript(request.Address, network)
	if err != nil {
		return err
	}
	for _, prevOutput := range request.Inputs {
		txHash, err := chainhash.NewHashFromStr(prevOutput.TxId)
		if err != nil {
			return err
		}
		outPoint := wire.NewOutPoint(txHash, prevOutput.VOut)
		in := wire.NewTxIn(outPoint, nil, nil)
		in.Sequence = bitcoin.DefaultSequenceNum
		tx.AddTxIn(in)
		pkScript, err := bitcoin.AddrToPkScript(prevOutput.Address, network)
		if err != nil {
			return err
		}
		if len(prevOutput.Data) == 0 {
			totalSenderAmount += btcutil.Amount(prevOutput.Amount)
		} else {
			for _, v1 := range prevOutput.Data {
				types[v1.AtomicalId] = v1.Type
				if v1.Type == FT {
					unspent[v1.AtomicalId] = unspent[v1.AtomicalId] + prevOutput.Amount
				} else {
					unspent[v1.AtomicalId] = 1
					totalSenderAmount += btcutil.Amount(prevOutput.Amount)
				}
			}
		}
		txOut := wire.NewTxOut(prevOutput.Amount, pkScript)
		tool.CommitTxPrevOutputFetcher.AddPrevOut(*outPoint, txOut)
	}
	cost := btcutil.Amount(0)
	for k, v := range unspent {
		total := v
		for _, v := range request.Outputs {
			addrPkScript, err := bitcoin.AddrToPkScript(v.Address, network)
			if err != nil {
				return err
			}
			tx.AddTxOut(wire.NewTxOut(v.Amount, addrPkScript))
			if len(v.Data) == 0 || v.Data[0].Type == NFT {
				cost += btcutil.Amount(v.Amount)
			}
			c := v.Amount
			if v.Data[0].Type == NFT {
				c = 1
			}
			if total < c {
				return errors.New(ErInsufficientBalance)
			}
			total -= c
		}
		if types[k] == FT {
			if total > minChangeValue {
				tx.AddTxOut(wire.NewTxOut(total, changePkScript))
			}
		} else {
			if total > 0 {
				cost += btcutil.Amount(minChangeValue)
				tx.AddTxOut(wire.NewTxOut(minChangeValue, changePkScript))
			}
		}
	}
	tx.AddTxOut(wire.NewTxOut(0, changePkScript))
	txForEstimate := wire.NewMsgTx(bitcoin.DefaultTxVersion)
	txForEstimate.TxIn = tx.TxIn
	txForEstimate.TxOut = tx.TxOut
	view, _ := request.Inputs.UtxoViewpoint(network)
	if err := bitcoin.Sign(txForEstimate, tool.CommitTxPrivateKeyList, tool.CommitTxPrevOutputFetcher); err != nil {
		return err
	}
	fee := btcutil.Amount(bitcoin.GetTxVirtualSizeByView(btcutil.NewTx(txForEstimate), view)) * btcutil.Amount(request.FeePerB)
	changeAmount := totalSenderAmount - cost - fee
	if int64(changeAmount) >= minChangeValue {
		tx.TxOut[len(tx.TxOut)-1].Value = int64(changeAmount)
	} else {
		tx.TxOut = tx.TxOut[:len(tx.TxOut)-1]
		txForEstimate.TxOut = txForEstimate.TxOut[:len(txForEstimate.TxOut)-1]
		feeWithoutChange := btcutil.Amount(bitcoin.GetTxVirtualSizeByView(btcutil.NewTx(txForEstimate), view)) * btcutil.Amount(request.FeePerB)
		if totalSenderAmount-cost-feeWithoutChange < 0 {
			tool.MustCommitTxFee = int64(cost + fee)
			return errors.New(ErInsufficientBalance)
		}
	}
	tool.CommitTx = tx
	return nil
}

func (tool *AtomicalTransferTool) signCommitTx() error {
	return bitcoin.Sign(tool.CommitTx, tool.CommitTxPrivateKeyList, tool.CommitTxPrevOutputFetcher)
}

func (tool *AtomicalTransferTool) GetCommitTxHex() (string, error) {
	return bitcoin.GetTxHex(tool.CommitTx)
}

func (tool *AtomicalTransferTool) CalculateFee() (int64, []int64) {
	commitTxFee := int64(0)
	for _, in := range tool.CommitTx.TxIn {
		commitTxFee += tool.CommitTxPrevOutputFetcher.FetchPrevOutput(in.PreviousOutPoint).Value
	}
	for _, out := range tool.CommitTx.TxOut {
		commitTxFee -= out.Value
	}
	return commitTxFee, make([]int64, 0)
}

func AtomicalTransfer(network *chaincfg.Params, request *AtomicalRequest) (*bitcoin.InscribeTxs, *Err) {
	tool, err := NewAtomicalTransferTool(network, request)
	if err != nil && err.ErrCode == ErrInsufficientBalance {
		return &bitcoin.InscribeTxs{
			CommitTx:    "",
			RevealTxs:   []string{},
			CommitTxFee: tool.MustCommitTxFee,
			CommitAddrs: tool.CommitAddrs,
		}, err
	}

	if err != nil {
		return nil, err
	}

	commitTx, e := tool.GetCommitTxHex()
	if e != nil {
		return nil, errUnknown
	}

	commitTxFee, revealTxFees := tool.CalculateFee()

	return &bitcoin.InscribeTxs{
		CommitTx:     commitTx,
		CommitTxFee:  commitTxFee,
		RevealTxFees: revealTxFees,
		CommitAddrs:  tool.CommitAddrs,
	}, nil
}
