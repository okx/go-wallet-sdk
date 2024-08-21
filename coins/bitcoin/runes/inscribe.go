package bitcoin

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"github.com/btcsuite/btcd/btcec/v2"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/okx/go-wallet-sdk/coins/bitcoin"
	"math/big"
	"sort"
	"strings"
)

var (
	ErrInvalidEdict = errors.New("ErrInvalidEdict")
)
var (
	TAG_BODY         = big.NewInt(0)
	TAG_Flags        = big.NewInt(2)
	TAG_Rune         = big.NewInt(4)
	TAG_Premine      = big.NewInt(6)
	TAG_Cap          = big.NewInt(8)
	TAG_Amount       = big.NewInt(10)
	TAG_HeightStart  = big.NewInt(12)
	TAG_HeightEnd    = big.NewInt(14)
	TAG_OffsetStart  = big.NewInt(16)
	TAG_OffsetEnd    = big.NewInt(118)
	TAG_Mint         = big.NewInt(20)
	TAG_Pointer      = big.NewInt(22)
	TAG_Cenotaph     = big.NewInt(126)
	TAG_Divisibility = big.NewInt(1)
	TAG_Spacers      = big.NewInt(3)
	TAG_Symbol       = big.NewInt(5)
	TAG_Nop          = big.NewInt(127)
)

type InscriptionData struct {
	ContentType string `json:"contentType,omitempty"`
	Body        []byte `json:"body"`
	RevealAddr  string `json:"revealAddr,omitempty"`
	Vout        uint   `json:"vout"` // default 1
}
type Output struct {
	Address string `json:"address"`
	Amount  int64  `json:"amount"`
}
type OpReturnData struct {
	Edicts          Edicts   `json:"edicts"`
	Etching         *Etching `json:"etching"`
	Burn            bool     `json:"burn"`
	IsDefaultOutput bool     `json:"isDefaultOutput"`
	DefaultOutput   int64    `json:"defaultOutput"`
	Mint            bool     `json:"mint"`
	MintNum         uint     `json:"mintNum"`
}

type Edict struct {
	Block  string `json:"block"`
	Id     string `json:"id"`
	Amount string `json:"amount"`
	Output uint   `json:"output"`
}
type Edicts []*Edict

func (s Edicts) RealEdicts(mint bool) (RealEdicts, error) {
	ss := make(RealEdicts, len(s))
	for k, v := range s {
		r, err := v.RealEdict(mint)
		if err != nil {
			return nil, err
		}
		ss[k] = r
	}
	return ss, nil
}

func (v Edicts) Fix() {
	for _, v := range v {
		v.Fix()
	}
}

func (v *Edict) Fix() {
	if v == nil {
		return
	}
	if len(v.Block) > 0 {
		return
	}
	param := strings.Split(v.Id, ":")
	v.Block = param[0]
	if len(param) > 1 {
		v.Id = param[1]
	}
}

func (v *Edict) RealEdict(mint bool) (*RealEdict, error) {
	b, ok := new(big.Int).SetString(v.Block, 10)
	if !ok {
		return nil, ErrInvalidEdict
	}
	id, ok := new(big.Int).SetString(v.Id, 10)
	if !ok {
		return nil, ErrInvalidEdict
	}

	amount, ok := new(big.Int).SetString(v.Amount, 10)
	if !ok && !mint {
		return nil, ErrInvalidEdict
	}
	if !mint && amount.Cmp(zero) <= 0 {
		return nil, ErrInvalidEdict
	}
	if (amount == nil || amount.Cmp(zero) <= 0) && mint {
		amount = big.NewInt(0)
	}
	return &RealEdict{
		Block:  b,
		Id:     id,
		Amount: amount,
		Output: v.Output,
	}, nil
}

type RealEdict struct {
	Block  *big.Int `json:"block"`
	Id     *big.Int `json:"id"`
	Amount *big.Int `json:"amount"`
	Output uint     `json:"output"`
}

type RealEdicts []*RealEdict

func (s RealEdicts) Less(i, j int) bool {
	if s[i].Block.Cmp(s[j].Block) == 0 {
		return s[i].Id.Cmp(s[j].Id) <= 0
	}
	return s[i].Block.Cmp(s[j].Block) <= 0
}

func (s RealEdicts) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s RealEdicts) Len() int {
	return len(s)
}

type Etching struct {
	Divisibility int64  `json:"divisibility"`
	Limit        string `json:"limit"`
	Rune         string `json:"rune"`
	Symbol       string `json:"symbol"`
	Term         int64  `json:"term"`
}

type PrevOutput struct {
	TxId       string `json:"txId"`
	VOut       uint32 `json:"vOut"`
	Amount     int64  `json:"amount"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

type PrevOutputs []*PrevOutput

func (s PrevOutputs) UtxoViewpoint(net *chaincfg.Params) (bitcoin.UtxoViewpoint, error) {
	view := make(bitcoin.UtxoViewpoint, len(s))
	for _, v := range s {
		h, err := chainhash.NewHashFromStr(v.TxId)
		if err != nil {
			return nil, err
		}
		changePkScript, err := bitcoin.AddrToPkScript(v.Address, net)
		if err != nil {
			return nil, err
		}
		view[wire.OutPoint{Index: v.VOut, Hash: *h}] = changePkScript
	}
	return view, nil
}

type InscriptionRequest struct {
	CommitTxPrevOutputList []*PrevOutput     `json:"commitTxPrevOutputList"`
	Outputs                []*Output         `json:"outputs"`
	CommitFeeRate          int64             `json:"commitFeeRate"`
	RevealFeeRate          int64             `json:"revealFeeRate"`
	InscriptionDataList    []InscriptionData `json:"inscriptionDataList"`
	RevealOutValue         int64             `json:"revealOutValue"`
	ChangeAddress          string            `json:"changeAddress"`
	MinChangeValue         int64             `json:"minChangeValue"`
	isMainnet              bool              `json:"isMainnet"`
}

type InscribeTxs struct {
	CommitTx     string   `json:"commitTx"`
	RevealTxs    []string `json:"revealTxs,omitempty"`
	CommitTxFee  int64    `json:"commitTxFee"`
	RevealTxFees []int64  `json:"revealTxFees,omitempty"`
	CommitAddrs  []string `json:"commitAddrs,omitempty"`
}

const (
	DefaultTxVersion      = 2
	DefaultSequenceNum    = 0xfffffffd
	DefaultRevealOutValue = int64(546)
	DefaultMinChangeValue = int64(546)

	MaxStandardTxWeight = 4000000 / 10
	WitnessScaleFactor  = 4
)

// todo
func CheckOpReturnData(data *OpReturnData) bool {

	return true
}

var (
	zero = big.NewInt(0)
	x7F  = big.NewInt(0x7F) //127
	x80  = big.NewInt(0x80) //128
)

func EncodeToVecV2(n *big.Int, buf *bytes.Buffer) {
	for n2 := new(big.Int).Rsh(n, 7); n2.Cmp(zero) > 0; {
		buf.Write(new(big.Int).Or(new(big.Int).And(n, x7F), x80).Bytes())
		n = n2
		n2 = new(big.Int).Rsh(n, 7)
	}
	v := new(big.Int).And(n, x7F)
	if r := v.Bytes(); len(r) > 0 {
		buf.Write(r)
	} else {
		buf.WriteByte(0)
	}
}

func EncodeToVec(n *big.Int) []int64 {
	i := 18
	out := make([]int64, 19)
	m := new(big.Int).SetInt64(0x7F)
	out[i] = (new(big.Int).And(n, m)).Int64()
	k := new(big.Int).SetInt64(128)
	one := new(big.Int).SetInt64(1)
	ff := new(big.Int).SetInt64(0xFF)
	for n.Cmp(m) > 0 {
		n.Div(n, k)
		n.Sub(n, one)
		i--
		out[i] = new(big.Int).And(n, ff).Int64() | 0x80
	}
	return out[i:]
}

func BuildOpReturnDataJson(a []byte) ([]byte, error) {
	o := &OpReturnData{}
	err := json.Unmarshal(a, &o)
	if err != nil {
		return nil, err
	}
	return BuildOpReturnData(o.Edicts, o.IsDefaultOutput, o.Mint, o.DefaultOutput)
}

func BuildOpReturnData(edicts Edicts, isDefaultOutput, mint bool, defaultOutput int64) ([]byte, error) {
	payload := &bytes.Buffer{}
	if len(edicts) == 0 {
		return nil, ErrInvalidEdict
	}
	es, err := edicts.RealEdicts(mint)
	if err != nil {
		return nil, ErrInvalidEdict
	}
	if mint && len(edicts) != 0 && edicts[0].Block != "" {
		EncodeToVecV2(TAG_Mint, payload)
		EncodeToVecV2(es[0].Block, payload) // only mint edicts[0].id
		EncodeToVecV2(TAG_Mint, payload)
		EncodeToVecV2(es[0].Id, payload) // only mint edicts[0].id
	}
	if isDefaultOutput {
		EncodeToVecV2(TAG_Pointer, payload)
		EncodeToVecV2(big.NewInt(defaultOutput), payload)
	}

	if es.Len() > 0 && !mint {
		EncodeToVecV2(TAG_BODY, payload)
		sort.Sort(es)
		id := big.NewInt(0)
		block := big.NewInt(0)
		for _, edict := range es {
			EncodeToVecV2(new(big.Int).Sub(edict.Block, block), payload)
			EncodeToVecV2(new(big.Int).Sub(edict.Id, id), payload)
			EncodeToVecV2(edict.Amount, payload)
			EncodeToVecV2(big.NewInt(int64(edict.Output)), payload)
			id = edict.Id
			block = edict.Block
		}
	}
	data := payload.Bytes()
	if len(data) > 80 {
		return nil, errors.New("The script is too long")
	}
	//fmt.Println(hex.EncodeToString(payload.Bytes()))
	builder := txscript.NewScriptBuilder()
	builder.AddOp(txscript.OP_RETURN).AddOp(txscript.OP_13).AddData(data)
	return builder.Script()
}

func Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	index := 0
	for j := 0; j < 8; j++ {
		if buf[j] != 0 {
			index = j
			break
		}
		if j >= 7 {
			index = 7
		}
	}
	return buf[index:]
}

func BuildOpReturnScriptForEdict(edicts []*Edict, isMainnet bool) ([]byte, error) {

	tagBody := new(big.Int).SetInt64(0)
	payload := []int64{}

	es := edicts
	if len(es) > 0 {
		payload = append(payload, EncodeToVec(tagBody)...)
	}

	sort.Slice(es, func(i, j int) bool {
		// return es[i].Id < es[j].Id
		a := es[i]
		b := es[j]
		idA, ok := new(big.Int).SetString(a.Id, 16)
		if !ok {
			return es[i].Id < es[j].Id
		}
		idB, ok := new(big.Int).SetString(b.Id, 16)
		if !ok {
			return es[i].Id < es[j].Id
		}
		return idA.Cmp(idB) <= 0
	})
	id := new(big.Int).SetInt64(0)
	for _, e := range es {
		idB, ok := new(big.Int).SetString(e.Id, 16)
		if !ok {
			return nil, errors.New("invalid edict id")
		}
		r := new(big.Int).Sub(idB, id)
		payload = append(payload, EncodeToVec(r)...)
		amountB, _ := new(big.Int).SetString(e.Amount, 10)
		payload = append(payload, EncodeToVec(amountB)...)
		output := new(big.Int).SetUint64(uint64(e.Output))
		payload = append(payload, EncodeToVec(output)...)
		id = idB
	}
	prefix := "R"
	if !isMainnet {
		prefix = "RUNE_TEST"
	}
	buf := &bytes.Buffer{}
	for _, v := range payload {
		if v < 0 {
			return nil, errors.New("invalid number")
		}
		buf.Write(Int64ToBytes(v))
	}

	inscriptionBuilder := txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData([]byte(prefix)).AddData(buf.Bytes())
	inscriptionScript, err := inscriptionBuilder.Script()
	if err != nil {
		return nil, err
	}
	return inscriptionScript, nil

	//return txscript.NullDataScript(buf.Bytes())
}

func BuildOpReturnScriptForRuneMainEdict(
	edicts []*Edict, isMainnet bool, isDefaultOutput bool, defaultOutput int64) ([]byte, error) {

	tagBody := new(big.Int).SetInt64(0)
	tagDefaultOupt := new(big.Int).SetInt64(12)
	payload := []int64{}
	// Adjustments may be necessary, depending on whether ordi defaults to 0 for transfer
	if isDefaultOutput {
		payload = append(payload, EncodeToVec(tagDefaultOupt)...)
		payload = append(payload, EncodeToVec(new(big.Int).SetInt64(defaultOutput))...)
	}

	es := edicts
	if len(es) > 0 {
		payload = append(payload, EncodeToVec(tagBody)...)
	}

	sort.Slice(es, func(i, j int) bool {
		// return es[i].Id < es[j].Id
		a := es[i]
		b := es[j]
		idA, ok := new(big.Int).SetString(a.Id, 16)
		if !ok {
			return es[i].Id < es[j].Id
		}
		idB, ok := new(big.Int).SetString(b.Id, 16)
		if !ok {
			return es[i].Id < es[j].Id
		}
		return idA.Cmp(idB) <= 0
	})
	id := new(big.Int).SetInt64(0)
	for _, e := range es {
		idB, ok := new(big.Int).SetString(e.Id, 16)
		if !ok {
			return nil, errors.New("invalid edict id")
		}
		r := new(big.Int).Sub(idB, id)
		payload = append(payload, EncodeToVec(r)...)
		amountB, _ := new(big.Int).SetString(e.Amount, 10)
		payload = append(payload, EncodeToVec(amountB)...)
		output := new(big.Int).SetUint64(uint64(e.Output))
		payload = append(payload, EncodeToVec(output)...)
		id = idB
	}
	prefix := "RUNE_TEST" // Currently even the Bitcoin mainnet only recognizes RUNE_TEST
	if !isMainnet {
		prefix = "RUNE_TEST"
	}
	buf := &bytes.Buffer{}
	for _, v := range payload {
		if v < 0 {
			return nil, errors.New("invalid number")
		}
		buf.Write(Int64ToBytes(v))
	}

	if len(buf.Bytes()) > 80 {
		return nil, errors.New("op-return payload exceeds length limit 80 ")
	}

	inscriptionBuilder := txscript.NewScriptBuilder().AddOp(txscript.OP_RETURN).AddData([]byte(prefix)).AddData(buf.Bytes())

	inscriptionScript, err := inscriptionBuilder.Script()
	if err != nil {
		return nil, err
	}
	return inscriptionScript, nil

	//return txscript.NullDataScript(buf.Bytes())
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
			script, err := bitcoin.PayToPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
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
				redeemScript, err := bitcoin.PayToWitnessPubKeyHashScript(btcutil.Hash160(pubKeyBytes))
				if err != nil {
					return err
				}
				in.SignatureScript = append([]byte{byte(len(redeemScript))}, redeemScript...)
			}
		}
	}

	return nil
}

func CalculateMintTxFee(receiver string, addr string, mintScript string, privateKey string, netParams *chaincfg.Params, feePerB int64, value int64) (int64, int64, error) {
	preTxBuild := bitcoin.NewTxBuild(2, netParams)
	preTxBuild.AddInput2("9f9ff5acc7b3966ccfc6acc77027209d62aab34e563a09180c58ef7296fca74b", 0, privateKey, addr, 0)
	preTxBuild.AddOutput(receiver, value)
	preTxBuild.AddOutput2("", mintScript, 0)
	tx, err := preTxBuild.Build()
	if err != nil {
		return 0, 0, err
	}
	vsize := bitcoin.GetTxVirtualSize(btcutil.NewTx(tx))
	return vsize * feePerB, vsize*feePerB + value, nil
}

func BuildMintL2Tx(txid string, vout uint32, amount int64, receiver, addr string, mintScript string, privateKey string, netParams *chaincfg.Params, feePerB int64, value int64) (*wire.MsgTx, error) {
	preTxBuild := bitcoin.NewTxBuild(2, netParams)
	preTxBuild.AddInput2(txid, vout, privateKey, addr, amount)
	preTxBuild.AddOutput(receiver, value)
	preTxBuild.AddOutput2("", mintScript, 0)
	return preTxBuild.Build()
}
